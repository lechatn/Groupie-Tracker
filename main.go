package main



import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

var port = ":8768"

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

type DatesAndArtists struct {
	Artist Artist
	Dates   Dates
	Locations Locations

}

var jsonList_Artists []Artist
var homeData map[string]interface{}
var jsonList_Location []Locations
var jsonList_Dates []Dates
var homeDates map[string][]Dates
var allLocation map[string][]Locations

// ///////////////////////////////////////////

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

	errUnmarshall := json.Unmarshal(body_General, &homeData)
	if errUnmarshall != nil {
		fmt.Println("Error3")
		return
	}

	////////////////////////////////////////////////////////////////////////

	url_Artists := homeData["artists"].(string)
	url_Locations := homeData["locations"].(string)
	url_Dates := homeData["dates"].(string)
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
	//fmt.Println(jsonList_Dates)

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

	listArtists := jsonList_Artists
	listDates := jsonList_Dates
	listLocations := jsonList_Location
	var Data []DatesAndArtists

	for i:=0; i<len(listArtists); i++{
		for j:=0; j<len(listDates); j++{
			if listArtists[i].IdArtists == listDates[j].IdDates{
				if listLocations[j].Id == listDates[j].IdDates{
					var inter DatesAndArtists
					inter.Artist = listArtists[i]
					inter.Dates = listDates[j]
					inter.Locations = listLocations[j]
					Data = append(Data, inter)
					//fmt.Println(inter)
				}	
			}
			continue
		}
	}


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tHome := template.Must(template.ParseFiles("./templates/home.html"))
		//fmt.Println(Data)
		tHome.Execute(w, Data)
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

	fmt.Println("http://localhost" + port) // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)

}
