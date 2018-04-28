package redis

import (
	"packetbeat/config"
	"packetbeat/protos"
)

type redisConfig struct {
	config.ProtocolCommon `config:",inline"`
}

var (
	defaultConfig = redisConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
			//统计时长 、 默认 10秒
			HowLong: protos.DefaultHowLong,
			//执行次数 默认 2次
			CountTimes: protos.DefaultCountTimes,
			//热key单位时间内出现的次数 默认 2000
			HotKeysCount: protos.DefaultHotKeysCount,
			//大value size大小 bytes 默认 1M
			BigValueSize: protos.DefaultBigValueSize,
		},
	}
)
