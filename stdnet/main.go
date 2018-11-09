package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	f, err := os.OpenFile("logFile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
    log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	if len(os.Args) < 3 {
		log.Fatal("Missing args.")
	}

	remoteIP := os.Args[1]
	remotePort := os.Args[2]
	remoteAddr := fmt.Sprintf("%s:%s", remoteIP, remotePort)

	log.Printf("Connecting to %s", remoteAddr)

	rAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}

	rConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		panic(err)
	}
	defer rConn.Close()

	tcp_con_handle(rConn)
}

// Handles TC connection and perform synchorinization:
// TCP -> Stdout and Stdin -> TCP
func tcp_con_handle(con net.Conn) {
	chan_to_stdout := stream_copy(con, os.Stdout)
	chan_to_remote := stream_copy(os.Stdin, con)
	select {
	case <-chan_to_stdout:
		log.Println("Remote connection is closed")
	case <-chan_to_remote:
		log.Println("Local program is terminated")
	}
}

// Performs copy operation between streams: os and tcp streams
func stream_copy(src io.Reader, dst io.Writer) <-chan int {
	sync_channel := make(chan int)
	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				log.Printf("Connection from %v is closed\n", con.RemoteAddr())
			}
			sync_channel <- 0 // Notify that processing is finished
		}()
		piper, pipew := io.Pipe()
		go func() {
			defer pipew.Close()
			io.Copy(pipew, src)
		}()
		io.Copy(dst, piper)
		piper.Close()
	}()
	return sync_channel
}
