import supertest from 'supertest';
import request from './request';

describe('posts', () => {
  let posts = request.extend({ prefixUrl: '/posts' });
  let createdPostId: string;

  // TODO: fix
  test.skip('GET /', async () => {
    let lastResponse = await posts
      .get('/')
      .query({
        ownerId: '1ae02f69-ea1a-4308-b825-0e5896e652e4',
        authorId: '1ae02f69-ea1a-4308-b825-0e5896e652e4',
      })
      .expect(200)

    expect(lastResponse.body).toMatchSnapshot();

    lastResponse = await posts
      .get('/')
      .query({
        ownerId: '1ae02f69-ea1a-4308-b825-0e5896e652e4',
        authorId: '1ae02f69-ea1a-4308-b825-0e5896e652e4',
        cursor: lastResponse.body.data.next,
      })
      .expect(200);

    expect(lastResponse.body).toMatchSnapshot();
  });

  test('POST /publish', async () => {
    let response = await posts
      .post('/publish')
      .send({
        ownerId: '1ae02f69-ea1a-4308-b825-0e5896e652e4',
        message: 'my test message',
      })
      .expect(200);

    expect(response.body).toStrictEqual(expect.objectContaining({
      data: {
        id: expect.any(String),
      },
      ok: true,
    }));
    createdPostId = response.body.data.id;
  });

  test('POST /remove', async () => {
    const response = await posts
      .post('/remove')
      .send({
        id: createdPostId,
      }).expect(200);

    expect(response.body).toStrictEqual(expect.objectContaining({
      data: {
        id: createdPostId,
      },
      ok: true,
    }));
  });

  describe('permissions', () => {
    describe('GET /:id, has permissions', () => {
      test('auth.user == post.author', async () => {
        let response = await posts
          .get(`/27b7bf17-38f9-4ed5-b0a8-501a90f7c8e7`)
          .expect(200);

        expect(response.body).toMatchSnapshot();
      });

      test('auth.user == post.owner', async () => {
        let response = await posts
          .get(`/27b7bf27-38f9-4ed5-b0a8-501a90f7c8e7`)
          .expect(200);

        expect(response.body).toMatchSnapshot();
      });
    });

    test('GET /:id, has no permissions', async () => {
      let response = await posts
        .get('/c7b7bf17-38f9-4ed5-b0a8-011a90f7c8e2')
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });
  });
});

