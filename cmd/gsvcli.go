package main

import (
	"github.com/ringbrew/gsvcli/subcmd"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("sub command is required")
	}

	subCmd, err := subcmd.New(subcmd.Name(os.Args[1]))
	if err != nil {
		log.Fatal("[error]" + err.Error())
	}

	if err := subCmd.Parse(os.Args[2:]); err != nil {
		log.Fatal("[error]" + err.Error())
	}

	if err := subCmd.Process(); err != nil {
		log.Fatal("[error]" + err.Error())
	}

	log.Println("[info]done!!!")
}
