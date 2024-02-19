package Server

import (
	"Groupie/structure"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func LoadArtistes(w http.ResponseWriter, r *http.Request) []structure.Artist {
	url_Artists := "https://groupietrackers.herokuapp.com/api/artists"

	var jsonList_Artists []structure.Artist
	response_Artists, err := http.Get(url_Artists)
	if err != nil {
		fmt.Println("Erreur de connexion a l'API :", err)
		os.Exit(0)
	}

	defer response_Artists.Body.Close()

	body_Artists, err := io.ReadAll(response_Artists.Body)
	if err != nil {
		fmt.Println("Erreur de lecture de l'API :", err)
		os.Exit(0)
	}
	errUnmarshall1 := json.Unmarshal(body_Artists, &jsonList_Artists)
	if errUnmarshall1 != nil {
		fmt.Println("Erreur de décodage de l'API :", errUnmarshall1)
		os.Exit(0)
	}
	return jsonList_Artists

}

func LoadLocation(w http.ResponseWriter, r *http.Request, data []structure.Artist) {
	//res := make([]map[string][]string, len(data))
	for i := 0; i < len(data); i++ {
		url_location := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(data[i].IdArtists)

		response_location, err := http.Get(url_location)
		if err != nil {
			fmt.Println("Erreur de connexion à l'API :", err)
			os.Exit(0)
		}

		defer response_location.Body.Close()

		body_location, err := io.ReadAll(response_location.Body)
		if err != nil {
			fmt.Println("Erreur de lecture de l'API :", err)
			os.Exit(0)
		}

		var json_location structure.Relations
		errUnmarshall2 := json.Unmarshal(body_location, &json_location)
		if errUnmarshall2 != nil {
			fmt.Println("Erreur de décodage de l'API :", errUnmarshall2)
			os.Exit(0)
		}

	}

}

func LoadRelation(w http.ResponseWriter, r *http.Request, id string, infos_artist structure.Artist) structure.Relations {

	var json_Relation structure.Relations

	url_Relations := "https://groupietrackers.herokuapp.com/api/relation/" + id

	response_Relations, err := http.Get(url_Relations)
	if err != nil {
		fmt.Println("Erreur de connexion à l'API :", err)
		os.Exit(0)
	}

	defer response_Relations.Body.Close()

	body_Relations, err := io.ReadAll(response_Relations.Body)
	if err != nil {
		fmt.Println("Erreur de lecture de l'API :", err)
		os.Exit(0)
	}

	//fmt.Println(body_Relations)

	errUnmarshall3 := json.Unmarshal(body_Relations, &json_Relation)
	if errUnmarshall3 != nil {
		fmt.Println("Erreur de décodage de l'API :", errUnmarshall3)
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
	//fmt.Println(data)
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
			fmt.Println("Erreur de connexion à l'API :", err)
			os.Exit(0)
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Ereur de lecture de l'API :", err)
			os.Exit(0)
		}

		var data []structure.Place
		errUnmarshall4 := json.Unmarshal(body, &data)
		if errUnmarshall4 != nil {
			fmt.Println("Erreur de décodage de l'API :", errUnmarshall4)
			os.Exit(0)
		}

		var inter []string
		inter = append(inter, data[0].Lat)
		inter = append(inter, data[0].Lon)
		res[city] = inter
		inter = nil

	}
	return res
}
