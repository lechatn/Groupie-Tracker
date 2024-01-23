package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

var port = ":8080"

type base struct {
	ApiData string
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
	// j := map[string]any{}
	for _, artist := range jsonList {
		artistData = append(artistData, artist.Images)
	}

	apiData := artistData[0]
	fmt.Println(apiData)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Lunch a new page for the lose condition
		tmpl := template.Must(template.ParseFiles("./templates/home.html")) // Read the home page
		data := base{
			apiData,
		}
		tmpl.Execute(w, data)
	})

	fmt.Println("http://localhost:8080") // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)
}
