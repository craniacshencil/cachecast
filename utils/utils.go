package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONWrapper struct {
	Data interface{}
}

func (m *JSONWrapper) MarshalBinary() (data []byte, err error) {
	marshalledData, err := json.Marshal(m.Data)
	if err != nil {
		return nil, err
	}
	return marshalledData, nil
}

func (m *JSONWrapper) UnmarshalBinary(data []byte) (err error) {
	err = json.Unmarshal(data, &m.Data)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ParseBody(message interface{}, v any) error {
	var body io.ReadCloser

	switch m := message.(type) {

	case *http.Request:
		if m.Body == nil {
			return errors.New("empty request body")
		}
		body = m.Body

	case *http.Response:
		if m.Body == nil {
			return errors.New("empty response body")
		}
		body = m.Body

	default:
		return errors.New("neither response nor request")
	}

	defer body.Close()
	if err := json.NewDecoder(body).Decode(v); err != nil {
		return errors.New("couldn't parse body")
	}
	return nil
}
