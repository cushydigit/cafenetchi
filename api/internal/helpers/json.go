package helpers

import (
	"cafenetchi-api/internal/types"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ReadJSON reads JSON data from an HTTP request body and decodes it into the provided data structure.
//
// Parameters:
// - w: The HTTP response writer to write any error responses to.
// - r: The HTTP request containing the JSON data to be read.
// - data: The data structure to decode the JSON data into.
//
// Returns:
// - error: An error if there was an issue reading or decoding the JSON data.
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON writes a JSON response to the HTTP response writer.
//
// Parameters:
// - w: The HTTP response writer to write the JSON response to.
// - status: The HTTP status code to set in the response.
// - data: The data to be marshaled into JSON.
// - headers: Optional custom headers to be included in the response.
//
// Returns:
// - error: An error if there was an issue writing the JSON response.
func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
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

func OK(w http.ResponseWriter, data any) error {
	return writeJSON(w, http.StatusOK, types.Response{
		Error: false,
		Data:  data,
	})
}

func Created(w http.ResponseWriter, data any) error {
	return writeJSON(w, http.StatusCreated, types.Response{
		Error: false,
		Data:  data,
	})
}

func Message(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, types.Response{
		Error:   false,
		Message: message,
	})
}

func Error(w http.ResponseWriter, err error) error {
	apiErr, ok := err.(types.APIError)
	if !ok {
		apiErr = types.ErrInternalServer
	}
	return writeJSON(w, apiErr.Status, types.Response{
		Error:   true,
		Message: apiErr.Message,
	})
}
