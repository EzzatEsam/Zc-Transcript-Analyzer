package main

import (
	"bytes"
	"github.com/ledongthuc/pdf"
)

type Reader struct {
	fileName string
}
func (rd Reader) readLines() (string, error) {
	f, r, err := pdf.Open(rd.fileName)
	// remember close file
	if err != nil {
		println(err.Error())
		return "", err
	}
    defer f.Close()
	
	var buf bytes.Buffer
    b, err := r.GetPlainText()
    if err != nil {
        return "", err
    }
    buf.ReadFrom(b)
	return buf.String(), nil
}