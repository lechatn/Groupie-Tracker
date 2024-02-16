var map = L.map('map').setView([40, -0.09], 1);



L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);


fetch('http://localhost:8768/locationForJs')
    .then(response => response.json())
    .then(location => {
        console.log(location);
        for (let j = 0; j<location.length; j++) {
            ville.push(location[j].locations);
        }
        for (let i = 0 ; i  <ville.length; i++) {
            for (let z = 0; z < ville[i].length; z++) {
                fetch('https://nominatim.openstreetmap.org/search?q=' +ville[i][z]+ '&format=json')
                    .then(response => response.json())
                    .then(data => {
                        var lat = data[0].lat;
                        var lon = data[0].lon;
                        var marker = L.marker([lat, lon]).addTo(map);
                        text = "En concert à " + ville[i][z];
                        marker.bindPopup(text);
                    });
            }
        }
    });