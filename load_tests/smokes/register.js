import http from 'k6/http';
import { check, sleep } from 'k6';
import { NOSEBOOK_URL } from '../const';

const vus = 10;
const iterationsPerVu = 200;
const alreadyExistingInDbSqlInitCount = 5;
export const options = {
  vus,
  iterations: vus * iterationsPerVu,
};

export const setup = () => {
  sleep(__VU);  
};

export default () => {
  const index = (alreadyExistingInDbSqlInitCount * 2 + ((__VU - 1) * iterationsPerVu) + __ITER);
  const res = http.post(`${NOSEBOOK_URL}/register`, JSON.stringify({
    firstName: 'Virtual',
    lastName: 'User ' + index,
    nick: 'virtual_user_' + index,
    password: '123123123',
  }));
  
  sleep(0.5);
  
  check(res, {
    'Status': () => res.status === 200,
  });
};

