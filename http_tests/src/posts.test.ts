import supertest from 'supertest';
import request from './request';

describe('posts', () => {
  let posts = request.extend('/posts');
  let createdPostId: string;

  test('GET /', async () => {
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

    expect(response.body).toEqual(expect.objectContaining({
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

    expect(response.body).toEqual(expect.objectContaining({
      data: {
        id: createdPostId,
      },
      ok: true,
    }));
  });
});

