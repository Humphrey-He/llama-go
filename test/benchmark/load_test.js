import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '1m30s', target: 50 },
    { duration: '20s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.1'],
  },
};

export default function () {
  const payload = JSON.stringify({
    model: 'tinyllama-chat',
    messages: [
      { role: 'user', content: 'Hello, how are you?' }
    ],
    temperature: 0.7,
    max_tokens: 256,
    stream: false,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer test-key',
    },
  };

  const res = http.post('http://localhost:8080/v1/chat/completions', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
