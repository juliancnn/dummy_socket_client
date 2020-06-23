package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Path to unix socket
var sockPath string
// True if you want to close the socket after sending the data
var notClose bool

func main() {

	// Number of conexion (And goroutines)
	var numGorun uint
	// Time before launch a new goroutines in ms
	var timeSleep uint
	// True if you want wait sigterm for close
	var waitTime uint

	flag.StringVar(&sockPath, "f", "./echo.sock", "Unix socket path")
	flag.UintVar(&numGorun, "n", 200, "Number of conexion (And goroutines)")
	flag.UintVar(&timeSleep, "t", 2, "Time before launch a new goroutines (in ms)")
	flag.BoolVar(&notClose, "u", false, "Don't close socket after send data")
	flag.UintVar(&waitTime, "w", 0, "Wait time between the data was sent and "+
		"the application closes (seconds)")
	flag.Parse()

	// Sem for sync
	semEndRoutine := make(chan bool, numGorun)
	semSendOk := make(chan bool, numGorun)

	// go go go!
	for i := uint(0); i < numGorun; i++ {
		go dial(i, semEndRoutine, semSendOk)
		time.Sleep(time.Duration(timeSleep) * time.Millisecond)
	}


	// Wait data send
	for i := uint(0); i < numGorun; i++ {
		<- semSendOk
	}
	fmt.Printf("All data was send\n")

	// Sleep before close
	time.Sleep(time.Duration(waitTime) * time.Second)

	// Send end signal to goroutines
	for i := uint(0); i < numGorun; i++ {
		semEndRoutine <- false
	}

	fmt.Printf("Bye!\n")
}

func dial(threadNumber uint, semEnd chan bool, semSendOk chan bool) {
	var strToSend string

	// Resolver and conect to socket
	addr, err := net.ResolveUnixAddr("unix", sockPath)
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}

	// Create dammy data and send
	strToSend = fmt.Sprintf("[goroutine N: %d ] Test message\n", threadNumber)
	if i, err := conn.Write([]byte(strToSend)); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success, sent %v bytes ** %s **\n", i, strings.TrimSuffix(strToSend, "\n"))
	}

	// Notify data send ok
	semSendOk <- false

	//Close socket?
	if false == notClose {
		conn.Close()
	}

	// Wait end
	<- semEnd
	// Force not close of go garb collector
	if true == notClose {
		conn.Write([]byte("Bye!"))
		conn.Close()
	}


}
