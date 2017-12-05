package handlers

import (
	"net/http"
	"github.com/jshodd/guestMapper/database"
	"fmt"
	"strings"
)


// ClearDB deletes all nodes and relationships from the database
func ClearDB(w http.ResponseWriter, r *http.Request) {
	database.Init()
	defer database.Close()
	database.ClearDatabase()
	http.Redirect(w,r,"/", 301)
}

// Add creates a Node and Relationship if they do not already exist
func Add(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	person := database.Person{Name: strings.Trim(fmt.Sprintf("%s",r.Form["name"]), "[]")}
	target := database.Person{Name: strings.Trim(fmt.Sprintf("%s",r.Form["target"]), "[]")}
	relation := database.Relationship{person, target, strings.Trim(fmt.Sprintf("%s",r.Form["relation"]),"[]")}
	database.Init()
	defer database.Close()
	database.CreateNode(person)
	database.CreateNodeRelationship(relation)
	http.Redirect(w,r,"/", 301)
}

// GenerateTest loads a set of test data into the database
func GenerateTest(w http.ResponseWriter, r *http.Request){
	database.Init()
	defer database.Close()
	database.GenerateTestData()
	http.Redirect(w,r,"/", 301)
}

// GraphDB returns a json export of the graph for alchemy.js to visualize
func GraphDB(w http.ResponseWriter, r *http.Request){
	database.Init()
	defer database.Close()
	database.ExportGraph(w, r)
}