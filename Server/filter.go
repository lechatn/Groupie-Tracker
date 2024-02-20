//Define the package Server
package Server
//Import the used packages
import (
	"Groupie/structure"
	"net/http"
	"sort"
	"strings"
)

func SearchArtist(w http.ResponseWriter, r *http.Request, jsonList_Artists []structure.Artist, originalData []structure.Artist, lettre string) []structure.Artist {
	//Function to search an artist by the firsts letters of his group name
	var new_data []structure.Artist
	if lettre == "" {
		return jsonList_Artists
	}
	if len(originalData) != len(jsonList_Artists) { // If the user has already made a search, we reset the data
		jsonList_Artists = originalData
	}
	if strings.ToUpper(lettre) == "ALL" { // If the user wants to see all the artists, he can type "all" in the search bar or click on the "view all artists" button	
		return originalData
	}
	for i := 0; i < len(jsonList_Artists); i++ {
		for j := 0; j < len(lettre); j++ {
			if strings.ToUpper(string(jsonList_Artists[i].Name[j])) == strings.ToUpper(string(lettre[j])) { // We compare the firts letter of the group name with the first letter of the client search
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
	return new_data // We return the new data with the artists that match with the client search
}

func SortData(w http.ResponseWriter, r *http.Request, jsonList_Artists []structure.Artist) []structure.Artist {
	//Function to sort the data with 6 different possibilities
	order1 := r.FormValue("alpha") // Sort in alphabetical order
	order2 := r.FormValue("unalpha") // Sort in reverse alphabetical order
	order3 := r.FormValue("firstalbum") // Sort by the date of the first album
	order4 := r.FormValue("CreationDate") // Sort by the creation date of the group
	order5 := r.FormValue("mostArtists") // Sort by the greatest number of members
	order6 := r.FormValue("lessArtists") // Sort by the smallest number of members
	if order1 == "" && order2 == "" && order3 == "" && order4 == "" && order5 == "" && order6 == ""{ // If the client doesn't want to sort the data, we return the data as it is
		return jsonList_Artists
	}
	if order1 != "" { // If the client choose the firts option, we sort the data in alphabetical order
		sort.Slice(jsonList_Artists, func(i, j int) bool { // We use the sort.Slice function to sort the data
			return strings.ToUpper(jsonList_Artists[i].Name) < strings.ToUpper(jsonList_Artists[j].Name)
		})
	} else if order2 != "" { // If the client choose the second option, we sort the data in reverse alphabetical order
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			return strings.ToUpper(jsonList_Artists[i].Name) > strings.ToUpper(jsonList_Artists[j].Name)
		})
	} else if order3 != "" { // If the client choose the third option, we sort the data by the date of the first album
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			if jsonList_Artists[i].FirstAlbum[6:] == jsonList_Artists[j].FirstAlbum[6:] { // Firstly, we compare the year of the first album
				if jsonList_Artists[i].FirstAlbum[3:5] == jsonList_Artists[j].FirstAlbum[3:5] { // Secondly, we compare the month of the first album
					return jsonList_Artists[i].FirstAlbum[:2] < jsonList_Artists[j].FirstAlbum[:2] // Finally, we compare the day of the first album
				}
				return jsonList_Artists[i].FirstAlbum[3:5] < jsonList_Artists[j].FirstAlbum[3:5]
			}
			return jsonList_Artists[i].FirstAlbum[6:] < jsonList_Artists[j].FirstAlbum[6:]
		})
	} else if order4 != "" { // If the client choose the fourth option, we sort the data by the creation date of the group
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			return jsonList_Artists[i].CreationDate < jsonList_Artists[j].CreationDate
		})
	} else if order5 != "" { // If the client choose the fifth option, we sort the data by the greatest number of members
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			return len(jsonList_Artists[i].Members) > len(jsonList_Artists[j].Members)
		})
	} else if order6 != "" { // If the client choose the sixth option, we sort the data by the smallest number of members
		sort.Slice(jsonList_Artists, func(i, j int) bool {
			return len(jsonList_Artists[i].Members) < len(jsonList_Artists[j].Members)
		})
	}
	return jsonList_Artists
}
