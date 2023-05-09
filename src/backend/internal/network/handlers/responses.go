package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type ErrorModel struct {
	Error          string
	HTTPStatusCode int
}

func sendErrorResponse(w http.ResponseWriter, e *ErrorModel) {
	w.WriteHeader(e.HTTPStatusCode)
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(e)
	w.Write(reqBodyBytes.Bytes())
}

func sendResponse(w http.ResponseWriter, s int, event int) {
	w.WriteHeader(s)
	data, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func sendMapResponse(w http.ResponseWriter, s int, event map[string]any) {
	w.WriteHeader(s)
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(event)
	data, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}
