// remote-typograf is a service which prepares a text for the web by quering
// the webservice made and maintaning by ArtLebedev Studio.
//
// See https://www.artlebedev.ru/tools/typograf/about/ for details.
//
// The service should be strict about timeouts from the original webservice
// because it's undocumented and untested.
//
// Usage for a client:
//   * send a text with a specified encoding to the service endpoint
//   * wait for a fixed time interval
//   * receive the text (processed or the original one)
//
// The service uses the python-client available at https://www.artlebedev.ru/tools/typograf/webservice/.
//
// Author: Ihar Suvorau.
// Time: Fri Feb 24 20:48:14 MSK 2017.
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	url  = "http://typograf.artlebedev.ru/webservices/typograf.asmx"
	resp = new(response)
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("listening at :8080...")
	http.ListenAndServe(":8080", nil)
}

// a request must be in JSON and conaint only a text
// response must be in JSON
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	}

	if r.Method == "POST" {
		var err error
		var data = make(map[string]string)

		if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
			respond(err, "", 500, w)
			return
		}
		defer r.Body.Close()

		if data["encoding"] == "" {
			data["encoding"] = "UTF-8"
		}

		var buf *bytes.Buffer
		buf, err = prepareRequest(data)
		if err != nil {
			respond(err, "", 500, w)
			return
		}

		client := new(http.Client)
		client.Timeout = time.Millisecond * 500

		var req *http.Request
		req, err = http.NewRequest("POST", url, buf)
		req.Header.Add("Content-Type", "text/xml")
		req.Header.Add("Content-Length", string(buf.Len()))
		req.Header.Add("SOAPAction", "http://typograf.artlebedev.ru/webservices/ProcessText")

		var start = time.Now()
		var tresp *http.Response
		if tresp, err = client.Do(req); err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				respond(err, "typograf service took too long to respond", 503, w)
				return
			} else {
				respond(err, "", 500, w)
				return
			}
		}
		log.Printf("origin request made in %v\n", time.Now().Sub(start))
		defer tresp.Body.Close()

		var result = new(Envelope)
		if err = xml.NewDecoder(tresp.Body).Decode(&result); err != nil {
			respond(err, "", 500, w)
			return
		}

		data["data"] = result.Body.ProcessTextResponse.ProcessTextResult

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if err = json.NewEncoder(w).Encode(data); err != nil {
			respond(err, "", 500, w)
			return
		}
	}

}

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xsd     string   `xml:"xmlns:xsd,attr"`
	Soap    string   `xml:"xmlns:soap,attr"`
	Body    *Body
}

type Body struct {
	XMLName             xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	ProcessText         *ProcessText
	ProcessTextResponse struct {
		ProcessTextResult string
	}
}

type ProcessText struct {
	XMLName    xml.Name `xml:"http://typograf.artlebedev.ru/webservices/ ProcessText"`
	Text       string   `xml:"text"`
	EntityType int      `xml:"entityType"`
	UseBr      int      `xml:"useBr"`
	UseP       int      `xml:"useP"`
	MaxNobr    int      `xml:"maxNobr"`
}

func prepareRequest(data map[string]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	var head = fmt.Sprintf("<?xml version=\"1.0\" encoding=\"%s\"?>", data["encoding"])
	buf.Write([]byte(head))

	data["data"] = strings.Replace(data["data"], "&", "&amp;", -1)
	data["data"] = strings.Replace(data["data"], "<", "&lt;", -1)
	data["data"] = strings.Replace(data["data"], ">", "&gt;", -1)

	var s = &Envelope{
		Xsi:  "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:  "http://www.w3.org/2001/XMLSchema",
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: &Body{
			ProcessText: &ProcessText{
				Text:       data["data"],
				EntityType: 4,
				UseBr:      1,
				UseP:       1,
				MaxNobr:    3,
			},
		},
	}
	b, err := xml.Marshal(s)
	buf.Write(b)
	return &buf, err
}

// response type is used during error responding.
type response struct {
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	Data       string `json:"data,omitempty"`
}

// respond responds in JSON.
// Must be sure to call return after writing to the ResponseWriter.
func respond(err error, msg string, code int, w http.ResponseWriter) {
	log.Println("[error]", err)
	resp.Error = err.Error()
	resp.Message = msg
	resp.StatusCode = code
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println(err)
	}
}
