// WSL
import http from 'k6/http';

Date.prototype.addDays = function(days) {
    var date = new Date(this.valueOf());
    date.setDate(date.getDate() + days);
    return date;
}


export default function () {
    const startDate = new Date();
    startDate.setFullYear(2013, 10, 10);

    const endDate = new Date();
    startDate.setFullYear(2014, 10, 10);

    const baseUrl = 'http://localhost:5000/hotels/';

    for (let i = 0; i < 101; i++) {

        const indate = new Date(startDate);
        const outdate = new Date(endDate);

        // Modify indate and outdate by 1 day
        indate.addDays(i);
        outdate.addDays(i);

        const url = `${baseUrl}?indate=${indate.toISOString()}&outdate=${outdate.toISOString()}`;

        http.get(url);
    }
}




// p(95) i p(99.9)
//k6 run --iterations=100 --vus=10 --summary-trend-stats="med,p(95),p(99.9)" script.js
//