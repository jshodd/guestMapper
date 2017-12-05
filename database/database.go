package database

import (
	"encoding/json"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io"
	"log"
	"net/http"
)

var driver bolt.Driver
var conn bolt.Conn
var err error

// D3Response is the graph response
type D3Response struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"edges"`
}

// Node is the graph response node
type Node struct {
	Id     int    `json:"id"`
	Title  string `json:"caption"`
}

// Link is the graph response link
type Link struct {
	Source int    `json:"source"`
	Target int    `json:"target"`
	Label  string `json:"caption"`
}
type Person struct {
	Name   string
}
type Relationship struct {
	P1       Person
	P2       Person
	Relation string
}

func ClearDatabase() {
	_, err := conn.ExecNeo("MATCH (n) DETACH DELETE n", nil)
	handleError(err)
	fmt.Println("Database Cleared")
}

func CreateNode(person Person) {
	result, _, _, err := conn.QueryNeoAll("MATCH (n:PERSON {name: {name}}) RETURN n", map[string]interface{}{"name": person.Name})
	handleError(err)
	if len(result) == 0 {
		stmt := fmt.Sprintf("CREATE (:PERSON {name: {name}})")
		_, err := conn.ExecNeo(stmt, map[string]interface{}{"name": person.Name})
		handleError(err)
		fmt.Printf("You Created: %s\n", person.Name)
		return
	}
	fmt.Println("That node already exists!")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateNodeRelationship(relation Relationship) {
	result, _, _, err := conn.QueryNeoAll("MATCH (p1:PERSON {name: {p1}})-[r]-(P2:PERSON {name: {P2}}) RETURN r", map[string]interface{}{"p1": relation.P1.Name, "P2": relation.P2.Name})
	handleError(err)

	if len(result) == 0 {

		stmt := fmt.Sprintf("MATCH (p1:PERSON {name: {p1}}) MATCH (P2:PERSON {name: {P2}}) CREATE (p1)-[:%s]->(P2)", relation.Relation)
		_, err := conn.ExecNeo(stmt, map[string]interface{}{"p1": relation.P1.Name, "P2": relation.P2.Name})
		handleError(err)
		fmt.Printf("You Created Relation \"%s\" Between %s and %s\n", relation.Relation, relation.P1.Name, relation.P2.Name)
		return
	}
	fmt.Println("That Relation already exists!")
}

func GenerateTestData() {
	tasha := Person{"Natasha Pacifico"}
	CreateNode(tasha)

	jake := Person{"Jacob Shodd"}
	CreateNode(jake)

	CreateNodeRelationship(Relationship{tasha, jake, "Engaged"})

	mike := Person{"Mike Shodd"}
	CreateNode(mike)
	CreateNodeRelationship(Relationship{mike, jake, "Father"})

	sara := Person{"Sara Shodd"}
	CreateNode(sara)
	CreateNodeRelationship(Relationship{sara, jake, "Mother"})

	angie := Person{"Angela Pacifico"}
	CreateNode(angie)
	CreateNodeRelationship(Relationship{angie, tasha, "Mother"})

	jerry := Person{"Jerry Pacifico"}
	CreateNode(jerry)
	CreateNodeRelationship(Relationship{jerry, tasha, "Father"})

	sam := Person{"Samantha Shodd"}
	CreateNode(sam)
	CreateNodeRelationship(Relationship{sam, jake, "Sister"})

	niko := Person{"Niko Pacifico"}
	CreateNode(niko)
	CreateNodeRelationship(Relationship{niko, tasha, "Brother"})
}

func Init() {
	driver = bolt.NewDriver()
	conn, err = driver.OpenNeo("bolt://localhost:7687")
	handleError(err)
}

func Close() {
	conn.Close()
}

func ExportGraph(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cypher := `
	MATCH
		(p:PERSON)-[r]->(q:PERSON)
	RETURN
		p.name,q.name,type(r),p.family,q.family`

	stmt, err := conn.PrepareNeo(cypher)
	if err != nil {
		log.Println("error preparing graph:", err)
		w.WriteHeader(500)
		w.Write([]byte("An error occurred querying the DB"))
		return
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(nil)
	if err != nil {
		log.Println("error querying graph:", err)
		w.WriteHeader(500)
		w.Write([]byte("An error occurred querying the DB"))
		return
	}

	d3Resp := D3Response{}
	row, _, err := rows.NextNeo()
	for row != nil && err == nil {
		p1 := row[0].(string)
		p2 := row[1].(string)
		r := row[2].(string)
		fmt.Printf("Processing %s and %s\n", p1, p2)
		check1 := -1
		check2 := -1
		for i, node := range d3Resp.Nodes {
			if p1 == node.Title {
				check1 = i
			} else if p2 == node.Title {
				check2 = i
			}
		}
		if check1 == -1 {
			if p1 == "Jacob Shodd" || p1 == "Natasha Pacifico" {
			}
			check1 = len(d3Resp.Nodes)
			d3Resp.Nodes = append(d3Resp.Nodes, Node{Id: check1, Title: p1})

		}
		if check2 == -1 {
			if p2 == "Jacob Shodd" || p2 == "Natasha Pacifico" {
			}
			check2 = len(d3Resp.Nodes)
			d3Resp.Nodes = append(d3Resp.Nodes, Node{Id: check2, Title: p2})

		}
		if (p2 == "Jacob Shodd" || p2 == "Natasha Pacifico") && (p1 == "Jacob Shodd" || p1 == "Natasha Pacifico") {
			r = "Engaged"
		}
		d3Resp.Links = append(d3Resp.Links, Link{Source: check1, Target: check2, Label: r})
		row, _, err = rows.NextNeo()
	}

	if err != nil && err != io.EOF {
		log.Println("error querying graph:", err)
		w.WriteHeader(500)
		w.Write([]byte("An error occurred querying the DB"))
		return
	} else if len(d3Resp.Nodes) == 0 {
		w.WriteHeader(404)
		return
	}

	err = json.NewEncoder(w).Encode(d3Resp)
	if err != nil {
		log.Println("error writing graph response:", err)
		w.WriteHeader(500)
		w.Write([]byte("An error occurred writing response"))
	}
}
