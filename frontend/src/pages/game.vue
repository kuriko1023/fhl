<template>
  <view>
    <image class="background" src="/static/game_background_scaled.jpg" ></image>
    <image class="background" :src="staticRes('game_background.png')" ></image>
    <view style="padding: 15px">
    <view id="subject" style="margin-bottom: 20px">
      <subject-block :mode="mode" :subject="subject">
      </subject-block>
    </view>
    <view id="answering" style="display: flex; margin: 0 15px">
      <image class="circle" :style="'width: 28px; height: 28px; ' + (((side === 1) ^ isHost) ? '' : 'display: none')" :src="hostAvatar" mode="widthFix"></image>
      <image class="circle" :style="'width: 28px; height: 28px; ' + (((side === 0) ^ isHost) ? '' : 'display: none')" :src="guestAvatar" mode="widthFix"></image>
      <view style="display: flex; width: 100%; overflow: scroll; height: 28px">
        <view v-for="word in answer" :key="word.word">
          <p class="kati" :style="{color: answerColor[word.highlight]}">{{word.word}}</p>
        </view>
      </view>
    </view>
      <view class="count_down">
        <count-down :active="timersActive" :color="playerColor[side]" :total="TURN_TIMER_MAX" :start="timerStart1" :end="timerEnd1"></count-down>
        <count-down :active="timersActive" color="#65d4e5" :total="CUMULATIVE_TIMER_MAX" :start="timerEnd1" :end="timerEnd2"></count-down>
      </view>
    </view>
    <image src="/static/history_background_scaled.jpg" class="history_background"></image>
    <image :src="staticRes('history_background.png')" class="history_background"></image>
    <view id="answerHistory">
      <history-block :data="history"/>
    </view>
    <view id="submit">
      <form>
        <view style="display: flex">
          <input :placeholder="side === 0 ? '可用标点分隔多句，至多 21 字' : '轮到对方，等待中'" @confirm='onSubmitAnswer' placeholder-style="color: #bac3ab; font-size: 12px" name="myAnswer"  class="input" adjust-position="false" maxlength="24" v-model='inputAnswer' :disabled='answerSendTimer !== -1' />
          <view v-if='side === 0' form-type="submit" :class="'btn' + (answerSendTimer !== -1 ? ' disabled' : '')" @click='onSubmitAnswer'>发送</view>
        </view>
      </form>
    </view>
    <uni-popup ref="popup" type="message">
      <uni-popup-message type="warn" :message="popMessage"/>
    </uni-popup>
  </view>
</template>

