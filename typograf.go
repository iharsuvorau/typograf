// Package typograf is a library for the ArtLebedev Studio's webwervice
// which prepares texts for the Web by replacing characters for the
// typographicaly correct ones.
// See https://www.artlebedev.ru/tools/typograf/about/ for details.
package typograf

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"time"
)

// ServiceURL is the address of the SOAP web service.
var ServiceURL = "http://typograf.artlebedev.ru/webservices/typograf.asmx"

// Envelope is the SOAP Envelope. Contains the whole request to the SOAP WS.
type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xsd     string   `xml:"xmlns:xsd,attr"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    *Body
}

// Body contains an input text and output text.
type Body struct {
	XMLName             xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	ProcessText         *ProcessText
	ProcessTextResponse struct {
		ProcessTextResult string
	}
}

// ProcessText is a type for the input text.
type ProcessText struct {
	XMLName    xml.Name `xml:"http://typograf.artlebedev.ru/webservices/ ProcessText"`
	Text       string   `xml:"text"`
	EntityType int      `xml:"entityType"`
	UseBr      int      `xml:"useBr"`
	UseP       int      `xml:"useP"`
	MaxNobr    int      `xml:"maxNobr"`
}

func replaceSymbols(in string) string {
	in = strings.Replace(in, "&", "&amp;", -1)
	in = strings.Replace(in, "<", "&lt;", -1)
	return strings.Replace(in, ">", "&gt;", -1)
}

func newEnvelope(s string) *Envelope {
	return &Envelope{
		Xsi:  "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:  "http://www.w3.org/2001/XMLSchema",
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: &Body{
			ProcessText: &ProcessText{
				Text:       s,
				EntityType: 4,
				UseBr:      1,
				UseP:       1,
				MaxNobr:    3,
			},
		},
	}
}

// PrepareRequestBody preprocesses an input string and returns a reader with the data for a SOAP WS
// for the following request.
func PrepareRequestBody(in string) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	s := replaceSymbols(in)
	b, err := xml.Marshal(newEnvelope(s))
	buf.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	buf.Write(b)
	return
}

// DoRequest sends a SOAP request to the webservice and returns the output string.
func DoRequest(body *bytes.Buffer) (out string, err error) {
	// making a request
	client := new(http.Client)
	client.Timeout = time.Millisecond * 500

	req, err := http.NewRequest("POST", ServiceURL, body)
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("Content-Length", string(body.Len()))
	req.Header.Add("SOAPAction", "http://typograf.artlebedev.ru/webservices/ProcessText")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	log.Printf("origin request made in %v\n", time.Since(start))
	defer resp.Body.Close()

	// decoding a response
	var result = new(Envelope)
	if err = xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	// returning the processed text only
	return result.Body.ProcessTextResponse.ProcessTextResult, nil
}

// Typogrify gets a piece of text and typography it to the one with correct typographical characters.
// It preprocesses the given text and makes a request returning a processed text and an error.
func Typogrify(in string) (out string, err error) {
	reqBody, err := PrepareRequestBody(in)
	if err != nil {
		return
	}

	if out, err = DoRequest(reqBody); err != nil {
		return
	}

	return
}
