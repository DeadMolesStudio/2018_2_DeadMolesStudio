package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/models"
)

type TestCase struct {
	// request
	RequestMethod  string
	RequestHeaders []string
	RequestBody    string
	// response
	Response        string
	ResponseHeaders []string
	ResponseBody    string
	StatusCode      int
}

func TestLoginFail(t *testing.T) {
	cases := []TestCase{
		TestCase{ // [0] empty body
			RequestMethod: "POST",
			StatusCode:    400,
		},
		TestCase{ // [1] wrong json format
			RequestMethod: "POST",
			RequestBody: `{
				"email": "username"
				"password": "password"
			}`,
			StatusCode: 400,
		},
		TestCase{ // [2] wrong json format
			RequestMethod: "POST",
			RequestBody: `{
				"abacaba": "bacbac",
				"isTrue": 1
			}`,
			StatusCode: 400,
		},
		TestCase{ // [3] pair <user & password> is incorrect
			RequestMethod: "POST",
			RequestBody: `{
				"email": "sdfadf",
				"password": "asfdaf"
			}`,
			StatusCode: 403,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(sessionHandler))
	defer ts.Close()

	tsURL := ts.URL
	for i, c := range cases {
		req, _ := http.NewRequest(c.RequestMethod, tsURL, strings.NewReader(c.RequestBody))
		w := httptest.NewRecorder()

		sessionHandler(w, req)

		if w.Code != c.StatusCode {
			t.Errorf("[%d] wrong StatusCode:\ngot: %d\nexpected: %d", i, w.Code, c.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != c.Response {
			t.Errorf("[%d] wrong Response:\ngot: %+v\nexpected: %+v", i, bodyStr, c.Response)
		}
	}
}

func TestLoginLogout(t *testing.T) {
	users = append(users, models.UserProfile{
		UserPassword: models.UserPassword{
			UserID:   1,
			Email:    "test",
			Password: "test",
		},
		Nickname: "tester",
		Stats: models.Stats{
			Record: 111,
		},
	})
	c := []TestCase{
		TestCase{ // login
			RequestMethod: "POST",
			RequestBody: `{
			"email": "test",
			"password": "test"
		}`,
			StatusCode: 200,
		},
		TestCase{ // login again
			RequestMethod: "POST",
			RequestBody: `{
			"email": "test",
			"password": "test"
		}`,
			StatusCode: 200,
		},
		TestCase{ // logout
			RequestMethod: "DELETE",
			StatusCode:    200,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(sessionHandler))
	defer ts.Close()

	tsURL := ts.URL

	// login
	req, _ := http.NewRequest(c[0].RequestMethod, tsURL, strings.NewReader(c[0].RequestBody))
	w := httptest.NewRecorder()

	sessionHandler(w, req)

	if w.Code != c[0].StatusCode {
		t.Errorf("login-logout case: wrong StatusCode:\ngot: %d\nexpected: %d", w.Code, c[0].StatusCode)
	}

	resp := w.Result()
	heads := resp.Header
	sessionCookie, ok := heads["Set-Cookie"]
	if !ok {
		t.Errorf("login-logout case: wrong Response: no session cookie got")
	}

	if w.Code != c[0].StatusCode {
		t.Errorf("login-logout case: wrong Response from first login:\ngot: %d\nexpected: %d", w.Code, c[1].StatusCode)
	}

	// second login
	req, _ = http.NewRequest(c[0].RequestMethod, tsURL, strings.NewReader(c[0].RequestBody))
	req.Header.Set("Cookie", sessionCookie[0])
	w = httptest.NewRecorder()

	sessionHandler(w, req)

	if w.Code != c[1].StatusCode {
		t.Errorf("login-logout case: wrong Response from repeated login:\ngot: %d\nexpected: %d", w.Code, c[1].StatusCode)
	}

	// logout
	tslogout := httptest.NewServer(http.HandlerFunc(sessionHandler))

	req, _ = http.NewRequest(c[0].RequestMethod, tslogout.URL, strings.NewReader(c[0].RequestBody))
	req.Header.Set("Cookie", sessionCookie[0])
	w = httptest.NewRecorder()

	sessionHandler(w, req)

	if w.Code != c[2].StatusCode {
		t.Errorf("login-logout case: wrong Response from repeated login:\ngot: %d\nexpected: %d", w.Code, c[1].StatusCode)
	}

	resp = w.Result()
	heads = resp.Header
	sessionCookie, ok = heads["Set-Cookie"]
	if !ok {
		t.Errorf("login-logout case: wrong Response: no session cookie got")
	}
}

func TestLogoutWhenAlreadyLoggedOut(t *testing.T) {
	cases := []TestCase{
		TestCase{ // [0] logged out
			RequestMethod: "GET",
			StatusCode:    200,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(sessionHandler))
	defer ts.Close()

	tsURL := ts.URL
	for i, c := range cases {
		req, _ := http.NewRequest(c.RequestMethod, tsURL, strings.NewReader(c.RequestBody))
		w := httptest.NewRecorder()

		sessionHandler(w, req)

		if w.Code != c.StatusCode {
			t.Errorf("[%d] wrong StatusCode:\ngot: %d\nexpected: %d", i, w.Code, c.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != c.Response {
			t.Errorf("[%d] wrong Response:\ngot: %+v\nexpected: %+v", i, bodyStr, c.Response)
		}
	}
}
