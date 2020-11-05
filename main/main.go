package main

import (
	"fmt"
	"log"
	"time"
	"github.com/fatemeh-al/DS_CA1_sync/main/broker"
)

type myMessage string

func recieveMessage(channel <-chan broker.Message, b broker.Broker){
	// time.Sleep(time.Second)
	for m := range channel {
		fmt.Printf("got message: %s\n", m)
		break
	}
	if err := b.Publish("ch1", fmt.Sprintf("send next message\n")); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var b broker.Broker

	// Trying the in-memory broker.
	b = broker.NewMemoryBroker()
	// b = broker.NewRedisBroker()

	// subCh is a readony channel that we will
	// receive messages published on "ch1".
	subCh, err := b.CreateChannel("ch1")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer b.Close()

		i := 0
		for {
			i++
			j := 0
			for {
				j++
				if err := b.Publish("ch1", fmt.Sprintf("message %d", i)); err != nil {
					log.Fatalln(err)
				}
				recieveMessage(subCh, b)
				time.Sleep(time.Second)
				select {
				case x := <-subCh:
					j = 3
				default:
					fmt.Println("No message recieved.")
				}
				if j == 3{
					break
				}
			}
			if i == 5 {
				if err := b.DeleteChannel("ch1"); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}
	}()
	
}
