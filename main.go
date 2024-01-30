package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	//"net/url"
	//"sort"
	"strings"
	"text/template"
	"sort"
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
	Id 	  int      					 `json:"id"`
	DateLocation map[string][]string `json:"datesLocations"`	
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

var artist_create = false
var originalData []Artist

// ///////////////////////////////////////////

func main() {
	fmt.Println()
	css := http.FileServer(http.Dir("style"))                // For add css to the html pages
	http.Handle("/style/", http.StripPrefix("/style/", css)) // For add css to the html pages
	img := http.FileServer(http.Dir("images"))               // For add css to the html pages
	http.Handle("/images/", http.StripPrefix("/images/", img))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tHome := template.Must(template.ParseFiles("./templates/home.html"))
		tHome.Execute(w, nil)
	})

	http.HandleFunc("/artistes", func(w http.ResponseWriter, r *http.Request) {
		if artist_create == false {
			jsonList_Artists = loadArtistes(w, r)
			artist_create = true
			originalData = jsonList_Artists
		}
		
		tArtistes := template.Must(template.ParseFiles("./templates/artistes.html")) // Read the artists page
		jsonList_Artists = SearchArtist(w, r, jsonList_Artists, originalData)
		jsonList_Artists = SortData(w, r, jsonList_Artists)
		tArtistes.Execute(w, jsonList_Artists)
	})

	http.HandleFunc("/dates", func(w http.ResponseWriter, r *http.Request) {
		loadDates(w, r)
	})

	http.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		loadLocation(w, r)
	})

	http.HandleFunc("/relation", func(w http.ResponseWriter, r *http.Request) {
		loadRelation(w, r)
	})

	fmt.Println("http://localhost:8768") // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)

}


func loadArtistes(w http.ResponseWriter, r *http.Request) []Artist {
	url_Artists := "https://groupietrackers.herokuapp.com/api/artists"

	var jsonList_Artists []Artist
	response_Artists, err := http.Get(url_Artists)
	if err != nil {
		fmt.Println("Error1")
		os.Exit(1)
	}

	defer response_Artists.Body.Close()

	body_Artists, err := io.ReadAll(response_Artists.Body)
	if err != nil {
		fmt.Println("Error5")
		os.Exit(1)
	}
	errUnmarshall1 := json.Unmarshal(body_Artists, &jsonList_Artists)
	if errUnmarshall1 != nil {
		fmt.Println("Error6")
		os.Exit(1)
	}
	return jsonList_Artists
	//fmt.Println(jsonList_Artists)

}

func loadDates(w http.ResponseWriter, r *http.Request) {
	url_Dates := "https://groupietrackers.herokuapp.com/api/dates"
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

	errUnmarshall2 := json.Unmarshal(body_Dates, &homeDates)
	if errUnmarshall2 != nil {
		fmt.Println("Error9")
		return
	}
	jsonList_Dates = homeDates["index"]

	tDates := template.Must(template.ParseFiles("./templates/dates.html")) // Read the dates page
	tDates.Execute(w, jsonList_Dates)

}

func loadLocation(w http.ResponseWriter, r *http.Request) {
	url_Locations := "https://groupietrackers.herokuapp.com/api/locations"

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
	//fmt.Println(jsonList_Location)

	tLocation := template.Must(template.ParseFiles("./templates/location.html")) // Read the location page	
	tLocation.Execute(w, jsonList_Location)
}

func loadRelation(w http.ResponseWriter, r *http.Request) {
	url_Relations := "https://groupietrackers.herokuapp.com/api/relation"

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
	fmt.Println(jsonList_Relations)

	tRelation := template.Must(template.ParseFiles("./templates/relation.html")) // Read the relation page
	tRelation.Execute(w, jsonList_Relations)
}

func SearchArtist(w http.ResponseWriter, r *http.Request, jsonList_Artists []Artist, originalData []Artist ) []Artist {
		var new_data []Artist
        lettre := r.FormValue("Check")
        fmt.Println(lettre)
		if lettre == "" {
			return jsonList_Artists
		}
		if len(originalData) != len(jsonList_Artists){
			jsonList_Artists = originalData
		}
		if strings.ToUpper(lettre) == "ALL" {
			return originalData
		}
		for i := 0; i < len(jsonList_Artists); i++ {
			for j := 0; j < len(lettre); j++ {
				if strings.ToUpper(string(jsonList_Artists[i].Name[j])) == strings.ToUpper(string(lettre[j])){
					if j == len(lettre)-1{
						new_data = append(new_data, jsonList_Artists[i])
					} else {
						if j == len(jsonList_Artists[i].Name)-1{
							new_data = append(new_data, jsonList_Artists[i])
							break
						} else {
							continue
						}	
					}
				} else {
					break
				}
			}
		}
		fmt.Println(new_data)
		return new_data
	}


	func SortData(w http.ResponseWriter, r *http.Request, jsonList_Artists []Artist) []Artist {
		order1 := r.FormValue("alpha")
		order2 := r.FormValue("unalpha")
		order3 := r.FormValue("firstalbum")
		if order1 == "" && order2 == "" && order3 == ""{
			return jsonList_Artists
		}
		if order1 != "" {
			sort.Slice(jsonList_Artists, func(i, j int) bool {
		 	return jsonList_Artists[i].Name < jsonList_Artists[j].Name })
		} else if order2 != "" {
			sort.Slice(jsonList_Artists, func(i, j int) bool {
		 	return jsonList_Artists[i].Name > jsonList_Artists[j].Name })
		} else if order3 != "" {
			sort.Slice(jsonList_Artists, func(i, j int) bool {
			return jsonList_Artists[i].FirstAlbum[6:] < jsonList_Artists[j].FirstAlbum[6:] })	
		}
		return jsonList_Artists
		
	}
