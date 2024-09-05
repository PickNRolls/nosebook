import * as ws from 'ws';
import { HOST, SESSION_HEADER_KEY, TEST_TESTER_SESSION } from './request';

export class WebSocket {
  private original: ws.WebSocket;

  public constructor(session?: string) {
    this.original = new ws.WebSocket(`${HOST}/ws`, {
      headers: {
        [SESSION_HEADER_KEY]: session ?? TEST_TESTER_SESSION,
      }
    });
  }

  public unwrap() {
    return this.original;
  }
}

