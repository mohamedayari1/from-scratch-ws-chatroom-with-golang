package main

import (
	"fmt"
	"log"
	"net"

	"boot.theprimeagen.tv/internal/request"
)

	


func main(){
	// readedFile, err := os.Open("messages.txt")
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
		
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		
			r, err := request.RequestFromReader(conn)
			if err != nil  {
				log.Fatal("error", "error", err)
			}
			fmt.Println("Received Request:")	
			fmt.Printf("Request Method: %s\n", r.RequestLine.Method)
			fmt.Printf("Request Target: %s\n", r.RequestLine.RequestTarget)
			fmt.Printf("HTTP Version: %s\n", r.RequestLine.HttpVersion)
		
		
		}

		
	}




