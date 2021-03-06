package redis

import (
	"bytes"
	"time"

	"libbeat/beat"
	"libbeat/common"
	"libbeat/logp"
	"libbeat/monitoring"

	"packetbeat/procs"
	"packetbeat/protos"
	"packetbeat/protos/applayer"
	"packetbeat/protos/tcp"
)

type stream struct {
	applayer.Stream
	//解析方法
	parser   parser
	tcptuple *common.TCPTuple
}

type redisConnectionData struct {
	//stream类型数组，长度为2
	streams   [2]*stream
	requests  messageList
	responses messageList
}

type messageList struct {
	head *redisMessage
	tail *redisMessage
}

// Redis 协议 插件
type redisPlugin struct {
	// config
	ports        []int
	sendRequest  bool
	sendResponse bool
	transactionTimeout time.Duration
	//解析后的結果
	results protos.Reporter
}

var (
	debugf  = logp.MakeDebug("redis")
	isDebug = true
)

var (
	unmatchedResponses = monitoring.NewInt(nil, "redis.unmatched_responses")
)

//入口
func init() {
	//把 redis名字注册
	protos.Register("redis", New)
}
//注册后执行New方法，创建对象
func New(
	testMode bool,
	results protos.Reporter,
	cfg *common.Config,
) (protos.Plugin, error) {
	p := &redisPlugin{}
	config := defaultConfig
	if !testMode {
		if err := cfg.Unpack(&config); err != nil {
			return nil, err
		}
	}

	if err := p.init(results, &config); err != nil {
		return nil, err
	}
	return p, nil
}
//初始化
func (redis *redisPlugin) init(results protos.Reporter, config *redisConfig) error {
	redis.setFromConfig(config)

	redis.results = results
	isDebug = logp.IsDebug("redis")

	return nil
}
//获取配置文件中的配置
func (redis *redisPlugin) setFromConfig(config *redisConfig) {
	redis.ports = config.Ports
	redis.sendRequest = config.SendRequest
	redis.sendResponse = config.SendResponse
	redis.transactionTimeout = config.TransactionTimeout
}

func (redis *redisPlugin) GetPorts() []int {
	return redis.ports
}

func (s *stream) PrepareForNewMessage() {
	parser := &s.parser
	s.Stream.Reset()
	parser.reset()
}

func (redis *redisPlugin) ConnectionTimeout() time.Duration {
	return redis.transactionTimeout
}

func (redis *redisPlugin) Parse(
	pkt *protos.Packet,
	tcptuple *common.TCPTuple,
	dir uint8,
	private protos.ProtocolData,
) protos.ProtocolData {
	defer logp.Recover("ParseRedis exception")

	conn := ensureRedisConnection(private)
	conn = redis.doParse(conn, pkt, tcptuple, dir)
	if conn == nil {
		return nil
	}
	return conn
}

func ensureRedisConnection(private protos.ProtocolData) *redisConnectionData {
	if private == nil {
		return &redisConnectionData{}
	}

	priv, ok := private.(*redisConnectionData)
	if !ok {
		logp.Warn("redis connection data type error, create new one")
		return &redisConnectionData{}
	}
	if priv == nil {
		logp.Warn("Unexpected: redis connection data not set, create new one")
		return &redisConnectionData{}
	}

	return priv
}

func (redis *redisPlugin) doParse(
	conn *redisConnectionData,
	pkt *protos.Packet,
	tcptuple *common.TCPTuple,
	dir uint8,
) *redisConnectionData {

	st := conn.streams[dir]
	if st == nil {
		st = newStream(pkt.Ts, tcptuple)
		conn.streams[dir] = st
		if isDebug {
			debugf("new stream: %p (dir=%v, len=%v)", st, dir, len(pkt.Payload))
		}
	}

	if err := st.Append(pkt.Payload); err != nil {
		if isDebug {
			debugf("%v, dropping TCP stream: ", err)
		}
		return nil
	}
	if isDebug {
		debugf("stream add data: %p (dir=%v, len=%v)", st, dir, len(pkt.Payload))
	}

	for st.Buf.Len() > 0 {
		if st.parser.message == nil {
			st.parser.message = newMessage(pkt.Ts)
		}

		ok, complete := st.parser.parse(&st.Buf)
		if !ok {
			// drop this tcp stream. Will retry parsing with the next
			// segment in it
			conn.streams[dir] = nil
			if isDebug {
				debugf("Ignore Redis message. Drop tcp stream. Try parsing with the next segment")
			}
			return conn
		}

		if !complete {
			// wait for more data
			break
		}

		msg := st.parser.message
		if isDebug {
			if msg.isRequest {
				debugf("REDIS (%p) request message: %s", conn, msg.message)
			} else {
				debugf("REDIS (%p) response message: %s", conn, msg.message)
			}
		}

		// all ok, go to next level and reset stream for new message
		redis.handleRedis(conn, msg, tcptuple, dir)
		st.PrepareForNewMessage()
	}

	return conn
}
/*
新stream
*/
func newStream(ts time.Time, tcptuple *common.TCPTuple) *stream {
	s := &stream{
		tcptuple: tcptuple,
	}
	s.parser.message = newMessage(ts)
	s.Stream.Init(tcp.TCPMaxDataInStream)
	return s
}
/*
新redis消息
添加时间
*/
func newMessage(ts time.Time) *redisMessage {
	return &redisMessage{ts: ts}
}

