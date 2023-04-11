package main

import (
	"helpers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var user helpers.User

/*
func main() {
	//Openning a connection to mysql database
	db, err := helpers.Init()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()


	http.HandleFunc("/register", user.Register_user_handler)
	http.HandleFunc("/Login", user.Login_handler)
	log.Println("Listening on HTTP port 8080")
	http.ListenAndServe("localhost:8080", nil)

}*/

func main() {
	r := mux.NewRouter()

	// Allow preflight requests
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	// Register your loginHandler for POST requests
	r.HandleFunc("/Login", user.Login_handler).Methods(http.MethodPost)

	// Start the HTTP server
	http.ListenAndServe(":8080", r)
}
