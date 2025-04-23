package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"sync"
	"time"

	param "UDP2TCP/struct"
)

// 全局声明一个消息队列，负责检测发送消息的超时情况
// 添加同步互斥锁，为了处理多线程竞争同一资源读写问题
var (
	congestionWindow []param.Message
	lock             sync.Mutex
	wg               sync.WaitGroup // 用于等待所有goroutine完成
)

func main() {
	//启动超时检测
	go startTimeoutChecker()

	//启动响应接收端口
	go startResponseServer()

	// 发送消息
	for i := 0; i < 10; i++ {
		wg.Add(1)
		msg := send("Hello, server!", i)
		congestionWindow = append(congestionWindow, msg)
	}

	// 等待所有goroutine完成
	wg.Wait()
	// 发送结束消息 并确保对方收到结束消息
	wg.Add(1)
	msg := send("Compose, server!", 10)
	congestionWindow = append(congestionWindow, msg)
	wg.Wait()
	log.Printf("Send Finished")
}

// 发送消息，通过独立的监听器接收响应
func send(content string, id int) param.Message {
	//随机生成一个端口
	clientAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 0,
	}

	// 目标端口地址
	targetAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 9050,
	}

	conn, err := net.ListenUDP("udp", clientAddr)
	if err != nil {
		log.Fatal(err)
	}

	// 创建消息结构体
	msg := param.Message{
		ID:      id,
		Content: content,
		Time:    time.Now().Unix(),
	}

	// 序列化消息
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(msg)
	if err != nil {
		log.Fatal(err)
	}

	// 发送消息
	_, err = conn.WriteToUDP(buffer.Bytes(), targetAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sent message: %+v\n", msg)

	return msg
}

// 启动响应接收端口
// 监听9051端口，接收响应消息
func startResponseServer() {
	responseAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 9051,
	}

	conn, err := net.ListenUDP("udp", responseAddr)
	if err != nil {
		log.Fatal(err)
	}

	//defer保证函数执行结束后释放资源
	defer conn.Close()

	log.Printf("Response server listening on %v\n", responseAddr)

	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		var receivedMsg param.Message
		dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
		err = dec.Decode(&receivedMsg)
		if err != nil {
			log.Printf("Error decoding message: %v\n", err)
			continue
		}
		log.Printf("Received message from %v:\n %+v\n", addr, receivedMsg)
		lock.Lock()

		// 处理接收到的消息 （待修改）
		var tempCongestionWindow []param.Message
		for _, msg := range congestionWindow {
			if msg.ID == receivedMsg.ID {
				log.Printf("Received response for message %d\n", msg.ID)
				// 删除已确认的消息
				wg.Done()
			} else {
				tempCongestionWindow = append(tempCongestionWindow, msg)
			}
		}
		congestionWindow = tempCongestionWindow
		lock.Unlock()
	}
}

// 启动超时检测器，定时检查消息是否超时
// 超时消息会被重新发送
func startTimeoutChecker() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	//循环监听 ticker.C 通道，通道每隔两秒触发一次，循环执行一次
	for range ticker.C {
		lock.Lock()

		log.Printf("Congestion window: %v\n", congestionWindow)

		currentTime := time.Now().Unix()

		tempCongestionWindow := []param.Message{}

		for _, msg := range congestionWindow {
			if currentTime-msg.Time > 2 {
				log.Printf("Message %d timeout, resend\n", msg.ID)
				sendMsg := send(msg.Content, msg.ID)
				tempCongestionWindow = append(tempCongestionWindow, sendMsg)
			} else {
				// 保留未超时的消息
				tempCongestionWindow = append(tempCongestionWindow, msg)
			}
		}

		congestionWindow = tempCongestionWindow

		lock.Unlock()
	}
}
