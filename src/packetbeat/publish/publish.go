package publish

import (
	"libbeat/beat"
	"libbeat/common"
	"libbeat/logp"
	"libbeat/processors"
	"time"
	"sync"
	packetConfig "packetbeat/config"
)

type TransactionPublisher struct {
	done      chan struct{}
	pipeline  beat.Pipeline
	canDrop   bool
	processor transProcessor
}

type transProcessor struct {
	ignoreOutgoing bool
	localIPs       []string
	name           string
}

var debugf = logp.MakeDebug("publish")

func NewTransactionPublisher(
	name string,
	pipeline beat.Pipeline,
	ignoreOutgoing bool,
	canDrop bool,
) (*TransactionPublisher, error) {
	localIPs, err := common.LocalIPAddrsAsStrings(false)
	if err != nil {
		return nil, err
	}

	p := &TransactionPublisher{
		done:     make(chan struct{}),
		pipeline: pipeline,
		canDrop:  canDrop,
		processor: transProcessor{
			localIPs:       localIPs,
			name:           name,
			ignoreOutgoing: ignoreOutgoing,
		},
	}
	return p, nil
}

func (p *TransactionPublisher) Stop() {
	close(p.done)
}

func (p *TransactionPublisher) CreateReporter(
	config *common.Config,
) (func(beat.Event), error) {
	//初始化redis配置
	redisConf := packetConfig.ProtocolCommon{}
	if err := config.Unpack(&redisConf); err != nil {
		return nil, err
	}
	// load and register the module it's fields, tags and processors settings
	meta := struct {
		Event      common.EventMetadata    `config:",inline"`
		Processors processors.PluginConfig `config:"processors"`
	}{}
	if err := config.Unpack(&meta); err != nil {
		return nil, err
	}
	processors, err := processors.New(meta.Processors)
	if err != nil {
		return nil, err
	}

	clientConfig := beat.ClientConfig{
		EventMetadata: meta.Event,
		Processor:     processors,
	}

	if p.canDrop {
		clientConfig.PublishMode = beat.DropIfFull
	}

	client, err := p.pipeline.ConnectWith(clientConfig)
	if err != nil {
		return nil, err
	}
	// 启动 worker, so post-processing and processor-pipeline
	// 可以同时监听嗅探器获取新beat.Event
	ch := make(chan beat.Event, 3)
	go p.worker(ch, client,redisConf)
	return func(event beat.Event) {
		select {
		case ch <- event:
		case <-p.done:
			// stop serving more send requests
			ch = nil
		}
	}, nil
}

func (p *TransactionPublisher) worker(ch chan beat.Event, client beat.Client,redisConf packetConfig.ProtocolCommon) {
	//定时器 5 秒执行一次 ,config config.ProtocolCommon
	t := time.NewTimer(5*time.Second)
	//创建一个map
	m := new(sync.Map)
	for {
		select {
		case <-p.done:
			return
		case <-t.C:
			// range
			findHotKeysBigValues(m,client,redisConf)
			m = new(sync.Map)
			//重新设置5秒过期
			t.Reset(5*time.Second)
		case event := <-ch:
			pub, _ := p.processor.Run(&event)
			if pub != nil {
				//client.Publish(*pub)
				//对列信息进行count和大values查询
				suspectedHotkeyStore(m, pub.Fields)
			}
		}
	}
}

func (p *transProcessor) Run(event *beat.Event) (*beat.Event, error) {
	//filter 验证
	//if err := validateEvent(event); err != nil {
	//	logp.Warn("Dropping invalid event: %v", err)
	//	return nil, nil
	//}
	//重新设置地址信息
	if !p.normalizeTransAddr(event.Fields) {
		return nil, nil
	}
	return event, nil
}

/*
filter 验证，是否有@timestamp、type
如果验证不通过则返回 error.
*/
func validateEvent(event *beat.Event) error {
	//fields := event.Fields
	//
	//if event.Timestamp.IsZero() {
	//	return errors.New("missing '@timestamp'")
	//}
	//
	//_, ok := fields["@timestamp"]
	//if ok {
	//	return errors.New("duplicate '@timestamp' field from event")
	//}
	//
	//t, ok := fields["type"]
	//if !ok {
	//	return errors.New("missing 'type' field from event")
	//}
	//
	//_, ok = t.(string)
	//if !ok {
	//	return errors.New("invalid 'type' field from event")
	//}

	return nil
}
/*
  重新设置、 删除源地址信息、删除目标地址信息
*/
func (p *transProcessor) normalizeTransAddr(event common.MapStr) bool {
	debugf("normalize address for: %v", event)
	src, ok := event["src"].(*common.Endpoint)
	debugf("has src: %v", ok)
	if ok {
		event["client_ip"] = src.IP
		event["client_port"] = src.Port
		delete(event, "src")
	}

	dst, ok := event["dst"].(*common.Endpoint)
	debugf("has dst: %v", ok)
	if ok {
		event["ip"] = dst.IP
		event["port"] = dst.Port
		delete(event, "dst")
	}
	return true
}