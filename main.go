package main

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

// Define all the struct and some variables

var jsonList_Artists []structure.Artist

var port = ":8768"

var artist_create = false
var originalData []structure.Artist

var location []structure.Locations

var data_artist structure.Relations

var infos_artist structure.Artist

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
		if !artist_create {
			jsonList_Artists = Server.LoadArtistes(w, r)
			artist_create = true
			originalData = jsonList_Artists
		}
		tArtistes := template.Must(template.ParseFiles("./templates/artistes.html")) // Read the artists page
		if r.FormValue("Check") != "" {
			lettre := r.FormValue("Check")
			jsonList_Artists = Server.SearchArtist(w, r, jsonList_Artists, originalData, lettre)
			lettre = ""
		}
		if r.FormValue("Search_artist") != "" {
			//fmt.Println("test")
			lettre := r.FormValue("Search_artist")
			jsonList_Artists = Server.SearchArtist(w, r, jsonList_Artists, originalData, lettre)
			lettre = ""
		}
		jsonList_Artists = Server.SortData(w, r, jsonList_Artists)
		tArtistes.Execute(w, jsonList_Artists)
	})

	http.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		Server.LoadLocation(w, r, Server.LoadArtistes(w, r))
	})

	http.HandleFunc("/relation", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		//fmt.Println(id)
		id_int, _ := strconv.Atoi(id)
		sort.Slice(originalData, func(i, j int) bool {
			return originalData[i].IdArtists < originalData[j].IdArtists
		})
		if len(originalData) == len(jsonList_Artists) {
			//fmt.Println("cas1")
			infos_artist = jsonList_Artists[id_int-1]
		} else {
			//fmt.Println("cas2")
			infos_artist = originalData[id_int-1]
		}
		data_artist = Server.LoadRelation(w, r, id, infos_artist)
	})

	http.HandleFunc("/relationForJs", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(data_artist)
		data_artist = structure.Relations{}
	})

	http.HandleFunc("/locationForJs", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(location)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("http://localhost:8768") // Creat clickable link in the terminal
	http.ListenAndServe(port, nil)

}
