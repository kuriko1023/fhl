import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

Vue.prototype.sendMessage = function (msg) {
  uni.sendSocketMessage({
    data: msg
  })
}

let messageListener = null;
const messageQueue = [];

Vue.prototype.registerSocketMessageListener = function () {
  messageListener = this;
  if (messageQueue.length > 0)
    messageListener.onSocketMessage();
};

Vue.prototype.peekSocketMessage = () => messageQueue[0];
Vue.prototype.popSocketMessage = (type) =>
  (type === undefined || messageQueue[0].type === type) ?
  messageQueue.shift() : undefined;

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
        let index = j + ''
        if(parse[1].indexOf(index) !== -1){
            tmpObject.highlight = 1
        }
        if(parse.length > 2 && parse[2].indexOf(index) !== -1){
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

App.mpType = 'app'

const app = new Vue({
  ...App
})
app.$mount()
