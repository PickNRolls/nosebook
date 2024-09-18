import { check, sleep } from 'k6';
import http from 'k6/http';
import { HOST } from '../const';
import * as chats from '../lib/chats';
import { connectWebSocket } from '../lib/connect-websocket';
import { faker } from '@faker-js/faker';
import { SharedArray } from 'k6/data';
import json from './data.json';
import execution from 'k6/execution';

export const options = {
  vus: 4, // Key for Smoke test. Keep it at 2, 3, max 5 VUs
  duration: '1m'
};

const shared = new SharedArray("data", () => {
  return json;
});

export default () => {
  const data = shared[__VU - 1];
  
  const nick = `virtual_user_${__VU - 1}`
  const loginJson = http.post(`${HOST}/login`, JSON.stringify({
    nick,
    password: '123123123',
  })).json();
  
  const sessionId = loginJson.data.session.sessionId;
  const userId = loginJson.data.user.id;
  const shouldMessageFirst = data.first;
  const interlocutorId = data.interlocutorId;

  const res = connectWebSocket({ sessionId }, socket => {
    socket.on('open', () => {
      console.log(`VU ${userId} connects websocket.`);

      sleep(1);

      if (shouldMessageFirst) {
        console.log(`VU ${userId} is sending message to ${interlocutorId} first.`);
        chats.sendMessage({
          interlocutorId,
          text: faker.lorem.sentences(),
        }, { sessionId });
      }
    });

    socket.on('message', (msg) => {
      if (execution.scenario.progress === 1) {
        socket.close()
        return
      }
      
      const message = JSON.parse(msg);

      if (message.type === 'new_message' && message.payload.author.id === interlocutorId) {
        console.log(`VU ${userId} received message.`);
        sleep(4);

        console.log(`VU ${userId} is sending message.`);
        chats.sendMessage({
          interlocutorId,
          text: faker.lorem.sentences(),
        }, { sessionId });
      }
    });

    socket.on('close', () => {
      console.log(`VU ${userId} disconnects websocket.`);
    });
  });

  check(res, {
    'Status is 101': (r) => r && r.status === 101,
  });
};

