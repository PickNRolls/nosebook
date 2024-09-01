import supertest, { Test } from 'supertest';
import TestAgent from 'supertest/lib/agent';

const HOST = 'http://backend:8080';
export const SESSION_HEADER_KEY = 'X-Auth-Session-Id';

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
      test.set(SESSION_HEADER_KEY, 'bb23af03-be50-4bce-b729-b259b2e02e54');
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


