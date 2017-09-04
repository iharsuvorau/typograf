// typografd is a service which prepares a text for the web by quering
// the webservice made and maintaning by ArtLebedev Studio.
// See https://www.artlebedev.ru/tools/typograf/about/ for details.
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/iharsuvorau/typograf"
)

var resp = new(response)

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
		// reading input text
		data := make(map[string]string)
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			respond(err, "", 500, w)
			return
		}
		defer r.Body.Close()

		buf, err := typograf.PrepareRequestBody(data["data"])
		if err != nil {
			respond(err, "", 500, w)
			return
		}

		if data["data"], err = typograf.DoRequest(buf); err != nil {
			respond(err, "", 500, w)
			return
		}

		// responding back
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
