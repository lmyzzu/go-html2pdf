package main

import (
	"fmt"
	"os"
	"runtime"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"bytes"
	"bufio"
	"net/url"
	"github.com/gorilla/schema"
	"io/ioutil"
	"log"
)

var decoder  = schema.NewDecoder()

type tplVars struct {
	Title string `schema:"title"`
}

func getBinaryPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if runtime.GOOS == "darwin" {
		return dir + "/bin/wkhtmltopdf_macOs"
	}
	return dir + "/bin/wkhtmltopdf_linux_amd64"
}

func billToPDF(paramsString url.Values) []byte {
	var vars tplVars
	// turning GET params into template varables
	err := decoder.Decode(&vars, paramsString)
	if err != nil {
		log.Println("Error in GET parameters : ", err)
	}
	// set wkhtmltopdf binary path
	wkhtmltopdf.SetPath(getBinaryPath())

	fmt.Printf("Using path: %v\n", wkhtmltopdf.GetPath())

	pdf, _ := wkhtmltopdf.NewPDFGenerator()
	pdf.PageSize.Set(wkhtmltopdf.PageSizeA4)

	text, _ := ioutil.ReadFile("templates/bill.html")
	tpl, err := template.New("template pdf").Parse(string(text))
	if err != nil {
		panic(err)
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	tpl.Execute(buf, vars)
	r := bufio.NewReader(buf)

	page := wkhtmltopdf.NewPageReader(r)
	page.NoImages.Set(true)
	page.DisableJavascript.Set(true)
	page.DisableExternalLinks.Set(true)
	page.DisableInternalLinks.Set(true)
	pdf.AddPage(page)
	//fmt.Printf("args: %v\n", pdf.Args())
	
	if err := pdf.Create(); err != nil {
		panic(err)
	}
	//ioutil.WriteFile("out/result.pdf", result, 0777)
	return pdf.Bytes()
}

func toPDF(fileData []byte) []byte {

	// set binary path
	wkhtmltopdf.SetPath(getBinaryPath())
	fmt.Printf("Using path: %v\n", wkhtmltopdf.GetPath())

	pdf, _ := wkhtmltopdf.NewPDFGenerator()
	pdf.PageSize.Set(wkhtmltopdf.PageSizeA4)

	buf := bytes.NewBuffer(fileData)
	r := bufio.NewReader(buf)
	page := wkhtmltopdf.NewPageReader(r)
	//page.NoImages.Set(true)
	page.DisableJavascript.Set(true)
	page.DisableExternalLinks.Set(true)
	page.DisableInternalLinks.Set(true)
	pdf.AddPage(page)
	fmt.Printf("args: %v\n", pdf.Args())

	if err := pdf.Create(); err != nil {
		panic(err)
	}

	return pdf.Bytes()
}

