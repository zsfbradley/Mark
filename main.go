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

	log.Println(s.Start())
}