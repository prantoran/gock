package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/prantoran/gock/lis"
	pub "github.com/prantoran/gock/pub"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
var (
	addr = "192.168.0.167:4201"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	maxID, err := strconv.Atoi(*flag.String("maxid", "array", "type of publishing to use"))
	if err != nil {
		log.Println("maxId not integer")
	}
	turns, err := strconv.Atoi(*flag.String("turns", "pub", "publisher name"))
	if err != nil {
		log.Println("turns not integer")
	}

	if err := lis.CreateListeners(addr, 1, maxID); err != nil {
		log.Println("Could not create listeners, err:", err)
	}

	s := pub.NewPublisher(addr, "/ws", "pub", maxID, turns)
	if err := s.Run(); err != nil {
		log.Println("sender run error:", err)
	}

	forever := make(chan bool)
	<-forever
}
