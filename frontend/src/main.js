import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

Vue.prototype.retrieveServerProfile = function (callback) {
  if (getApp().globalData.my) {
    callback();
    return;
  }
  const req = () => uni.login({success: (res) => uni.request({
    // url: 'https://flyhana.starrah.cn/profile/!kuriko1023',
    url: 'https://flyhana.starrah.cn/profile/' + res.code,
    success: (res) => {
      const obj = res.data;
      if (!obj || !obj.id) {
        req();
        return;
      }
      getApp().globalData.my = {
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

Vue.prototype.requestLocalProfile = function (callback) {
  console.log(getApp().globalData.my)
  if (getApp().globalData.my.nickname === null) {
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
    callback()
  }
};

let socketTask = null;

let messageListener = null;
const messageQueue = [];

Vue.prototype.connectSocket = function (obj) {
  const conn = () => {
    messageQueue.splice(0);
    socketTask = uni.connectSocket(obj);
  };
  if (socketTask) {
    socketTask.onClose(conn);
    socketTask.close();
  } else {
    conn();
  }
};

Vue.prototype.closeSocket = function () {
  socketTask.close();
  socketTask.onClose(() => socketTask = null);
};

Vue.prototype.registerSocketMessageListener = function () {
  messageListener = this;
  if (messageQueue.length > 0)
    messageListener.onSocketMessage();
};

Vue.prototype.sendSocketMessage =
  (msg) => uni.sendSocketMessage({data: JSON.stringify(msg)});

Vue.prototype.peekSocketMessage = () => messageQueue[0];
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

uni.onSocketClose(() => {
  console.log('socket closed!');
});

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
        tmpObject.word = parse[0][j]
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
      subject.subject1 = []
      subject.subject1.push(text)
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
