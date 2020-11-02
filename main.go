package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"./controllers"
	"./database"
	"./middlewares"
)

func init() {
	database.DbConnect()
}

func main() {
	defer database.ClientDisconnect()
	defer database.ConCancel()

	//routes
	r := mux.NewRouter()

	//middlewares
	r.Use(middlewares.AddAccessHeader)
	r.Use(middlewares.CheckAuth)

	//static file server
	r.PathPrefix("/uploads/images/").Handler(http.StripPrefix("/uploads/images/", http.FileServer(http.Dir("./uploads/images"))))

	//user routes
	r.HandleFunc("/api/users", controllers.GetUsers).Methods("GET")
	r.HandleFunc("/api/users/signup", controllers.SignUp).Methods("POST")
	// r.HandleFunc("/api/users/login", ).Methods("POST")

	// //places routes
	r.HandleFunc("/api/places/user/{uid}", controllers.GetPlacesByUserID).Methods("GET")
	r.HandleFunc("/api/places/{pid}", controllers.GetPlaceByID).Methods("GET")
	r.HandleFunc("/api/places", controllers.CreatePlace).Methods("POST")
	r.HandleFunc("/api/places/{pid}", controllers.UpdatePlace).Methods("PATCH")
	r.HandleFunc("/api/places/{pid}", controllers.DeletePlace).Methods("DELETE")
	http.ListenAndServe(":5000", r)
}
