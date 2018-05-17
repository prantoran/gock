package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"magic.pathao.com/pinku/socktest/conf"
	"magic.pathao.com/pinku/socktest/tok"
)

var (
	reqcnt    = 0
	listeners = map[int]*Listener{}
	maxTimes  = []time.Duration{}
	pmaxTimes = []time.Duration{}
)


// CreateListeners create (endID-startID+1) listeners by calling NewListener.
// NewListener calls a goroutine, so calling CreateListeners as a gorouting is unadvised.
func CreateListeners(l map[int]*Listener, tokens []string) error {
	var wg sync.WaitGroup
	wg.Add(len(tokens))
	for i, tok := range tokens {
		fmt.Println("creating lis:", i)
		listener := NewListener(conf.Addr, "/ws", tok)
		l[i] = listener // saving reference so that the timer chan can be used
		listener.Run(&wg)
	}
	wg.Wait()
	fmt.Println("all listener running")
	return nil
}

func publish(pubID, userID, msgID int, wg *sync.WaitGroup) {
	s := NewPublisher(conf.Addr, "/ws", "pub")
	if err := s.Connect(); err != nil {
		log.Println("publichser connect:", err)
		return
	}
	defer func() {
		s.Conn.Close()
		wg.Done()
	}()

	st := time.Now().String()

	body := map[string]interface{}{
		"type":    "relay",
		"user_id": int32(userID),
		"msg":     fmt.Sprintf("%v:%d:%v", pubID, msgID, st),
	}
	// ensuring that chan for listener.Ack[body.Msg] exists
	// listeners[userID].EnsureAckChan(body.Msg)
	log.Println("publisher:", pubID, " receiver:", userID, " msg:", body["msg"])
	if err := s.Conn.WriteJSON(body); err != nil {
		log.Printf("publisher:%v write: %v", pubID, err)
		return
	}
	reqcnt++
}

// RunPublisher sends msgs at random to listeners with uids in the range [1,maxID]
func RunPublisher(pubID int, wg *sync.WaitGroup) {
	psttime := time.Now()
	defer func() {
		insertDuration(time.Now().Sub(psttime), pmaxTimes)
		wg.Done()
	}()
	swg := sync.WaitGroup{}
	r := rand.New(rand.NewSource(99))
	for i := 1; i <= conf.Pub.Turns; i++ {
		swg.Add(1)
		id := r.Intn(conf.Lis.Workers) + 1
		publish(pubID, id, i, &swg)
	}
	swg.Wait()
	fmt.Println("maxtimes:", maxTimes)
}

// LaunchPublishers setup waitgroup, executes pubs number of publishers
func LaunchPublishers() {
	wg := sync.WaitGroup{}
	st := time.Now()
	for i := 1; i <= conf.Pub.Workers; i++ {
		wg.Add(1)
		go RunPublisher(i, &wg)
	}
	wg.Wait()
	fmt.Println("publisher completion max times:", pmaxTimes)
	fmt.Println("diff:", time.Now().Sub(st))
}

func main() {
	if err := conf.Load(); err != nil {
		log.Println("Could not load config:", err)
		return
	}
	fmt.Println("conf lis:", conf.Lis, " pub:", conf.Pub)
	// tokens := tok.GetRangeTokens(1, maxID)
	tokens := tok.GetFixedTokens()
	if err := CreateListeners(listeners, tokens); err != nil {
		log.Println("Main could not create listeners, err:", err)
		return
	}
	LaunchPublishers()

	forever := make(chan bool)
	<-forever
}
