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

  describe('GET /', () => {
    test('filter friends for userId', async () => {
      let response = await friends
        .get('/')
        .query({
          userId: 'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
          accepted: true,
        })
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });

    test('filter friends with limit', async () => {
      let response = await friends
        .get('/')
        .query({
          userId: 'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
          accepted: true,
          limit: 1,
        })
        .expect(200);

      expect(response.body).toMatchSnapshot();

      response = await friends
        .get('/')
        .query({
          userId: 'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
          limit: 1,
          accepted: true,
          next: response.body.data.next,
        })
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });

    test('filter online friends', async () => {
      let response = await friends
        .get('/')
        .query({
          userId: 'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
          onlyOnline: true,
          accepted: true,
        })
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });

    test('filter incoming pending requests', async () => {
      let response = await friends
        .get('/')
        .query({
          userId: 'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
          onlyIncoming: true,
          accepted: false,
        })
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });

    test('filter outcoming pending requests', async () => {
      let response = await friends
        .get('/')
        .query({
          userId: 'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
          onlyOutcoming: true,
          accepted: false,
        })
        .expect(200);

      expect(response.body).toMatchSnapshot();
    });
  });
});

