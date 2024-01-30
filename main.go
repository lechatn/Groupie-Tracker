package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//"net/url"
	//"sort"
	//"strings"
	"text/template"
	//"sort"
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
	Id 	  int      `json:"id"`
	Locations []string `json:"locations"`
}

type Relations struct {
	Id 	  int      					 `json:"id"`
	DateLocation map[string][]string `json:"datesLocations"`	
}

type DatesAndArtists struct {
	Artist Artist
	Dates   Dates
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tHome := template.Must(template.ParseFiles("./templates/home.html"))
		tHome.Execute(w, nil)
	})

	http.HandleFunc("/artistes", func(w http.ResponseWriter, r *http.Request) {
		loadArtistes(w, r)
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


func loadArtistes(w http.ResponseWriter, r *http.Request) {
	url_Artists := "https://groupietrackers.herokuapp.com/api/artists"

	var jsonList_Artists []Artist
	response_Artists, err := http.Get(url_Artists)
	if err != nil {
		fmt.Println("Error1")
		return
	}

	defer response_Artists.Body.Close()

	body_Artists, err := io.ReadAll(response_Artists.Body)
	if err != nil {
		fmt.Println("Error5")
		return
	}
	errUnmarshall1 := json.Unmarshal(body_Artists, &jsonList_Artists)
	if errUnmarshall1 != nil {
		fmt.Println("Error6")
		return
	}

	//fmt.Println(jsonList_Artists)

	tArtistes := template.Must(template.ParseFiles("./templates/artistes.html")) // Read the artists page
	tArtistes.Execute(w, jsonList_Artists)

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

/*func SearchArtist(w http.ResponseWriter, r *http.Request, Data []DatesAndArtists, originalData []DatesAndArtists) []DatesAndArtists {
		if len(Data) != len(originalData){
			Data = originalData
		}
		var new_data []DatesAndArtists
        lettre := r.FormValue("Check")
        fmt.Println(lettre)
		if lettre == "" {
			return Data
		}
		if strings.ToUpper(lettre) == "ALL"{
			return originalData
		}
		for i := 0; i < len(Data); i++ {
			for j := 0; j < len(lettre); j++ {
				if strings.ToUpper(string(Data[i].Artist.Name[j])) == strings.ToUpper(string(lettre[j])){
					if j == len(lettre)-1{
						new_data = append(new_data, Data[i])
					} else {
						if j == len(Data[i].Artist.Name)-1{
							new_data = append(new_data, Data[i])
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
		return new_data
	}


	func SortData(w http.ResponseWriter, r *http.Request, Data []DatesAndArtists) []DatesAndArtists {
		order1 := r.FormValue("alpha")
		order2 := r.FormValue("unalpha")
		order3 := r.FormValue("firstalbum")
		if order1 == "" && order2 == "" && order3 == ""{
			return Data
		}
		if order1 != "" {
			sort.Slice(Data, func(i, j int) bool {
		 	return Data[i].Artist.Name < Data[j].Artist.Name })
		} else if order2 != "" {
			sort.Slice(Data, func(i, j int) bool {
		 	return Data[i].Artist.Name > Data[j].Artist.Name })
		} else if order3 != "" {
			sort.Slice(Data, func(i, j int) bool {
			return Data[i].Artist.FirstAlbum[6:] < Data[j].Artist.FirstAlbum[6:] })	
		}
		return Data
		
	}
*/