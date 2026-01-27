import http from 'k6/http';
import { check } from 'k6';

export default function () {
    // Test against the local httpbin server
    let response = http.get('http://localhost:5454/get');

    check(response, {
        'status is 200': (r) => r.status === 200,
    });

    // Test POST endpoint
    response = http.post('http://localhost:5454/post', {
        key: 'value'
    });

    check(response, {
        'POST status is 200': (r) => r.status === 200,
    });
}
