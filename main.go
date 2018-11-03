package main

import (
	"log"

	"github.com/alaija/gocodelabru/api"
)

func main() {
	a := api.New(":9111")
	log.Fatal(a.Start())
}
