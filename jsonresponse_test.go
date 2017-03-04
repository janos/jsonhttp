// Copyright (c) 2015, 2016 Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonresponse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultRespond(t *testing.T) {
	w := httptest.NewRecorder()

	Respond(w, 0, nil)

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, statusCode)
	}

	m := &MessageResponse{}

	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if m.Code != 0 {
		t.Errorf("expected message code %d, got %d", 0, m.Code)
	}

	if m.Message != "" {
		t.Errorf("expected message message \"\", got \"%s\"", m.Message)
	}
}

func TestStatusCodeRespond(t *testing.T) {
	w := httptest.NewRecorder()

	Respond(w, http.StatusForbidden, nil)

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusForbidden {
		t.Errorf("expected status code %d, got %d", http.StatusForbidden, statusCode)
	}

	m := &MessageResponse{}

	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if m.Code != http.StatusForbidden {
		t.Errorf("expected message code %d, got %d", http.StatusForbidden, m.Code)
	}

	message := http.StatusText(http.StatusForbidden)
	if m.Message != message {
		t.Errorf("expected message message \"%s\", got \"%s\"", message, m.Message)
	}
}

func TestMessageRespond(t *testing.T) {
	w := httptest.NewRecorder()

	message := "test message"
	Respond(w, 0, MessageResponse{
		Message: message,
	})

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, statusCode)
	}

	m := &MessageResponse{}

	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if m.Code != 0 {
		t.Errorf("expected message code %d, got %d", 0, m.Code)
	}

	if m.Message != message {
		t.Errorf("expected message message \"%s\", got \"%s\"", message, m.Message)
	}
}

func TestMessageAndStatusCodeRespond(t *testing.T) {
	w := httptest.NewRecorder()

	message := "test custom forbidden message"
	Respond(w, http.StatusForbidden, MessageResponse{
		Message: message,
	})

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusForbidden {
		t.Errorf("expected status code %d, got %d", http.StatusForbidden, statusCode)
	}

	m := &MessageResponse{}

	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if m.Code != http.StatusForbidden {
		t.Errorf("expected message code %d, got %d", http.StatusForbidden, m.Code)
	}

	if m.Message != message {
		t.Errorf("expected message message \"%s\", got \"%s\"", message, m.Message)
	}
}

func TestCustomCodeRespond(t *testing.T) {
	w := httptest.NewRecorder()

	code := 1001
	Respond(w, 0, MessageResponse{
		Code: code,
	})

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, statusCode)
	}

	m := &MessageResponse{}

	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if m.Code != code {
		t.Errorf("expected message code %d, got %d", code, m.Code)
	}

	if m.Message != "" {
		t.Errorf("expected message message \"\", got \"%s\"", m.Message)
	}
}

func TestMessageAndCustomCodeRespond(t *testing.T) {
	w := httptest.NewRecorder()

	message := "test custom message"
	code := 1001
	Respond(w, http.StatusNotFound, MessageResponse{
		Message: message,
		Code:    code,
	})

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, statusCode)
	}

	m := &MessageResponse{}

	if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
		t.Errorf("json unmarshal response body: %s", err)
	}

	if m.Code != code {
		t.Errorf("expected message code %d, got %d", code, m.Code)
	}

	if m.Message != message {
		t.Errorf("expected message message \"%s\", got \"%s\"", message, m.Message)
	}
}

func TestPanicRespond(t *testing.T) {
	w := httptest.NewRecorder()

	defer func() {
		err := recover()
		if _, ok := err.(*json.UnsupportedTypeError); !ok {
			t.Errorf("exppected error from recover json.UnsupportedTypeError, got %#v", err)
		}
	}()

	Respond(w, http.StatusNotFound, map[bool]string{
		true: "",
	})
}

