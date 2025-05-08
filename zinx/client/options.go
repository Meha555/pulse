package client

import "time"

type ClientOptions func(*Client)

func WithName(name string) ClientOptions {
	return func(cli *Client) {
		cli.Name = name
	}
}

func WithIPVersion(network string) ClientOptions {
	return func(cli *Client) {
		cli.IPVersion = network
	}
}

func WithExitTimeout(timeout int) ClientOptions {
	return func(cli *Client) {
		cli.exitTimeout = time.Duration(timeout)
	}
}

func WithHeartBeatInterval(interval int) ClientOptions {
	return func(cli *Client) {
		cli.heartBeatInterval = time.Duration(interval)
	}
}
