// const apiServer = 'https://fhxl.ayuu.ink';
// const wsServer = 'wss://fhxl.ayuu.ink';
const apiServer = 'http://localhost:2310';
const wsServer = 'ws://localhost:2310';

const staticRes = (file) => `${apiServer}/static/${file}`;

export {
  apiServer, wsServer, staticRes,
}
