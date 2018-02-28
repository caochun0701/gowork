package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"libbeat/beat"
	"libbeat/common"
	"libbeat/logp"
	"libbeat/paths"
	"libbeat/publisher/pipeline/stress"
	"libbeat/service"

	// import queue types
	_ "libbeat/publisher/queue/memqueue"

	// import outputs
	_ "libbeat/outputs/console"
	_ "libbeat/outputs/elasticsearch"
	_ "libbeat/outputs/fileout"
	_ "libbeat/outputs/logstash"
)

var (
	duration   time.Duration // -duration <duration>
	overwrites = common.SettingFlag(nil, "E", "Configuration overwrite")
)

type config struct {
	Path    paths.Path
	Logging logp.Logging
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	info := beat.Info{
		Beat:     "stresser",
		Version:  "0",
		Name:     "stresser.test",
		Hostname: "stresser.test",
	}

	flag.DurationVar(&duration, "duration", 0, "Test duration (default 0)")
	flag.Parse()

	files := flag.Args()
	cfg, err := common.LoadFiles(files...)
	if err != nil {
		return err
	}

	service.BeforeRun()
	defer service.Cleanup()

	if err := cfg.Merge(overwrites); err != nil {
		return err
	}

	config := config{}
	if err := cfg.Unpack(&config); err != nil {
		return err
	}

	if err := paths.InitPaths(&config.Path); err != nil {
		return err
	}
	if err = logp.Init("test", time.Now(), &config.Logging); err != nil {
		return err
	}
	logp.SetStderr()

	return stress.RunTests(info, duration, cfg, nil)
}

func startHTTP(bind string) {
	go func() {
		http.ListenAndServe(bind, nil)
	}()
}
