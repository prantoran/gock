package main

import (
	"flag"
	"log"
	"strconv"
)

const (
	Register = "register"
	Admin    = "admin"
	Send     = "send"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")

var addr = "192.168.0.167:4200"
var drpCnt = 0
var crtCnt = 0
var incorrect = 0

func main() {
	flag.Parse()
	log.SetFlags(0)

	st, _ := strconv.Atoi(flag.Arg(0))
	nd, _ := strconv.Atoi(flag.Arg(1))

	// createcnt := 0
	// dropcnt := 0
	// clichan := make(chan int)

	for i := st; i <= nd; i++ {

		listener := NewSocListener(addr, "/ws", strconv.Itoa(i))
		listener.Run()

	}

	forever := make(chan bool)
	<-forever

}
