package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
	"bytes"
)

func text2pdf(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
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
	//name := strings.Split(header.Filename, ".")
	//fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// Reset the buffer to reduce memory allocation
	defer Buf.Reset()
	//contents := Buf.String()
	//fmt.Println(contents)
	//	pdfBytes := toPDF(Buf.Bytes())
	pdfBytes := toPDF(Buf.Bytes())

	// Write the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(pdfBytes)

}

func main() {

	port := 8080
	http.HandleFunc("/text2pdf", text2pdf)
	http.HandleFunc("/file2pdf", file2pdf)
	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}
