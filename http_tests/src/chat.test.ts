import request, { ASSER_ID, ASSER_SESSION } from "./request";
import { WebSocket } from "./websocket";

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

  const conv = request.extend({ prefixUrl: '/conversations' });

  test('POST /send-message', async () => {
    const websocket = new WebSocket(ASSER_SESSION).unwrap();

    const message = new Promise(res => {
      websocket.on('message', (data) => {
        const message = JSON.parse(data.toString());
        if (message.type === 'new_message') {
          res(message);
        }
      });
    });

    let response = await conv
      .post('/send-message')
      .send({
        recipientId: ASSER_ID,
        text: 'Hello asser'
      })
      .expect(200);

    expect(response.body).toMatchSnapshot();

    const m = await message as any
    m.payload.id = 'test id';
    m.payload.chatId = 'test chat id';
    expect(m).toMatchSnapshot();

    websocket.terminate();
  });
});

