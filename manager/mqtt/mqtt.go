// Copyright © 2018 Beau Brewer <beaubrewer@gmail.com>

package mqtt

import (
	"fmt"
	"os"
	"sync"

	"github.com/beaubrewer/bellmanv2/config"
	"github.com/beaubrewer/bellmanv2/manager/audio"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Topic is the MQTT topic bellman subscribes to
// Any device that will trigger bellman should publish to this topic
const Topic = "bellman/doorbell"

var quit chan struct{}
var waitgroup sync.WaitGroup

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	f := audio.GetAudio(string(msg.Payload()))
	audio.Play(f)
}

// Start the MQTT consumer
func Start() {
	quit = make(chan struct{}, 1)
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()
		//create a ClientOptions struct setting the broker address, clientid, turn
		//off trace output and set the default message handler
		opts := MQTT.NewClientOptions().AddBroker("tcp://" + config.GetMQTTHost())
		opts.SetClientID("bellman")
		opts.SetDefaultPublishHandler(f)
		opts.AutoReconnect = true
		opts.OnConnectionLost = func(client MQTT.Client, e error) {
			fmt.Println("Connection Lost")
		}

		//create and start a client using the above ClientOptions
		c := MQTT.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		//subscribe to the topic and request messages to be delivered
		//at a maximum qos of zero, wait for the receipt to confirm the subscription
		if token := c.Subscribe(Topic, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
		//from the server after sending each message
		// for i := 0; i < 5; i++ {
		// 	text := fmt.Sprintf("this is msg #%d!", i)
		// 	token := c.Publish(Topic, 0, false, text)
		// 	token.Wait()
		// }

		<-quit

		//unsubscribe from /go-mqtt/sample
		if token := c.Unsubscribe(Topic); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		fmt.Println("Disconnecting MQTT client")
		c.Disconnect(250)
	}()
}

// Stop closes the MQTT consumer
func Stop() {
	quit <- struct{}{}
	waitgroup.Wait()
}
