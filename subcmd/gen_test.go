package subcmd

import (
	"testing"
)

func TestGen(t *testing.T) {
	g := NewGenGrpc()

	if err := g.Setup(); err != nil {
		t.Error(err)
		return
	}

	if err := g.Gen("proto/example.proto", "delivery"); err != nil {
		t.Error(err)
		return
	}

	//if err := g.autoCommit("example", "proto/example.proto"); err != nil {
	//	t.Error(err)
	//	return
	//}
}
