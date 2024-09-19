import { sleep } from 'k6';
import { getCurrentStageIndex } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import http from 'k6/http';
import { AUTH_SESSION_HEADER, NOSEBOOK_URL } from '../const';
import * as chats from '../lib/chats';
import * as random from '../lib/random';

const stages = [
  { duration: 1000 * 60 * 1, target: 100 },
  { duration: 1000 * 60 * 5, target: 100 },
  { duration: 1000 * 60 * 1, target: 0 },
]
export const options = {
  stages: stages.map(stage => {
    return {
      ...stage,
      duration: `${stage.duration / 1000}s`,
    }
  }),
};

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
  const stageIndex = getCurrentStageIndex();
  const duration = stages[stageIndex].duration;

  const offset = 200;
  const userIndex = offset + __VU;
  const nick = `virtual_user_${userIndex}`
  // do not login on every stage, only first one
  if (loginJson == null) {
    sleep(random.intBetween(1, 4));
    
    loginJson = http.post(`${NOSEBOOK_URL}/login`, JSON.stringify({
      nick,
      password: '123123123',
    })).json();
  }

  const sessionId = loginJson.data.session.sessionId;
  const userId = loginJson.data.user.id;
  const shouldMessageFirst = userIndex % 2 === 0;
  const interlocutorNick = shouldMessageFirst ? `virtual_user_${userIndex - 1}` : `virtual_user_${userIndex + 1}`;

  if (usersJson == null) {
    sleep(random.intBetween(1, 4));

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
    duration,
    logging: false,
    userId,
    shouldMessageFirst,
    interlocutorId: interlocutor.id,
  }, { sessionId });
};

