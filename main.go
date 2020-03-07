package main

import (
	"./mNet"
	"log"
	"os"
)

func main() {
	sc := mNet.ServerConfig{
		Name:    "Frank",
		Host:    "0.0.0.0",
		Port:    "8888",
		Network: "tcp",
	}

	s, err := mNet.NewServer(sc)
	if err != nil {
		log.Println(err)
		return
	}

	if err := s.Start(); err != nil {
		log.Println(err)
		return
	}

	log.Printf("[%d] pid is over" , os.Getpid())
}
