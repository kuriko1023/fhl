import Taro from '@tarojs/taro';

const G = {};

const redirect = (url) => {
  Taro.redirectTo({
    url,
  });
};

const retrieveServerProfile = function (callback) {
  if (G.my) {
    callback();
    return;
  }
  const req = () => Taro.login({success: (res) => Taro.request({
    url: 'http://123.57.21.143/profile/' + res.code,
    success: (res) => {
      const obj = res.data;
      if (!obj || !obj.id) {
        req();
        return;
      }
      G.my = {
        id: obj.id,
        avatar: obj.avatar,
        nickname: obj.nickname,
      };
      callback();
    },
    fail: () => req(),
  })});
  req();
};

const requestLocalProfile = function (callback) {
  if (G.my.nickname === null) {
    Taro.getUserProfile({
      desc: '用于向其他玩家展示头像和昵称',
      success: (res) => {
        G.my.avatar = res.userInfo.avatarUrl
        G.my.nickname = res.userInfo.nickName
        G.my.create = true
        callback()
      },
      fail: () => {
        console.log('getUserProfile() failed')
      },
    })
  } else {
    callback()
  }
};

let socketTask = null;
let socketCloseCallback = null;
let socketUrl = null;

let messageListener = null;
const messageQueue = [];

Taro.onSocketClose(() => {
  console.log('socket closed!');
  socketTask = null;
  messageQueue.splice(0);
  if (socketCloseCallback) {
    socketCloseCallback();
    socketCloseCallback = null;
  } else {
    // Unexpected disconnection
    console.log('unexpected disconnection');
    const conn = () => socketUrl().then((url) => socketTask = Taro.connectSocket({
      url: url,
      success: () => {
        console.log('reconnected! ' + url)
      },
      fail: () => setTimeout(conn, 1000),
    }));
    conn();
  }
});

const connectSocket = function (obj) {
  socketUrl = obj.url;
  const conn = () =>
    socketUrl().then((url) =>
      socketTask = Taro.connectSocket(Object.assign(obj, {url: url})));
  if (socketTask) {
    socketCloseCallback = conn;
    socketTask.close();
  } else {
    conn();
  }
};

const closeSocket = function () {
  if (!socketTask) return;
  socketCloseCallback = () => socketTask = null;
  socketTask.close();
};

const registerSocketMessageListener = function () {
  messageListener = this;
  if (messageQueue.length > 0)
    messageListener.onSocketMessage();
};

const sendSocketMessage =
  (msg) => {
    console.log('send', msg)
    Taro.sendSocketMessage({data: JSON.stringify(msg), success: () => console.log('sent', msg), fail: () => console.log('fail', msg)});
  }

const peekSocketMessage = () => messageQueue[0];
const tryPeekSocketMessage = (type) => {
  for (let i = 0; i < messageQueue.length; i++)
    if (messageQueue[i].type === type) {
      messageQueue.splice(0, i);
      return true;
    }
  return false;
};
const tryPopSocketMessage = (type) =>
  (type === undefined || messageQueue[0].type === type) ?
  messageQueue.shift() : {_none: true};
const popSocketMessage = (types) => {
  if (typeof types === 'string') types = [types];
  while (messageQueue.length > 0) {
    const msg = messageQueue.shift();
    if (types === undefined || types.indexOf(msg.type) !== -1) return msg;
  }
  return {_none: true};
};

Taro.onSocketMessage((res) => {
  let payload = res.data;
  if (typeof payload !== 'string') return;
  try {
    payload = JSON.parse(payload);
  } catch (e) {
    return;
  }

  if (payload.error || true) console.log(payload);
  messageQueue.push(payload);

  if (messageListener !== null)
    messageListener.onSocketMessage();
});

const historySentenceParse = function(str) {
    let sentence = []
    let parse = str.split('/')
    for(let j = 0; j < parse[0].length; j++){
        let tmpObject = {}
        tmpObject.word = (parse[0][j] === ' ' ? '　' : parse[0][j])
        let index = j.toString(36)
        if(parse[1].indexOf(index) !== -1){
            tmpObject.highlight = 1
        }
        else if(parse.length > 2 && parse[2].indexOf(index) !== -1){
            tmpObject.highlight = 2
        }
        else{
            tmpObject.highlight = 0
        }
        sentence.push(tmpObject)
    }
    return sentence
}

const getHistory = function(){

}

const parseSubject = function (mode, text) {
  let subject = {}
  switch (mode) {
    case 'A': {
      subject.subject1 = text.split('')
      break
    }
    case 'B': {
      // 1 要用的字   2 当前的字
      let parse = text.split('/')
      subject.subject1 = []
      for (let i = 0; i < parse[0].length; i++) {
        let tmpObject = {}
        tmpObject.value = parse[0][i]
        if(parseInt(parse[1]) === i){
          tmpObject.show = 1
        }
        else  tmpObject.show = 0
        subject.subject1.push(tmpObject)
      }
      break
    }
    case 'C': {
      // 1 未使用  0 已使用
      let parse = text.split('/')
      subject.subject1 = parse[0].split(' ')
      let tmpArray = parse[1].split(' ')
      subject.subject2 = []
      for (let i = 0; i < tmpArray.length; i++) {
        let tmpObject = {}
        tmpObject.value = tmpArray[i]
        if(parse[2] !== undefined) {
          tmpObject.show = 1 - parseInt(parse[2][i])
          console.log("gamePage")
        }
        else {
          tmpObject.show = 1
          console.log("endPage")
        }
        subject.subject2.push(tmpObject)
      }
      break
    }
    case 'D': {
      // 1 未使用  0 已使用
      let parse = text.split('/')
      let tmpArray1 = parse[0].split(' ')
      let tmpArray2 = parse[1].split(' ')
      subject.subject1 = []
      subject.subject2 = []
      for (let i = 0; i < tmpArray1.length; i++) {
        let tmpObject1 = {}
        tmpObject1.value = tmpArray1[i]
        if(parse[2] !== undefined) {
          tmpObject1.show = 1 - parseInt(parse[2][i])
        }
        else tmpObject1.show = 1
        subject.subject1.push(tmpObject1)
        let tmpObject2 = {}
        tmpObject2.value = tmpArray2[i]
        if(parse[3] !== undefined) {
          tmpObject2.show = 1 - parseInt(parse[3][i])
        }
        else tmpObject2.show = 1
        subject.subject2.push(tmpObject2)
      }
      break
    }
  }
  return subject
}

export {
  G,
  redirect,
  retrieveServerProfile,
  requestLocalProfile,
  connectSocket,
  closeSocket,
  registerSocketMessageListener,
  sendSocketMessage,
  peekSocketMessage,
  tryPeekSocketMessage,
  tryPopSocketMessage,
  historySentenceParse,
  parseSubject,
};