func TestStandardHTTPResponds(t *testing.T) {
	for _, test := range []struct {
		f    func(w http.ResponseWriter, response interface{})
		code int
	}{
		{f: Continue, code: http.StatusContinue},
		{f: SwitchingProtocols, code: http.StatusSwitchingProtocols},
		{f: OK, code: http.StatusOK},
		{f: Created, code: http.StatusCreated},
		{f: Accepted, code: http.StatusAccepted},
		{f: NonAuthoritativeInfo, code: http.StatusNonAuthoritativeInfo},
		{f: ResetContent, code: http.StatusResetContent},
		{f: PartialContent, code: http.StatusPartialContent},
		{f: MultipleChoices, code: http.StatusMultipleChoices},
		{f: MovedPermanently, code: http.StatusMovedPermanently},
		{f: Found, code: http.StatusFound},
		{f: SeeOther, code: http.StatusSeeOther},
		{f: NotModified, code: http.StatusNotModified},
		{f: UseProxy, code: http.StatusUseProxy},
		{f: TemporaryRedirect, code: http.StatusTemporaryRedirect},
		{f: PermanentRedirect, code: http.StatusPermanentRedirect},
		{f: BadRequest, code: http.StatusBadRequest},
		{f: Unauthorized, code: http.StatusUnauthorized},
		{f: PaymentRequired, code: http.StatusPaymentRequired},
		{f: Forbidden, code: http.StatusForbidden},
		{f: NotFound, code: http.StatusNotFound},
		{f: MethodNotAllowed, code: http.StatusMethodNotAllowed},
		{f: NotAcceptable, code: http.StatusNotAcceptable},
		{f: ProxyAuthRequired, code: http.StatusProxyAuthRequired},
		{f: RequestTimeout, code: http.StatusRequestTimeout},
		{f: Conflict, code: http.StatusConflict},
		{f: Gone, code: http.StatusGone},
		{f: LengthRequired, code: http.StatusLengthRequired},
		{f: PreconditionFailed, code: http.StatusPreconditionFailed},
		{f: RequestEntityTooLarge, code: http.StatusRequestEntityTooLarge},
		{f: RequestURITooLong, code: http.StatusRequestURITooLong},
		{f: UnsupportedMediaType, code: http.StatusUnsupportedMediaType},
		{f: RequestedRangeNotSatisfiable, code: http.StatusRequestedRangeNotSatisfiable},
		{f: ExpectationFailed, code: http.StatusExpectationFailed},
		{f: Teapot, code: http.StatusTeapot},
		{f: UpgradeRequired, code: http.StatusUpgradeRequired},
		{f: PreconditionRequired, code: http.StatusPreconditionRequired},
		{f: TooManyRequests, code: http.StatusTooManyRequests},
		{f: RequestHeaderFieldsTooLarge, code: http.StatusRequestHeaderFieldsTooLarge},
		{f: UnavailableForLegalReasons, code: http.StatusUnavailableForLegalReasons},
		{f: InternalServerError, code: http.StatusInternalServerError},
		{f: NotImplemented, code: http.StatusNotImplemented},
		{f: BadGateway, code: http.StatusBadGateway},
		{f: ServiceUnavailable, code: http.StatusServiceUnavailable},
		{f: GatewayTimeout, code: http.StatusGatewayTimeout},
		{f: HTTPVersionNotSupported, code: http.StatusHTTPVersionNotSupported},
	} {
		w := httptest.NewRecorder()
		test.f(w, nil)
		m := &MessageResponse{}

		if err := json.Unmarshal(w.Body.Bytes(), m); err != nil {
			t.Errorf("json unmarshal response body: %s", err)
		}

		if m.Code != test.code {
			t.Errorf("expected message code %d, got %d", test.code, m.Code)
		}

		if m.Message != http.StatusText(test.code) {
			t.Errorf("expected message message \"%s\", got \"%s\"", http.StatusText(test.code), m.Message)
		}
	}
}

func TestNewMessage(t *testing.T) {
	message := "testing message"
	m := NewMessage(message)
	if m.Message != message {
		t.Errorf("expected message %q, got %q", message, m.Message)
	}
	if m.Code != 0 {
		t.Errorf("expected message %d, got %d", 0, m.Code)
	}
}
