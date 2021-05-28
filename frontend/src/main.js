import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

Vue.prototype.sendMessage = function (msg) {
  uni.sendSocketMessage({
        data: msg
      })
}

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
