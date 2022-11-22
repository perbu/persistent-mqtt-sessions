package main

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/perbu/persistent-mqtt-sessions/config"
	"log"
	"net/url"
	"os"
	"time"
)

type Message struct {
	Topic   string
	Payload []byte
}

func NewConnection(c config.Config, p *persistentSession) (*autopaho.ConnectionManager, error) {
	brokerUrl, err := url.Parse(c.Broker)
	if err != nil {
		return nil, fmt.Errorf("error parsing broker url: %w", err)
	}
	log.Println("Connecting to broker", brokerUrl)
	logger := log.New(os.Stdout, "paho ", log.LstdFlags)
	debugLogger := log.New(os.Stdout, "paho-debug ", log.LstdFlags)
	pahoConfig := autopaho.ClientConfig{
		BrokerUrls: []*url.URL{brokerUrl},
		//		TlsCfg:            nil,

		KeepAlive:         5,
		ConnectRetryDelay: 5 * time.Second,
		ConnectTimeout:    5 * time.Second,
		OnConnectError:    func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
		Debug:             logger,
		PahoDebug:         debugLogger,
		ClientConfig: paho.ClientConfig{
			ClientID:    c.ClientID,
			Persistence: p,
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
				p.HandleMessage(m)
			}),
			EnableManualAcknowledgment: true,
		},
	}

	cm, err := autopaho.NewConnection(context.TODO(), pahoConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection manager: %w", err)
	}
	err = cm.AwaitConnection(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error awaiting connection: %w", err)
	}
	log.Println("Connected to broker. Subscribing to topic", c.Topic)
	subReq := &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{c.Topic: {QoS: 1}},
	}
	_, err = cm.Subscribe(context.TODO(), subReq)
	if err != nil {
		return nil, fmt.Errorf("error subscribing: %w", err)
	}
	log.Println("Subscribed to topic", c.Topic)

	return cm, nil
}

func realMain() error {
	mqttConfig, ok := config.GetConfig("producer")
	if !ok {
		return fmt.Errorf("config not found")
	}
	per := &persistentSession{}
	conn, err := NewConnection(mqttConfig, per)
	if err != nil {
		return fmt.Errorf("error creating connection: %w", err)
	}
	defer conn.Disconnect(context.TODO())
	time.Sleep(10 * time.Second)

	return nil
}

func main() {
	err := realMain()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
