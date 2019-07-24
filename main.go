package main

import (
	"log"
	"os"
	"xcore/auth"
	"xcore/cmd"
	"xcore/config"
	"xcore/world"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("please specify which server to start [auth | world]")
		return
	}

	switch args[1] {
	case "auth":
		go runAuthServer()
	case "world":
		go runWorldServer()
	default:
		log.Println("please specify server type to run [auth | world]")
		return
	}

	cmd.StartCLI()
}

func runAuthServer() {
	s, err := auth.NewServer(config.Current())
	if err != nil {
		log.Panic(err)
	}
	if err := s.Run(); err != nil {
		log.Panic(err)
	}
}

func runWorldServer() {
	s, err := world.NewServer(config.Current())
	if err != nil {
		log.Panic(err)
	}

	if err := s.Run(); err != nil {
		log.Panic(err)
	}
}
