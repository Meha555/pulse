package zutils

// viper库更好，还可以监听文件变动

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type zServerConf struct {
	Name          string `json:"name"`
	Host          string `json:"host"`
	Port          uint16 `json:"port"`
	MaxConn       uint   `json:"max_conn"`
	MaxPacketSize uint32 `json:"max_packet_size"`
}

type zConfig struct {
	Server zServerConf `json:"server"`
}

func (c *zConfig) Reload(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

const confFile = "zinx.json"

var Conf *zConfig

// 这里用init()函数，保证Conf不会是空指针
func init() {
	// 设置 Conf 的默认值
	Conf = &zConfig{
		Server: zServerConf{
			Host:          "0.0.0.0",
			Name:          "MY-ZINX Server",
			Port:          3333,
			MaxConn:       10,
			MaxPacketSize: 4096,
		},
	}
}

// 这里不用init()函数的原因是，不想让zinx的客户端用这个库也会加载服务端的配置文件
func Init() {
	path := os.Getenv("MY_ZINK_CONFIG_PATH")
	if path == "" {
		path, _ = os.Getwd()
	}
	// 使用 path/filepath 库的 Join 函数来拼接路径，避免路径分隔符数量不对的问题
	confPath := filepath.Join(path, confFile)
	if err := Conf.Reload(confPath); err != nil {
		log.Printf("Loading config from %s failed", confPath)
	}
	log.Println(*Conf)
}
