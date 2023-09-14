package events

import (
	"log"
	"sync"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

// DataChannel is a channel which can accept an DataEvent
type DataChannel chan DataEvent // DataChannelSlice is a slice of DataChannels
type DataChannelSlice []DataChannel

type EventBus struct {
	subscribers map[string]DataChannelSlice
	mutex       sync.RWMutex
}

func (eb *EventBus) Subscribe(topic string) DataChannel {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	ch := make(DataChannel)
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	return ch
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	if chans, found := eb.subscribers[topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				log.Printf("Publish: Sending data to channel %s", topic)
				select {
				case ch <- data:
					// send data to the channel if possible
					log.Printf("Publish: Sent data to channel %s", topic)
				default:
					// noop
				}
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
}

func NewBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string]DataChannelSlice),
	}
}
