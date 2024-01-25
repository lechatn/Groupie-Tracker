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
	ApiLocation []string
}

type dataApi struct {
	Images       string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

var jsonList_Artists []dataApi
var artistData []string
var generalData map[string]interface{}

func main() {
	css := http.FileServer(http.Dir("style"))                // For add css to the html pages
	http.Handle("/style/", http.StripPrefix("/style/", css)) // For add css to the html pages
	img := http.FileServer(http.Dir("images"))               // For add css to the html pages
	http.Handle("/images/", http.StripPrefix("/images/", img))

	url_General := "https://groupietrackers.herokuapp.com/api"
	

	////////////////////////////////////////////////////////////////////////////

	response_General, err := http.Get(url_General)
	if err != nil {
		fmt.Println("Error4")
		return
	}
	defer response_General.Body.Close()

	body_General, err := ioutil.ReadAll(response_General.Body)
	if err != nil {
		fmt.Println("Error5")
		return
	}

	errUnmarshall := json.Unmarshal(body_General, &generalData)
	if errUnmarshall != nil {
		fmt.Println("Error6")
		return
	}

	////////////////////////////////////////////////////////////////////////

	url_Artists := generalData["artists"].(string)
	//url_Locations := generalData["locations"].(string)
	//url_Dates := generalData["dates"].(string)
	//url_Relations := generalData["relation"].(string)


/// ////////////////////////////////////////////////////////////////////Partie artsites
	response_Artists, err := http.Get(url_Artists)
	if err != nil {
		fmt.Println("Error1")
		return
	}
	defer response_Artists.Body.Close()



	body_Artists, err := ioutil.ReadAll(response_Artists.Body)
	if err != nil {
		fmt.Println("Error2")
		return
	}
	errUnmarshall2 := json.Unmarshal(body_Artists, &jsonList_Artists)
	if errUnmarshall2 != nil {
		fmt.Println("Error3")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Lunch a new page for the lose condition
		tHome := template.Must(template.ParseFiles("./templates/home.html")) // Read the home page
		tHome.Execute(w, jsonList_Artists)
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
