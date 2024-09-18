import { AUTH_SESSION_HEADER, HOST } from "../const";
import ws from "k6/ws";

export const connectWebSocket = (auth, callback) => {
  return ws.connect(`ws://localhost:8080/ws`, {
    headers: {
      [AUTH_SESSION_HEADER]: auth.sessionId,
    }
  }, callback);
};

