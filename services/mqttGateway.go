package services

import (
	"bytes"
	"strings"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

const MQTT_BROKER = "tcp://192.168.0.106:1883"

func createMQTTClient(brokerAddr, clientId, username, password string) *MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker(brokerAddr)
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	client := MQTT.NewClient(opts)
	return client
}

func subscribe(client *MQTT.Client, id string, sub chan<-MQTT.Message) {
	fmt.Println("Start MQTT subsribing..")
	
	topic := "templar87@github/"
//	topic = strConcatenateByBuffer(topic, id)
	
	subToken := client.Subscribe(
		topic,
		0,
		func(client *MQTT.Client, msg MQTT.Message) {
			sub <- msg
		})

	if subToken.Wait() && subToken.Error() != nil {
//		fmt.Println(subToken.Error())
//		os.Exit(1)
		panic(subToken.Error())
	}
}


func publish(client *MQTT.Client, input string) {
	token := client.Publish("templar87@github/", 0, true, input)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}


func input(pub chan<-string) {
	for {
		var input string
		fmt.Scanln(&input)
		
		pub <- input
	}
}


func InitMQTT() {
	go startMQTTClient()
}

func startMQTTClient() {
	id := "server"
	
	
	client := createMQTTClient(MQTT_BROKER, id, "test", "test")
	defer client.Disconnect(250)
	
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	
	sub := make(chan MQTT.Message)
	go subscribe(client, id, sub)
	pub := make(chan string)
//	go input(pub)
	for {
		select {
			case s := <-sub:
				msg := string(s.Payload())
				fmt.Printf("\nmsg: %s\n", msg)

			case p := <-pub:
				publish(client, p)
		}
	}
}

func strConcatenateByJoin(str1, str2 string) string {
	s := []string{}
	s = append(s, str1)
	s = append(s, str2)
	
	return strings.Join(s, "")
}

func strConcatenateByBuffer(str1, str2 string) string {
	var buffer bytes.Buffer
	
	buffer.Write([]byte(str1))
	buffer.Write([]byte(str2))
	
	return buffer.String()
}













