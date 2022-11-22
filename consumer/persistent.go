package main

import (
	"github.com/eclipse/paho.golang/packets"
	"github.com/eclipse/paho.golang/paho"
	"log"
)

// local implementation of the persistent session interface

type persistentSession struct {
}

func (p *persistentSession) HandleMessage(msg *paho.Publish) {
	topic := msg.Topic
	payload := msg.Payload
	log.Println("Received message on topic", topic, "with payload", string(payload))

}

func (p *persistentSession) Open() {
	//TODO implement me
	panic("implement me")
}

func (p *persistentSession) Put(u uint16, packet packets.ControlPacket) {
	//TODO implement me
	panic("implement me")
}

func (p *persistentSession) Get(u uint16) packets.ControlPacket {
	//TODO implement me
	panic("implement me")
}

func (p *persistentSession) All() []packets.ControlPacket {
	//TODO implement me
	panic("implement me")
}

func (p *persistentSession) Delete(u uint16) {
	//TODO implement me
	panic("implement me")
}

func (p *persistentSession) Close() {
	//TODO implement me
	panic("implement me")
}

func (p *persistentSession) Reset() {
	//TODO implement me
	panic("implement me")
}
