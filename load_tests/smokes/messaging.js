import { sleep } from 'k6';
import http from 'k6/http';
import { NOSEBOOK_URL } from '../const';
import * as chats from '../lib/chats';
import * as random from '../lib/random';
import { connectWebSocket } from '../lib/connect-websocket';
import { faker } from '@faker-js/faker';
import { SharedArray } from 'k6/data';
import json from './data.json';
import { setTimeout } from 'k6/timers';

const duration = 1 * 60 * 1000;
export const options = {
  vus: 4, // Key for Smoke test. Keep it at 2, 3, max 5 VUs
  iterations: 4,
  duration: duration / 1000 + 's',
};

const shared = new SharedArray("data", () => {
  return json;
});

export default () => {
  const data = shared[__VU - 1];

  const nick = `virtual_user_${__VU - 1}`
  const loginJson = http.post(`${NOSEBOOK_URL}/login`, JSON.stringify({
    nick,
    password: '123123123',
  })).json();

  const sessionId = loginJson.data.session.sessionId;
  const userId = loginJson.data.user.id;
  const shouldMessageFirst = data.first;
  const interlocutorId = data.interlocutorId;

  const ws = connectWebSocket({ sessionId });

  ws.onopen = () => {
    console.log(`VU ${userId} connects websocket.`);
    setTimeout(() => {
       ws.close();
    }, duration);

    sleep(random.intBetween(2, 10));

    if (shouldMessageFirst) {
      console.log(`VU ${userId} is sending message to ${interlocutorId} first.`);
      chats.sendMessage({
        interlocutorId,
        text: faker.lorem.sentences(),
      }, { sessionId });
    }
  };

  ws.onmessage = event => {
    const message = JSON.parse(event.data);

    if (message.type === 'new_message' && message.payload.author.id === interlocutorId) {
      console.log(`VU ${userId} received message.`);
      sleep(random.intBetween(1, 10));

      console.log(`VU ${userId} is sending message.`);
      chats.sendMessage({
        interlocutorId,
        text: faker.lorem.sentences(),
      }, { sessionId });
    }
  };

  ws.onclose = () => {
    console.log(`VU ${userId} disconnects websocket.`);
  };
};

