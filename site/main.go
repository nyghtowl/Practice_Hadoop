// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type handler struct {
	hf http.HandlerFunc
}


type RadioButton struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type templateParams struct {
	PageTitle       string
	Date 			string
	PageRadioButtons []RadioButton
	Answer          string

}


func main() {
	//Start the web server, set the port to listen to 8080. Without assumes localhost.
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
		log.Print("Defaulting to port %s ", port)
	}

	//Code adjust from https://github.com/campoy/go-web-workshop/blob/master/section02/README.md
/*	r := mux.NewRouter()*/

	//This method takes in the URL path "/" and a function that takes in a response writer, and a http request.
	/*r.HandleFunc("/", indexHandler)*/

	//// match only GET requests on /media/
	//r.HandleFunc("/media/", listMedia).Methods("GET")
	//
	//// match only POST requests on /media/
	//r.HandleFunc("/media/", addMedia).Methods("POST")
	//
	//// match GET regardless of mediaID
	//r.HandleFunc("/media/{mediaID}", getMedia)

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver
	http.Handle("/", handler{indexHandler})
	http.Handle("/selected", handler{UserSelected})
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func (h *handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	// Just redirect to the first one
	http.Redirect(w, r, "/site/", http.StatusTemporaryRedirect)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	current := r.URL.Path
	/*	w.Header().Set("Content-Type", "text/plain")*/
	if current == "/" || current == "/selected" {
		h.hf(w, r)
		return
	} else {
		http.NotFound(w, r)
		return
	}

}
func indexHandler(w http.ResponseWriter, r *http.Request) {

	params := templateParams {}
	params.PageTitle = "Flip th Script"
	params.PageRadioButtons = []RadioButton{
		RadioButton{"animalselect", "cats", false, false, "Cats"},
		RadioButton{"animalselect", "dogs", false, false, "Dogs"},}
    params.Date = time.Now().Format("02-01-2006")


	t := template.Must(template.ParseFiles("content/index.html"))
/*	if err := getBQData(r); err != nil {
		log.Fatalf("Error getting data: %v\n", err)
	}
*/
	//execute the template and pass it the params struct to fill in the gaps
	if err := t.ExecuteTemplate(w, "index.html", params); err !=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print("Form execute error: ", err)
	}
}
/*https://blog.scottlogic.com/2017/02/28/building-a-web-app-with-go.html*/
func UserSelected(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	// r.Form is now either
	// map[animalselect:[cats]] OR
	// map[animalselect:[dogs]]
	// so get the animal which has been selected
	youranimal := r.Form.Get("animalselect")

	params := templateParams {}
	params.PageTitle = "Your preferred animal"
	params.Answer = youranimal


	// generate page by passing page variables into template

	t := template.Must(template.ParseFiles("content/index.html"))

	//execute the template and pass it the params struct to fill in the gaps
	if err := t.ExecuteTemplate(w, "index.html", params); err !=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print("Post execute error: ", err)
	}
}