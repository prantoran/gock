package lis

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"magic.pathao.com/pinku/socktest/pub"
)

const (
	timeLayout = "2006-01-02 15:04:05.999999999 -0700 MST"
)

var (
	// variables used to count active and dropped conns,
	drpCnt = 0
	crtCnt = 0
	// to count the number of times wrong clients received msgs
	incorrect = 0
)

// Listener stores the vars req by a websocket listener
type Listener struct {
	Token string
	Path  string
	Host  string
	Conn  *websocket.Conn
}

// NewListener returns a new Listener
func NewListener(host, path, tok string) *Listener {
	s := Listener{}
	s.Token = tok
	s.Host = host
	s.Path = path

	return &s
}

// LisMsg represents tokens received by Listener
type LisMsg struct {
	Type  string `json:"type,omitempty"`
	Token string `json:"token,omitempty"`
}

// Connect creates a websocket conn
func (s *Listener) Connect() error {
	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	s.Conn = c
	return err
}

// Register registers listener to websocket server
func (s *Listener) Register() error {
	// registering token, useless
	return s.Conn.WriteJSON(LisMsg{
		Token: s.Token,
		Type:  "register",
	})
}

// Run creates websocket url, dials and listens to websock
func (s *Listener) Run(wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
	}()
	if err := s.Connect(); err != nil {
		return fmt.Errorf("listener Run Connect err: %v", err)
	}
	if err := s.Register(); err != nil {
		return fmt.Errorf("listener register err: %v", err)
	}

	go func() {
		defer func() {
			drpCnt++
			fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		}()
		crtCnt++
		fmt.Println("tot alive connection", crtCnt, "tot droped connection", drpCnt)
		for {
			msg := pub.Msg{}
			err := websocket.ReadJSON(s.Conn, &msg)
			if err != nil {
				log.Println("listener:", s.Token, " read:", err)
				return
			}
			var st, pubID, msgID string
			if _, err := fmt.Sscanf(msg.Msg, "%v:%v:%v", &pubID, &msgID, &st); err != nil {
				log.Println("listener could not parse msg:", msg.Msg)
				return
			}
			stime, err := time.Parse(timeLayout, st)
			d := time.Now().Sub(stime)
			if err != nil {
				log.Println("listener received time not in correct format")
				return
			}
			pub.InsertDuration(d, pub.MaxTimes)
		}
	}()
	return nil
}
