package main

import (
	message "UDP2TCP/struct"
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"net"
	"time"
	"sync"
)

var(
	congestionWindow []message.Message
	lock             sync.Mutex
	wq               sync.WaitGroup
	windowCounter    int 		//创建一个拥塞窗口计数器，确保顺序交付
)

func main() {
	wq.Add(1)
	windowCounter = 0

	go startUDPServer()
	go startWindowReader()

	wq.Wait()
}

func startUDPServer() {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 9050,
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("UDP server started on :9050")
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		// 反序列化接收到的消息
		var receivedMsg message.Message
		dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
		err = dec.Decode(&receivedMsg)
		if err != nil {
			log.Printf("Error decoding message: %v\n", err)
			continue
		}
		log.Printf("Received message from %v:\n %+v\n", addr, receivedMsg)

		// 10%概率丢包
		if rand.Float64() < 0.2 {
			log.Printf("Simulated packet loss, skipping response")
			continue
		}

		lock.Lock()

		//模拟拥塞窗口固定长度，超出长度则丢弃
		if len(congestionWindow) >= 5 {
			log.Printf("Congestion window full, skipping message %d\n", receivedMsg.ID)
			lock.Unlock()
			continue
		}

		// 模拟顺序交付
		if windowCounter == receivedMsg.ID {
			// 清空拥塞窗口, 并重置计数器
			if receivedMsg.Content == "Compose, server!" {
				log.Printf("Clear congestion window")
				windowCounter = 0
				congestionWindow = []message.Message{}
			} else {
				log.Printf("Delivering message %d\n", receivedMsg.ID)
				congestionWindow = append(congestionWindow, receivedMsg)
				windowCounter++
				log.Printf("Congestion window: %v\n", congestionWindow)
			}

		} else {
			log.Printf("Skipping message %d\n", receivedMsg.ID)
			lock.Unlock()
			continue
		}

		lock.Unlock()

		// 创建响应消息
		responseMsg := message.Message{
			ID:      receivedMsg.ID,
			Content: "Hello, client!",
			Time:    time.Now().Unix(),
		}

		// 序列化响应消息
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(responseMsg)
		if err != nil {
			log.Printf("Error encoding response: %v\n", err)
			continue
		}

		responseAddr := &net.UDPAddr{
			IP:   addr.IP,
			Port: 9051,
		}

		_, err = conn.WriteToUDP(buffer.Bytes(), responseAddr)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Sent response to %v\n", responseAddr)
	}
}

// 模拟读取窗口内数据
func startWindowReader() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C{
		lock.Lock()
		if len(congestionWindow) > 0 {
			congestionWindow = congestionWindow[1:]
			
		}
		lock.Unlock()
	}
}
