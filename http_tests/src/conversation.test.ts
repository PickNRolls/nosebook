import request, { ASSER_ID, ASSER_SESSION } from "./request";
import { WebSocket } from "./websocket";

describe('conversation', () => {
  const conv = request.extend({ prefixUrl: '/conversations' });

  test('POST /send-message', async () => {
    const websocket = new WebSocket(ASSER_SESSION).unwrap();

    const message = new Promise(res => {
      websocket.once('message', (data) => {
        res(data.toString());
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
    expect(await message).toMatchSnapshot();

    websocket.terminate();
  });
});

