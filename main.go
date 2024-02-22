// Define the package main 
package main

// Import the used packages
import (
	"Groupie/Server"
	"Groupie/structure"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
)


var jsonList_Artists []structure.Artist

var port = ":8768"

var id string

var artist_create = false

var originalData []structure.Artist

var data_artist structure.Relations

var infos_artist structure.Artist


func main() {
	css := http.FileServer(http.Dir("style"))                
	http.Handle("/style/", http.StripPrefix("/style/", css)) 
	img := http.FileServer(http.Dir("images"))              
	http.Handle("/images/", http.StripPrefix("/images/", img))
	js := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Define the default path
		tHome := template.Must(template.ParseFiles("./templates/home.html"))
		tHome.Execute(w, nil)
	})

	http.HandleFunc("/artistes", func(w http.ResponseWriter, r *http.Request) { // Define the path for the artists page
		if !artist_create { // If the artist list is not created, we create it. If it's already created, we don't create it again to avoid to overload the API
			jsonList_Artists = Server.LoadArtistes(w, r)
			artist_create = true
			originalData = jsonList_Artists
		}
		tArtistes := template.Must(template.ParseFiles("./templates/artistes.html"))
		if r.FormValue("Check") != "" { // If the client type a letter in the search bar, we call the function SearchArtist to filter the data
			lettre := r.FormValue("Check")
			jsonList_Artists = Server.SearchArtist(w, r, jsonList_Artists, originalData, lettre)
			lettre = ""
			if len(jsonList_Artists) == 0 {
				terror := template.Must(template.ParseFiles("./templates/error.html"))
				terror.Execute(w,nil)
				jsonList_Artists = originalData
				return
			}
		}
		if r.FormValue("Search_artist") != "" { // If the client type a letter in the search bar, we call the function SearchArtist to filter the data
			lettre := r.FormValue("Search_artist")
			jsonList_Artists = Server.SearchArtist(w, r, jsonList_Artists, originalData, lettre)
			lettre = ""
			if len(jsonList_Artists) == 0 {
				terror := template.Must(template.ParseFiles("./templates/error.html"))
				terror.Execute(w, nil)
				jsonList_Artists = originalData
				return
			}
		}
		jsonList_Artists = Server.SortData(w, r, jsonList_Artists)
		tArtistes.Execute(w, jsonList_Artists)
	})
	http.HandleFunc("/loading", func(w http.ResponseWriter, r *http.Request) { // Define the path for the loading page
		tloading := template.Must(template.ParseFiles("./templates/loading.html"))
		id = r.URL.Query().Get("id") // We get the id selected by the client
		redirect := false 
		go func() { // We use a goroutine to load the data and start the loading page in the same time
			id_int, _ := strconv.Atoi(id)
			sort.Slice(originalData, func(i, j int) bool {
				return originalData[i].IdArtists < originalData[j].IdArtists
			})
			if len(originalData) == len(jsonList_Artists) {
				infos_artist = jsonList_Artists[id_int-1]
			} else {
				infos_artist = originalData[id_int-1]
			}
			data_artist = Server.LoadRelation(w, r, id, infos_artist)
			redirect = true
		}()
		if !redirect {
			tloading.Execute(w, id)
		}
	})

	http.HandleFunc("/relation", func(w http.ResponseWriter, r *http.Request) { // Define the path for the relation page
		Server.LoadRelation(w, r, id, infos_artist)
	})

	http.HandleFunc("/relationForJs", func(w http.ResponseWriter, r *http.Request) { // Define the path for the recup of the data by the javascript file
		json.NewEncoder(w).Encode(data_artist) // We encode the data to send it to the javascript file
		data_artist = structure.Relations{} // We reset the data to aavoid some problems
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) 

	fmt.Println("http://localhost:8768") // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)

}
