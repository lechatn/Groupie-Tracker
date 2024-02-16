package Server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"Groupie/structure"
)

func LoadArtistes(w http.ResponseWriter, r *http.Request) []structure.Artist {
	url_Artists := "https://groupietrackers.herokuapp.com/api/artists"

	var jsonList_Artists []structure.Artist
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

func LoadLocation(w http.ResponseWriter, r *http.Request, data []structure.Artist) {

	//A faire

}

func LoadRelation(w http.ResponseWriter, r *http.Request, id string, infos_artist structure.Artist) structure.Relations {

	var json_Relation structure.Relations

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

	data := structure.Relations{}
	data.Id = json_Relation.Id
	data.DateLocation = json_Relation.DateLocation
	data.Infos = infos_artist
	data.Coordonnees = SearchLatLon(data.DateLocation)

	tRelation := template.Must(template.ParseFiles("./templates/relation.html")) // Read the relation page
	//fmt.Println(data)
	tRelation.Execute(w, data)

	return data
}


func SearchLatLon(relation map[string][]string) map[string][]string {
	res := make(map[string][]string, len(relation))
	url := ""
	for city := range relation {
		city = strings.ReplaceAll(city, "-", ",")
		fmt.Println(city)
		if city == "willemstad,netherlands_antilles" {
			url = "https://nominatim.openstreetmap.org/search?q=willemstad&format=json"
		} else {
			url = "https://nominatim.openstreetmap.org/search?q=" + city + "&format=json"
		}

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

		var data []structure.Place
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