package main

import (
    "fmt"
    "log"
    "net/http"
	"os"
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

func fileReadHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path !="/fileRead"{
		http.Error(w, "404 nope", http.StatusNotFound)
		return
	}

 	data, err := os.ReadFile("access.log")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(data)

    fmt.Fprintf(w, string(data))

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