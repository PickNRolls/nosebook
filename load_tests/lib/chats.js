import http from "k6/http";
import { sleep } from "k6";
import { AUTH_SESSION_HEADER, NOSEBOOK_URL } from "../const";
import * as random from './random';
import { connectWebSocket } from './connect-websocket';
import { faker } from "@faker-js/faker";
import execution from "k6/execution";

export const sendMessage = (opts, auth) => {
  return http.post(`${NOSEBOOK_URL}/conversations/send-message`, JSON.stringify({
    recipientId: opts.interlocutorId,
    text: opts.text,
  }), {
    headers: {
      [AUTH_SESSION_HEADER]: auth.sessionId,
    }
  });
};

const RESEND_INTERVAL = 1000 * 20;

export const runMessagingScenario = (opts, auth) => {
  const { sessionId } = auth;
  const { duration, logging = true, userId, interlocutorId, shouldMessageFirst } = opts;

  const ws = connectWebSocket({ sessionId });

  let timerId = null;
  if (shouldMessageFirst) {
    timerId = setInterval(() => {
      logging && console.log(`VU ${__VU} userId:${userId} is resending message, because he did not get answer.`);
      sendMessage({
        interlocutorId,
        text: faker.lorem.sentences(),
      }, { sessionId });
    }, RESEND_INTERVAL);
  }

  ws.onopen = () => {
    logging && console.log(`VU ${__VU} userId:${userId} connects websocket.`);
    setTimeout(() => {
      ws.close();
    }, duration);

    sleep(random.intBetween(2, 10));

    if (shouldMessageFirst) {
      logging && console.log(`VU ${__VU} userId:${userId} is sending message to ${interlocutorId} first.`);
      sendMessage({
        interlocutorId,
        text: faker.lorem.sentences(),
      }, { sessionId });
    }
  };

  ws.onmessage = event => {
    try {
      const message = JSON.parse(event.data);

      if (message.type === 'new_message' && message.payload.author.id === interlocutorId) {
        if (timerId != null) {
          clearInterval(timerId);
          timerId = null;
        }

        logging && console.log(`VU ${__VU} userId:${userId} received message.`);
        sleep(random.intBetween(1, 10));

        logging && console.log(`VU ${__VU} userId:${userId} is sending message.`);
        sendMessage({
          interlocutorId,
          text: faker.lorem.sentences(),
        }, { sessionId });
      }
    } catch (e) {
      console.log('Failed to handle message');
      console.log(event.data);
      execution.test.abort();
    }
  };

  ws.onclose = () => {
    logging && console.log(`VU ${__VU} userId:${userId} disconnects websocket.`);
  };
};

