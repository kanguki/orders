package kafka

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var KAFKA_SERVERS=os.Getenv("KAFKA_SERVERS") //host:port strings, separated by comma. eg: "localhost:9092,localhost:9093"
var KAFKA_MAIN_TOPIC=os.Getenv("KAFKA_MAIN_TOPIC")
var KAFKA_ID=os.Getenv("KAFKA_ID")
var KAFKA_GROUP_ID=os.Getenv("KAFKA_GROUP_ID")
var KafkaProducer *kafka.Producer = nil
var KafkaConsumer *kafka.Consumer = nil

func NewProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":     KAFKA_SERVERS,
		"request.required.acks": -1,
		"message.max.bytes":     1000000,
	})
	if err != nil {
		log.Fatal("Error connecting kafka: ", err)
	}
	return producer
}

func NewConsumer(topics string, groupId string) *kafka.Consumer { //eg: "test1,test2"
	gId := KAFKA_GROUP_ID
	if groupId != "" { //zero value
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
	consumer.SubscribeTopics(strings.Split(topics,",") , nil)
	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else {
				// The client will automatically try to recover from all errors.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()
	return consumer
}
