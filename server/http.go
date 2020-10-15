package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/jadams/o365-attack-toolkit-us/api"
	"github.com/jadams/o365-attack-toolkit-us/model"

	"github.com/gorilla/mux"
)

func StartIntServer(config model.Config) {

	log.Printf("Starting Internal Server on 127.0.0.1:%d \n", config.Server.InternalPort)

	route := mux.NewRouter()

	route.HandleFunc("/", GetUsers).Methods("GET")
	route.HandleFunc(model.IntAbout, GetAbout).Methods("GET")

	// Routes for Users
	route.HandleFunc(model.IntGetAll, GetUsers).Methods("GET")

	// Route for files
	route.HandleFunc(model.IntUserFiles, GetUserFiles).Methods("GET")
	route.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("downloads/"))))

	//Route for emails
	route.HandleFunc(model.IntUserEmails, GetUserEmails).Methods("GET")
	route.HandleFunc(model.IntUserEmails, SearchUserEmails).Methods("POST") //  For searching

	route.HandleFunc(model.IntAllEmails, GetAllEmails).Methods("GET")
	route.HandleFunc(model.IntAllEmails, SearchEmails).Methods("POST") // For Searching

	route.HandleFunc(model.IntUserEmail, GetEmail).Methods("GET")

	// The route for the file downloads.

	route.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Server.Host, config.Server.InternalPort),
		Handler: route,
	}
	server.ListenAndServe()

}

func StartExtServer(config model.Config) {

	log.Printf("Starting External Server on %s:%d \n", config.Server.Host, config.Server.ExternalPort)
	route := mux.NewRouter()
	route.HandleFunc(model.ExtTokenPage, GetToken).Methods("GET")
	route.PathPrefix(model.ExtMainPage).Handler(http.FileServer(http.Dir("./static/")))
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Server.Host, config.Server.ExternalPort),
		Handler: route,
	}
	//server.ListenAndServeTLS(config.Server.Certificate,config.Server.Key)
	server.ListenAndServe()
}

// GetToken will handle the request and initilize the thing with the code
func GetToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	r.ParseForm()
	jsonData := api.GetAllTokens(r.FormValue("code"))
	if jsonData != nil {
		authResponse := model.AuthResponse{}
		json.Unmarshal(jsonData, &authResponse)
		log.Println("Authentication Successful")
		go api.InitializeProfile(authResponse.AccessToken, authResponse.RefreshToken)
		// Do the redirect here
		w.Write([]byte("{\"status\":\"OK\"}"))
	}
}
