package main

import (
	"fmt"
	"log"
	"unsafe"
)

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

var userdata uint64

func main() {
	s := "data"
	SetData(s)
	data := GetData()
	fmt.Println(data)
}

func SetData(data interface{}) {
	dataPtr := unsafe.Pointer(&data)
	userdata = uint64(uintptr(dataPtr))
}

func GetData() interface{} {
	dataPtr := unsafe.Pointer(uintptr(userdata))
	return *(*interface{})(dataPtr)
}
