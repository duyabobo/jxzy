package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MySQL MysqlConf `json:",optional"`
}

type MysqlConf struct {
	DataSource string
}
