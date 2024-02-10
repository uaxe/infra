package main

import (
	"log"
	"sync"
	"time"

	"github.com/uaxe/infra/telnet"
	"github.com/uaxe/infra/telnet/linereader"
)

func main() {
	svr := telnet.NewServer(":9999", telnet.HandleFunc(exampleHandler), telnet.NAWSOption)
	_ = svr.ListenAndServe()
}

func exampleHandler(c *telnet.Connection) {
	log.Printf("Connection received: %s", c.RemoteAddr())
	lr := linereader.New()
	go func() { _ = lr.ReadLines(c) }()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range lr.C {
			log.Printf("Received line: %v", line)
		}
	}()
	time.Sleep(time.Millisecond)
	nh := c.OptionHandlers[telnet.NAWS].(*telnet.NAWSHandler)
	log.Printf("Client width: %d, height: %d", nh.Width, nh.Height)
	wg.Wait()
	log.Printf("Goodbye %s!", c.RemoteAddr())
}
