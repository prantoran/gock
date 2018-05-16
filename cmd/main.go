package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	// addr = flag.String("addr", "localhost:8080", "http service address")
	addr   = "192.168.0.167:4201"
	reqcnt = 0
)

var listeners = map[int]*Listener{}

var maxTimes = []time.Duration{}

func insertDuration(t time.Duration) {
	l := len(maxTimes)
	cur := t
	for i := 0; i < l && i < 10; i++ {
		u := maxTimes[i]
		// fmt.Println("\tu:", u, " t:", t)
		if u < cur {
			for j := i; j < 10 && j < l; j++ {
				u = maxTimes[j]
				maxTimes[j] = cur
				cur = u
			}
			if l < 10 {
				maxTimes = append(maxTimes, cur)
			}
			return
		}
	}
	if l < 10 {
		maxTimes = append(maxTimes, t)
	}
}

// CreateListeners create (endID-startID+1) listeners by calling NewListener.
// NewListener calls a goroutine, so calling CreateListeners as a gorouting is unadvised.
func CreateListeners(addr string, startID, endID int, l map[int]*Listener) error {
	fmt.Println("creating listeners st:", startID, " nd:", endID)
	if endID < startID {
		return fmt.Errorf("StartID greater than endID")
	}
	var wg sync.WaitGroup
	wg.Add(endID - startID + 1)

	for i := startID; i <= endID; i++ {
		fmt.Println("creating lis:", i)
		listener := NewListener(addr, "/ws", strconv.Itoa(i))
		l[i] = listener // saving reference so that the timer chan can be used
		listener.Run(&wg)
	}

	wg.Wait()
	return nil
}

func publish(pubID, userID, msgID int) {
	s := NewPublisher(addr, "/ws", "pub")
	if err := s.Connect(); err != nil {

		log.Println("publichser run: connect:", err)
		return
	}
	body := PubMsg{
		Type:   "relay",
		UserID: int32(userID),
		Msg:    fmt.Sprintf("%v:%d", pubID, msgID),
	}
	log.Println("pub bod:", body, " conn")

	// ensuring that chan for listener.Ack[body.Msg] exists
	listeners[userID].mu.Lock()
	if listeners[userID].Ack[body.Msg] == nil {
		log.Println("created nu chan userID:", userID, " body.Msg:", body.Msg)
		listeners[userID].Ack[body.Msg] = make(chan bool)
	}
	listeners[userID].mu.Unlock()

	startTime := time.Now()
	if err := s.Conn.WriteJSON(body); err != nil {
		log.Printf("write: %v", err)
		return
	}
	<-listeners[userID].Ack[body.Msg]
	d := time.Now().Sub(startTime)
	insertDuration(d)
	reqcnt++
	log.Println("pubid:", pubID, " listener:", userID, " received msg, duration:", d, " reqcnt:", reqcnt)
	defer func() {
		// time.Sleep(time.Second)

		s.Conn.Close()
	}()
}

// RunPublisher sends msgs at random to listeners with uids in the range [1,maxID]
func RunPublisher(pubID, maxID, turns int, wg *sync.WaitGroup) {
	r := rand.New(rand.NewSource(99))
	fmt.Println("publid:", pubID)
	for i := 1; i <= turns; i++ {
		id := r.Intn(maxID) + 1
		publish(pubID, id, i)
	}
	fmt.Println("maxtimes:", maxTimes)
	wg.Done()
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	maxID := *flag.Int("maxid", 500, "maxid indicates max number of clients")
	turns := *flag.Int("turns", 300, "number of websocket requests by a publisher")
	pubs := *flag.Int("pubs", 5, "number of publisher")

	if err := CreateListeners(addr, 1, maxID, listeners); err != nil {
		log.Println("Could not create listeners, err:", err)
	}
	wg := sync.WaitGroup{}
	st := time.Now()
	for i := 1; i <= pubs; i++ {
		wg.Add(1)
		go RunPublisher(i, maxID, turns, &wg)
	}
	wg.Wait()
	fmt.Println("diff:", time.Now().Sub(st))
	// fmt.Println("the end")

	forever := make(chan bool)
	<-forever
}
