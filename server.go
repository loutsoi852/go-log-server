package main

import (
    "fmt"
    "log"
    "net/http"
	"os"
	"io"
	
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/hello" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }

    fmt.Fprintf(w, "Hello!")
}

func fileAppendHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path !="/fileAppend"{
		http.Error(w, "404 nope", http.StatusNotFound)
		return
	}

    // If the file doesn't exist, create it, or append to the file
    f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write([]byte("appended some data\n")); err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }

    fmt.Fprintf(w, "Done!")

}

func getLastLineWithSeek(filepath string) string {
    fileHandle, err := os.Open(filepath)

    if err != nil {
        panic("Cannot open file")
        os.Exit(1)
    }
    defer fileHandle.Close()

    line := ""
    var cursor int64 = 0
    stat, _ := fileHandle.Stat()
    filesize := stat.Size()
    for { 
        cursor -= 1
        fileHandle.Seek(cursor, io.SeekEnd)

        char := make([]byte, 1)
        fileHandle.Read(char)

        if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
            break
        }

        line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

        if cursor == -filesize { // stop if we are at the begining
            break
        }
    }

    return line
}

func fileReadHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path !="/fileRead"{
		http.Error(w, "404 nope", http.StatusNotFound)
		return
	}

 	// data, err := os.ReadFile("access.log")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// os.Stdout.Write(data)

    // fmt.Fprintf(w, string(data))

    fmt.Fprintf(w, string(getLastLineWithSeek("access.log")))
    fmt.Fprintf(w, "Done!")
    fmt.Fprintf(w, "Done!")


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

 func main() {
	fileServer := http.FileServer(http.Dir("./static")) // New code
    http.Handle("/", fileServer) // New code
    http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/fileAppend", fileAppendHandler)
	http.HandleFunc("/fileRead", fileReadHandler)

	fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
        log.Fatal(err)
    }
}