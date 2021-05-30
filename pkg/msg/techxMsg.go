package msg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func NewKafkaRequest(topic string, metadata string, data interface{}) (msg string, msgId string) {
	if os.Getenv("MSG_TEMPLATE") == "TECHX" {
		return NewTechxRequestMsg(topic, metadata, data)
	} else {
		log.Fatal("MSG_TEMPLATE is not set. exit.")
		return "", ""
	}
}

func GetMsgData(msg string) (string, string) {
	if os.Getenv("MSG_TEMPLATE") == "TECHX" {
		techxMsg := (&TechxMsg{}).Deserialize(msg)
		if techxMsg.MessageType == "RESPONSE" {
			return Serialize(techxMsg.Data), techxMsg.TransactionId
		}
		return "techx-msg-but-not-response", ""
	} else {
		log.Fatal("MSG_TEMPLATE is not set. exit.")
		return "", ""
	}
}

func NewTechxRequestMsg(topic string, uri string, data interface{}) (string, string) {
	msg, id := NewTechxRequest(topic, uri, data)
	return Serialize(msg), id
}

func NewTechxRequest(topic string, uri string, data interface{}) (*TechxMsg, string) {
	return CreateTechxMsg(topic, uri, "REQUEST_RESPONSE", data, "REQUEST")
}

func CreateTechxMsg(topic string, uri string, responseUri string, data interface{}, messageType string) (*TechxMsg, string) {
	topicResponse := fmt.Sprintf("%v.%v.response.%v", APPLICATION_NAME, topic, NODE_ID)
	dataStr := Serialize(data)
	now := time.Now().Unix()
	return &TechxMsg{
		MessageType:   messageType,
		SourceId:      APPLICATION_NAME,
		MessageId:     now,
		TransactionId: strconv.FormatInt(time.Now().Unix(), 10),
		Uri:           uri,
		ResponseDestination: TechxMsgResponseDestination{
			Topic: topicResponse,
			Uri:   responseUri,
		},
		Data: dataStr,
	}, strconv.FormatInt(time.Now().Unix(), 10)

}

type TechxMsg struct {
	MessageType         string                      `json:"messageType"`
	SourceId            interface{}                 `json:"sourceId,omitempty"`
	MessageId           interface{}                 `json:"messageId,omitempty"`
	TransactionId       string                      `json:"transactionId,omitempty"`
	Uri                 string                      `json:"uri"`
	ResponseDestination TechxMsgResponseDestination `json:"responseDestination,omitempty"`
	Data                interface{}                 `json:"data"`
}

type TechxMsgResponseDestination struct {
	Topic string `json:"topic"`
	Uri   string `json:"uri"`
}

func (techxMsg *TechxMsg) Deserialize(data string) *TechxMsg {
	//https://stackoverflow.com/questions/14289256/cannot-convert-data-type-interface-to-type-string-need-type-assertion
	techxMsg, ok := deserialize(data, techxMsg).(*TechxMsg)
	if ok == false {
		log.Fatal("Error casting msg to TechxMsg.", techxMsg)
	}
	return techxMsg
}

func Serialize(object interface{}) string {
	byte, err := json.Marshal(object)
	if err != nil {
		log.Printf("error serializing string: %v", err)
		return ""
	}
	return string(byte)
}

func deserialize(fromString string, toObject interface{}) interface{} {
	err := json.Unmarshal([]byte(fromString), toObject)
	if err != nil {
		log.Printf("error deserializing object: %v", err)
		return nil
	}
	return toObject
}
