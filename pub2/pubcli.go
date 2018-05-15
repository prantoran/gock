package main

import (
	"flag"
	"log"
	"strconv"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
var addr = "192.168.0.167:4201"

func main() {
	flag.Parse()
	log.SetFlags(0)

	token := flag.Arg(0)

	n, _ := strconv.Atoi(flag.Arg(1))
	m, _ := strconv.Atoi(flag.Arg(2))
	
	s := NewSocSender(addr, "/ws", token, n, m)
	if err := s.Run(); err != nil {
		log.Println("sender run error:", err)
	}

	forever := make(chan bool)
	<-forever
}