<script>
// import uniPopUp from '@dcloudio/uni-ui/lib/uni-popup/uni-popup.vue';
import uniPopupMessage from '@dcloudio/uni-ui/lib/uni-popup-message/uni-popup-message.vue';
import {uniPopup} from '@dcloudio/uni-ui'
// import {uniPopupMessage} from '@dcloudio/uni-ui'
import Subject from "@/components/Subject";
import History from "@/components/History";
import CountDown from "@/components/CountDown";
export default {
  name: "GamePage",
  components: {
    "subject-block": Subject,
    "history-block": History,
    "count-down": CountDown,
    uniPopup,
    uniPopupMessage,
  },
  data() {
    return {
      isHost: false,
      hostAvatar: '',
      guestAvatar: '',
      answerColor: {
        0: '#000000',
        1: '#ad2b29',
        2: '#5e270d',
      },
      playerColor: {
        0: '#535353',
        1: '#fcecbb',
      },
      TURN_TIMER_MAX: 60,
      CUMULATIVE_TIMER_MAX: 60,
      timerActive: false,
      timerStart1: 0,
      timerEnd1: 0,
      timerEnd2: 0,
      side: 0,  // 0 -- 我方  1 -- 对方
      myExtraTime: 29.99,
      sideExtraTime: 29.99,
      info: 'counting...',
      popMessage: '无效的答案',
      mode: "C",
      // subject: {
      //   "subject1": ["雁", "古", "梦"],
      //   // '长' 舟 送 寄 事 神 不 生 西风 多少
      //   "subject2": [
      //     {
      //       'value':'长',
      //       'show': 1,
      //     },
      //     {
      //       'value': '舟',
      //       'show': 1,
      //     },
      //     {
      //       'value': '送',
      //       'show': 1,
      //     },
      //     {
      //       'value': '寄',
      //       'show': 1,
      //     },
      //     {
      //       'value': '多少',
      //       'show': 0,
      //     },
      //   ],
      // },
      subject: {},
      answer:  [],
      history: [],
      inputAnswer: '',
      answerSendTimer: -1,
    }
  },
  onLoad() {
    window.history.pushState({}, '', window.location.pathname + window.location.search);
    this.isHost = getApp().globalData.isHost;
    this.side = (this.isHost ? 0 : 1);
    this.hostAvatar = getApp().globalData.hostAvatar;
    this.guestAvatar = getApp().globalData.guestAvatar;
    this.registerSocketMessageListener();
  },
  methods:{
    pop(){
      this.$refs.popup.open()
    },
    //change count down side
    changeSide(){
      this.history.push(this.answer)
      this.answer = []

      let extraTime
      if(this.side) {
        console.log("side opponent -> my")
        this.side = 0
        extraTime = this.myExtraTime
      }
      else {
        console.log("side my -> opponent")
        this.side = 1
        extraTime = this.sideExtraTime
      }
      const now = Date.now() / 1000
      this.timerStart1 = now
      this.timerEnd1 = now + this.TURN_TIMER_MAX
      this.timerEnd2 = now + this.TURN_TIMER_MAX + extraTime
      this.timersActive = true
    },
    clearAnswerSendTimer() {
      if (this.answerSendTimer !== -1) {
        clearTimeout(this.answerSendTimer)
        this.answerSendTimer = -1
        return true
      } else {
        return false
      }
    },
    onSubmitAnswer(e){
      if (this.side !== 0) return;
      if (this.answerSendTimer !== -1) return;
      this.answerSendTimer = setTimeout(() => {
        if (this.clearAnswerSendTimer()) {
          this.popMessage = '【断线】请检查网络连接并重试'
          this.pop()
        }
      }, 5000)
      const normalizedAnswer = this.inputAnswer
        .replace(/[ ，。？！\/,.?!]+/g, ' ')
        .replace(/《》“”「」『』—/g, '')
        .trim()
      this.sendSocketMessage({
        'type': 'answer',
        'text': normalizedAnswer,
      })
    },
    updateTimers (msg) {
      const guestTimer = msg.guest_timer / 1000;
      const hostTimer = msg.host_timer / 1000;
      if (getApp().globalData.isHost) {
        this.myExtraTime = hostTimer;
        this.sideExtraTime = guestTimer;
      } else {
        this.myExtraTime = guestTimer;
        this.sideExtraTime = hostTimer;
      }
      const now = Date.now() / 1000
      this.timerStart1 = now
      this.timerEnd1 = now + msg.turn_timer / 1000
      this.timerEnd2 = now + msg.turn_timer / 1000 +
        (this.side === 0 ? this.myExtraTime : this.sideExtraTime)
      this.timersActive = true
    },

    onSocketMessage() {
      if (this.tryPeekSocketMessage('end_status')) {
        uni.redirectTo({
          url: '/pages/finish'
        });
        return;
      }

      const msg = this.popSocketMessage(['game_status', 'game_update', 'invalid_answer']);
      switch (msg.type){
        case 'game_status': {
          console.log('game_status', msg);
          this.mode = msg.mode
          this.subject = this.parseSubject(msg.mode, msg.subject)
          this.history = []
          for(let i = 0; i < msg.history.length; i++){
            let sentence = this.historySentenceParse(msg.history[i])
            this.history.push(sentence)
          }
          this.side = (this.history.length % 2) ^ (this.isHost ? 0 : 1);
          this.updateTimers(msg)
          break
        }
        case 'game_update':{
          console.log('game_update', msg);
          this.answer = this.historySentenceParse(msg.text)
          this.timersActive = false
          if (this.side === 0) this.inputAnswer = ''
          setTimeout(() => {
            this.clearAnswerSendTimer()
            switch (this.mode){
              case 'B': {
                let index = parseInt(msg.update)
                for (let i = 0; i < this.subject.subject1.length; i++)
                  this.subject.subject1[i].show =
                    (i < index ? 2 : (i === index ? 1 : 0));
                break
              }
              case 'C': {
                let index = parseInt(msg.update)
                this.subject.subject2[index].show = 0
                break
              }
              case 'D': {
                let index = msg.update.split(',')
                this.subject.subject1[parseInt(index[0])].show = 0
                this.subject.subject2[parseInt(index[1])].show = 0
                break
              }
            }
            this.updateTimers(msg)
            this.changeSide()
          }, 2000);
          break
        }
        case 'invalid_answer':{
          this.clearAnswerSendTimer()
          let popupText = '【' + msg.reason + '】　'
          switch (msg.reason) {
            case '不审题': popupText += '不匹配剩余的关键词'; break;
            case '复读机': popupText += '与此前的答案重复'; break;
            case '没背熟': popupText += '存在相似但不一致的诗句'; break;
            case '大文豪': popupText += '没有找到相似的诗句'; break;
            case '捣浆糊': popupText += '总字数少于四字'; break;
            case '碎碎念': popupText += '总字数多于二十一字'; break;
          }
          this.popMessage = popupText
          this.pop()
        }
      }
    }
  }
}
</script>

<style scoped>
  .circle {
    width: 10%;
    border-radius: 50%;
    margin-right: 5px;
  }
  #popup {
    background: white;
    padding: 20px;
  }
  .background{
    position: absolute;
    height: 100%;
    width: 100%;
    z-index: -1;
  }
  #answerHistory{
    position: absolute;
    height: 30%;
    width: 80%;
    margin: 10px 10% 10px 10%;
    overflow: scroll;
  }
  .history_background{
    position: absolute;
    height: 30%;
    width: 80%;
    margin: 10px 10% 10px 10%;
  }
  .input {
    border: 1px solid #bac3ab;
    border-radius: 7px;
    background: white;
    padding: 3px;
    width: 70%;
    margin: 0 5%;
  }
  #submit{
    position: fixed;
    bottom: 5vh;
    left: calc(50% - 30vh);
    width: 100%;
    max-width: 60vh;
  }
  .kati {
    font-family: Kai;
    font-size: 21px;
    padding: 0 1px;
  }

  .btn{
    background-color: #366440;
    font-family: Kai;
    color: white;
    border-radius: 6px;
    font-size: 14px;
    padding: 6px 10px;
    width: 10%;
    text-align: center;
    vertical-align: center;
    margin-right: 5%;
  }

  .btn.disabled {
    background-color: #aaccbb;
  }

  .count_down{
    margin: 10px 5px;
  }
</style>
