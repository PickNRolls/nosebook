import supertest from 'supertest';
import request from './request';

describe('comments', () => {
  let comments = request.extend({ prefixUrl: '/comments' });
  let createdCommentId: string;

  test('GET /', async () => {
    const postId = 'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7';

    let lastResponse = await comments
      .get('/')
      .query({
        postId,
      })
      .expect(200)

    expect(lastResponse.body).toMatchSnapshot();

    lastResponse = await comments
      .get('/')
      .query({
        postId,
        next: lastResponse.body.data.next,
      })
      .expect(200);

    expect(lastResponse.body).toMatchSnapshot();
  });

  test('POST /publish-on-post', async () => {
    let response = await comments
      .post('/publish-on-post')
      .send({
        id: 'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
        message: 'my test comment',
      })
      .expect(200);

    expect(response.body).toStrictEqual(expect.objectContaining({
      data: {
        id: expect.any(String),
      },
      ok: true,
    }));
    createdCommentId = response.body.data.id;
  });

  test('POST /remove', async () => {
    const response = await comments
      .post('/remove')
      .send({
        id: createdCommentId,
      }).expect(200);

    expect(response.body).toStrictEqual(expect.objectContaining({
      data: {
        id: createdCommentId,
      },
      ok: true,
    }));
  });

  describe('permissions', () => {
    test('GET /:id, has permissions', async () => {
      let response = await comments
        .get(`/120c79b8-3927-48b7-a308-1ffd3db6036f`)
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });

    test('GET /:id, has no permissions', async () => {
      let response = await comments
        .get('/d0023f4d-8d7f-4907-9438-d2ed2a9661f0')
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });
  });
});

