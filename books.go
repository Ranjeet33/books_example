package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	if r.Method == "GET" {
		Books, err := json.Marshal(books)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(Books)
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid ID"))
		} else {
			sbook, err := json.Marshal(Book)
			if err != nil {
				log.Println(err)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(sbook)
		}
	}
}

func createNewBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		var newBook book
		err = json.Unmarshal(reqBody, &newBook)
		if err != nil {
			log.Println(err)
			return
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
				log.Println(err)
				return
			}
			w.Write(book)
		} else {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("ALREADY_EXISTS"))
		}
	}
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "PUT" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		var newBook book
		index := 0
		err = json.Unmarshal(reqBody, &newBook)
		if err != nil {
			log.Println(err)
			return
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
				log.Println(err)
				return
			}
			w.Write(book)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("ID not found"))
		}
	} else if r.Method == "PATCH" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		var newBook book
		index := 0
		err = json.Unmarshal(reqBody, &newBook)
		if err != nil {
			log.Println(err)
			return
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
				log.Println(err)
				return
			}
			w.Write(book)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("ID not found"))
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

func collectionBook(w http.ResponseWriter, r *http.Request) {
	getAllBooks(w, r)
	createNewBook(w, r)
}
func manageBook(w http.ResponseWriter, r *http.Request) {

	getBook(w, r)
	deleteBook(w, r)
	updateBook(w, r)
}

func main() {
	http.HandleFunc("/books", collectionBook)
	http.HandleFunc("/books/", manageBook)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Println(err)
	}
}
