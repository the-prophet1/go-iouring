package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var params IoUringParams
	ring, err := NewIoUringQueueParams(4, &params)
	must(err)
	data, _ := json.Marshal(ring)
	paramsData, _ := json.Marshal(params)
	fmt.Println(string(data))
	fmt.Println(string(paramsData))
}
