package store

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

type Controller struct {
	Repository Repository
}

//Middleware to handle aunthentication
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := authorizationHeader

			token, error := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte("secret"), nil
			})
			if error != nil {
				json.NewEncoder(w).Encode(Exception{Message: error.Error()})
				return
			}
			if token.Valid {
				log.Println("Token was valid")
				context.Set(req, "decoded", token.Claims)
				next(w, req)
			} else {
				json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
			}

		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

func (c *Controller) GetToken(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	log.Println("Username: " + user.Username)
	log.Println("Password: " + user.Password)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

// GET/
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	products := c.Repository.GetProducts() // List of all products
	data, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddProduct POST/
func (c *Controller) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // Read the body of the request

	if err != nil {
		log.Fatalln("Error adding product", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("1")
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error adding product", err)
	}

	if err := json.Unmarshal(body, &product); err != nil { // Unmarshall body contents as a type Candidate
		w.WriteHeader(422) // Unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error adding product unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	log.Println(product)
	success := c.Repository.AddProduct(product) // Adding product to DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}
