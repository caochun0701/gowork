package protos

import (
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"log"
	"encoding/json"
)

type ProtocolPlugin func(
	testMode bool,
	results Reporter,
	cfg *common.Config,
) (Plugin, error)

// Reporter is used by plugin instances to report new transaction events.
type Reporter func(beat.Event)

// Functions to be exported by a protocol plugin
//扩展协议接口，主要是获取ports
type Plugin interface {
	// Called to return the configured ports
	GetPorts() []int
}
//tcp协议接口
type TCPPlugin interface {
	Plugin

	// Called when TCP payload data is available for parsing.
	//解析协议
	Parse(pkt *Packet, tcptuple *common.TCPTuple,
		dir uint8, private ProtocolData) ProtocolData

	// Called when the FIN flag is seen in the TCP stream.
	//处理tcp链接断开
	ReceivedFin(tcptuple *common.TCPTuple, dir uint8,
		private ProtocolData) ProtocolData

	// Called when a packets are missing from the tcp stream.
	// 处理空包 丢包
	GapInStream(tcptuple *common.TCPTuple, dir uint8, nbytes int,
		private ProtocolData) (priv ProtocolData, drop bool)

	// ConnectionTimeout returns the per stream connection timeout.
	// Return <=0 to set default tcp module transaction timeout.
	//处理超时时间
	ConnectionTimeout() time.Duration
}
//udp协议接口
type UDPPlugin interface {
	Plugin

	// ParseUDP is invoked when UDP payload data is available for parsing.
	//udp协议是无状态协议，所以只做协议解析
	ParseUDP(pkt *Packet)
}

// Protocol identifier.
type Protocol uint16

// Protocol constants.
const (
	UnknownProtocol Protocol = iota
)

// Protocol names
var protocolNames = []string{
	"unknown",
}

func (p Protocol) String() string {
	if int(p) >= len(protocolNames) {
		return "impossible"
	}
	return protocolNames[p]
}

var (
	protocolPlugins = map[Protocol]ProtocolPlugin{}
	protocolSyms    = map[string]Protocol{}
)

func Lookup(name string) Protocol {
	if p, exists := protocolSyms[name]; exists {
		return p
	}
	return UnknownProtocol
}

func Register(name string, plugin ProtocolPlugin) {
	proto := Protocol(len(protocolNames))
	if p, exists := protocolSyms[name]; exists {
		// keep symbol table entries if plugin gets overwritten
		proto = p
	} else {
		protocolNames = append(protocolNames, name)
		protocolSyms[name] = proto
	}

	protocolPlugins[proto] = plugin
}
