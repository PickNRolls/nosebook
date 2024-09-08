import request, { ASSER_SESSION } from './request';
import { WebSocket } from "./websocket";

describe('like', () => {
  let like = request.extend({
    prefixUrl: '/like'
  });

  test('POST /post', async () => {
    const websocket = new WebSocket(ASSER_SESSION).unwrap();

    const message = new Promise(res => {
      websocket.once('message', (data) => {
        res(JSON.parse(data.toString()));
      });
    });

    async function doLike() {
      return await like
        .post('/post')
        .send({
          id: 'c7b7bf17-38f9-4ed5-b0a8-501a90f7c829',
        })
        .expect(200);
    }

    expect((await doLike()).body).toMatchSnapshot();
    expect((await doLike()).body).toMatchSnapshot();

    expect(await message).toMatchSnapshot();

    websocket.terminate();
  });

  test('POST /comment', async () => {
    const websocket = new WebSocket(ASSER_SESSION).unwrap();

    const message = new Promise(res => {
      websocket.once('message', (data) => {
        res(JSON.parse(data.toString()));
      });
    });

    async function doLike() {
      return await like
        .post('/comment')
        .send({
          id: '620c79b7-3927-48b7-a308-1ffd3db6036f',
        })
        .expect(200);
    }

    expect((await doLike()).body).toMatchSnapshot();
    expect((await doLike()).body).toMatchSnapshot();

    expect(await message).toMatchSnapshot();

    websocket.terminate();
  });
});

