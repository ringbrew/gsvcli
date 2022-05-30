package subcmd

import (
	"errors"
	"log"
)

type InitSubCmd struct {
	ProjectName string
}

func (sc *InitSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid project name")
	}
	sc.ProjectName = args[0]
	return nil
}

func (sc *InitSubCmd) Process() error {
	log.Printf("init project[%s] running\n", sc.ProjectName)
	p := NewInitProject(sc.ProjectName)
	if err := p.Check(); err != nil {
		return err
	}

	if err := p.SetGoEnv(); err != nil {
		return err
	}

	if err := p.GetTemplate(); err != nil {
		return err
	}

	if err := p.Render(); err != nil {
		return err
	}

	return p.Complete()
}
