var map = L.map('map').setView([40, -0.09], 1);



L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);


fetch('http://localhost:8768/relationForJs')
    .then(response => response.json())
    .then(relation => {
        console.log(relation["Coordonnees"]);
        for (let city in relation["Coordonnees"]) {
            var lat = relation["Coordonnees"][city][0];
            var lon = relation["Coordonnees"][city][1];
            var marker = L.marker([lat, lon]).addTo(map);
            text = "En concert à " + noHypen(city);
            marker.bindPopup(text);
        }
        
        
    }
    );

function noHypen (place) {
    place = place.replace("-", " ");
    place = place.replace("_", " ");
    return place
}
