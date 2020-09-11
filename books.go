package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type book struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

var arr = []book{}

func get(e *book) string {
	return (e.ID + e.Name + e.Description + e.Date)
}

func post(reqbody []byte) string {
	var t book
	err := json.Unmarshal(reqbody, &t)
	if err != nil {
		panic(err)
	}
	arr = append(arr, t)
	return t.ID
}

func put(reqbody []byte) string {
	var t book
	err := json.Unmarshal(reqbody, &t)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(arr); i++ {
		if t.ID == arr[i].ID {
			arr[i] = t
		}
	}
	return "data replaced"
}
func patch(reqbody []byte) string {
	var t book
	err := json.Unmarshal(reqbody, &t)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(arr); i++ {
		if t.ID == arr[i].ID {
			if t.Name != arr[i].Name {
				arr[i].Name = t.Name
			} else if t.Description != arr[i].Description {
				arr[i].Description = t.Description
			} else if t.Date != arr[i].Date {
				arr[i].Date = t.Date
			}
		}
	}
	return "patch updated"
}

func addCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		key := "id"
		val := r.URL.Query().Get(key)
		for i := 0; i < len(arr); i++ {
			if val == arr[i].ID {
				io.WriteString(w, get(&arr[i]))
			}
		}
	} else if r.Method == "POST" {
		reqbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, post(reqbody))
	} else if r.Method == "PUT" {
		reqbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, put(reqbody))
	} else if r.Method == "PATCH" {
		reqbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, patch(reqbody))
	} else {
		io.WriteString(w, "404 page not found")
	}
	addCookie(w, "TestCookieName", "TestValue", 30*time.Minute)

}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
