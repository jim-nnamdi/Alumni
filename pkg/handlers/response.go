package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"
)

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Err     string      `json:"err,omitempty"`
	Success bool        `json:"success"`
	TTL     int         `json:"ttl"`
}

func GetSuccessResponse(data interface{}, ttl int) []byte {
	resp := &response{
		Success: true,
		TTL:     ttl,
	}

	if data == nil || (reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) {
		resp.Data = nil
	} else {
		resp.Data = data
	}
	responseBytes, _ := json.Marshal(resp)
	return responseBytes
}

func GetErrorResponseBytes(data interface{}, ttl int, err error) []byte {
	resp := &response{
		Success: false,
		Err:     "",
		TTL:     ttl,
	}
	if data == nil || (reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) {
		resp.Data = nil
	} else {
		resp.Data = data
	}
	if err != nil {
		resp.Err = err.Error()
	}

	responseBytes, _ := json.Marshal(resp)

	return responseBytes
}
func createSessionKey(patronID string, t time.Time) string {
	var (
		hash = md5.New()
		err  error
	)
	if _, err = io.WriteString(hash, patronID); err != nil {
		return err.Error()
	}
	if _, err = io.WriteString(hash, t.Format(time.RFC3339)); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
