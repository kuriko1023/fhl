import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

import {
  apiServer, wsServer, staticRes,
} from 'utils';

Vue.mixin({
  methods: { staticRes }
});

if (uni.getSystemInfoSync().uniPlatform === 'mp-weixin') {
  Vue.prototype.adaptedLogin = (options) => {
    uni.login({
      success: (resp) => options.success('wx_login:' + resp.code),
      fail: () => options.fail && options.fail(),
    });
  }
} else {
  let uid = null;
  Vue.prototype.adaptedLogin = (options) => {
    if (uid === null) {
      do {
        uid = prompt('你的名字~（不超过 20 字）');
      } while (uid === null || uid.length > 20);
      /* if (uid === null) {
        if (options.fail) options.fail();
        return;
      } */
    }
    const hex = [];
    for (const n of new TextEncoder().encode(uid))
      hex.push(n.toString(16).padStart('0', 2))
    options.success('web_login:' + hex.join(''));
  };
}

Vue.prototype.retrieveServerProfile = function (callback) {
  if (getApp().globalData.my) {
    callback();
    return;
  }
  const req = () => this.adaptedLogin({success: (code) => uni.request({
    url: `${apiServer}/profile/${code}`,
    success: (res) => {
      const obj = res.data;
      if (!obj || !obj.id) {
        req();
        return;
      }
      getApp().globalData.my = {
        id: obj.id,
        nickname: obj.nickname,
      };
      callback();
    },
    fail: () => req(),
  })});
  req();
};

Vue.prototype.requestLocalProfile = function (callback) {
  console.log(getApp().globalData.my)
  if (!getApp().globalData.my.nickname) {
    if (uni.getSystemInfoSync().uniPlatform === 'mp-weixin') {
      uni.getUserProfile({
        desc: '用于向其他玩家展示头像和昵称',
        success: (res) => {
          console.log(res)
          getApp().globalData.my.avatar = res.userInfo.avatarUrl
          getApp().globalData.my.nickname = res.userInfo.nickName
          getApp().globalData.my.create = true
          callback()
        },
        fail: () => {
          console.log('getUserProfile() failed')
        },
      })
    } else {
      alert('!!!')
    }
  } else {
    callback()
  }
};

let socketTask = null;
let socketCloseCallback = null;
let socketUrl = null;

let messageListener = null;
const messageQueue = [];

uni.onSocketClose(() => {
  console.log('socket closed!');
  socketTask = null;
  messageQueue.splice(0);
  if (socketCloseCallback) {
    socketCloseCallback();
    socketCloseCallback = null;
  } else {
    // Unexpected disconnection
    console.log('unexpected disconnection');
    const conn = () => socketUrl().then((url) => socketTask = uni.connectSocket({
      url: url,
      success: () => {
        console.log('reconnected! ' + url)
      },
      fail: () => setTimeout(conn, 1000),
    }));
    conn();
  }
});

Vue.prototype.connectSocket = function (obj) {
  socketUrl = obj.url;
  const conn = () =>
    socketUrl().then((url) =>
      socketTask = uni.connectSocket(Object.assign(obj, {url: url})));
  if (socketTask) {
    socketCloseCallback = conn;
    socketTask.close();
  } else {
    conn();
  }
};

Vue.prototype.closeSocket = function () {
  if (!socketTask) return;
  socketCloseCallback = () => socketTask = null;
  socketTask.close();
};

Vue.prototype.registerSocketMessageListener = function () {
  messageListener = this;
  if (messageQueue.length > 0)
    messageListener.onSocketMessage();
};

Vue.prototype.sendSocketMessage =
  (msg) => {
    console.log('send', msg)
    uni.sendSocketMessage({data: JSON.stringify(msg), success: () => console.log('sent', msg), fail: () => console.log('fail', msg)});
  }

Vue.prototype.peekSocketMessage = () => messageQueue[0];
Vue.prototype.tryPeekSocketMessage = (type) => {
  for (let i = 0; i < messageQueue.length; i++)
    if (messageQueue[i].type === type) {
      messageQueue.splice(0, i);
      return true;
    }
  return false;
};
Vue.prototype.tryPopSocketMessage = (type) =>
  (type === undefined || messageQueue[0].type === type) ?
  messageQueue.shift() : {_none: true};
Vue.prototype.popSocketMessage = (types) => {
  if (typeof types === 'string') types = [types];
  while (messageQueue.length > 0) {
    const msg = messageQueue.shift();
    if (types === undefined || types.indexOf(msg.type) !== -1) return msg;
  }
  return {_none: true};
};

uni.onSocketMessage((res) => {
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

Vue.prototype.historySentenceParse = function(str) {
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

Vue.prototype.getHistory = function(){

}

Vue.prototype.parseSubject = function (mode, text) {
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

App.mpType = 'app'

const app = new Vue({
  ...App
})
app.$mount()
