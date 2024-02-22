// Define the package structure
package structure

// Define the structs

// Define DataLocation struct
type DataLocation struct {
	Locations []string `json:"locations"`
	Id        int      `json:"id"`
	Dates     string   `json:"dates"`
}

// Define Dates struct
type Dates struct {
	Index   []string `json:"index"`
	IdDates int      `json:"id"`
	Dates   []string `json:"dates"`
}

// Define Locations struct
type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
}

// Define Relations struct
type Relations struct {
	Id           int                 `json:"id"`
	DateLocation map[string][]string `json:"datesLocations"`
	Coordonnees  map[string][]string
	Infos        Artist
}

// Define Place struct
type Place struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// Define Artist struct
type Artist struct {
	IdArtists    int      `json:"id"`
	Images       string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}