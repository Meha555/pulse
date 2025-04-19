package znet

import (
	"fmt"
	"log"
	"my-zinx/zinx/ziface"
	"my-zinx/zinx/zutils"
	"net"
	"time"
)

type Server struct {
	Name      string
	IPVersion string
	Ip        string
	Port      uint16
	// 映射请求到具体的API回调
	ApiMapper ziface.IApiMapper
	// 工作协程池
	wokerPool *WokerPool

	maxConnCount uint // 最大连接数
}

func NewServer() *Server {
	// 消息队列（worker协程从中取数据）mq容量和worker数量相同。mq容量更大没意义
	mq := NewMsgQueue(int(zutils.Conf.Server.MaxWorkerPoolSize))
	mapper := NewApiMapper()
	return &Server{
		Name:         zutils.Conf.Server.Name,
		IPVersion:    "tcp4",
		Ip:           zutils.Conf.Server.Host,
		Port:         zutils.Conf.Server.Port,
		ApiMapper:    mapper,
		wokerPool:    NewWokerPool(mq.Cap(), mq, mapper),
		maxConnCount: zutils.Conf.Server.MaxConn,
	}
}

func (s *Server) Listen() {
	log.Printf("Server Start with config: %+v\n", zutils.Conf)
	doListern(s)
}

func (s *Server) Serve() {
	log.Println("Server Serve")
	doServe(s)
}

func (s *Server) Shutdown() {
	log.Println("Server Stop")
	doShutdown(s)
}

func (s *Server) ListenAndServe() {
	s.Listen()
	s.Serve()
}

// 确保 Server 实现了 ziface.IServer 的所有方法（让编译器帮我们检查）
var _ ziface.IServer = (*Server)(nil)

func doListern(s *Server) {
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
	log.Println(s.Name, " Listening...")

	// 启动协程池
	s.wokerPool.Start()

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
			dealConn := NewConnection(peer, s.wokerPool)
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

func doShutdown(s *Server) {
	// TODO 将其它需要清理的连接信息或其他信息一并停止或清理

	s.wokerPool.Stop()
}
