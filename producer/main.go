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

func NewConnection(c config.Config) (*autopaho.ConnectionManager, error) {
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
			ClientID: c.ClientID,
		},
	}
	cm, err := autopaho.NewConnection(context.TODO(), pahoConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection manager: %w", err)
	}
	time.Sleep(time.Second)
	log.Println("Waiting for connection to come up")
	err = cm.AwaitConnection(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error awaiting connection: %w", err)
	}
	log.Println("Connected to broker")
	return cm, nil
}

func realMain() error {
	mqttConfig, ok := config.GetConfig("producer")
	if !ok {
		return fmt.Errorf("config not found")
	}
	conn, err := NewConnection(mqttConfig)
	if err != nil {
		return fmt.Errorf("error creating connection: %w", err)
	}
	defer conn.Disconnect(context.TODO())
	i := 0
	for {
		i++
		message := fmt.Sprintf("Hello world %d", i)
		_, err := conn.Publish(context.TODO(), &paho.Publish{
			Topic:   mqttConfig.Topic,
			QoS:     1,
			Payload: []byte(message),
		})
		if err != nil {
			return fmt.Errorf("error publishing message: %w", err)
		}
		log.Printf("Published message: %s", message)
		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	err := realMain()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
