package znet

import (
	"fmt"
	"log"
	"my-zinx/zinx/ziface"
	"net"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	Name      string
	IPVersion string
	Ip        string
	Port      uint16
	Router    ziface.IRouter

	maxConnCount uint // 最大连接数
}

func NewServer(name string, port uint16) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		Ip:        "0.0.0.0",
		Port:      port,
		// Router:       make([]ziface.IRouter, 0),
		Router:       nil,
		maxConnCount: 100,
	}
}

func (s *Server) Start() {
	log.Println("Server Start")
	doStart(s)
}

func (s *Server) Serve() {
	log.Println("Server Serve")
	doServe(s)
}

func (s *Server) Stop() {
	log.Println("Server Stop")
	doStop(s)
}

func (s *Server) AddRouter(router ziface.IRouter) {
	// s.Router = append(s.Router, router)
	s.Router = router
}

// 确保 Server 实现了 ziface.IServer 的所有方法（让编译器帮我们检查）
var _ ziface.IServer = (*Server)(nil)

func doStart(s *Server) {
	endpoint, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Println("ResolveTCPAddr error:", err)
		return
	}
	listener, err := net.ListenTCP(s.IPVersion, endpoint)
	if err != nil {
		log.Println("ListenTCP error:", err)
		return
	}
	log.Println("Start Zinx server success, ", s.Name, " Listening...")

	// 启用单独的协程来处理客户端连接
	// 这是go语言的风格，能用异步一般用异步。这样主协程接下来还可以做其他工作，比如后面的Serve()方法
	go func() {
		var connCount uint = 0
		for {
			peer, err := listener.AcceptTCP()
			if err != nil {
				log.Println("AcceptTCP error:", err)
				continue
			}
			connCount++
			if connCount > s.maxConnCount {
				log.Println("Too many connections, close this new connection")
				connCount--
				peer.Close()
				continue
			}

			// TODO 处理该新连接请求的业务方法, 此时应该有 handler, 它和 conn 是绑定的
			// TODO 也许这个处理业务的dealConn服务端应该记录在map中，不然的话ConnID()没有生成的意义
			dealConn := NewConnection(peer, uuid.New(), s.Router)
			// 启动子协程处理业务
			go dealConn.Open()
		}
	}()
}

func doServe(s *Server) {
	// TODO 是否在启动服务的时候, 还需要做其它事情呢? 比如自定义 Logger 或加入鉴权中间件等
	// 阻塞, 否则 main goroutine 退出, listenner 也将会随之退出
	for {
		time.Sleep(10 * time.Second)
	}
}

func doStop(s *Server) {
	// TODO 将其它需要清理的连接信息或其他信息一并停止或清理
}
