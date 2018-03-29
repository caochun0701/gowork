package fileout

import (
	"testing"
	"log"
	"libbeat/beat"
	"github.com/satori/go.uuid"
)

func TestFileOutputConfig(t *testing.T) {

	var defaultConfig = &config{
		NumberOfFiles: 10,
		RotateEveryKb: 10 * 1024,
		Permissions:   0600,
	}
	err := defaultConfig.Validate()
	if err != nil{
		log.Println(err)
	}
	beat := beat.Info{"caocu","adfad","0.0.1","adfad","adfaf",uuid.UUID{},}

	log.Println(beat)


}
