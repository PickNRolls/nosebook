import { WebSocket } from "k6/experimental/websockets";
import { AUTH_SESSION_HEADER, NOSEBOOK_HOST } from "../const";

export const connectWebSocket = (auth) => {
  const ws = new WebSocket(`ws://${NOSEBOOK_HOST}/ws`, null, {
    headers: {
      [AUTH_SESSION_HEADER]: auth.sessionId,
    }
  });

  return ws;
};

