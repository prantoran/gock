package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	// addr = flag.String("addr", "localhost:8080", "http service address")
	addr = "192.168.0.167:4201"
)

/* Function to run the groutine to run for stdin read */
func read(r io.Reader) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lines <- scan.Text()
		}
	}()
	return lines
}

// Publisher sends objects to websocket
type Publisher struct {
	Name string
	Path string
	Host string
	Type string
}

// NewPublisher creates new instances of Publisher
func NewPublisher(host, path, name, tp string) *Publisher {
	s := Publisher{}
	s.Name = name
	s.Host = host
	s.Path = path
	s.Type = tp
	return &s
}

// MsgBody represents an object sent to websocket
type MsgBody struct {
	Type      string  `json:"type,omitempty"`
	Msg       string  `json:"msg,omitempty"`
	UserID    int32   `json:"user_id,omitempty"`
	DriverID  int32   `json:"driver_id,omitempty"`
	CurLat    float32 `json:"current_latitude,omitempty"`
	CurLong   float32 `json:"current_longitude,omitempty"`
	PromptAct string  `json:"prompt_action,omitempty"`
}

// Run creates websocket url, dials and listens to websock
func (s *Publisher) Run() error {
	u := url.URL{Scheme: "ws", Host: s.Host, Path: s.Path, RawQuery: "token=" + s.Name}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	switch t := s.Type; t {
	case "array":
		return s.ParseArrayofStrings(c)
	case "json":
		return s.ParseJSONString(c)
	case "fixed":
		return s.PublishFixedObj(c)
	}
	return nil
}

// ParseJSONString inputs json string, unmarshals and sends object to websocket
func (s *Publisher) ParseJSONString(c *websocket.Conn) error {
	mes := read(os.Stdin) //Reading from Stdin
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case msg := <-mes:
			body := MsgBody{}
			json.Unmarshal([]byte(msg), &body)
			if err := c.WriteJSON(body); err != nil {
				return fmt.Errorf("write: %v", err)
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("write close: %v", err)
			}
			return nil
		}
	}
}

// PublishFixedObj sends a fixed object to websocket
func (s *Publisher) PublishFixedObj(c *websocket.Conn) error {
	fixedObj := MsgBody{
		UserID:    12022,
		DriverID:  12,
		CurLat:    23.89778678,
		CurLong:   90.87764554,
		Type:      "relay",
		PromptAct: "update_ride",
	}
	if err := c.WriteJSON(fixedObj); err != nil {
		return fmt.Errorf("write: %v", err)
	}
	return nil
}

// ParseArrayofStrings parses words from string input in console and writes object to websocket
func (s *Publisher) ParseArrayofStrings(c *websocket.Conn) error {
	mes := read(os.Stdin) //Reading from Stdin
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case msg := <-mes:
			w := strings.Fields(msg)
			wlen := len(w)
			tp := w[0]
			uid, _ := strconv.Atoi(w[1])
			body := MsgBody{
				Type:   tp,
				UserID: int32(uid),
			}
			if wlen > 2 {
				body.Msg = w[2]
			}
			fmt.Printf("w: %v\n", w)
			fmt.Println("body tp:", body.Type, " userid:", body.UserID, " msg:", body.Msg)
			if err := c.WriteJSON(body); err != nil {
				return fmt.Errorf("write: %v", err)
			}
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("write close: %v", err)
			}
			return nil
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	tp := flag.String("type", "array", "type of publishing to use")
	name := flag.String("name", "pub", "publisher name")

	fmt.Println("tp:", *tp, " name:", *name)
	s := NewPublisher(addr, "/ws", *name, *tp)
	if err := s.Run(); err != nil {
		log.Println("sender run error:", err)
	}
}
