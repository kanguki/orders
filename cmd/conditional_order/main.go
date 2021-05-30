package main

import (
	"log"

	"mo.io/conditional_order/pkg/cache"
	"mo.io/conditional_order/pkg/db"
	mokafka "mo.io/conditional_order/pkg/mq/kafka"
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
	 mokafka.NewDriver("test", "test")
	// go func(){
	// 	response := driver.ProduceAndGetResponse("tes1", "")
	// 	if response == "" {
	// 		log.Println("Ignore nonsense msg, not standard format")
	// 	} else {
	// 		log.Printf("Receive msg %v", response)
	// 	}	
	// }()

	select {}
}
