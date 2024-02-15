package Server

import (
	"net/http"
	"sort"
	"strings"
	"Groupie/structure"
	
)

func SearchArtist(w http.ResponseWriter, r *http.Request, jsonList_Artists []structure.Artist, originalData []structure.Artist, lettre string) []structure.Artist {
	var new_data []structure.Artist
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

func SortData(w http.ResponseWriter, r *http.Request, jsonList_Artists []structure.Artist) []structure.Artist {
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