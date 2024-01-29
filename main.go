package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

var port = ":8768"

type dataApi struct {
	IdArtists    int      `json:"id"`
	Images       string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type datesData struct {
	Index   []string `json:"index"`
	IdDates int      `json:"id"`
	Dates   []string `json:"dates"`
}

type base struct {
	Art []dataApi
	Dat []datesData
}

// var art []dataApi
// var dat []datesData

var jsonList_Artists []dataApi
var generalData map[string]interface{}

// ///////////////////////////////////////////

var jsonList_Dates []datesData
var dataDates map[string][]datesData

func main() {
	css := http.FileServer(http.Dir("style"))                // For add css to the html pages
	http.Handle("/style/", http.StripPrefix("/style/", css)) // For add css to the html pages
	img := http.FileServer(http.Dir("images"))               // For add css to the html pages
	http.Handle("/images/", http.StripPrefix("/images/", img))

	url_General := "https://groupietrackers.herokuapp.com/api"

	////////////////////////////////////////////////////////////////////////////

	response_General, err := http.Get(url_General)
	if err != nil {
		fmt.Println("Error1")
		return
	}
	defer response_General.Body.Close()

	body_General, err := io.ReadAll(response_General.Body)
	if err != nil {
		fmt.Println("Error2")
		return
	}

	errUnmarshall := json.Unmarshal(body_General, &generalData)
	if errUnmarshall != nil {
		fmt.Println("Error3")
		return
	}

	////////////////////////////////////////////////////////////////////////

	url_Artists := generalData["artists"].(string)
	//url_Locations := generalData["locations"].(string)
	url_Dates := generalData["dates"].(string)
	//url_Relations := generalData["relation"].(string)

	/// ////////////////////////////////////////////////////////////////////Partie artsites
	response_Artists, err := http.Get(url_Artists)
	if err != nil {
		fmt.Println("Error4")
		return
	}

	defer response_Artists.Body.Close()

	body_Artists, err := io.ReadAll(response_Artists.Body)
	if err != nil {
		fmt.Println("Error6")
		return
	}
	errUnmarshall2 := json.Unmarshal(body_Artists, &jsonList_Artists)
	if errUnmarshall2 != nil {
		fmt.Println("Error7")
		return
	}
	////////////////////////////////////////////////////////////////////////////////////
	response_Dates, err := http.Get(url_Dates)
	if err != nil {
		fmt.Println("Error4")
		return
	}

	defer response_Dates.Body.Close()

	body_Dates, err := io.ReadAll(response_Dates.Body)
	if err != nil {
		fmt.Println("Error6")
		return
	}

	errUnmarshall3 := json.Unmarshal(body_Dates, &dataDates)
	if errUnmarshall3 != nil {
		fmt.Println("Error7")
		return
	}

	////////////////////////////////////////////////////////////////////////////////////

	jsonList_Dates = dataDates["index"]
	art := jsonList_Artists
	dat := jsonList_Dates

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Lunch a new page for the lose condition
		tHome := template.Must(template.ParseFiles("./templates/hom.html")) // Read the home page
		data := base{
			art,
			dat,
		}
		tHome.Execute(w, data)
	})

	http.HandleFunc("/artistes", func(w http.ResponseWriter, r *http.Request) {
		tArtistes := template.Must(template.ParseFiles("./templates/artistes.html")) // Read the artists page
		tArtistes.Execute(w, nil)
	})

	http.HandleFunc("/dates", func(w http.ResponseWriter, r *http.Request) {
		tDates := template.Must(template.ParseFiles("./templates/dates.html")) // Read the dates page
		tDates.Execute(w, nil)
	})

	http.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		tLocation := template.Must(template.ParseFiles("./templates/location.html")) // Read the location page
		tLocation.Execute(w, nil)
	})

	fmt.Println("http://localhost:8768") // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)

}
