package console

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"

	"libbeat/beat"
	"libbeat/common"
	"libbeat/logp"
	"libbeat/outputs"
	"libbeat/outputs/codec"
	"libbeat/outputs/codec/json"
	"libbeat/publisher"
)

type console struct {
	out      *os.File
	observer outputs.Observer
	writer   *bufio.Writer
	codec    codec.Codec
	index    string
}

type consoleEvent struct {
	Timestamp time.Time `json:"@timestamp" struct:"@timestamp"`

	// Note: stdlib json doesn't support inlining :( -> use `codec: 2`, to generate proper event
	Fields interface{} `struct:",inline"`
}

func init() {
	outputs.RegisterType("console", makeConsole)
}

func makeConsole(
	beat beat.Info,
	observer outputs.Observer,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	err := cfg.Unpack(&config)
	if err != nil {
		return outputs.Fail(err)
	}

	var enc codec.Codec
	if config.Codec.Namespace.IsSet() {
		enc, err = codec.CreateEncoder(beat, config.Codec)
		if err != nil {
			return outputs.Fail(err)
		}
	} else {
		enc = json.New(config.Pretty, beat.Version)
	}

	index := beat.Beat
	c, err := newConsole(index, observer, enc)
	if err != nil {
		return outputs.Fail(fmt.Errorf("console output initialization failed with: %v", err))
	}

	// check stdout actually being available
	if runtime.GOOS != "windows" {
		if _, err = c.out.Stat(); err != nil {
			err = fmt.Errorf("console output initialization failed with: %v", err)
			return outputs.Fail(err)
		}
	}

	return outputs.Success(config.BatchSize, 0, c)
}

func newConsole(index string, observer outputs.Observer, codec codec.Codec) (*console, error) {
	c := &console{out: os.Stdout, codec: codec, observer: observer, index: index}
	c.writer = bufio.NewWriterSize(c.out, 8*1024)
	return c, nil
}

func (c *console) Close() error { return nil }
func (c *console) Publish(batch publisher.Batch) error {
	st := c.observer
	events := batch.Events()
	st.NewBatch(len(events))

	dropped := 0
	for i := range events {
		ok := c.publishEvent(&events[i])
		if !ok {
			dropped++
		}
	}

	c.writer.Flush()
	batch.ACK()

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return nil
}

var nl = []byte("\n")

func (c *console) publishEvent(event *publisher.Event) bool {
	serializedEvent, err := c.codec.Encode(c.index, &event.Content)
	if err != nil {
		if !event.Guaranteed() {
			return false
		}

		logp.Critical("Unable to encode event: %v", err)
		return false
	}

	if err := c.writeBuffer(serializedEvent); err != nil {
		c.observer.WriteError(err)
		logp.Critical("Unable to publish events to console: %v", err)
		return false
	}

	if err := c.writeBuffer(nl); err != nil {
		c.observer.WriteError(err)
		logp.Critical("Error when appending newline to event: %v", err)
		return false
	}

	c.observer.WriteBytes(len(serializedEvent) + 1)
	return true
}

func (c *console) writeBuffer(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := c.writer.Write(buf[written:])
		if err != nil {
			return err
		}

		written += n
	}
	return nil
}
