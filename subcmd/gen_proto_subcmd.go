package subcmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
)

type GenProtoSubCmd struct {
	Domain string
}

func (sc *GenProtoSubCmd) Parse(args []string) error {
	if len(args) == 0 {
		return errors.New("invalid domain")
	}
	sc.Domain = args[0]
	return nil
}

func (sc *GenProtoSubCmd) Process() error {
	log.Printf("[INFO] gen proto command is running\n")

	p := NewGenHttp(sc.Domain)
	return p.Process()
}

func NewGenProtoSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "proto",
		Short: "generator for proto file message. usage: gsv gen proto {domain}",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal(errors.New("invalid domain").Error())
			}
			domain := args[0]

			log.Printf("gen proto command is running domain[%s]\n", domain)

			p := NewGenProto(domain)
			if err := p.Process(); err != nil {
				log.Fatal(err.Error())
			}
		},
	}
}
