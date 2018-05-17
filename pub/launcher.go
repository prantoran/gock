package pub

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"magic.pathao.com/pinku/socktest/conf"
)

// LaunchPublishers setup waitgroup, executes pubs number of publishers
func LaunchPublishers() ([]time.Duration, []time.Duration) {
	wg := sync.WaitGroup{}
	Reqcnt = 0
	MaxTimes = []time.Duration{}
	PmaxTimes = []time.Duration{}
	st := time.Now()
	for i := 1; i <= conf.Pub.Workers; i++ {
		wg.Add(1)
		go RunPublisher(i, &wg)
	}
	wg.Wait()
	fmt.Println("publisher completion max times:", PmaxTimes)
	fmt.Println("diff:", time.Now().Sub(st))
	return MaxTimes, PmaxTimes
}

// RunPublisher sends msgs at random to listeners with uids in the range [1,maxID]
func RunPublisher(pubID int, wg *sync.WaitGroup) {
	psttime := time.Now()
	defer func() {
		InsertDuration(time.Now().Sub(psttime), PmaxTimes)
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
	fmt.Println("maxtimes:", MaxTimes)
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
	Reqcnt++
}
