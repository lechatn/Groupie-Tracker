var map = L.map('map').setView([40, -0.09], 13);
console.log('test');

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);



////////////////////////////////////////////////////////////////////////////////
//console.log('Value from Go:', location);
    //var ville = data.vill
    //fetch('https://nominatim.openstreetmap.org/search?q=' + ville + '&format=json')
        //.then(function(response) {
        //    return response.json();
       // })
        //.then(function(data) {
        //    var lat = data[0].lat;
        //    var lon = data[0].lon;

       //     L.marker([lat, lon]).addTo(map)
       //     map.setView([lat, lon], 9)
      //  });
////////////////////////////////////////////////////////////////////////////////