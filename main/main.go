package main

import (
	"fmt"
	"log"
	"time"
	"github.com/pkg/errors"
	// "github.com/fatemeh-al/DS_CA1_sync/main/broker"
)

type Message interface{}

type Subscriber interface {
	CreateChannel(channel string) (<-chan Message, error)
	DeleteChannel(channel string) error
	Subscribe(channel string) (<-chan Message, error)
}

type Publisher interface {
	Publish(channel string, m Message) error
}

type Broker interface {
	Subscriber
	Publisher
	Close() error
}

type memoryBroker struct {
	chans map[string]chan Message
}

func NewMemoryBroker() *memoryBroker {
	return &memoryBroker{
		chans: make(map[string]chan Message),
	}
}

func (b *memoryBroker) CreateChannel(channel string) (<-chan Message, error) {
	ch := make(chan Message, 3)
	b.chans[channel] = ch
	return ch, nil
}

func (b *memoryBroker) Subscribe(channel string) (<-chan Message, error) {
	ch := b.chans[channel]
	return ch, nil
}

func (b *memoryBroker) DeleteChannel(channel string) error {
	ch, ok := b.chans[channel]
	if !ok {
		return errors.Errorf("cannot find channel %s", channel)
	}
	close(ch)
	delete(b.chans, channel)
	return nil
}

func (b *memoryBroker) Publish(channel string, m Message) error {
	ch, ok := b.chans[channel]
	if !ok {
		return errors.Errorf("cannot find channel %s", channel)
	}
	if len(ch) < 3 {
		ch <- m
		return nil
	}
	return nil;
}

func (b *memoryBroker) Close() error {
	for name, ch := range b.chans {
		close(ch)
		delete(b.chans, name)
	}
	return nil
}

type myMessage string

func recieveMessage(channel <-chan Message, b Broker){
	// time.Sleep(time.Second)
	for m := range channel {
		fmt.Printf("got message: %s\n", m)
		break
	}
	// if err := b.Publish("ch1", fmt.Sprintf("send next message\n")); err != nil {
	// 	log.Fatalln(err)
	// }
}

func main() {
	var b Broker

	// Trying the in-memory broker.
	b = NewMemoryBroker()
	// b = broker.NewRedisBroker()

	// subCh is a readony channel that we will
	// receive messages published on "ch1".
	subCh, err := b.CreateChannel("ch1")
	if err != nil {
		log.Fatalln(err)
	}

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
			case x := <-subCh:{
				fmt.Printf("%s", x)
				j = 3
			}
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
	
}
