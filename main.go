package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

// Client struct holds info about the user
type Client struct {
	dataType   int    // 0 = single key press, 1 = typing
	menuCurr   int    // 0 = main, 1 = edit sig
	username   string // Grab from DOOR.SYS
	userid     int    // Grab from DOOR.SYS
	nodeNumber int    // Grab from DOOR.SYS
}

// NewClient allows us to update user info
func NewClient() *Client {
	return &(Client{})
}

func dropData() {

	file, err := os.Open("/home/robbiew/bbs/temp/1/door.sys")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	count := 0
	for _, eachLn := range text {

		if count == 3 {
			fmt.Printf("Node: %v\r\n", eachLn)
		}
		if count == 35 {
			fmt.Printf("User: %v\r\n", eachLn)
		}
		if count == 51 {
			break
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// Main input loop
func readWrapper(dataChan chan []byte, errorChan chan error) {

	for {
		buf := make([]byte, 1024)
		reqLen, err := os.Stdin.Read(buf)
		if err != nil {
			errorChan <- err
			return
		}
		dataChan <- buf[:reqLen]
	}
}

func main() {

	go dropData()
	time.Sleep(100 * time.Millisecond)

	errorChan := make(chan error)
	dataChan := make(chan []byte)

	go readWrapper(dataChan, errorChan)
	r := bytes.NewBuffer(make([]byte, 0, 1024))

	fmt.Printf("---------------------\r\n")
	fmt.Printf("Cmd? ")
	for {
		select {
		case data := <-dataChan:
			// Get input from user
			fmt.Printf(string(data))
			if bytes.Equal(data, []byte("\r\n")) || bytes.Equal(data, []byte("\r")) {
				fmt.Printf("you typed: %q\r\n", r.String())
				r.Reset()
				break
			}
			// ESC aborts and returns to BBS
			if bytes.Equal(data, []byte("\033")) || bytes.Equal(data, []byte("\033\r\n")) || bytes.Equal(data, []byte("\033\r")) || bytes.Equal(data, []byte("\033\n")) {
				fmt.Printf("\r\nAborted!\r\n")
				r.Reset()
				break
			}
			r.Write(data)
			// otherwise continue printing menu for invalid submissions
			continue

		case err := <-errorChan:
			log.Println("An error occured:", err.Error())
			return
		}
		break
	}

	fmt.Printf("\r\nReturning...\r\n")
	time.Sleep(500 * time.Millisecond)

}
