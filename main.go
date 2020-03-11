package main

import (
	"./mNet"
	"log"
)

func main() {
	s := mNet.NewServer()

	err := s.Load()
	if err != nil {
		log.Println(err)
		return
	}

	err = mNet.InitializesIDG("mark", "6.0")
	if err != nil {
		log.Println(err)
		return
	}

	err = s.RegisterEncFunc(entrance)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(s.Start())

	for {
	}
}

func entrance() error {
	log.Println("================== Entrance Function ==================")
	return nil
}
