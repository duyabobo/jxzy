package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MySQL    MysqlConf          `json:",optional"`
	LLMRpc   zrpc.RpcClientConf `json:",optional"`
	BsLlmRpc zrpc.RpcClientConf `json:",optional"`
	BsRagRpc zrpc.RpcClientConf `json:",optional"`
}

type MysqlConf struct {
	DataSource string
}
