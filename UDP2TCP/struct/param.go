package param

import (
	"net"
	"sync"
	"time"
)

// Message 定义了客户端和服务器之间传输的消息结构
// 其余的都是ai生成的，被我重写废弃了
type Message struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type ResponseListener struct {
	Conn       *net.UDPConn
	Responses  chan Message
	MsgQueue   *MessageQueue
	ServerAddr *net.UDPAddr
}

type MessageQueue struct {
	Mutex    sync.Mutex
	Messages map[int]*QueueItem
}

type QueueItem struct {
	Msg      Message
	SendTime time.Time
	Retries  int
}
