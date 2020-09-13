package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type book struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

var books []book

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	Books, err := json.Marshal(books)
	if err != nil {
		panic(err)
	}
	w.Write(Books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var Book book
		key := "id"
		val := r.URL.Query().Get(key)
		index := 0
		for i := 0; i < len(books); i++ {
			if val == books[i].ID {
				Book = books[i]
				index = i
			}
		}
		if val != books[index].ID {
			w.Write([]byte("400 INVALID_ARGUMENT</br>Invalid ID"))
		} else {
			sbook, err := json.Marshal(Book)
			if err != nil {
				panic(err)
			}
			w.Write(sbook)
		}
	}
}

func createNewBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var newBook book
		err = json.Unmarshal(reqBody, &newBook)
		if err != nil {
			panic(err)
		}
		var con = true
		for i := 0; i < len(books); i++ {
			if newBook.ID == books[i].ID {
				con = false
			}
		}
		if con == true {
			books = append(books, newBook)
			book, err := json.Marshal(newBook)
			if err != nil {
				panic(err)
			}
			w.Write(book)
		} else {
			w.Write([]byte("409 ALREADY_EXISTS"))
		}
	}
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var newBook book
		index := 0
		err = json.Unmarshal(reqBody, &newBook)
		if err != nil {
			panic(err)
		}

		key := "id"
		val := r.URL.Query().Get(key)
		for i := 0; i < len(books); i++ {
			if val == books[i].ID {
				books[i] = newBook
				val = books[i].ID
				index = i
			}
		}
		if val == books[index].ID {
			book, err := json.Marshal(newBook)
			if err != nil {
				panic(err)
			}
			w.Write(book)
		} else {
			w.Write([]byte("404 NOT_FOUND</br>ID not found"))
		}
	} else if r.Method == "PATCH" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var newBook book
		index := 0
		err = json.Unmarshal(reqBody, &newBook)
		if err != nil {
			panic(err)
		}

		key := "id"
		val := r.URL.Query().Get(key)
		for i := 0; i < len(books); i++ {
			if val == books[i].ID {
				if newBook.ID != books[i].ID {
					if newBook.ID != "" {
						books[i].ID = newBook.ID
					}
				}
				if newBook.Name != books[i].Name {
					if newBook.Name != "" {
						books[i].Name = newBook.Name
					}
				}
				if newBook.Description != books[i].Description {
					if newBook.Description != "" {
						books[i].Description = newBook.Description
					}
				}
				if newBook.Date != books[i].Date {
					if newBook.Date != "" {
						books[i].Date = newBook.Date
					}
				}

				val = books[i].ID
				index = i
			}
		}
		if val == books[index].ID {
			book, err := json.Marshal(books[index])
			if err != nil {
				panic(err)
			}
			w.Write(book)
		} else {
			w.Write([]byte("404 NOT_FOUND</br>ID not found"))
		}
	}
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		key := "id"
		val := r.URL.Query().Get(key)

		for index, book := range books {
			if val == book.ID {
				books = append(books[:index], books[index+1:]...)
			}
		}
	}
}

func main() {
	http.HandleFunc("/books", getAllBooks)
	http.HandleFunc("/books/", getBook)
	http.HandleFunc("/book", createNewBook)
	http.HandleFunc("/books/d", deleteBook)
	http.HandleFunc("/book/", updateBook)
	http.ListenAndServe(":8000", nil)
}