func (redis *redisPlugin) handleRedis(
	conn *redisConnectionData,
	m *redisMessage,
	tcptuple *common.TCPTuple,
	dir uint8,
) {
	m.tcpTuple = *tcptuple
	//方向
	m.direction = dir
	m.cmdlineTuple = procs.ProcWatcher.FindProcessesTuple(tcptuple.IPPort())
	if m.isRequest {
		// wait for response
		conn.requests.append(m)
	} else {
		conn.responses.append(m)
		redis.correlate(conn)
	}
}
/*
相关联的
*/
func (redis *redisPlugin) correlate(conn *redisConnectionData) {
	// drop responses with missing requests
	if conn.requests.empty() {
		for !conn.responses.empty() {
			debugf("Response from unknown transaction. Ignoring")
			unmatchedResponses.Add(1)
			conn.responses.pop()
		}
		return
	}
	//将请求与响应合并到事务中
	for !conn.responses.empty() && !conn.requests.empty() {
		requ := conn.requests.pop()
		resp := conn.responses.pop()
		if redis.results != nil {
			event := redis.newTransaction(requ, resp)
			redis.results(event)
		}
	}
}

func (redis *redisPlugin) newTransaction(requ, resp *redisMessage) beat.Event {
	error := common.OK_STATUS
	if resp.isError {
		error = common.ERROR_STATUS
	}
	/*
	 取消返回数据，以减少数据大小
	*/
	//var returnValue map[string]common.NetString
	//if resp.isError {
	//	returnValue = map[string]common.NetString{
	//		"error": resp.message,
	//	}
	//} else {
	//	returnValue = map[string]common.NetString{
	//		"return_value": resp.message,
	//	}
	//}
	/*
	源地址信息
	*/
	src := &common.Endpoint{
		IP:   requ.tcpTuple.SrcIP.String(),
		Port: requ.tcpTuple.SrcPort,
		//Proc: string(requ.cmdlineTuple.Src),
	}
	/*
	目标地址信息
	*/
	dst := &common.Endpoint{
		IP:   requ.tcpTuple.DstIP.String(),
		Port: requ.tcpTuple.DstPort,
		//Proc: string(requ.cmdlineTuple.Dst),
	}
	if requ.direction == tcp.TCPDirectionReverse {
		src, dst = dst, src
	}
	// resp_time in milliseconds
	responseTime := int32(resp.ts.Sub(requ.ts).Nanoseconds() / 1e6)
	//格式化解析结果为map
	fields := common.MapStr{
		"type":         "redis",
		"status":       error,
		"responsetime": responseTime,
		//"redis":        returnValue,
		"method":       common.NetString(bytes.ToUpper(requ.method)),
		"key_name":     requ.path,
		"query":        requ.message,
		"bytes_in":     uint64(requ.size),
		"bytes_out":    uint64(resp.size),
		"src":          src,
		"dst":          dst,
	}
	if redis.sendRequest {
		fields["request"] = requ.message
	}
	if redis.sendResponse {
		fields["response"] = resp.message
	}
	return beat.Event{
		//Timestamp: requ.ts,
		Fields:    fields,
	}
}
//空包丢包处理
func (redis *redisPlugin) GapInStream(tcptuple *common.TCPTuple, dir uint8,
	nbytes int, private protos.ProtocolData) (priv protos.ProtocolData, drop bool) {

	// tsg: being packet loss tolerant is probably not very useful for Redis,
	// because most requests/response tend to fit in a single packet.

	return private, true
}
//断开链接的处理
func (redis *redisPlugin) ReceivedFin(tcptuple *common.TCPTuple, dir uint8,
	private protos.ProtocolData) protos.ProtocolData {

	// TODO: check if we have pending data that we can send up the stack

	return private
}

func (ml *messageList) append(msg *redisMessage) {
	if ml.tail == nil {
		ml.head = msg
	} else {
		ml.tail.next = msg
	}
	msg.next = nil
	ml.tail = msg
}
//清空messageList操作
func (ml *messageList) empty() bool {
	return ml.head == nil
}

func (ml *messageList) pop() *redisMessage {
	if ml.head == nil {
		return nil
	}

	msg := ml.head
	ml.head = ml.head.next
	if ml.head == nil {
		ml.tail = nil
	}
	return msg
}

func (ml *messageList) last() *redisMessage {
	return ml.tail
}
