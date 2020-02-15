// Copyright (c) 2015 Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnmarshalRequestBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", strings.NewReader("{}"))

	var m StatusResponse

	err := UnmarshalRequestBody(w, r, m)

	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, statusCode)
	}

	if m.Code != 0 {
		t.Errorf("expected message code %d, got %d", 0, m.Code)
	}

	if m.Message != "" {
		t.Errorf("expected message message \"\", got %q", m.Message)
	}
}

func TestUnmarshalRequestBody_message(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", strings.NewReader(`{"message":"test message"}`))

	var m StatusResponse

	if err := UnmarshalRequestBody(w, r, &m); err != nil {
		t.Errorf("unexpected error: %#v", err)
	}

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, statusCode)
	}

	if m.Code != 0 {
		t.Errorf("expected message code %d, got %d", 0, m.Code)
	}

	if m.Message != "test message" {
		t.Errorf("expected message message \"test message\", got %q", m.Message)
	}
}

func TestUnmarshalRequestBody_emptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", nil)

	err := UnmarshalRequestBody(w, r, nil)

	if err != ErrEmptyRequestBody {
		t.Errorf("expected error %#v, got %#v", ErrEmptyRequestBody, err)
	}

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, statusCode)
	}

	var m StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Errorf("json unmarshal response body: %#v", err)
	}

	if m.Code != http.StatusBadRequest {
		t.Errorf("expected message code %d, got %d", http.StatusBadRequest, m.Code)
	}

	if m.Message != ErrEmptyRequestBody.Error() {
		t.Errorf("expected message message %q, got %q", ErrEmptyRequestBody.Error(), m.Message)
	}
}

func TestUnmarshalRequestBody_contentLength0(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", strings.NewReader("{}"))
	r.Header.Add("Content-Length", "0")

	err := UnmarshalRequestBody(w, r, nil)

	if err != ErrEmptyRequestBody {
		t.Errorf("expected error %#v, got %#v", ErrEmptyRequestBody, err)
	}

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, statusCode)
	}

	var m StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Errorf("json unmarshal response body: %#v", err)
	}

	if m.Code != http.StatusBadRequest {
		t.Errorf("expected message code %d, got %d", http.StatusBadRequest, m.Code)
	}

	if m.Message != ErrEmptyRequestBody.Error() {
		t.Errorf("expected message message %q, got %q", ErrEmptyRequestBody.Error(), m.Message)
	}
}

func TestUnmarshalRequestBody_jsonSyntaxError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", strings.NewReader("{1}"))

	err := UnmarshalRequestBody(w, r, nil)

	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected error json.SyntaxError, got %#v", err)
	}

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, statusCode)
	}

	var m StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Errorf("json unmarshal response body: %#v", err)
	}

	if m.Code != http.StatusBadRequest {
		t.Errorf("expected message code %d, got %d", http.StatusBadRequest, m.Code)
	}

	message := "invalid character '1' looking for beginning of object key string (offset 2)"
	if m.Message != message {
		t.Errorf("expected message message %q, got %q", message, m.Message)
	}
}

func TestUnmarshalRequestBody_jsonUnmarshalTypeError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", strings.NewReader(`{"code":"invalid code"}`))

	err := UnmarshalRequestBody(w, r, &StatusResponse{})

	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Errorf("expected error json.UnmarshalTypeError, got %#v", err)
	}

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, statusCode)
	}

	var m StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Errorf("json unmarshal response body: %#v", err)
	}

	if m.Code != http.StatusBadRequest {
		t.Errorf("expected message code %d, got %d", http.StatusBadRequest, m.Code)
	}

	message := "expected json int value but got string (offset 22)"
	if m.Message != message {
		t.Errorf("expected message message %q, got %q", message, m.Message)
	}
}
