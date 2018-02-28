package format

import (
	"testing"

	"libbeat/beat"
	"libbeat/common"
	"libbeat/common/fmtstr"
)

func TestFormatStringWriter(t *testing.T) {
	t.SkipNow()

	format := fmtstr.MustCompileEvent("test %{[msg]}")
	expectedValue := "test message"

	codec := New(format)
	output, err := codec.Encode("test", &beat.Event{Fields: common.MapStr{"msg": "message"}})

	if err != nil {
		t.Errorf("Error during event write %v", err)
	} else {
		if string(output) != expectedValue {
			t.Errorf("Expected value (%s) does not equal with output %s", expectedValue, output)
		}
	}
}
