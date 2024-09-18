import http from "k6/http";
import { AUTH_SESSION_HEADER, NOSEBOOK_HOST, NOSEBOOK_URL } from "../const";

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

