const apiServer = 'https://fhxl.ayuu.ink';
const wsServer = 'wss://fhxl.ayuu.ink';

const staticRes = (file) => `${apiServer}/static/${file}`;

export {
  apiServer, wsServer, staticRes,
}
