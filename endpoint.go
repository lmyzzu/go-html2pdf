package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
	"bytes"
	"os"
	"strconv"
)

func text2pdf(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	pdfBytes := toPDF([]byte(r.PostForm.Get("htmlstring")))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(pdfBytes)
}

func file2pdf(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer
	file, header, err := r.FormFile("file")
	_ = header
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%v\n", err)))
	}
	defer file.Close()

	// Copy the file data to a buffer
	io.Copy(&Buf, file)

	// Reset the buffer to reduce memory allocation
	defer Buf.Reset()

	pdfBytes := toPDF(Buf.Bytes())

	// Write the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(pdfBytes)

}

func main() {
	var port uint16
	if len(os.Args) < 2 {
		// default port
		port = 8080
	} else if i, err := strconv.ParseUint(os.Args[1], 10, 16); err == nil {
		// is it really the correct way to do this ?!
		port = uint16(i)
	}

	// 2 endpoints
	http.HandleFunc("/text2pdf", text2pdf)
	http.HandleFunc("/file2pdf", file2pdf)
	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}
