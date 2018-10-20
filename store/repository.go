package store

import (
	"log"
	"fmt"

	"gopkg.in/mgo.v2"
)

type Repository struct{}

const SERVER = "mongodb://jastin:jastin10@ds137483.mlab.com:37483/db-store"
const DBNAME = "db-store"
const COLLECTION = "store"

// GetProduct
func (r Repository) GetProducts() Products {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to connect to DB", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Products{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to fetch products", err)
	}

	return results
}

// AddProduct
func (r Repository) AddProduct(product Product) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	session.DB(DBNAME).C(COLLECTION).Insert(product)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("New product added to DB")

	return true
}
