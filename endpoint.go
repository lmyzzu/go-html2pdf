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

func bill2pdf(w http.ResponseWriter, r *http.Request) {

	pdfBytes := billToPDF(r.URL.Query())
	// Write the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(pdfBytes)
}

func main() {
	// Default port
	var port uint16 = 8080
	if len(os.Args) >= 2 {
		// is it really the correct way to do this ?!
		if i, err := strconv.ParseUint(os.Args[len(os.Args) - 1], 10, 16); err == nil {
			port = uint16(i)
		}
	}

	// endpoints
	http.HandleFunc("/text2pdf", text2pdf)
	http.HandleFunc("/file2pdf", file2pdf)
	http.HandleFunc("/bill2pdf", bill2pdf)
	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}
