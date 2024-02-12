var map = L.map('map').setView([40, -0.09], 13);



L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);


fetch('http://localhost:8768/relationForJs')
    .then(response => response.json())
    .then(relation => {
        var ville = Object.keys(relation)
        for (let i = 0 ; i<ville.length; i++) {
            fetch('https://nominatim.openstreetmap.org/search?q=' +ville[i]+ '&format=json')
                .then(response => response.json())
                .then(data => {
                    var lat = data[0].lat;
                    var lon = data[0].lon;
                    var marker = L.marker([lat, lon]).addTo(map);
                    text = "En concert à " + ville[i] + " le : " + relation[ville[i]];
                    marker.bindPopup(text);
                });
            }
    });