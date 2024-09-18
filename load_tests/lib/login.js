import http from "k6/http";
import { HOST } from "../const";

export const login = (nickname, password) => {
  return http.post(`${HOST}/login`, JSON.stringify({
    nick: nickname,
    password: password,
  }));
};

