package config

import (
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MySQL    MysqlConf          `json:",optional"`
	BsRagRpc zrpc.RpcClientConf `json:",optional"`
    BsLlmRpc zrpc.RpcClientConf `json:",optional"`
    Cache    cache.CacheConf    `json:",optional"`
}

type MysqlConf struct {
	DataSource string
}
