import supertest, { Test } from 'supertest';
import TestAgent from 'supertest/lib/agent';

export const HOST = 'http://backend:8080';
export const SESSION_HEADER_KEY = 'X-Auth-Session-Id';
export const ASSER_ID = 'ed1a3fd0-4d0b-4961-b4cd-cf2123577666';
export const ASSER_SESSION = 'bb23af03-be50-4bce-b729-b259b2e02e56';
export const TEST_TESTER_SESSION = 'bb23af03-be50-4bce-b729-b259b2e02e54'
export const DRUGTESTER_ID = '1ae02f69-ea1a-4308-b825-0e5896e652e4';

class TestAgentWrapper {
  private host: string;
  private agent: TestAgent;
  private auth: boolean;

  public constructor(host: string, auth?: boolean) {
    this.host = host;
    this.agent = supertest(host);
    this.auth = auth ?? true;
  }

  public unwrap(): TestAgent {
    return this.agent;
  }

  public extend(opts?: {
    prefixUrl?: string;
    auth?: boolean;
  }): TestAgentWrapper {
    let url = this.host;
    if (opts?.prefixUrl) {
      url += opts.prefixUrl;
    }
    return new TestAgentWrapper(url, opts?.auth ?? this.auth);
  }

  private hook<M extends 'get' | 'post'>(method: M, url: string): Test {
    const test = this.agent[method](url);

    if (this.auth) {
      test.set(SESSION_HEADER_KEY, TEST_TESTER_SESSION);
    }

    return test;
  }

  public get(url: string): Test {
    return this.hook('get', url);
  }

  public post(url: string): Test {
    return this.hook('post', url);
  }
};

export default new TestAgentWrapper(HOST);


