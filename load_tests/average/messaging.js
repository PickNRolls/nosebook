import http from 'k6/http';
import { AUTH_SESSION_HEADER, NOSEBOOK_URL } from '../const';
import * as chats from '../lib/chats';
import { rampUpWs, rampUpWsOptions } from '../lib/ws';

const stages = [
  { duration: 1000 * 60 * 3, target: 1000 },
  { duration: 1000 * 60 * 3, target: 1000 },
  { duration: 1000 * 60 * 3, target: 0 },
];

export const options = rampUpWsOptions(stages);

let loginJson = null;
let usersJson = null;

export const teardown = () => {
  if (loginJson == null) {
    return;
  }

  http.post(`${NOSEBOOK_URL}/logout`, null, {
    headers: {
      [AUTH_SESSION_HEADER]: loginJson.data.session.sessionId,
    }
  });
};

export default () => {
  const { vuIndex, duration } = rampUpWs(stages);
  
  const offset = 300;
  const userIndex = offset + vuIndex;
  const nick = `virtual_user_${userIndex}`
  if (loginJson == null) {
    loginJson = http.post(`${NOSEBOOK_URL}/login`, JSON.stringify({
      nick,
      password: '123123123',
    })).json();
  }

  const sessionId = loginJson.data.session.sessionId;
  const userId = loginJson.data.user.id;
  const shouldMessageFirst = userIndex % 2 === 1;
  const interlocutorNick = shouldMessageFirst ? `virtual_user_${userIndex - 1}` : `virtual_user_${userIndex + 1}`;

  if (usersJson == null) {
    usersJson = http.get(`${NOSEBOOK_URL}/users?text=${interlocutorNick}`, {
      headers: {
        [AUTH_SESSION_HEADER]: sessionId,
      },
    }).json();

    if (usersJson.data.totalCount !== 1) {
      throw new Error("Can't find 1-to-1 interlocutor by nickname");
    }
  }
  
  const interlocutor = usersJson.data.data[0];

  chats.runMessagingScenario({
    duration: duration,
    logging: false,
    userId,
    shouldMessageFirst,
    interlocutorId: interlocutor.id,
  }, { sessionId });
};

