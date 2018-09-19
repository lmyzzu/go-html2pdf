package main

import (
	"fmt"
	"os"
	"runtime"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"bytes"
	"bufio"
)

type tplVars struct {
	Title string
}

func guard(err error) {
	if err != nil {
		panic(err)
	}
}


func toPDFWithTemplate(fileData []byte) []byte {

	dir, err := os.Getwd()
	guard(err)

	fmt.Printf("Running on %s", runtime.GOOS)
	var path string
	if runtime.GOOS == "darwin" {
		path = dir + "/bin/wkhtmltopdf_macOs"
	} else {
		path = dir + "/bin/wkhtmltopdf_linux_amd64"
	}

	// set wkhtmltopdf binary path
	wkhtmltopdf.SetPath(path)

	fmt.Printf("Using path: %v\n", wkhtmltopdf.GetPath())

	pdf, _ := wkhtmltopdf.NewPDFGenerator()
	pdf.PageSize.Set(wkhtmltopdf.PageSizeA4)

	//text, _ := ioutil.ReadFile("assets/index.html")
	tpl, err := template.New("template pdf").Parse(string(fileData))
	guard(err)

	var b []byte
	buf := bytes.NewBuffer(b)
	vars := tplVars{"yooo"}
	tpl.Execute(buf, vars)
	r := bufio.NewReader(buf)

	page := wkhtmltopdf.NewPageReader(r)
	page.NoImages.Set(true)
	page.DisableJavascript.Set(true)
	page.DisableExternalLinks.Set(true)
	page.DisableInternalLinks.Set(true)
	pdf.AddPage(page)

	//fmt.Printf("p: %T\n", r)
	fmt.Printf("args: %v\n", pdf.Args())

	err = pdf.Create()
	guard(err)
	//ioutil.WriteFile("out/result.pdf", result, 0777)
	//os.Remove("tmp.html")
	//fmt.Println("Ok.")
	return pdf.Bytes()
}

func toPDF(fileData []byte) []byte {

	dir, err := os.Getwd()
	guard(err)

	fmt.Println(runtime.GOOS) // darwin
	var path string
	if runtime.GOOS == "darwin" {
		path = dir + "/bin/wkhtmltopdf_macOs"
	} else {
		path = dir + "/bin/wkhtmltopdf_linux_amd64"
	}
	// set binary path
	wkhtmltopdf.SetPath(path)
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
	err = pdf.Create()
	guard(err)
	return pdf.Bytes()
}

