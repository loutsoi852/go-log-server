package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var logfile string = "access.log"
var logfile2 string = "access-1.log"

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", http.FileServer(http.Dir("./form")))
	//http.HandleFunc("/hello", helloHandler)
	http.Handle("/form", fileServer)
	http.HandleFunc("/sendForm", formHandler)
	http.HandleFunc("/send", fileAppendHandler)
	http.HandleFunc("/read", fileReadHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}

func fileAppendHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/send" {
		http.Error(w, "404 nope", http.StatusNotFound)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	// anonymous struct type: handy for one-time use
	t := struct {
		Log *string `json:"log"` // pointer so we can test for field absence
	}{}

	err := d.Decode(&t)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Log == nil {
		http.Error(w, "missing field 'log' from JSON object", http.StatusBadRequest)
		return
	}

	// optional extra check
	if d.More() {
		http.Error(w, "extraneous data after JSON object", http.StatusBadRequest)
		return
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// if _, err := f.Write([]byte("appended some data\n")); err != nil {
	// 	log.Fatal(err)
	// }
	if _, err := f.Write([]byte(*t.Log)); err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		log.Fatal(err)
	}

	stat, _ := f.Stat()
	filesize := stat.Size()

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	//do something when log file size is too big
	if filesize > 1000 {
		e := os.Rename(logfile, logfile2)
		if e != nil {
			log.Fatal(e)
		}

		newF, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		if err := newF.Close(); err != nil {
			log.Fatal(err)
		}

	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// got the input we expected: no more, no less
	fmt.Println(*t.Log)
	fmt.Println(filesize)

	fmt.Fprintf(w, "Success")
}

func getLastLineWithSeek(filepath string, totalLines int) string {
	fileHandle, err := os.Open(filepath)

	if err != nil {
		panic("Cannot open file")
		os.Exit(1)
	}
	defer fileHandle.Close()

	lines := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()

	var count int = 0

	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)
		char := make([]byte, 1)
		fileHandle.Read(char)
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			count += 1
			if count == totalLines {
				break
			}
		}
		lines = fmt.Sprintf("%s%s", string(char), lines)
		if cursor == -filesize {
			break
		}
	}

	// fmt.Println("io.SeekEnd", io.SeekEnd)
	// fmt.Printf("slice len=%d cap=%d %v\n", len(slice), cap(slice), slice)
	// fmt.Println("count", count)
	// fmt.Println("cursor", cursor)
	// fmt.Println("lines", lines)

	return lines
}

func fileReadHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/fileRead" {
	// 	http.Error(w, "404 nope", http.StatusNotFound)
	// 	return
	// }

	var totalLines string
	lines, ok := r.URL.Query()["lines"]
	if !ok || len(lines[0]) < 1 {
		totalLines = "10"
	} else {
		totalLines = lines[0]
	}

	var i int
	if _, err := fmt.Sscanf(totalLines, "%5d", &i); err == nil {
		fmt.Fprintf(w, string(getLastLineWithSeek(logfile, i)))
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/hello" {
// 		http.Error(w, "404 not found.", http.StatusNotFound)
// 		return
// 	}

// 	if r.Method != "GET" {
// 		http.Error(w, "Method is not supported.", http.StatusNotFound)
// 		return
// 	}

// 	fmt.Fprintf(w, "Hello!")
// }
