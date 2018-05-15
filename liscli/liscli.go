package main

import (
	"flag"
	"log"
)

const (
	Register = "register"
	Admin    = "admin"
	Send     = "send"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")

var addr = "192.168.0.167:4201"

func main() {
	flag.Parse()
	log.SetFlags(0)

	token := flag.Arg(0)

	listener := NewSocListener(addr, "/ws", token)
	listener.Run()

	// forever := make(chan bool)
	// <-forever
}

/*
00wO3ztoWQ16AfsGaCS5ddaUz6pkY2FpeHavP2RB
*/
