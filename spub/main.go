package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"magic.pathao.com/pinku/socktest/conf"
)

var userIDs = []int{
	999,
	998,
	997,
	996,
	995,
	993,
	992,
	991,
	990,
	989,
	988,
	987,
	985,
	982,
	978,
	977,
	974,
	969,
	966,
	965,
	962,
	960,
	956,
	955,
	954,
	953,
	952,
	951,
	945,
	940,
	939,
	938,
	937,
	92,
	916,
	915,
	912,
	910,
	904,
	903,
	902,
	900,
	899,
	898,
	897,
	896,
	895,
	894,
	893,
	892,
	889,
	888,
	884,
	883,
	882,
	878,
	876,
	873,
	872,
	870,
	869,
	868,
	862,
	856,
	855,
	854,
	853,
	852,
	850,
	849,
	846,
	845,
	844,
	843,
	842,
	840,
	838,
	833,
	831,
	828,
	824,
	823,
	821,
	814,
	813,
	812,
	811,
	810,
	809,
	808,
	807,
	806,
	805,
	802,
	801,
	8,
	786,
	784,
	783,
	782,
}

var cnt = 0
var sucCnt = 0

func pub(host string, userID int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	u := url.URL{Scheme: "ws", Host: host, Path: ""}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("connection error", err)
		return
	}

	defer c.Close()

	cnt++

	writeCnt := cnt

	log.Println("writing ", writeCnt)
	body := map[string]interface{}{
		"type":    "relay",
		"user_id": userID,
		"id_test": userID,
		"cnt":     cnt,
		"value":   "test publish",
	}

	if err := c.WriteJSON(body); err != nil {
		log.Printf("write error : %v", err)
		return
	}

	sucCnt++
	fmt.Println("successfully write done", writeCnt)

}

func pubtoall(pwg *sync.WaitGroup, addr string) {
	defer func() {
		pwg.Done()
	}()

	wg := sync.WaitGroup{}

	for i := range userIDs {
		wg.Add(1)
		go pub(addr, userIDs[i], &wg)
		// break
	}

	wg.Wait()
}

func main() {
	if err := conf.Load(); err != nil {
		log.Println("Could not load config:", err)
		return
	}
	fmt.Println("conf.Addr:", conf.Addr, "turns:", conf.Pub.Turns)

	wg := sync.WaitGroup{}
	for i := 0; i < conf.Pub.Turns; i++ {
		wg.Add(1)
		go pubtoall(&wg, conf.Addr)
	}

	wg.Wait()

	fmt.Println("success tot write", sucCnt)

	forever := make(chan bool)
	<-forever
}
