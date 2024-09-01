import supertest, { Test } from 'supertest';
import TestAgent from 'supertest/lib/agent';

const HOST = 'http://backend:8080';
const SESSION_HEADER_KEY = 'X-Auth-Session-Id';

class TestAgentWrapper {
  private host: string;
  private agent: TestAgent;

  public constructor(host: string) {
    this.host = host;
    this.agent = supertest(host);
  }

  public unwrap(): TestAgent {
    return this.agent;
  }

  public extend(prefixUrl: string): TestAgentWrapper {
    return new TestAgentWrapper(this.host + prefixUrl);
  }

  private hook<M extends 'get' | 'post'>(method: M, url: string): Test {
    const test = this.agent[method](url);

    test.set(SESSION_HEADER_KEY, 'bb23af03-be50-4bce-b729-b259b2e02e54');

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


