import http from 'k6/http';
import { check, sleep } from 'k6';
import { faker } from '@faker-js/faker';
import { NOSEBOOK_URL } from '../const';

export const options = {
  vus: 3, // Key for Smoke test. Keep it at 2, 3, max 5 VUs
  duration: '1m'
};

export default () => {
  const res = http.post(`${NOSEBOOK_URL}/register`, JSON.stringify({
    firstName: faker.person.firstName(),
    lastName: faker.person.lastName(),
    nick: faker.internet.userName() + __VU + __ITER,
    password: faker.internet.password({
      length: 12,
    }),
  }), {
    headers: 'application/json',
  });
  
  sleep(1);
  
  check(res, {
    'Status': () => res.status === 200,
  });
};

