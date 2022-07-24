package subcmd

import (
	"errors"
	"github.com/spf13/cobra"
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

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "generator for init project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal(errors.New("invalid init param"))
			}

			log.Printf("init project[%s] running\n", args[0])
			p := NewInitProject(args[0])
			if err := p.Check(); err != nil {
				log.Fatal(err.Error())
			}

			if err := p.SetGoEnv(); err != nil {
				log.Fatal(err.Error())
			}

			if err := p.GetTemplate(); err != nil {
				log.Fatal(err.Error())
			}

			if err := p.Render(); err != nil {
				log.Fatal(err.Error())
			}

			if err := p.Complete(); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
}
