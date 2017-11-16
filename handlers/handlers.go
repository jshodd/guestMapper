package handlers

import (
	"net/http"
	"github.com/jshodd/guestMapper/database"
)

func ClearDB(w http.ResponseWriter, r *http.Request) {
	database.Init()
	defer database.Close()
	database.ClearDatabase()
	http.Redirect(w,r,"/", 301)
}

func AddPerson(w http.ResponseWriter, r *http.Request){}

func AddRelationship(w http.ResponseWriter, r *http.Request){}

func GenerateTest(w http.ResponseWriter, r *http.Request){
	database.Init()
	defer database.Close()
	database.GenerateTestData()
	http.Redirect(w,r,"/", 301)
}

func GraphDB(w http.ResponseWriter, r *http.Request){
	database.Init()
	defer database.Close()
	database.ExportGraph(w, r)
}