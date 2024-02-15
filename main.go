package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"sort"
	"strconv"
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
	Coordonnees  map[string][]string
	Infos        Artist
}

type Place struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

var jsonList_Artists []Artist

var allLocation map[string][]Locations

var json_Relation Relations

var port = ":8768"

var artist_create = false
var originalData []Artist

var relation map[string][]string
var location []Locations

var data_artist Relations

var infos_artist Artist

// ///////////////////////////////////////////

func main() {
	css := http.FileServer(http.Dir("style"))                // For add css to the html pages
	http.Handle("/style/", http.StripPrefix("/style/", css)) // For add css to the html pages
	img := http.FileServer(http.Dir("images"))               // For add css to the html pages
	http.Handle("/images/", http.StripPrefix("/images/", img))
	js := http.FileServer(http.Dir("js")) // For add css to the html pages
	http.Handle("/js/", http.StripPrefix("/js/", js))

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
		if r.FormValue("Check") != "" {
			lettre := r.FormValue("Check")
			jsonList_Artists = SearchArtist(w, r, jsonList_Artists, originalData, lettre)
			lettre = ""
		}
		jsonList_Artists = SortData(w, r, jsonList_Artists)
		tArtistes.Execute(w, jsonList_Artists)
	})

	http.HandleFunc("/dates", func(w http.ResponseWriter, r *http.Request) {
		loadDates(w, r)
	})

	http.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		location = loadLocation(w, r)
	})

	http.HandleFunc("/relation", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		fmt.Println(id)
		id_int, _ := strconv.Atoi(id)
		sort.Slice(originalData, func(i, j int) bool {
			return originalData[i].IdArtists < originalData[j].IdArtists
		})
		if len(originalData) == len(jsonList_Artists) {
			fmt.Println("cas1")
			infos_artist = jsonList_Artists[id_int-1]
		} else {
			fmt.Println("cas2")
			infos_artist = originalData[id_int-1]
		}
		data_artist = loadRelation(w, r, id, infos_artist)
	})

	http.HandleFunc("/relationForJs", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(data_artist)
		data_artist = Relations{}
	})

	http.HandleFunc("/locationForJs", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(location)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
		fmt.Println("Error60")
		os.Exit(1)
	}
	return jsonList_Artists

}

func loadDates(w http.ResponseWriter, r *http.Request) {
	var jsonList_Dates []Dates
	var homeDates map[string][]Dates

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
		fmt.Println("Error15")
		return
	}
	jsonList_Dates = homeDates["index"]

	tDates := template.Must(template.ParseFiles("./templates/dates.html")) // Read the dates page
	tDates.Execute(w, jsonList_Dates)

}

func loadLocation(w http.ResponseWriter, r *http.Request) []Locations {

	var allLocation map[string][]Locations

	url_Locations := "https://groupietrackers.herokuapp.com/api/locations"
	response_Location, err := http.Get(url_Locations)
	if err != nil {
		fmt.Println("Error7")
		os.Exit(0)
	}

	defer response_Location.Body.Close()

	body_Location, err := io.ReadAll(response_Location.Body)
	if err != nil {
		fmt.Println("Error8")
		os.Exit(0)
	}

	errUnmarshall4 := json.Unmarshal(body_Location, &allLocation)
	if errUnmarshall4 != nil {
		fmt.Println("Error9")
		os.Exit(0)
	}

	tLocation := template.Must(template.ParseFiles("./templates/location.html")) // Read the location page
	tLocation.Execute(w, nil)

	return allLocation["index"]
}

func loadRelation(w http.ResponseWriter, r *http.Request, id string, infos_artist Artist) Relations {

	var json_Relation Relations

	url_Relations := "https://groupietrackers.herokuapp.com/api/relation/" + id

	response_Relations, err := http.Get(url_Relations)
	if err != nil {
		fmt.Println("Error4")
		os.Exit(0)
	}

	defer response_Relations.Body.Close()

	body_Relations, err := io.ReadAll(response_Relations.Body)
	if err != nil {
		fmt.Println("Error5")
		os.Exit(0)
	}

	//fmt.Println(body_Relations)

	errUnmarshall3 := json.Unmarshal(body_Relations, &json_Relation)
	if errUnmarshall3 != nil {
		fmt.Println(errUnmarshall3)
		os.Exit(0)
	}
	
	data := Relations{}
	data.Id = json_Relation.Id
	data.DateLocation = json_Relation.DateLocation
	data.Infos = infos_artist
	data.Coordonnees = SearchLatLon(data.DateLocation)

	tRelation := template.Must(template.ParseFiles("./templates/relation.html")) // Read the relation page
	//fmt.Println(data)
	tRelation.Execute(w, data)

	return data
}

