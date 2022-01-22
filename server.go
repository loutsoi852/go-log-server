package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"strings"
)

var logFileA string = "logA.log"
var logFileB string = "logB.log"
var fileSizeLimit int64 = 100

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/send", fileAppendHandler)
	r.HandleFunc("/read/{lines}", fileReadHandler)
	http.Handle("/", r)


 	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./form/index.html")
    })



	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}

func getFileDetails(file string)(m int64, s int64, f *os.File){

        f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatal(err)
        }
		stat, _ := f.Stat()
        modTime := stat.ModTime().Unix()
		size := stat.Size()
		//fmt.Printf("Type: %T \n", f)

		return modTime, size, f
}

func getTruncFile(file string) (f *os.File){
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func closeFile(f *os.File){
	fmt.Printf("Closing: %v\n",  f.Name())
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
func getLatestFile(truncIt bool)(f *os.File){

	//get latest ModTime.unix file
	//is file size less than X
		//yes return file
	
	//no then truncate the other file and return it sd


	modTimeA, sizeA, fA := getFileDetails(logFileA)
	modTimeB, sizeB, fB := getFileDetails(logFileB)
	// fmt.Printf("Type: %T Value: %v\n", modTimeA, modTimeA)
	// fmt.Printf("Type: %T Value: %v\n", sizeA, sizeA)
	// fmt.Printf("Type: %T Value: %v\n", fA, fA.Name())
	// fmt.Printf("Type: %T Value: %v\n", fB, fB.Name())

	// fmt.Println(modTimeA)
	// fmt.Println(modTimeB)
	// fmt.Println("sizeA",sizeA)
	// fmt.Println("sizeB", sizeB)

	if(modTimeA >= modTimeB){
		if(truncIt && sizeA > fileSizeLimit){
			closeFile(fA)
			closeFile(fB)
			truncF := getTruncFile(logFileB)
			fmt.Printf("TRUNC Type: %T Value: %v\n", truncF, truncF.Name())
			return truncF
		}
		closeFile(fB)
		return fA
	}else{
		if(truncIt && sizeB > fileSizeLimit){
			closeFile(fB)
			closeFile(fA)
			truncF := getTruncFile(logFileA)
			fmt.Printf("TRUNC Type: %T Value: %v\n", truncF, truncF.Name())
			return truncF
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


	fp := getLatestFile(true)
	fmt.Printf("FP Type: %T Value: %v\n", fp, fp.Name())
	if _, err := fp.Write([]byte(*t.Log)); err != nil {
		log.Fatal(err)
	}
	if _, err := fp.Write([]byte("\n")); err != nil {
		log.Fatal(err)
	}
 	closeFile(fp)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Success")

	return


	// // If the file doesn't exist, create it, or append to the file
	// f, err := os.OpenFile(logFileA, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// //f, err := os.OpenFile(logFileA, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // if _, err := f.Write([]byte("appended some data\n")); err != nil {
	// // 	log.Fatal(err)
	// // }
	// if _, err := f.Write([]byte(*t.Log)); err != nil {
	// 	log.Fatal(err)
	// }
	// if _, err := f.Write([]byte("\n")); err != nil {
	// 	log.Fatal(err)
	// }

	// stat, _ := f.Stat()
	// filesize := stat.Size()

	// if err := f.Close(); err != nil {
	// 	log.Fatal(err)
	// }

	// //do something when log file size is too big
	// if filesize > 1000 {
	// 	e := os.Rename(logFileA, logFileB)
	// 	if e != nil {
	// 		log.Fatal(e)
	// 	}

	// 	newF, err := os.OpenFile(logFileA, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if err := newF.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}

	// }

	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	// // got the input we expected: no more, no less
	// fmt.Println(*t.Log)
	// fmt.Println(filesize)

	// fmt.Fprintf(w, "Success")
}





func getLastLineWithSeek(filepath string, lineLimit int) (string, int) {
	fmt.Println("filepath",filepath)

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
		if(filesize ==0 ){
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
	//get latest file
	//fetch lines
		//return if is enough lines
		//else
			//fetch more lines from other file 


	fp := getLatestFile(false)
	var firstFile string
	var secondFile string
	if (fp.Name()==logFileA){
		firstFile = logFileA
		secondFile = logFileB
	}else{
		firstFile = logFileB
		secondFile = logFileA
	}
	//closeFile(fp)

	fmt.Printf("Type: %T Value: %v\n", firstFile, firstFile)
	fmt.Println("fp.Name",fp.Name())

	fmt.Println("firstFile",firstFile)
	fmt.Println("secondFile",secondFile)

	s, count := getLastLineWithSeek(firstFile,lineLimit)
	fmt.Println("count", count)
	fmt.Println("lineLimit", lineLimit)
	
	var secondLineLimit int
	if(count<lineLimit){
		secondLineLimit = lineLimit - count
		fmt.Println("secondLineLimit", secondLineLimit)
		s2, _ := getLastLineWithSeek(secondFile,secondLineLimit)
		return s2+s
	}else{
		return s

	}




}






func fileReadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------start---------------")

	vars := mux.Vars(r)
    lines := vars["lines"]
 
	var lineLimit int
	if _, err := fmt.Sscanf(lines, "%5d", &lineLimit); err != nil {
		panic(err)
	}
	s := fetchLines(lineLimit)

	ss := strings.Split(s, "\n")
	// fmt.Println(ss)

	//fmt.Println("s", s)
	j, err := json.Marshal(ss)
    if err != nil {
        fmt.Printf("Error: %s", err.Error())
    }  
	fmt.Fprintf(w, "{\"data\":"+string(j)+"}")
	fmt.Println("-------------end---------------")
}