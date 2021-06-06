package msg

import (
	"log"
	"os"
)

func GetResData(msg string) (string, string) {
	if os.Getenv("MSG_TEMPLATE") == "TECHX" {
		techxMsg := (&TechxMsg{}).Deserialize(msg)
		if techxMsg != nil {
			if techxMsg.MessageType == "RESPONSE" {
				return Serialize(techxMsg.Data), techxMsg.TransactionId
			}
			return "techx-msg-but-not-response", ""
		}
		return "not-techx-msg", ""
	} else {
		log.Fatal("MSG_TEMPLATE is not set. exit.")
		return "", ""
	}
}

func GetReqData(msg string) (uri string, data string, responseTopic string, id string) {
	if os.Getenv("MSG_TEMPLATE") == "TECHX" {
		techxMsg := (&TechxMsg{}).Deserialize(msg)
		if techxMsg != nil {
			if techxMsg.MessageType == "REQUEST" || techxMsg.MessageType == "MESSAGE" {
				return techxMsg.Uri, Serialize(techxMsg.Data), techxMsg.ResponseDestination.Topic, techxMsg.TransactionId
			}
			return "", "techx-msg-but-not-request", "", ""
		}
		return "", "not-techx-msg", "", ""
	} else {
		log.Fatal("MSG_TEMPLATE is not set. exit.")
		return "", "", "", ""
	}
}
