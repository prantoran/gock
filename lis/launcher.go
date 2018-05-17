package lis

import (
	"fmt"
	"sync"

	"magic.pathao.com/pinku/socktest/conf"
)

// LaunchListeners create (endID-startID+1) listeners by calling NewListener.
// NewListener calls a goroutine, so calling CreateListeners as a gorouting is unadvised.
func LaunchListeners(l map[int]*Listener, tokens []string) error {
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
