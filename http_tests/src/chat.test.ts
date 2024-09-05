import request from './request';

describe('posts', () => {
  let chats = request.extend({ prefixUrl: '/chats' });


  describe('GET /', () => {
    test('result list of 2 chats', async () => {
      let response = await chats
        .get('')
        .expect(200)

      expect(response.body).toMatchSnapshot();
    });

    test('filter limit = 1', async () => {
      let response = await chats
        .get('')
        .query({
          limit: 1,
        })
        .expect(200)

      expect(response.body).toMatchSnapshot();

      response = await chats
        .get('')
        .query({
          limit: 1,
          next: response.body.data.next,
        })
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });
  });
});

