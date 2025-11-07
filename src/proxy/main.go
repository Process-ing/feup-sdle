package main

import (
	"fmt"
	"github.com/zeromq/goczmq"
)

func main() {
	major := goczmq.CZMQVersionMajor
	minor := goczmq.CZMQVersionMinor
	fmt.Printf("ZeroMQ Version: %d.%d\n", major, minor)
}