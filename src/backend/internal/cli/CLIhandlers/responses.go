package CLIhandlers

import (
	"bytes"
	"encoding/json"
	"log"
)

type ErrorModel struct {
	Error string
}

func ErrorResponse(e *ErrorModel) string {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(e)
	if err != nil {
		log.Fatal(err)
	}
	return string(reqBodyBytes.Bytes()[:])
}

func Response(event int) string {
	data, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return string(data[:])
}

func MapResponse(event map[string]any) string {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(event)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return string(data[:])
}
