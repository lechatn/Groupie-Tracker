var map = L.map('mapid');

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);

////////////////////////////////////////////////////////////////////////////////

fetch('http://localhost:8768/JavaScript')
    .then(response => response.json())
    .then(data => {
        console.log('Value from Go:', data.vill);
        var ville = data.vill
        fetch('https://nominatim.openstreetmap.org/search?q=' + ville + '&format=json')
            .then(function(response) {
                return response.json();
            })
            .then(function(data) {
                var lat = data[0].lat;
                var lon = data[0].lon;

                L.marker([lat, lon]).addTo(map)
                map.setView([lat, lon], 9)
            });

    })


////////////////////////////////////////////////////////////////////////////////