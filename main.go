package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robbiew/autosig/kilo"
)

var (
	name   string
	id     int
	menu   = 1
	newSig string

	reset = "\u001b[0m"

	// Foreground ANSI colors
	fgBlack   = "\u001b[30m"
	fgRed     = "\u001b[31m"
	fgGreen   = "\u001b[32m"
	fgYellow  = "\u001b[33m"
	fgBlue    = "\u001b[34m"
	fgMagenta = "\u001b[35m"
	fgCyan    = "\u001b[36m"
	fgWhite   = "\u001b[37m"

	// Foreground ANSI colors, bright
	fgBlackBr   = "\u001b[30;1m"
	fgRedBr     = "\u001b[31;1m"
	fgGreenBr   = "\u001b[32;1m"
	fgYellowBr  = "\u001b[33;1m"
	fgBlueBr    = "\u001b[34;1m"
	fgMagentaBr = "\u001b[35;1m"
	fgCyanBr    = "\u001b[36;1m"
	fgWhiteBr   = "\u001b[37;1m"

	// Background ANSU colors
	bgBlack   = "\u001b[40m"
	bgRed     = "\u001b[41m"
	bgGreen   = "\u001b[42m"
	bgYellow  = "\u001b[43m"
	bgBlue    = "\u001b[44m"
	bgMagenta = "\u001b[45m"
	bgCyan    = "\u001b[46m"
	bgWhite   = "\u001b[47m"
)

// User struct for the database
type User struct {
	uid    int
	attrib string
	value  string
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
	// catch db errors
}

func replaceColors(currentSig string) string {

	r := strings.NewReplacer(
		"|00", fgBlack,
		"|01", fgBlue,
		"|02", fgGreen,
		"|03", fgCyan,
		"|04", fgRed,
		"|05", fgMagenta,
		"|06", fgYellow,
		"|07", fgWhite,
		"|08", fgBlackBr,
		"|09", fgBlueBr,
		"|10", fgGreenBr,
		"|11", fgCyanBr,
		"|12", fgRedBr,
		"|13", fgMagentaBr,
		"|14", fgYellowBr,
		"|15", fgWhiteBr,
		"\r", "\r\n")

	return r.Replace(currentSig)

}

func sigWithPipes(currentSig string) string {

	r := strings.NewReplacer(
		"\r", "\r\n")

	return r.Replace(currentSig)

}

func getUsers(db *sql.DB, id2 int) User {
	rows, err := db.Query(`select * from details where attrib = 'signature'`)
	checkError(err)
	for rows.Next() {
		var tempUser User
		err =
			rows.Scan(&tempUser.uid, &tempUser.attrib, &tempUser.value)

		checkError(err)
		if tempUser.uid == id2 {
			return tempUser
		}
	}
	return User{}
}

func dropFileData() {

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

		if count == 35 {
			// fmt.Printf("Name: %v\r\n", eachLn)
			name = eachLn

		}
		if count == 25 {
			// fmt.Printf("Id: %v\r\n", eachLn)
			idInt, err := strconv.Atoi(eachLn)
			if err != nil {
				fmt.Println(err)
			}
			id = idInt
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

	go dropFileData()
	time.Sleep(100 * time.Millisecond)

	db, _ := sql.Open("sqlite3", "/home/robbiew/bbs/data/users.sqlite3") // Open the SQLite File

	currentSig := getUsers(db, id)

	fmt.Println("\033[H\033[2J")
	fmt.Printf("Your current Auto Signature:\r\n\n")

	sigPipes := sigWithPipes(currentSig.value)

	sigEscapes := replaceColors(currentSig.value)
	fmt.Println(sigEscapes)

	fmt.Println("\u001b[0m")

	errorChan := make(chan error)
	dataChan := make(chan []byte)

	go readWrapper(dataChan, errorChan)
	// r := bytes.NewBuffer(make([]byte, 0, 1024))

	fmt.Printf("(E) Edit\r\n")
	fmt.Printf("(Q) Quit\r\n\n")
	fmt.Printf("Cmd? ")

	for {
		select {
		case data := <-dataChan:
			if menu == 1 {
				t := strings.TrimSuffix(strings.TrimSuffix(string(data), "\r\n"), "\n")
				switch t {
				// default:
				// 	fmt.Println("client hit invalid key...")
				case "Q", "q":
					kilo.Start(sigPipes)
					fmt.Println("\033[H\033[2J")
					fmt.Printf("Your current Auto Signature:\r\n\n")
					replaced := replaceColors(currentSig.value)
					fmt.Println(replaced)

					fmt.Println("\u001b[0m")

					fmt.Printf("(E) Edit\r\n")
					fmt.Printf("(Q) Quit\r\n\n")
					fmt.Printf("Cmd? ")
					menu = 2
				case "E", "e":
					kilo.Start(sigPipes)
					fmt.Println("\033[H\033[2J")
					fmt.Printf("Your current Auto Signature:\r\n\n")
					replaced := replaceColors(currentSig.value)
					fmt.Println(replaced)

					fmt.Println("\u001b[0m")

					fmt.Printf("(E) Edit\r\n")
					fmt.Printf("(Q) Quit\r\n\n")
					fmt.Printf("Cmd? ")
					menu = 1
				}
				continue
			}
			if menu == 2 {
				fmt.Printf("\r\n\nReturning...\r\n")
				time.Sleep(200 * time.Millisecond)
				os.Exit(0)
			}

		case err := <-errorChan:
			log.Println("An error occured:", err.Error())
			return
		}
		break
	}
	fmt.Printf("\r\n\nReturning...\r\n")
	time.Sleep(200 * time.Millisecond)

}
