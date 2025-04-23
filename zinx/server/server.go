package server

import (
	"context"
	"fmt"
	"log"
	"my-zinx/zinx/core/job"
	"my-zinx/zinx/core/message"
	"my-zinx/zinx/core/session"
	iface "my-zinx/zinx/interface"
	"my-zinx/zinx/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Name      string
	IPVersion string
	Ip        string
	Port      uint16
	// 连接管理器
	connMgr iface.ISessionMgr
	// 映射请求到具体的API回调
	JobRouter iface.IJobRouter
	// 工作协程池
	wokerPool *job.WokerPool
}

func NewServer() *Server {
	// 消息队列（worker协程从中取数据）mq容量和worker数量相同。mq容量更大没意义
	mq := message.NewMsgQueue(int(utils.Conf.Server.MaxWorkerPoolSize))
	router := job.NewJobRouter()
	return &Server{
		Name:      utils.Conf.Server.Name,
		IPVersion: "tcp4",
		Ip:        utils.Conf.Server.Host,
		Port:      utils.Conf.Server.Port,
		connMgr:   session.NewConnMgr(),
		JobRouter: router,
		wokerPool: job.NewWokerPool(mq.Cap(), mq, router),
	}
}

func (s *Server) Listen() {
	log.Printf("Server Start with config: %s\n", utils.Conf)

	// 忽略信号
	// 在某些系统中，syscall.SIGCHLD 可能未定义，这里仅忽略 SIGPIPE 信号
	signal.Ignore(syscall.SIGPIPE)

	endpoint, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Println("ResolveTCPAddr error:", err)
		return
	}
	listener, err := net.ListenTCP(s.IPVersion, endpoint) // FIXME listener没有Close啊
	if err != nil {
		log.Println("ListenTCP error:", err)
		return
	}
	log.Println(s.Name, " Listening...")

	// 注册心跳路由
	s.JobRouter.AddJob(job.HeartBeatTag, &job.HeartBeatJob{})
	// 启动协程池
	s.wokerPool.Start()

	// 启用单独的协程来处理客户端连接
	// 这是go语言的风格，能用异步一般用异步。这样主协程接下来还可以做其他工作，比如后面的Serve()方法
	go func() {
		for {
			peer, err := listener.AcceptTCP()
			if err != nil {
				log.Println("AcceptTCP error:", err)
				continue
			}
			if s.connMgr.Count() > utils.Conf.Server.MaxConnCount {
				log.Println("Too many connections, close this new connection")
				peer.Close()
				continue
			}
			log.Printf("New connection from %s", peer.RemoteAddr())

			dealConn := session.NewConnection(peer, context.Background(), s.wokerPool)
			s.connMgr.Add(dealConn)
			// 启动子协程处理业务
			go dealConn.Open()
		}
	}()
}

func (s *Server) Serve() {
	log.Println("Server Serve")

	// TODO 是否在启动服务的时候, 还需要做其它事情呢? 比如自定义 Logger 或加入鉴权中间件等

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quitCh := make(chan os.Signal, 1)                      // REVIEW 这里为啥必须是1的buffered chan
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM) // 订阅SIGINT和SIGTERM
	<-quitCh                                               // 从管道中读
	s.Shutdown()
}

func (s *Server) Shutdown() {
	log.Println("Server Shutdown")

	// TODO 将其它需要清理的连接信息或其他信息一并停止或清理

	s.connMgr.Clear()
	s.wokerPool.Stop()
}

func (s *Server) ListenAndServe() {
	s.Listen()
	s.Serve()
}

// 确保 Server 实现了 iface.IServer 的所有方法（让编译器帮我们检查）
var _ iface.IServer = (*Server)(nil)
