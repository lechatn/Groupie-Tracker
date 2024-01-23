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
var artistData []string
var nameArtist []string

func main() {
	css := http.FileServer(http.Dir("style"))                // For add css to the html pages
	http.Handle("/style/", http.StripPrefix("/style/", css)) // For add css to the html pages
	img := http.FileServer(http.Dir("images"))               // For add css to the html pages
	http.Handle("/images/", http.StripPrefix("/images/", img))
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
	errUnmarshall := json.Unmarshal(body, &jsonList)
	if errUnmarshall != nil {
		fmt.Println("Error")
		return
	}
	// j := map[string]any{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Lunch a new page for the lose condition
		tHome := template.Must(template.ParseFiles("./templates/home.html")) // Read the home page^
		tHome.Execute(w, jsonList)
	})

	http.HandleFunc("/artistes", func(w http.ResponseWriter, r *http.Request) {
		tArtistes := template.Must(template.ParseFiles("./templates/artistes.html")) // Read the home page
		tArtistes.Execute(w, nil)
	})

	http.HandleFunc("/dates", func(w http.ResponseWriter, r *http.Request) {
		tDates := template.Must(template.ParseFiles("./templates/dates.html")) // Read the home page
		tDates.Execute(w, nil)
	})

	http.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		tLocation := template.Must(template.ParseFiles("./templates/location.html")) // Read the home page
		tLocation.Execute(w, nil)
	})

	fmt.Println("http://localhost:8888") // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)
}
