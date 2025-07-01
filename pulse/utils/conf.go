package utils

// viper库更好，还可以监听文件变动

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type zServerConf struct {
	Name              string `json:"name"`
	Host              string `json:"host"`
	Port              uint16 `json:"port"`
	HeartBeatTick     uint   `json:"heartbeat_tick"`
	ConnTimeout       uint   `json:"conn_timeout"`
	MaxConnCount      uint   `json:"max_conn_count"`
	MaxMsgQueueSize   uint   `json:"max_msg_queue_size"`
	MaxPacketSize     uint32 `json:"max_packet_size"`
	MaxWorkerPoolSize uint   `json:"max_worker_pool_size"`
	RequestPoolMode   bool   `json:"request_pool_mode"`
}

type zLogConf struct {
	Level      int    `json:"level"`
	Format     string `json:"format"`
	File       string `json:"file,omitempty"`
	Path       string `json:"path,omitempty"`
	MaxLogSize int64  `json:"max_log_size,omitempty"`
}

type zConfig struct {
	Server zServerConf `json:"server"`
	Log    zLogConf    `json:"log"`
}

func (c *zConfig) Reload(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

func (c *zConfig) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

const confFile = "pulse.json"

var Conf *zConfig

// 这里用init()函数，保证Conf不会是空指针
func init() {
	// 设置 Conf 的默认值
	Conf = &zConfig{
		Server: zServerConf{
			Host:              "0.0.0.0",
			Name:              "pulse Server@" + uuid.New().String(),
			Port:              3333,
			HeartBeatTick:     3,
			ConnTimeout:       60,
			MaxConnCount:      10,
			MaxMsgQueueSize:   50,
			MaxPacketSize:     4096,
			MaxWorkerPoolSize: 10,
			RequestPoolMode:   false,
		},
		Log: zLogConf{
			Level:  2,
			Format: "[%t] [%c %l] [%f:%L:%g] %m",
		},
	}

	path := os.Getenv("PULSE_CONFIG_PATH")
	if path == "" {
		path, _ = os.Getwd()
	}
	// 使用 path/filepath 库的 Join 函数来拼接路径，避免路径分隔符数量不对的问题
	confPath := filepath.Join(path, confFile)
	if err := Conf.Reload(confPath); err != nil {
		log.Printf("loading config from %s failed: %v", confPath, err)
	}
}
