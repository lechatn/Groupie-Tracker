// Define the package Server
package Server

//Import the used packages
import (
	"Groupie/structure"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func LoadArtistes(w http.ResponseWriter, r *http.Request) []structure.Artist {
	//Function to load the page "artists" from the API
	url_Artists := "https://groupietrackers.herokuapp.com/api/artists" // We define the good URL

	var jsonList_Artists []structure.Artist
	response_Artists, err := http.Get(url_Artists)
	// Management of the error
	if err != nil {
		fmt.Println("Erreur de connexion a l'API 'artists' :", err)
		os.Exit(0)
	}

	defer response_Artists.Body.Close() // We wait for the end of the function to close the body

	body_Artists, err := io.ReadAll(response_Artists.Body) // We read the body of the API
	// Management of the error
	if err != nil {
		fmt.Println("Erreur de lecture de l'API 'artists' :", err)
		os.Exit(0)
	}
	errUnmarshall1 := json.Unmarshal(body_Artists, &jsonList_Artists) // We decode the result with function Unmarshall
	// Management of the error
	if errUnmarshall1 != nil {
		fmt.Println("Erreur de décodage de l'API 'artists' :", errUnmarshall1)
		os.Exit(0)
	}
	return jsonList_Artists

}

func LoadRelation(w http.ResponseWriter, r *http.Request, id string, infos_artist structure.Artist) structure.Relations {
	// Function to load the page "relation" from the API
	var json_Relation structure.Relations

	url_Relations := "https://groupietrackers.herokuapp.com/api/relation/" + id // We define the good URL with the id selected by the client

	response_Relations, err := http.Get(url_Relations)
	// Management of the error
	if err != nil {
		fmt.Println("Erreur de connexion à l'API 'relation' :", err)
		os.Exit(0)
	}

	defer response_Relations.Body.Close()

	body_Relations, err := io.ReadAll(response_Relations.Body)
	// Management of the error
	if err != nil {
		fmt.Println("Erreur de lecture de l'API 'relation' :", err)
		os.Exit(0)
	}

	//fmt.Println(body_Relations)

	errUnmarshall3 := json.Unmarshal(body_Relations, &json_Relation)
	// Management of the error
	if errUnmarshall3 != nil {
		fmt.Println("Erreur de décodage de l'API 'relation' :", errUnmarshall3)
		os.Exit(0)
	}

	data := structure.Relations{} // We create a new structure to stock the data
	data.Id = json_Relation.Id
	data.DateLocation = json_Relation.DateLocation
	data.Infos = infos_artist
	data.Coordonnees = SearchLatLon(data.DateLocation) // We call the function SearchLatLon to get the latitude and longitude of all the cities of DateLocation

	tRelation := template.Must(template.ParseFiles("./templates/relation.html")) // Read the relation page
	tRelation.Execute(w, data)
	return data
}

func SearchLatLon(relation map[string][]string) map[string][]string {
	res := make(map[string][]string, len(relation)) // We create a map to stock the latitude and longitude for the javascript file
	url := ""
	for city := range relation {
		city = strings.ReplaceAll(city, "-", ",")      // We replace the "-" by "," to avoid the error of the API nominatim
		if city == "willemstad,netherlands_antilles" { // We have to make a special case for the city "willemstad" because the API doesn't recognize the city "willemstad,netherlands_antilles"
			url = "https://nominatim.openstreetmap.org/search?q=willemstad&format=json"
		} else {
			url = "https://nominatim.openstreetmap.org/search?q=" + city + "&format=json" // Else, we define the good URL by using the city
		}

		response, err := http.Get(url)
		// Management of the error
		if err != nil {
			fmt.Println("Erreur de connexion à l'API nominatim :", err)
			os.Exit(0)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		// Management of the error
		if err != nil {
			fmt.Println("Ereur de lecture de l'API 'nominatim' :", err)
			os.Exit(0)
		}

		var data []structure.Place
		errUnmarshall4 := json.Unmarshal(body, &data)
		// Management of the error
		if errUnmarshall4 != nil {
			fmt.Println("Erreur de décodage de l'API 'nominatim' :", errUnmarshall4)
			os.Exit(0)
		}

		var inter []string
		inter = append(inter, data[0].Lat) // We stock the latitude and longitude in an intermediate slice
		inter = append(inter, data[0].Lon)
		res[city] = inter // We stock the intermediate slice in the result map
		inter = nil // We reset the intermediate slice

	}
	return res
}
