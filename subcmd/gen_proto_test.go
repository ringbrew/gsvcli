package subcmd

import (
	"testing"
)

func TestNewGenProto(t *testing.T) {
	gp := NewGenProto("test/example")
	if err := gp.Process(); err != nil {
		t.Error(err.Error())
		return
	}
}
