package main

import (
	"fmt"
	"log"

	"mo.io/conditional_order/pkg/cache"
	"mo.io/conditional_order/pkg/db"
	mokafka "mo.io/conditional_order/pkg/mq/kafka"
	momsg "mo.io/conditional_order/pkg/msg"
)

func main() {
	//db
	repo := db.GetStopOrderRepo()
	var records []db.StopOrder
	err := repo.Get(&records)
	if err != nil {
		log.Fatal("Error querying all records in stop order repo: ", err)
	}
	log.Printf("Records: %v", records)
	//cache
	redis := cache.GetCache()
	mo := redis.Get("mo")
	if mo == nil {
		redis.Set("mo", "3")
		mo = redis.Get("mo")
	}
	log.Printf("mo: %v", mo)
	//message queue
	driver := mokafka.NewDriver("test", "groupTest")
	go driver.ConsumeRes()
	go driver.ConsumeReq(handler)
	select {}
}
func handler(driver *mokafka.Driver, msg string) {
	uri, data, responseTopic, id := momsg.GetReqData(string(msg))
	// log.Printf("uri - data extracted: %v - %v. after processing, need to send to topic: %v with id %v", uri, data, responseTopic, id)
	if responseTopic == "" {
		log.Println("this msg doesnt need a response")
		return
	}
	if uri == "/trigger" {
		response := driver.ProduceAndGetResponse("middleware", "crap", fmt.Sprintf("%v, craphihihi", data))
		if response == "" {
			log.Println("Ignore nonsense msg, not standard format")
		} else {
			log.Printf("Msg From middleware: %v", response)
			driver.ProduceResponse(id, responseTopic, fmt.Sprintf("msg processed: ~~~~%v~~~~", response))
		}
	}
}
