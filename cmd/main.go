package main

import (
	"fmt"
	"log"

	"magic.pathao.com/pinku/socktest/conf"
	"magic.pathao.com/pinku/socktest/lis"
	"magic.pathao.com/pinku/socktest/pub"
	"magic.pathao.com/pinku/socktest/tok"
)

var (
	listeners = map[int]*lis.Listener{}
)

func main() {
	if err := conf.Load(); err != nil {
		log.Println("Could not load config:", err)
		return
	}
	fmt.Println("conf lis:", conf.Lis, " pub:", conf.Pub)
	// tokens := tok.GetRangeTokens(1, maxID)

	tokens := tok.GetFixedTokens()
	if err := lis.LaunchListeners(listeners, tokens); err != nil {
		log.Println("Main could not create listeners, err:", err)
		return
	}
	mxTimes, pMaxTimes := pub.LaunchPublishers()

	fmt.Println("mxtimes per publisher:", mxTimes)
	fmt.Println("pmaxtimes for all publishers:", pMaxTimes)
	forever := make(chan bool)
	<-forever
}
