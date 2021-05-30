package mokafka //import "mo.io/conditional_order/pkg/mq/kafka"

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	constant "mo.io/conditional_order/conf"
	momsg "mo.io/conditional_order/pkg/msg"
)

var NODE_ID = os.Getenv("NODE_ID")
var APPLICATION_NAME = os.Getenv("APPLICATION_NAME")
var KAFKA_SERVERS = os.Getenv("KAFKA_SERVERS") //host:port strings, separated by comma. eg: "localhost:9092,localhost:9093"
var KAFKA_MAIN_TOPIC = os.Getenv("KAFKA_MAIN_TOPIC")
var KAFKA_GROUP_ID = os.Getenv("KAFKA_GROUP_ID")
var KAFKA_REQUEST_TIMEOUT = os.Getenv("KAFKA_REQUEST_TIMEOUT")

type Driver struct {
	topic      string
	dataHolder map[string]chan string
	p_worker   *kafka.Producer
	c_worker   *kafka.Consumer
}

func NewDriver(topic string, groupId string) *Driver {
	m := map[string]chan string{}
	return &Driver{
		topic:      topic,
		dataHolder: m,
		p_worker:   NewProducer(),
		c_worker:   NewConsumer(fmt.Sprintf("%v.%v.response.%v", APPLICATION_NAME, topic, NODE_ID), groupId, &m),
	}
}

func (driver *Driver) ProduceAndGetResponse(uri string, data interface{}) string {
	//TODO: wrap data with metadata about id, type, topic response...
	msgToSend, msgId := momsg.NewKafkaRequest(driver.topic, uri, data)
	log.Println("Sent msg: ", msgToSend)
	deliveryChan := deliveryHandler()
	driver.p_worker.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &driver.topic, Partition: kafka.PartitionAny},
		Value:          []byte(msgToSend),
	}, *deliveryChan)
	driver.dataHolder[msgId] = make(chan string)
	select {
	case response := <-driver.dataHolder[msgId]:
		log.Println("Chan received: ", response)
		return response
	case <-time.After(setTimeout()):
		return constant.REQUEST_TIMEOUT
	}
}

func NewProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":     KAFKA_SERVERS,
		"request.required.acks": -1,
		"message.max.bytes":     1000000,
	})
	if err != nil {
		log.Fatal("Error connecting kafka: ", err)
	}
	log.Println("New kafka producer ready!")
	return producer
}

func NewConsumer(topic string, groupId string, m *map[string]chan string) *kafka.Consumer {
	topics := strings.Split(topic, ",") //just to make configuration happy, as it has to be a []string
	gId := KAFKA_GROUP_ID
	if groupId != "" {
		gId = groupId
	}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        KAFKA_SERVERS,
		"client.id":                uuid.New().String(),
		"group.id":                 gId,
		"auto.offset.reset":        "latest",
		"enable.auto.offset.store": true,
	})
	if err != nil {
		log.Fatal("Error connecting kafka: ", err)
	}
	consumer.SubscribeTopics(topics, nil)
	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				data, msgId := momsg.GetMsgData(string(msg.Value))
				log.Println("data extracted: ", data, "id: ", msgId)
				// if <- (*m)[msgId] != "" {}
				// (*m)[msgId] <- data
			} else {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()
	log.Printf("New kafka consumer listening on topic: %v", topics)
	return consumer
}

func setTimeout() time.Duration {
	timeout, err := strconv.Atoi(KAFKA_REQUEST_TIMEOUT)
	if err != nil {
		log.Fatal("Error converting kafka request timeout to number")
	}
	return time.Duration(timeout) * time.Second
}

func deliveryHandler() *chan kafka.Event { //handle sent message, whether it's successful or not
	deliveryChan := make(chan kafka.Event)
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &deliveryChan
}
