import supertest from 'supertest';
import request, { SESSION_HEADER_KEY } from './request';

describe('auth', () => {
  let auth = request.extend({
    auth: false,
  });

  test('GET /whoami', async () => {
    let response = await auth
      .get('/whoami')
      .expect(403)

    expect(response.body).toMatchSnapshot();
  });

  test('POST /register', async () => {
    let response = await auth
      .post('/register')
      .send({
        firstName: 'test',
        lastName: 'test',
        nick: 'test_nick2',
        password: '123123123',
      });

    expect(response.body).toMatchObject({
      ok: true,
      data: {
        user: {
          id: expect.any(String),
          firstName: 'test',
          lastName: 'test',
          nick: 'test_nick2',
          passhash: expect.any(String),
          createdAt: expect.any(String),
        },
        session: {
          sessionId: expect.any(String),
          userId: expect.any(String),
          createdAt: expect.any(String),
          expiresAt: expect.any(String),
        },
      }
    });
  });

  test('POST /login', async () => {
    let response = await auth
      .post('/login')
      .send({
        nick: 'test_tester',
        password: '123123123',
      });

    expect(response.body).toMatchObject({
      ok: true,
      data: {
        user: {
          id: expect.any(String),
          firstName: 'Test',
          lastName: 'Tester',
          nick: 'test_tester',
          passhash: expect.any(String),
          createdAt: expect.any(String),
          lastActivityAt: expect.any(String),
        },
        session: {
          sessionId: expect.any(String),
          userId: expect.any(String),
          createdAt: expect.any(String),
          expiresAt: expect.any(String),
        },
      }
    });
  });

  test('POST /logout', async () => {
    let response = await auth
      .post('/logout')
      .set(SESSION_HEADER_KEY, 'bb23af03-be50-4bce-b729-b259b2e02e55');

    expect(response.body).toMatchObject({
      ok: true,
      data: {
        sessionId: 'bb23af03-be50-4bce-b729-b259b2e02e55',
        userId: expect.any(String),
        createdAt: expect.any(String),
        expiresAt: expect.any(String),
      }
    });
  });
});