func SearchArtist(w http.ResponseWriter, r *http.Request, jsonList_Artists []Artist, originalData []Artist, lettre string) []Artist {
	var new_data []Artist
	//fmt.Println(lettre)
	if lettre == "" {
		return jsonList_Artists
	}
	if len(originalData) != len(jsonList_Artists) {
		//fmt.Println(originalData)
		jsonList_Artists = originalData
	}
	if strings.ToUpper(lettre) == "ALL" {
		return originalData
	}
	for i := 0; i < len(jsonList_Artists); i++ {
		for j := 0; j < len(lettre); j++ {
			if strings.ToUpper(string(jsonList_Artists[i].Name[j])) == strings.ToUpper(string(lettre[j])) {
				if j == len(lettre)-1 {
					new_data = append(new_data, jsonList_Artists[i])
				} else {
					if j == len(jsonList_Artists[i].Name)-1 {
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
	return new_data
}

func SortData(w http.ResponseWriter, r *http.Request, jsonList_Artists []Artist) []Artist {
	order1 := r.FormValue("alpha")
	order2 := r.FormValue("unalpha")
	order3 := r.FormValue("firstalbum")
	order4 := r.FormValue("CreationDate")
	if order1 == "" && order2 == "" && order3 == "" && order4 == "" {
		return jsonList_Artists
	}
	if order1 != "" {
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			return strings.ToUpper(jsonList_Artists[i].Name) < strings.ToUpper(jsonList_Artists[j].Name)
		})
	} else if order2 != "" {
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			return strings.ToUpper(jsonList_Artists[i].Name) > strings.ToUpper(jsonList_Artists[j].Name)
		})
	} else if order3 != "" {
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			if jsonList_Artists[i].FirstAlbum[6:] == jsonList_Artists[j].FirstAlbum[6:] {
				if jsonList_Artists[i].FirstAlbum[3:5] == jsonList_Artists[j].FirstAlbum[3:5] {
					return jsonList_Artists[i].FirstAlbum[:2] < jsonList_Artists[j].FirstAlbum[:2]
				}
				return jsonList_Artists[i].FirstAlbum[3:5] < jsonList_Artists[j].FirstAlbum[3:5]
			}
			return jsonList_Artists[i].FirstAlbum[6:] < jsonList_Artists[j].FirstAlbum[6:]
		})
	} else if order4 != "" {
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			return jsonList_Artists[i].CreationDate < jsonList_Artists[j].CreationDate
		})
	}
	return jsonList_Artists
}

func SearchLatLon(relation map[string][]string) map[string][]string {
	res := make(map[string][]string, len(relation))
	for city := range relation {
		city = strings.ReplaceAll(city, "-", ",")
		url := "https://nominatim.openstreetmap.org/search?q=" + city + "&format=json"
		response, err := http.Get(url)

		if err != nil {
			fmt.Println("Error1")
			os.Exit(0)
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error2")
			os.Exit(0)
		}

		var data []Place
		errUnmarshall := json.Unmarshal(body, &data)
		if errUnmarshall != nil {
			fmt.Println("Error3")
			os.Exit(0)
		}

		var inter []string
		inter = append(inter, data[0].Lat)
		inter = append(inter, data[0].Lon)
		res[city] = inter
		inter = nil

	}
	//fmt.Println(res)
	return res
}

/*		if r.FormValue("Search_artist") != "" {
		//fmt.Println("test")
		lettre := r.FormValue("Search_artist")
		jsonList_Artists = SearchArtist(w, r, jsonList_Artists, originalData, lettre)
		lettre = ""
	}*/
