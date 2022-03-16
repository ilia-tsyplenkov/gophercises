package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/ilia-tsyplenkov/gophercises/nc"
)

const addr = ":8080"

func main() {
	server := nc.NewServer(addr, "tcp")
	_, err := server.Listen()
	if err != nil {
		panic(err)
	}
	log.Printf("start listening: %s.....\n", addr)
	_, err = server.Accept()
	if err != nil {
		panic(err)
	}
	for {
		buffer := bufio.NewReader(server)
		received, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Print(received)
	}
}
