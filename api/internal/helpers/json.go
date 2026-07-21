package helpers

import (
	"cafenetchi-api/internal/types"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ReadJSON decodes a JSON request body into data.
//
// The request body is limited to 1 MiB. it reject empty bodies and multiple JSON values.
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var httpMaxByteError *http.MaxBytesError
		switch {
		case errors.Is(err, io.EOF):
			return types.ErrInvalidRequest
		case errors.Is(err, io.ErrUnexpectedEOF):
			return types.ErrInvalidRequest
		case errors.As(err, &httpMaxByteError):
			return types.ErrInvalidRequest
		case errors.As(err, &syntaxError):
			return types.ErrInvalidRequest
		case errors.As(err, &unmarshalTypeError):
			return types.ErrInvalidRequest
		default:
			return types.ErrInvalidRequest
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return types.ErrInvalidRequest
	}

	return nil
}

// WriteJSON marshals data as JSON and writes it to the response i
// with provided HTTP status code.
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		// add or overwrite with custom headers
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(out); err != nil {
		return err
	}
	return nil
}

// Ok writes a 200 OK JSON response.
func OK(w http.ResponseWriter, data any) {
	_ = WriteJSON(w, http.StatusOK, types.Response{
		Data: data,
		Code: http.StatusOK,
	})
}

// Created writes a 201 Created JSON response
func Created(w http.ResponseWriter, data any) {
	_ = WriteJSON(w, http.StatusCreated, types.Response{
		Data: data,
		Code: http.StatusCreated,
	})
}

// Message writes a JSON response containing only a message.
func Message(w http.ResponseWriter, status int, message string) {
	_ = WriteJSON(w, status, types.Response{
		Message: message,
		Code:    status,
	})
}

// Error writes a API error response.
//
// if err is not an APIError, an Internal Server Error is returned.
func Error(w http.ResponseWriter, err error) {
	var apiErr types.APIError
	if !errors.As(err, &apiErr) {
		apiErr = types.ErrInternalServer
	}
	_ = WriteJSON(w, apiErr.Status, types.Response{
		Error:   true,
		Message: apiErr.Message,
		Code:    apiErr.Status,
	})
}
