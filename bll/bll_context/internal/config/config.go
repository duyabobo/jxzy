package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MySQL    MysqlConf          `json:",optional"`
	LLMRpc   zrpc.RpcClientConf `json:",optional"`
	BsLlmRpc zrpc.RpcClientConf `json:",optional"`
	BsRagRpc zrpc.RpcClientConf `json:",optional"`
	Bailian  BailianConfig      `json:",optional"`
}

type MysqlConf struct {
	DataSource string
}

type BailianConfig struct {
	APIKey string `json:"api_key"`
}
