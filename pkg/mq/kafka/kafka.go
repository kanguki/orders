package mokafka //import "mo.io/conditional_order/pkg/mq/kafka"

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	topic        string
	dataHolder   *map[string](chan string)
	p_worker     *kafka.Producer
	c_req_worker *kafka.Consumer
	c_res_worker *kafka.Consumer
}

func NewDriver(topic string, groupId string) *Driver {
	m := make(map[string](chan string))
	return &Driver{
		topic:        topic,
		dataHolder:   &m,
		p_worker:     NewProducer(),
		c_req_worker: NewConsumer(groupId),
		c_res_worker: NewConsumer(groupId),
	}
}

func (driver *Driver) ProduceAndGetResponse(topic string, uri string, data interface{}) string {
	msgToSend, msgId := momsg.NewKafkaRequest(driver.topic, uri, data)
	deliveryChan := deliveryHandler()
	driver.p_worker.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msgToSend),
	}, *deliveryChan)
	log.Printf("Sent msg: %v to topic %v", msgToSend, topic)
	(*driver.dataHolder)[msgId] = make(chan string)
	select {
	case response := <-(*driver.dataHolder)[msgId]:
		delete(*driver.dataHolder, msgId)
		return response
	case <-time.After(setTimeout()):
		delete(*driver.dataHolder, msgId)
		return constant.REQUEST_TIMEOUT
	}
}

func (driver *Driver) ProduceResponseNoProcess(msgId string, topicResponse string, uriResponse string, data interface{}) {
	//TODO: wrap data with metadata about id, type, topic response...
	msgToSend, msgId := momsg.NewKafkaRequest(driver.topic, uriResponse, data)
	log.Printf("Sent msg: %v to topic %v", msgToSend, topicResponse)
	deliveryChan := deliveryHandler()
	driver.p_worker.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicResponse, Partition: kafka.PartitionAny},
		Value:          []byte(msgToSend),
	}, *deliveryChan)
}

func (driver *Driver) ProduceResponse(msgId string, topicResponse string, data interface{}) {
	//TODO: wrap data with metadata about id, type, topic response...
	msgToSend, msgId := momsg.NewKafkaResponse(msgId, data)
	log.Printf("Sent msg: %v to topic %v", msgToSend, topicResponse)
	deliveryChan := deliveryHandler()
	driver.p_worker.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicResponse, Partition: kafka.PartitionAny},
		Value:          []byte(msgToSend),
	}, *deliveryChan)
}

func (driver *Driver) ConsumeReq(do func(driver *Driver, msg string)) {
	topics := strings.Split(driver.topic, ",")
	errC := driver.c_req_worker.SubscribeTopics(topics, nil)
	if errC != nil {
		log.Fatal("Error consuming topic: ", topics)
	}
	log.Printf("Listening on topic: %v", topics)
	for {
		msg, err := driver.c_req_worker.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Receive message %s: %s\n", msg.TopicPartition.Offset, string(msg.Value))
			do(driver, string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func (driver *Driver) ConsumeRes() {
	topics := strings.Split(fmt.Sprintf("%v.%v.response.%v", APPLICATION_NAME, driver.topic, NODE_ID), ",")
	errC := driver.c_res_worker.SubscribeTopics(topics, nil)
	if errC != nil {
		log.Fatal("Error consuming topic: ", topics)
	}
	log.Printf("Listening on topic: %v", topics)
	for {
		msg, err := driver.c_res_worker.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Receive response %s: %s\n", msg.TopicPartition.Offset, string(msg.Value))
			data, msgId := momsg.GetResData(string(msg.Value))
			if (*driver.dataHolder)[msgId] != nil {
				(*driver.dataHolder)[msgId] <- data
			} else {
				log.Println("data extracted: ", data, "id: ", msgId)
			}
			//then process data
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
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

func NewConsumer(groupId string) *kafka.Consumer {
	gId := KAFKA_GROUP_ID
	if groupId != "" {
		gId = groupId
	}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        KAFKA_SERVERS,
		"client.id":                NODE_ID,
		"group.id":                 gId,
		"auto.offset.reset":        "latest",
		"enable.auto.offset.store": true,
	})
	if err != nil {
		log.Fatal("Error connecting kafka: ", err)
	}
	log.Println("New kafka consumer ready!")
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
					fmt.Println(ev.TopicPartition)
				}
			}
		}
	}()
	return &deliveryChan
}
