package fileout

import (
	"os"
	"path/filepath"

	"libbeat/common"
	"libbeat/common/file"
	"libbeat/logp"
	"libbeat/outputs"
	"libbeat/outputs/codec"
	"libbeat/publisher"
)
func init() {
	outputs.RegisterType("file", makeFileout)
}

type fileOutput struct {
	//beat     beat.Info
	observer outputs.Observer
	rotator  *file.Rotator
	codec    codec.Codec
}

// makeFileout instantiates a new file output instance.
func makeFileout(
	//beat beat.Info,
	observer outputs.Observer,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}

	// disable bulk support in publisher pipeline
	cfg.SetInt("bulk_max_size", -1, -1)
	fo := &fileOutput{
		//beat:     beat,
		observer: observer,
	}
	if err := fo.init(config); err != nil {
		return outputs.Fail(err)
	}
	return outputs.Success(-1, 0, fo)
}


func (out *fileOutput) init(c config) error {
	//读取配置文件中的文件输出的路径和文件名
	var path string
	if c.Filename != "" {
		path = filepath.Join(c.Path, c.Filename)
	}
	var err error
	out.rotator, err = file.NewFileRotator(
		path,
		file.MaxSizeBytes(c.RotateEveryKb*1024),
		file.MaxBackups(c.NumberOfFiles),
		file.Permissions(os.FileMode(c.Permissions)),
	)
	if err != nil {
		return err
	}

	out.codec, err = codec.CreateEncoder(c.Codec)
	if err != nil {
		return err
	}

	logp.Info("Initialized file output. "+
		"path=%v max_size_bytes=%v max_backups=%v permissions=%v",
		path, c.RotateEveryKb*1024, c.NumberOfFiles, os.FileMode(c.Permissions))

	return nil
}

// Implement Outputer
func (out *fileOutput) Close() error {
	return out.rotator.Close()
}

func (out *fileOutput) Publish(
	batch publisher.Batch,
) error {
	defer batch.ACK()

	st := out.observer
	events := batch.Events()
	st.NewBatch(len(events))
	dropped := 0
	for i := range events {
		event := &events[i]
		//serializedEvent, err := out.codec.Encode(out.beat.Beat, &event.Content)
		serializedEvent, err := out.codec.EncodeFile(&event.Content)
		if err != nil {
			if event.Guaranteed() {
				logp.Critical("Failed to serialize the event: %v", err)
			} else {
				logp.Warn("Failed to serialize the event: %v", err)
			}

			dropped++
			continue
		}
		if _, err = out.rotator.Write(append(serializedEvent, '\n')); err != nil {
			st.WriteError(err)

			if event.Guaranteed() {
				logp.Critical("Writing event to file failed with: %v", err)
			} else {
				logp.Warn("Writing event to file failed with: %v", err)
			}

			dropped++
			continue
		}

		st.WriteBytes(len(serializedEvent) + 1)
	}

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return nil
}
