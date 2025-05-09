import http from 'k6/http';
import { check, sleep, group } from 'k6';
import exec from 'k6/execution';

// Options to configure the number of users and the execution time for each scenario
export const options = {
    scenarios: {
        'healthz': {
            executor: 'shared-iterations',
            vus: 1, // Number of virtual users
            iterations: 1, // total amount of calling this function
        },
        'shorten': {
            executor: 'shared-iterations',
            vus: 5, // Number of virtual users
            iterations: 15, // total amount of calling this function
        },
        'get': {
            executor: 'constant-vus',
            vus: 50,
            duration: '10s',
            startTime: '10s', // Delay the start of this scenario until the shorten scenario finishes
        },
    },
};

// This will hold the IDs generated by the shorten scenario
let generatedIds = [];

const MONITORING_URL = 'http://localhost:8001';
const REQUEST_URL = 'http://localhost:8002/api/v1';

// Function for shortening a URL and collecting IDs
function shortenUrl() {
    const shortenPayload = JSON.stringify({
        url: 'https://example.com',
    });

    const headers = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const shortenRes = http.post(`${REQUEST_URL}/shorten`, shortenPayload, headers);

    check(shortenRes, {
        'shorten status is 201': (r) => r.status === 201,
        'response has ID': (r) => !!r.json().id,
    });

    const id = shortenRes.json().id;
    generatedIds.push(id); // Collect the ID for later use
    sleep(1); // Simulate user delay
}

// Function for getting a shortened URL using an ID
function getShortenedUrl(id) {
    const getRes = http.get(`${REQUEST_URL}/shorten/${id}`);

    check(getRes, {
        'get status is 200': (r) => r.status === 200,
    });

    sleep(1); // Simulate user delay
}

// The default function will run for each virtual user (VU)
export default function () {
    const scenario = exec.scenario.name;

    if (scenario === 'healthz') {
        group("healthz", () => {
            let res = http.get(`${MONITORING_URL}/healthz/liveness`);
            check(res, {
                'liveness is status 200': (r) => r.status === 200,
            });

            res = http.get(`${MONITORING_URL}/healthz/readiness`);
            check(res, {
                'readiness is status 200': (r) => r.status === 200,
            });
        });
    } else if (scenario === 'shorten') {
        shortenUrl();  // Generate and collect IDs
    } else if (scenario === 'get' && generatedIds.length > 0) {
        let randomIndex = Math.floor(Math.random() * generatedIds.length);
        const randomID = generatedIds[randomIndex];
        getShortenedUrl(randomID);  // Fetch shortened URL
    }
}
