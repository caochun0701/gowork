package config

import (
	"time"

	"libbeat/common"
	"libbeat/processors"
	"packetbeat/procs"
)

type Config struct {
	Interfaces      InterfacesConfig          `config:"interfaces"`
	Flows           *Flows                    `config:"flows"`
	Protocols       map[string]*common.Config `config:"protocols"`
	ProtocolsList   []*common.Config          `config:"protocols"`
	Procs           procs.ProcsConfig         `config:"procs"`
	IgnoreOutgoing  bool                      `config:"ignore_outgoing"`
	ShutdownTimeout time.Duration             `config:"shutdown_timeout"`
}

type InterfacesConfig struct {
	Device       string `config:"device"`
	Type         string `config:"type"`
	File         string `config:"file"`
	WithVlans    bool   `config:"with_vlans"`
	BpfFilter    string `config:"bpf_filter"`
	Snaplen      int    `config:"snaplen"`
	BufferSizeMb int    `config:"buffer_size_mb"`
	TopSpeed     bool
	Dumpfile     string
	OneAtATime   bool
	Loop         int
}

type Flows struct {
	Enabled       *bool                   `config:"enabled"`
	Timeout       string                  `config:"timeout"`
	Period        string                  `config:"period"`
	EventMetadata common.EventMetadata    `config:",inline"`
	Processors    processors.PluginConfig `config:"processors"`
}

type ProtocolCommon struct {
	Ports              []int         `config:"ports"`
	SendRequest        bool          `config:"send_request"`
	SendResponse       bool          `config:"send_response"`
	TransactionTimeout time.Duration `config:"transaction_timeout"`
	//统计时长 、 默认 10秒
	HowLong int `config:"how_long"`
	//执行次数 默认 2次
	CountTimes int `config:"count_times"`
	//热key单位时间内出现的次数
	HotKeysCount int `config:"hot_keys_count"`
	//大value size大小 bytes
	BigValueSize uint64 `config:"big_value_size"`
}

func (f *Flows) IsEnabled() bool {
	return f != nil && (f.Enabled == nil || *f.Enabled)
}
