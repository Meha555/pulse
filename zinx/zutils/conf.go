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

func init() {
	// 设置 Conf 的默认值
	Conf = &zConfig{
		Server: zServerConf{
			Host:          "0.0.0.0",
			Name:          "MY-ZINX Server",
			Port:          8999,
			MaxConn:       10,
			MaxPacketSize: 4096,
		},
	}
	path := os.Getenv("MY_ZINK_CONFIG_PATH")
	if path == "" {
		path, _ = os.Getwd()
	}
	// 使用 path/filepath 库的 Join 函数来拼接路径，避免路径分隔符数量不对的问题
	confPath := filepath.Join(path, confFile)
	if err := Conf.Reload(confPath); err != nil {
		log.Panicf("Loading config from %s failed", confPath)
	}
	log.Println(*Conf)
}