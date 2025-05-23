package server

import (
	"fmt"

	"my-zinx/logging"
	"my-zinx/server/common"
	"my-zinx/server/job"
	"my-zinx/server/session"
	"my-zinx/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var logger = logging.NewStdLogger(logging.LevelInfo, "server", "[%t] [%c %l] [%f:%C:%L:%g] %m", false)

type Server struct {
	Name      string
	IPVersion string
	Ip        string
	Port      uint16
	// 连接管理器
	sessionMgr common.ISessionMgr
	// 映射请求到具体的API回调
	jobRouter job.IJobRouter
	// 工作协程池
	workerPool *job.WorkerPool
}

func NewServer() *Server {
	// 消息队列（worker协程从中取数据）mq容量和worker数量相同。mq容量更大没意义
	mq := utils.NewBlockingQueue[common.IRequest](int(utils.Conf.Server.MaxWorkerPoolSize))
	router := job.NewJobRouter()
	return &Server{
		Name:       utils.Conf.Server.Name,
		IPVersion:  "tcp4",
		Ip:         utils.Conf.Server.Host,
		Port:       utils.Conf.Server.Port,
		sessionMgr: session.NewSessionMgr(),
		jobRouter:  router,
		workerPool: job.NewWorkerPool(mq.Cap(), mq, router),
	}
}

func (s *Server) Route(tag uint16, job job.IJob) *Server {
	s.jobRouter.AddJob(tag, job)
	return s
}

func (s *Server) Listen() {
	logger.Infof("Server Start with config: %s\n", utils.Conf)

	// 忽略信号
	// 在某些系统中，syscall.SIGCHLD 可能未定义，这里仅忽略 SIGPIPE 信号
	signal.Ignore(syscall.SIGPIPE)

	endpoint, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		logger.Errorf("ResolveTCPAddr error: %v", err)
		return
	}
	listener, err := net.ListenTCP(s.IPVersion, endpoint) // FIXME listener没有Close啊
	if err != nil {
		logger.Errorf("ListenTCP error: %v", err)
		return
	}
	logger.Infof("%s Listening on %s:%d ...", s.Name, s.Ip, s.Port)

	// 注册心跳路由
	s.jobRouter.AddJob(job.HeartBeatTag, &job.HeartBeatJob{})
	// 启动协程池
	s.workerPool.Start()

	// 启用单独的协程来处理客户端连接
	// 这是go语言的风格，能用异步一般用异步。这样主协程接下来还可以做其他工作，比如后面的Serve()方法
	go func() {
		for {
			peer, err := listener.AcceptTCP()
			if err != nil {
				logger.Errorf("AcceptTCP error: %v", err)
				continue
			}
			if s.sessionMgr.Count() > utils.Conf.Server.MaxConnCount {
				logger.Warn("Too many connections, close this new connection")
				peer.Close()
				continue
			}
			logger.Debugf("New connection from %s", peer.RemoteAddr())

			clientSession := session.NewSession(peer, s.workerPool)
			s.sessionMgr.Add(clientSession)
			// 启动子协程处理业务
			go clientSession.Open()
		}
	}()
}

func (s *Server) Serve() {
	logger.Debug("Server Serve")

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh
	s.Shutdown()
}

func (s *Server) Shutdown() {
	logger.Debug("Server Shutdown")

	// 将其它需要清理的连接信息或其他信息一并停止或清理

	s.sessionMgr.Clear()
	s.workerPool.Stop()
}

func (s *Server) ListenAndServe() {
	s.Listen()
	s.Serve()
}

// 确保 Server 实现了 IServer 的所有方法（让编译器帮我们检查）
var _ IServer = (*Server)(nil)
