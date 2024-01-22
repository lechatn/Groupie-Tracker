package main

import (
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
)

var port = ":8080"

func main() {
    // css := http.FileServer(http.Dir("style"))                // For add css to the html pages
    // http.Handle("/style/", http.StripPrefix("/style/", css)) // For add css to the html pages
    url := "https://groupietrackers.herokuapp.com/api/artists"
    response, err := http.Get(url)
    if err != nil {
        fmt.Println("Error")
        return
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("Error")
        return
    }
    apiData := (body)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Lunch a new page for the lose condition
        tmpl := template.Must(template.ParseFiles("./templates/index.html")) // Read the home page
        data := struct {
			ApiData []string
		}{
			ApiData: []string{string(apiData)},
		}
        tmpl.Execute(w, data)
    })

    fmt.Println("http://localhost:8080") // Creat clickable link in the terminal
    http.ListenAndServe(port, nil)
}