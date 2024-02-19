var map = L.map('map').setView([40, -0.09], 1);
var markers = [];
var markerClusters = L.markerClusterGroup();

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);


fetch('http://localhost:8768/relationForJs')
    .then(response => response.json())
    .then(relation => {
        console.log(relation);
        var myIcon = L.icon({
            iconUrl: relation["Infos"]["image"],
            iconSize: [30, 30],
            iconAnchor: [12, 41],
            popupAnchor: [1, -34],
            shadowSize: [41, 41]
        });
        for (let city in relation["Coordonnees"]) {
            var lat = relation["Coordonnees"][city][0];
            var lon = relation["Coordonnees"][city][1];
            var marker = L.marker([lat, lon] , {icon: myIcon});
            cityForLocation = withHypen(city);
            text = "En concert à " + noHypen(city) + " le " + relation["datesLocations"][cityForLocation];
            marker.bindPopup(text);

            markerClusters.addLayer(marker); // Nous ajoutons le marqueur aux groupes
            markers.push(marker); // Nous ajoutons le marqueur à la liste des marqueurs 
            
        }
        var group = new L.featureGroup(markers); // Nous créons le groupe des marqueurs pour adapter le zoom
        map.fitBounds(group.getBounds().pad(0.5)); // Nous demandons à ce que tous les marqueurs soient visibles, et ajoutons un padding (pad(0.5)) pour que les marqueurs ne soient pas coupés
        map.addLayer(markerClusters);

    }
    );

function noHypen (place) {
    place = place.replace("-", " ");
    place = place.replace("_", " ");
    place = place.replace(",", " ");
    return place
}

function withHypen (place) {
    return place.replace(",", "-");
}