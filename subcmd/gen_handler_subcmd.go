package subcmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
)

type GenHandlerSubCmd struct {
	Domain string
}

func (sc *GenHandlerSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid domain")
	}
	sc.Domain = args[0]
	return nil
}

func (sc *GenHandlerSubCmd) Process() error {
	log.Printf("[INFO] gen domain command is running\n")

	p := NewGenHttp(sc.Domain)
	return p.Process()
}

func NewHandlerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "handler",
		Short: "generator for http code",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal(errors.New("invalid domain").Error())
			}
			domain := args[0]

			log.Printf("gen handler command is running domain[%s]\n", domain)

			p := NewGenHttp(domain)
			if err := p.Process(); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
}