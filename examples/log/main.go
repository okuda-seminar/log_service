package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"log_service/pkg/logger"
	"log_service/pkg/rabbitmq"
)

func init() {
	log.SetFlags(0)
}

func main() {
	conn, ch, err := rabbitmq.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	defer ch.Close()

	ctx := context.Background()
	w, err := logger.NewWritter(ctx, ch)
	if err != nil {
		panic(err)
	}
	log.SetOutput(w)

	num := 10
	var sumSec int64 = 0
	for i := 0; i < num; i++ {
		t := time.Now()
		log.Println("Hello, World!")
		sec := time.Since(t).Milliseconds()
		sumSec += sec
		fmt.Printf("%vms\n", sec)
	}

	fmt.Printf("Average: %vms\n", int(sumSec)/num)
}
