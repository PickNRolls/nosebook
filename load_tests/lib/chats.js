import http from "k6/http";
import { AUTH_SESSION_HEADER, HOST } from "../const";

export const sendMessage = (opts, auth) => {
  return http.post(`${HOST}/conversations/send-message`, JSON.stringify({
    recipientId: opts.interlocutorId,
    text: opts.text,
  }), {
    headers: {
      [AUTH_SESSION_HEADER]: auth.sessionId,
    }
  });
};

