package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const logFileA string = "logA.log"
const logFileB string = "logB.log"
const fileSizeLimit int64 = 300000
const maxWSConnections = 10

var upgrader = websocket.Upgrader{}

type Log struct {
	Time string `json:"time"`
	Log  string `json:"log"`
}

type wsConn struct {
	status bool
	conn   websocket.Conn
}

var wsCons [10]wsConn

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/send", fileAppendHandler)
	r.HandleFunc("/read/{lines}", fileReadHandler)
	http.HandleFunc("/liveLogs", liveLogs)
	http.Handle("/", r)

	r.PathPrefix("/test").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./form.html")
	})

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./wsClient.html")
	})

	fmt.Printf("Starting server at port 7777\n")
	if err := http.ListenAndServe("127.0.0.1:7777", nil); err != nil {
		log.Fatal(err)
	}
}

func closeConn(conn *websocket.Conn, index int) {
	conn.Close()
	wsCons[index].status = false
}

func liveLogs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}

	var index int
	for i, wsc := range wsCons {
		if !wsc.status {
			wsCons[i].status = true
			wsCons[i].conn = *conn
			index = i
			break
		}
	}

	defer closeConn(conn, index)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket:", err)
			break
		}
	}
}

func getFileDetails(file string) (m int64, s int64, f *os.File) {

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	stat, _ := f.Stat()
	modTime := stat.ModTime().UnixNano()
	size := stat.Size()
	return modTime, size, f
}

func getTruncFile(file string) (f *os.File) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
func getLatestFile(truncIt bool) (f *os.File) {
	modTimeA, sizeA, fA := getFileDetails(logFileA)
	modTimeB, sizeB, fB := getFileDetails(logFileB)

	if modTimeA >= modTimeB {
		if truncIt && sizeA > fileSizeLimit {
			closeFile(fA)
			closeFile(fB)
			return getTruncFile(logFileB)
		}
		closeFile(fB)
		return fA
	} else {
		if truncIt && sizeB > fileSizeLimit {
			closeFile(fB)
			closeFile(fA)
			return getTruncFile(logFileA)
		}
		closeFile(fA)
		return fB
	}

}

func fileAppendHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/send" {
		http.Error(w, "404 nope", http.StatusNotFound)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	t := struct {
		Log *string `json:"log"`
	}{}

	err := d.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Log == nil {
		http.Error(w, "missing field 'log' from JSON object", http.StatusBadRequest)
		return
	}

	if d.More() {
		http.Error(w, "extraneous data after JSON object", http.StatusBadRequest)
		return
	}

	fp := getLatestFile(true)
	now := time.Now().UTC().UnixNano()

	data := Log{
		Time: strconv.FormatInt(now, 10),
		Log:  *t.Log,
	}

	file, _ := json.Marshal(data)

	if _, err := fp.Write(file); err != nil {
		log.Fatal(err)
	}

	if _, err := fp.Write([]byte("\n")); err != nil {
		log.Fatal(err)
	}
	closeFile(fp)

	for _, wsc := range wsCons {
		if wsc.status {
			err = wsc.conn.WriteMessage(1, file)
			if err != nil {
				log.Println("write failed:", err)
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Success")

}

func getLastLineWithSeek(filepath string, lineLimit int) (string, int) {

	fileHandle, err := os.Open(filepath)

	if err != nil {
		panic("Cannot open file")
	}
	defer fileHandle.Close()

	lines := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()

	var count int = 0
	for {
		if filesize == 0 {
			break
		}
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)
		char := make([]byte, 1)
		fileHandle.Read(char)
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			count += 1
			if count == lineLimit {
				break
			}
		}
		lines = fmt.Sprintf("%s%s", string(char), lines)
		if cursor == -filesize {
			count += 1
			break
		}
	}

	return lines, count
}

func fetchLines(lineLimit int) string {
	fp := getLatestFile(false)
	var firstFile string
	var secondFile string
	if fp.Name() == logFileA {
		firstFile = logFileA
		secondFile = logFileB
	} else {
		firstFile = logFileB
		secondFile = logFileA
	}
	//closeFile(fp)

	s, count := getLastLineWithSeek(firstFile, lineLimit)

	var secondLineLimit int
	if count < lineLimit {
		secondLineLimit = lineLimit - count
		s2, _ := getLastLineWithSeek(secondFile, secondLineLimit)
		return s2 + s
	} else {
		return s

	}
}

func fileReadHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	lines := vars["lines"]

	var lineLimit int
	if _, err := fmt.Sscanf(lines, "%5d", &lineLimit); err != nil {
		panic(err)
	}
	s := fetchLines(lineLimit)
	ss := strings.Split(s, "\n")

	j, err := json.Marshal(ss)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "{\"data\":"+string(j)+"}")
}
