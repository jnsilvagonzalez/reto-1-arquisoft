import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
    // A number specifying the number of VUs to run concurrently.
    vus: 1000,
    // A string specifying the total duration of the test run.
    duration: '60s',
};


export default function () {
    const payload = JSON.stringify({
        "RqUID": "1213",
        "IdVehicle": "VC-2345",
        "Speed": 60,
        "Address": "cRA 11",
        "Latitude": "23.23",
        "Longitude": "22.34"
    });
    const headers = { 'Content-Type': 'application/json' };
    http.post('http://alb-receptor-235373314.us-east-1.elb.amazonaws.com:80/api/signal', payload, { headers });
    sleep(1);
}
