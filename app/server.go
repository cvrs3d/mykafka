package main

import (
	"fmt"
	"net"
	"os"
)

// // Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
// var _ = net.Listen
// var _ = os.Exit

func main() {
	
	// Bind to port 9092

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			continue
		}

		go func() {
			defer conn.Close()
			// Read 4 bytes from the connection
			buffer := make([]byte, 4)
			_, err := conn.Read(buffer) // Read 4 bytes from the connection
			if err != nil {
				fmt.Println("Failed to read from connection")
				return
			}
			// Convert the 4 bytes to a uint32
			length := uint32(buffer[0]) << 24 | uint32(buffer[1]) << 16 | uint32(buffer[2]) << 8 | uint32(buffer[3])
			// Read "length" bytes from the connection
			buffer = make([]byte, length)
			_, err = conn.Read(buffer)
			if err != nil {
				fmt.Println("Failed to read from connection")
				return
			}
			// Print the bytes read from the connection
			// first 4 bytes are the length of the message
			// next 2 bytes are the reqeust_api_key
			// next 2 bytes are the request_api_version
			// next 4 bytes are the correlation_id

			// read correlation_id
			correlationId := uint32(buffer[4]) << 24 | uint32(buffer[5]) << 16 | uint32(buffer[6]) << 8 | uint32(buffer[7])
			// encode correlation_id as bidendian
			var response = []byte{0,0,0,0,0,0,0,0}
			response[4] = byte(correlationId >> 24)
			response[5] = byte(correlationId >> 16)
			response[6] = byte(correlationId >> 8)
			response[7] = byte(correlationId)
			requestApiVersion := int16(buffer[2]) << 8 | int16(buffer[3])
			if requestApiVersion < 0 || requestApiVersion > 4 {
				response = append(response, 0, 35)
			} else {
				response = append(response, 0, 0)
			}
			conn.Write(response)
		} ()
	}
}
