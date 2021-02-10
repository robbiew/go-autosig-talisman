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
	updatedSig string
	sigPipes   string
	currentSig User
	exit       int

	reset = "\u001b[0m"

	// Foreground ANSI colors
	fgBlack   = reset + "\u001b[30m"
	fgRed     = reset + "\u001b[31m"
	fgGreen   = reset + "\u001b[32m"
	fgYellow  = reset + "\u001b[33m"
	fgBlue    = reset + "\u001b[34m"
	fgMagenta = reset + "\u001b[35m"
	fgCyan    = reset + "\u001b[36m"
	fgWhite   = reset + "\u001b[37m"

	// Foreground ANSI colors, bright
	fgBlackBr   = reset + "\u001b[30;1m"
	fgRedBr     = reset + "\u001b[31;1m"
	fgGreenBr   = reset + "\u001b[32;1m"
	fgYellowBr  = reset + "\u001b[33;1m"
	fgBlueBr    = reset + "\u001b[34;1m"
	fgMagentaBr = reset + "\u001b[35;1m"
	fgCyanBr    = reset + "\u001b[36;1m"
	fgWhiteBr   = reset + "\u001b[37;1m"

	// Background ANSI colors
	bgBlack   = reset + "\u001b[40m"
	bgRed     = reset + "\u001b[41m"
	bgGreen   = reset + "\u001b[42m"
	bgYellow  = reset + "\u001b[43m"
	bgBlue    = reset + "\u001b[44m"
	bgMagenta = reset + "\u001b[45m"
	bgCyan    = reset + "\u001b[46m"
	bgWhite   = reset + "\u001b[47m"
)

// User struct for the database
type User struct {
	uid    int
	attrib string
	value  string
}

func checkError(err error) {

	if err != nil {
		log.Println(err)
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

func addNewLines(signature string) string {
	r := strings.NewReplacer(
		" ", "|| char(10) ||",
	)
	return r.Replace(signature)
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
		"\n", "\r\n ")

	return r.Replace(currentSig)
}

func getUsers(db *sql.DB, id2 int) User {

	rows, err := db.Query(`select * from details where attrib = 'signature'`)
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		var tempUser User
		err =
			rows.Scan(&tempUser.uid, &tempUser.attrib, &tempUser.value)

		checkError(err)
		if tempUser.uid == id2 {
			return tempUser
		}
		rows.Close()
	}
	return User{}
}

func updateUser(db *sql.DB, id int, value string, attrib string) {

	stmt, _ := db.Prepare(`update details set value=? where uid=? AND attrib=?`)
	_, err := stmt.Exec(updatedSig, id, attrib)
	checkError(err)
}

func dropFileData() {

	node := os.Args[1]

	file, err := os.Open("/home/robbiew/bbs/temp/" + node + "/door.sys")
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

	buf := make([]byte, 1024)
	reqLen, err := os.Stdin.Read(buf)
	if err != nil {
		errorChan <- err
		return
	}
	dataChan <- buf[:reqLen]
}

func main() {

	go dropFileData()
	time.Sleep(100 * time.Millisecond)
	db, _ := sql.Open("sqlite3", "/home/robbiew/bbs/data/users.sqlite3") // Open the SQLite File
	errorChan := make(chan error)
	dataChan := make(chan []byte)

	menu = "main"
	dataType = "key"

	for {
		go readWrapper(dataChan, errorChan)
		currentSig := getUsers(db, id)
		sigPipes := currentSig.value
		sigEscapes := replaceColors(currentSig.value)
		fmt.Println("\033[H\033[2J")
		showArt("header")

		newSigEscapes := replaceColors(updatedSig)

		if len(newSigEscapes) > 0 {
			fmt.Printf(" NEW Auto Signature:\r\n")
			fmt.Println("\u001b[1C")
			fmt.Println(newSigEscapes + reset)
			fmt.Printf("\r\n")
			fmt.Printf(" " + fgRed + "(" + fgRedBr + "S" + fgRed + ")" + fgRedBr + " Save/Keep" + reset + "\r\n")
			fmt.Printf(" " + fgRed + "(" + fgRedBr + "Q" + fgRed + ")" + fgRedBr + " Quit/Abort" + reset + "\r\n")
			fmt.Printf("\033[?25l")
		} else {
			fmt.Printf(" Current Auto Signature:\r\n")
			fmt.Println("\u001b[1C")
			fmt.Println(sigEscapes + reset)
			fmt.Printf("\r\n")
			fmt.Printf(" " + fgRed + "(" + fgRedBr + "E" + fgRed + ")" + fgRedBr + " Edit" + reset + "\r\n")
			fmt.Printf(" " + fgRed + "(" + fgRedBr + "Q" + fgRed + ")" + fgRedBr + " Quit" + reset + "\r\n")
			fmt.Printf("\033[?25l")
		}

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
					case "S", "s":
						menu = "save"
					}
				}
			}
			if dataType == "typed" {
				if menu == "edit" {
					updatedSig = kilo.Start(sigPipes)
					dataType = "key"
					menu = "main"
				}
			}
			// fall through statement to close connection
			if menu == "quit" {
				break
			}
			if menu == "save" {
				fmt.Printf("\r\n %vSaving %v%v%v%v's Auto Signature...%v", fgRedBr, reset, fgRed, name, fgRedBr, reset)
				newLines := addNewLines(updatedSig)
				updateUser(db, id, newLines, "signature")
				time.Sleep(400 * time.Millisecond)
				log.Printf("%v (id: %v) updated signature", name, id)
				break
			}
			continue
		case err := <-errorChan:
			log.Println("An error occured:", err.Error())
			return
		case <-time.After(1 * time.Minute):
			fmt.Println("\r\n Timed out!")
			break
		}

		fmt.Printf("\r\n %vReturning to BBS...%v", fgCyan, reset)
		time.Sleep(1 * time.Second)
		break
	}
	db.Close()
	os.Exit(0)
}
