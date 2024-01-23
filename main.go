package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

var port = ":8888"

type base struct {
	ApiImage []string
	ApiName  []string
}

type dataApi struct {
	Images       string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDate  string   `json:"concertDate"`
	Relation     string   `json:"relation"`
}

var jsonList []dataApi
var artistImage []string
var artistName []string

func main() {
	// css := http.FileServer(http.Dir("style"))                // For add css to the html pages
	// http.Handle("/style/", http.StripPrefix("/style/", css)) // For add css to the html pages
	url := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error")
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error")
		return
	}
	json.Unmarshal(body, &jsonList)
	for _, artistImg := range jsonList {
		artistImage = append(artistImage, artistImg.Images)
	}
	for _, artistNa := range jsonList {
		artistName = append(artistName, artistNa.Images)
	}
	apiImage := artistImage[0:51]
	apiName := artistName[0:51]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Lunch a new page for the lose condition
		tmpl := template.Must(template.ParseFiles("./templates/home.html")) // Read the home page
		data := base{
			apiImage,
			apiName,
		}
		tmpl.Execute(w, data)
	})

	fmt.Println("http://localhost:8888") // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)
}
