package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robbiew/autosig/kilo"
)

var (
	name       string
	id         int
	menu       string
	dataType   string
	newSig     string
	sigPipes   string
	currentSig User
	exit       int

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

func showArt(menu string) {

	var b bytes.Buffer
	file := (menu + ".ans")
	art, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}
	b.Write([]byte(art))
	b.WriteTo(os.Stdout)
	fmt.Printf("\r\n")
	return
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
		"\r", "\r\n  ")

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
			name = eachLn
		}
		if count == 25 {

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

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func anyoneThere(t time.Time) {

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Size() > 0 {
		fmt.Println("\u001b[1;1HSomething here")
	} else {
		fmt.Println("\u001b[1;1Hempty")
	}

}

func main() {

	// go doEvery(1*time.Second, anyoneThere)
	go dropFileData()
	time.Sleep(100 * time.Millisecond)
	db, _ := sql.Open("sqlite3", "/home/robbiew/bbs/data/users.sqlite3") // Open the SQLite File
	errorChan := make(chan error)
	dataChan := make(chan []byte)

	menu = "main"
	dataType = "key"

	go readWrapper(dataChan, errorChan)

	for {
		currentSig := getUsers(db, id)
		sigPipes := sigWithPipes(currentSig.value)
		sigEscapes := replaceColors(currentSig.value)
		fmt.Println("\033[H\033[2J")
		showArt("header")
		fmt.Printf(" \u001b[30;1m\u001b[0m+-------------------------------------------------\u001b[0m+\r\n")
		fmt.Println("\u001b[2C")
		fmt.Println(sigEscapes)
		fmt.Printf("\u001b[1D\u001b[30;1m\u001b[0m+-------------------------------------------------\u001b[0m+\r\n\n")
		fmt.Printf(" \u001b[31m(\u001b[31;1mE\u001b[0m\u001b[31m) \u001b[31mEdit\u001b[0m\r\n")
		fmt.Printf(" \u001b[31m(\u001b[31;1mQ\u001b[0m\u001b[31m) \u001b[31mQuit\u001b[0m\r\n")
		fmt.Printf("\033[?25l")

		select {
		case data := <-dataChan:
			t := strings.TrimSuffix(strings.TrimSuffix(string(data), "\r\n"), "\n")
			if dataType == "key" {
				if menu == "main" {
					switch t {
					case "Q", "q":
						menu = "quit"
					case "E", "e":
						dataType = "typed"
						menu = "edit"
					}
				}
			}
			if dataType == "typed" {
				if menu == "edit" {
					kilo.Start(sigPipes)
					dataType = "key"
					menu = "main"
				}
			}
			// fall through statement to close connection
			if menu == "quit" {
				break
			}
			continue
		case err := <-errorChan:
			log.Println("An error occured:", err.Error())
			return
		}
		fmt.Println("\r\nClosing")
		time.Sleep(500 * time.Millisecond)
		break
	}
	os.Exit(0)
}
