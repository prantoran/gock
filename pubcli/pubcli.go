package main

import (
	"flag"
	"log"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
var addr = "192.168.0.167:4200"

func main() {
	flag.Parse()
	log.SetFlags(0)

	token := flag.Arg(0)

	s := NewSocSender(addr, "/ws", token)
	if err := s.Run(); err != nil {
		log.Println("sender run error:", err)
	}

}
