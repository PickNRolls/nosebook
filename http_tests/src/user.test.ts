import request from './request';

describe('users', () => {
  let users = request.extend({ prefixUrl: '/users' });

  describe('online', () => {
    test('GET /:id', async () => {
      const response = await users
        .get('/ed1a3fd0-4d0b-4961-b4cd-cf212357740d')
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });
  });

  describe('offline', () => {
    test('GET /:id', async () => {
      const response = await users
        .get('/1ae02f69-ea1a-4308-b825-0e5896e652e4')
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });
  });
});

