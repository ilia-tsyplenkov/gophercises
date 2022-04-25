package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ilia-tsyplenkov/gophercises/nc"
)

var serverMode bool
var udpMode bool
var connectionType string = "tcp"
var host, port string

func init() {
	flag.BoolVar(&serverMode, "l", false, "Listening for an incoming connection.")
	flag.StringVar(&host, "h", "", "destination")
	flag.StringVar(&port, "p", "", "port")
}

func main() {
	flag.Parse()
	addr := host + ":" + port
	if serverMode {
		server := nc.NewServer(addr, connectionType)
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

	} else {
		client := nc.NewClient(addr, connectionType)
		conn, err := client.Dial()
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(os.Stdin)
		var readErr error
		var msg string
		for {
			if msg, readErr = reader.ReadString('\n'); readErr == nil {
				fmt.Fprint(client, msg)
			} else {
				if readErr != io.EOF {
					log.Fatal(readErr)
				}
				break
			}
		}
	}

}
