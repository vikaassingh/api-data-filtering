package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	Address   Address
}

type Users struct {
	Users []User `json:"users"`
}

type Address struct {
	City string
}

type Recipe struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"userId"`
	Name   string `json:"name"`
}

type RecipePage struct {
	Recipes []Recipe `json:"recipes"`
}

type UserRecipe struct {
	UserID     uint   `json:"user_id"`
	RecipeName string `json:"recipe_name"`
	City       string `json:"city"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user-recipe", GetData)

	http.ListenAndServe(":8080", r)
}

func GetData(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("https://dummyjson.com/users")
	if err != nil {
		fmt.Printf("error dring get user api response: %s", err)
	}

	body := resp.Body
	defer body.Close()
	decoder := json.NewDecoder(body)
	var userData Users
	err = decoder.Decode(&userData)
	if err != nil {
		fmt.Printf("error dring decoding user api response: %s", err)
	}
	m := map[uint]User{}
	for _, usr := range userData.Users {
		m[usr.ID] = usr
	}
	// fmt.Println(m)
	// Get recipies

	resp, err = http.Get("https://dummyjson.com/recipes")
	if err != nil {
		fmt.Printf("error dring get recipes api response: %s", err)
	}

	rbody := resp.Body
	defer rbody.Close()
	decoder = json.NewDecoder(rbody)
	var recipeData RecipePage
	err = decoder.Decode(&recipeData)
	if err != nil {
		fmt.Printf("error dring decoding recipes api response: %s", err)
	}
	// fmt.Println(recipeData)
	var userRecipe []UserRecipe
	for _, rec := range recipeData.Recipes {
		if _, ok := m[rec.UserID]; ok {
			userRecipe = append(userRecipe, UserRecipe{
				UserID:     rec.UserID,
				RecipeName: rec.Name,
				City:       m[rec.UserID].Address.City,
			})

		}
	}
	// fmt.Println(userRecipe)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "Application/json")
	bytes, err := json.Marshal(userRecipe)
	if err != nil {
		fmt.Printf("error during marshal the user recipe data")
	}
	w.Write(bytes)
}
