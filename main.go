package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var articles = map[int]article{
	1: {
		"Simple REST API with Go",
		"Let's create a simple REST API with Go.",
	},
}

func index(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(articles)
}

func create(res http.ResponseWriter, req *http.Request) {
	newArticle := article{}
	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(reqBody, &newArticle)
	articles[len(articles)+1] = newArticle

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(articles)
}

func read(res http.ResponseWriter, req *http.Request) {
	articleParam := mux.Vars(req)["article"]
	articleID, err := strconv.Atoi(articleParam)

	if err != nil {
		log.Fatal(err)
	}

	if articles[articleID].Title == "" &&
		articles[articleID].Description == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(
		map[int]article{
			articleID: articles[articleID],
		},
	)
}

func replace(res http.ResponseWriter, req *http.Request) {
	articleParam := mux.Vars(req)["article"]
	articleID, err := strconv.Atoi(articleParam)

	if err != nil {
		log.Fatal(err)
	}

	if articles[articleID].Title == "" &&
		articles[articleID].Description == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	updatedArticle := article{}
	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(reqBody, &updatedArticle)
	articles[articleID] = updatedArticle

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(articles)
}

func modify(res http.ResponseWriter, req *http.Request) {
	articleParam := mux.Vars(req)["article"]
	articleID, err := strconv.Atoi(articleParam)

	if err != nil {
		log.Fatal(err)
	}

	if articles[articleID].Title == "" &&
		articles[articleID].Description == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	newArticle := article{}
	updatedArticle := article{}
	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(reqBody, &updatedArticle)

	if updatedArticle.Title != "" {
		newArticle.Title = updatedArticle.Title
	} else {
		newArticle.Title = articles[articleID].Title
	}

	if updatedArticle.Description != "" {
		newArticle.Description = updatedArticle.Description
	} else {
		newArticle.Description = articles[articleID].Description
	}

	articles[articleID] = newArticle

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(articles)
}

func remove(res http.ResponseWriter, req *http.Request) {
	articleParam := mux.Vars(req)["article"]
	articleID, err := strconv.Atoi(articleParam)

	if err != nil {
		log.Fatal(err)
	}

	if articles[articleID].Title == "" &&
		articles[articleID].Description == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the entry from the map.
	// In the real situation,
	// you need to delete it from the database.
	delete(articles, articleID)

	// Mostly in the delete operation,
	// return status ok is enough if succeed.
	res.WriteHeader(http.StatusOK)

	// Feel free to open this comment,
	// to check that the entry is deleted.
	// res.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(res).Encode(articles[articleID])
}

func main() {
	router := mux.NewRouter()
	version := "/v1"
	path := "/articles"

	router.HandleFunc(version+path, index).Methods("GET")
	router.HandleFunc(version+path, create).Methods("POST")
	router.HandleFunc(version+path+"/{article}", read).Methods("GET")
	router.HandleFunc(version+path+"/{article}", replace).Methods("PUT")
	router.HandleFunc(version+path+"/{article}", modify).Methods("PATCH")
	router.HandleFunc(version+path+"/{article}", remove).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":80", router))
}
