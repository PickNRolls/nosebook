import { WebSocket } from "./websocket";

describe('websocket', () => {
  test('GET /ws', async () => {
    const websocket = new WebSocket().unwrap();

    const onUpgrade = jest.fn();
    const onOpen = jest.fn();

    const wait = new Promise(resolve => {
      websocket.once('upgrade', onUpgrade);
      websocket.once('open', () => {
        onOpen();
        resolve(null);
      });
    });

    await wait;

    expect(onUpgrade).toHaveBeenCalled();
    expect(onOpen).toHaveBeenCalled();

    return new Promise(resolve => {
      websocket.once('close', resolve);
      websocket.terminate();
    });
  });
});

