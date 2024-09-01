import request from './request';

describe('friendship', () => {
  let auth = request.extend({
    prefixUrl: '/friendship'
  });

  test('POST /send-request', async () => {
    let response = await auth
      .post('/send-request')
      .send({
        responderId: '48683858-796c-45ad-a361-9e3d6d003354',
        message: 'test add',
      })
      .expect(200);

    expect(response.body).toMatchSnapshot();
  });
});

