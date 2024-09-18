import http from "k6/http";
import { NOSEBOOK_URL } from "../const";

export const login = (nickname, password) => {
  return http.post(`${NOSEBOOK_URL}/login`, JSON.stringify({
    nick: nickname,
    password: password,
  }));
};

