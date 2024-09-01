import request from './request';

describe('friendship', () => {
  let friends = request.extend({
    prefixUrl: '/friendship'
  });

  test('POST /send-request', async () => {
    let response = await friends
      .post('/send-request')
      .send({
        responderId: '48683858-796c-45ad-a361-9e3d6d003354',
        message: 'test add',
      })
      .expect(200);

    expect(response.body).toMatchSnapshot();
  });

  test('POST /accept-request', async () => {
    let response = await friends
      .post('/accept-request')
      .send({
        requesterId: '1ae02f69-ea1a-4308-b825-0e5896e652e4',
      })
      .expect(200);

    expect(response.body).toMatchSnapshot();
  });

  test('POST /deny-request', async () => {
    let response = await friends
      .post('/deny-request')
      .send({
        requesterId: 'baa0e8bc-385f-4314-9580-29855aff2229',
      })
      .expect(200);

    expect(response.body).toMatchSnapshot();
  });

  test('POST /remove-friend', async () => {
    let response = await friends
      .post('/remove-friend')
      .send({
        friendId: '2db640fd-7aa4-4bba-8ee6-3935b700297a',
      })
      .expect(200);

    expect(response.body).toMatchSnapshot();
  });
});

