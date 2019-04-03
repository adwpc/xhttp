package xhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("myrespdata"))

		if r.Method != http.MethodGet {
			t.Errorf("Except method [GET] got [%s]", r.Method)
		}

		if r.Header.Get("header") != "myheader" {
			t.Errorf("Except header [header: myheader] got [%s]", r.Header.Get("header"))
		}

		r.ParseForm()
		if r.Form.Get("param") != "myparam" {
			t.Errorf("Except param [param=myparam] got [param=%s]", r.Form.Get("p"))
		}

		body, _ := ioutil.ReadAll(r.Body)
		if string(body) != "mybody" {
			t.Errorf("Except body [mybody] got [%s]", body)
		}

		fmt.Printf("[request]: method:%s header:%s param:%s body:%s\n", r.Method, r.Header.Get("header"), r.Form.Get("param"), body)
	}))

	defer ts.Close()

	str, err := New().Get().AddHeader("header", "myheader").AddParam("param", "myparam").SetBody("mybody").RespToString(ts.URL)

	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("[response]: string:%s\n", str)
	}
}

func TestGetJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"key\":\"value\"}"))

		if r.Method != http.MethodGet {
			t.Errorf("Except method [GET] got [%s]", r.Method)
		}

		if r.Header.Get("header") != "myheader" {
			t.Errorf("Except header [header: myheader] got [%s]", r.Header.Get("header"))
		}

		r.ParseForm()
		if r.Form.Get("param") != "myparam" {
			t.Errorf("Except param [param=myparam] got [param=%s]", r.Form.Get("p"))
		}

		body, _ := ioutil.ReadAll(r.Body)
		if string(body) != "mybody" {
			t.Errorf("Except body [mybody] got [%s]", body)
		}

		fmt.Printf("[request]: method:%s header:%s param:%s body:%s\n", r.Method, r.Header.Get("header"), r.Form.Get("param"), body)
	}))

	defer ts.Close()

	type Json struct {
		Key string `json:"key"`
	}
	var j Json
	err := New().Get().AddHeader("header", "myheader").AddParam("param", "myparam").SetBody("mybody").RespToJson(ts.URL, &j)

	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("[response]: %v\n", j)
	}
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("myrespdata"))

		if r.Method != http.MethodPost {
			t.Errorf("Except method [POST] got [%s]", r.Method)
		}

		if r.Header.Get("header") != "myheader" {
			t.Errorf("Except header [header: myheader] got [%s]", r.Header.Get("header"))
		}

		r.ParseForm()
		if r.Form.Get("param") != "myparam" {
			t.Errorf("Except param [param=myparam] got [param=%s]", r.Form.Get("p"))
		}

		body, _ := ioutil.ReadAll(r.Body)
		if string(body) != "mybody" {
			t.Errorf("Except body [mybody] got [%s]", body)
		}

		fmt.Printf("[request]: method:%s header:%s param:%s body:%s\n", r.Method, r.Header.Get("header"), r.Form.Get("param"), body)
	}))

	defer ts.Close()

	str, err := New().Post().AddHeader("header", "myheader").AddParam("param", "myparam").SetBody("mybody").RespToString(ts.URL)

	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("[response]: str:%s\n", str)
	}
}
