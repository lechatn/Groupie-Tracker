package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

// Define all the struct and some variables

type Artist struct {
	IdArtists    int      `json:"id"`
	Images       string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type DataLocation struct {
	Locations []string `json:"locations"`
	Id        int      `json:"id"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Index   []string `json:"index"`
	IdDates int      `json:"id"`
	Dates   []string `json:"dates"`
}

type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Relations struct {
	Id           int                 `json:"id"`
	DateLocation map[string][]string `json:"datesLocations"`
}

type DatesAndArtists struct {
	Artist    Artist
	Dates     Dates
	Locations Locations
	Relations Relations
}

var homeData map[string]interface{}

var jsonList_Artists []Artist

var jsonList_Location []Locations
var allLocation map[string][]Locations

var jsonList_Dates []Dates
var homeDates map[string][]Dates

var jsonList_Relations []Relations
var allRelations map[string][]Relations

var port = ":8768"

// ///////////////////////////////////////////

func main() {
	fmt.Println()
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

	errUnmarshall := json.Unmarshal(body_General, &homeData)
	if errUnmarshall != nil {
		fmt.Println("Error3")
		return
	}

	////////////////////////////////////////////////////////////////////////

	url_Artists := homeData["artists"].(string)
	url_Locations := homeData["locations"].(string)
	url_Dates := homeData["dates"].(string)
	url_Relations := homeData["relation"].(string)

	/// ////////////////////////////////////////////////////////////////////Partie artsites
	response_Artists, err := http.Get(url_Artists)
	if err != nil {
		fmt.Println("Error4")
		return
	}

	defer response_Artists.Body.Close()

	body_Artists, err := io.ReadAll(response_Artists.Body)
	if err != nil {
		fmt.Println("Error5")
		return
	}
	errUnmarshall2 := json.Unmarshal(body_Artists, &jsonList_Artists)
	if errUnmarshall2 != nil {
		fmt.Println("Error6")
		return
	}
	////////////////////////////////////////////////////////////////////////////////////
	response_Dates, err := http.Get(url_Dates)
	if err != nil {
		fmt.Println("Error7")
		return
	}

	defer response_Dates.Body.Close()

	body_Dates, err := io.ReadAll(response_Dates.Body)
	if err != nil {
		fmt.Println("Error8")
		return
	}

	errUnmarshall3 := json.Unmarshal(body_Dates, &homeDates)
	if errUnmarshall3 != nil {
		fmt.Println("Error9")
		return
	}
	jsonList_Dates = homeDates["index"]

	////////////////////////////////////////////////////////////////////////////////////

	////////////////////////////////////////////////////////////////////////////////////
	response_Location, err := http.Get(url_Locations)
	if err != nil {
		fmt.Println("Error7")
		return
	}

	defer response_Location.Body.Close()

	body_Location, err := io.ReadAll(response_Location.Body)
	if err != nil {
		fmt.Println("Error8")
		return
	}

	errUnmarshall4 := json.Unmarshal(body_Location, &allLocation)
	if errUnmarshall4 != nil {
		fmt.Println("Error9")
		return
	}
	jsonList_Location = allLocation["index"]
	//fmt.Println(jsonList_Dates)

	////////////////////////////////////////////////////////////////////////////////////

	/// ////////////////////////////////////////////////////////////////////Partie artsites
	response_Relations, err := http.Get(url_Relations)
	if err != nil {
		fmt.Println("Error4")
		return
	}

	defer response_Relations.Body.Close()

	body_Relations, err := io.ReadAll(response_Relations.Body)
	if err != nil {
		fmt.Println("Error5")
		return
	}

	errUnmarshall5 := json.Unmarshal(body_Relations, &allRelations)
	if errUnmarshall5 != nil {
		fmt.Println("Error6")
		return
	}

	jsonList_Relations = allRelations["index"]
	////////////////////////////////////////////////////////////////////////////////////

	Data := createData(jsonList_Artists, jsonList_Dates, jsonList_Location, jsonList_Relations)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tHome := template.Must(template.ParseFiles("./templates/home.html"))
		new_data := SearchBarMembers(w, r, Data)
		tHome.Execute(w, new_data)
		new_data = Data
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
		tLocation.Execute(w, jsonList_Relations)
	})

	fmt.Println("http://localhost" + port) // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)

}

func createData(jsonList_Artists []Artist, jsonList_Dates []Dates, jsonList_Location []Locations, jsonRelations []Relations) []DatesAndArtists {
	var Data []DatesAndArtists

	for i := 0; i < len(jsonList_Artists); i++ {
		for j := 0; j < len(jsonList_Dates); j++ {
			if jsonList_Artists[i].IdArtists == jsonList_Dates[j].IdDates {
				if jsonList_Location[j].Id == jsonList_Dates[j].IdDates {
					if jsonList_Artists[i].IdArtists == jsonRelations[i].Id {
						var inter DatesAndArtists
						inter.Artist = jsonList_Artists[i]
						inter.Dates = jsonList_Dates[j]
						inter.Locations = jsonList_Location[j]
						inter.Relations = jsonList_Relations[i]
						Data = append(Data, inter)
					}

				}
			}
			continue
		}
	}
	return Data
}

func SearchBar(w http.ResponseWriter, r *http.Request, Data []DatesAndArtists) []DatesAndArtists {
	var new_data []DatesAndArtists
	lettre := r.FormValue("Check")
	if lettre == "" {
		return Data
	}
	for i := 0; i < len(Data); i++ {
		for j := 0; j < len(lettre); j++ {
			if strings.ToUpper(string(Data[i].Artist.Name[j])) == strings.ToUpper(string(lettre[j])) {
				if j == len(lettre)-1 {
					new_data = append(new_data, Data[i])
				} else {
					continue
				}
			} else {
				break
			}
		}
	}
	return new_data
}

func SearchBarMembers(w http.ResponseWriter, r *http.Request, Data []DatesAndArtists) []DatesAndArtists {
	var new_data []DatesAndArtists
	find := false
	lettre := r.FormValue("Check")
	if lettre == "" {
		return Data
	}
	for i := 0; i < len(Data); i++ {
		for k := 0; k < len(Data[i].Artist.Members); k++ {
			if find == true {
				find = false
				break
			}
			for j := 0; j < len(lettre); j++ {
				if strings.ToUpper(string(Data[i].Artist.Members[k][j])) == strings.ToUpper(string(lettre[j])) {
					if j == len(lettre)-1 {
						new_data = append(new_data, Data[i])
						find = true
					} else {
						continue
					}
				} else {
					break
				}
			}
		}
	}
	return new_data
}
