var map = L.map('map').setView([40, -0.09], 1); // Iniit the map
var markers = []; // Define the slice for the markers
var markerClusters = L.markerClusterGroup(); // Define the group of markers

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map); 


fetch('http://localhost:8768/relationForJs') // We take the data from the server
    .then(response => response.json()) // We transform the data into json
    .then(relation => {
        var myIcon = L.icon({ // We define our icon. The icon corresponds to the image of the artist in the API
            iconUrl: relation["Infos"]["image"],
            iconSize: [30, 30],
            iconAnchor: [12, 41],
            popupAnchor: [1, -34],
            shadowSize: [41, 41]
        });
        for (let city in relation["Coordonnees"]) {
            var lat = relation["Coordonnees"][city][0]; // We take the latitude and longitude of the city
            var lon = relation["Coordonnees"][city][1];
            var marker = L.marker([lat, lon] , {icon: myIcon});
            cityForLocation = withHypen(city); // We replace the spaces by hypens for find the date of the concert in  datesLocations
            text = "En concert à " + noHypen(city) + " le " + relation["datesLocations"][cityForLocation]; // We define the message in the popup
            marker.bindPopup(text);

            markerClusters.addLayer(marker); // Add the marker to the group
            markers.push(marker);
            
        }
        var group = new L.featureGroup(markers);
        map.fitBounds(group.getBounds().pad(0.5));
        map.addLayer(markerClusters);
    }
    );

function noHypen (place) { // Function to replace the hypens by spaces
    place = place.replace("-", " ");
    place = place.replace("_", " ");
    place = place.replace(",", " ");
    return place
}

function withHypen (place) { // Function to replace the spaces by hypens
    return place.replace(",", "-");
}