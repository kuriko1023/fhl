<template>
  <view>
    <image class="background" src="/static/game_background.png" ></image>
    <view style="padding: 15px">
    <view id="subject" style="margin-bottom: 20px">
      <subject-block :mode="mode" :text="subject">
      </subject-block>
    </view>
    <view id="answering" style="display: flex; margin: 0 15px">
      <image class="circle" src="/static/picture.png" mode="widthFix"></image>
      <view style="display: flex;">
        <view v-for="word in answer" :key="word.word">
          <text class="kati" :style="{color: answerColor[word.highlight]}">{{word.word}}</text>
        </view>
      </view>
    </view>
<!--      <progress percent="80" stroke-width="3"/>-->
<!--      <progress percent="100" stroke-width="3"/>-->
      <view class="count_down">
      <count-down :active="active1" :color="playerColor[side]" :time="TURN_TIMER_MAX" :current="current1" @finish="onFinish" ref="countdown1"></count-down>
      <count-down :active="active2" color="#65d4e5" :time="CUMULATIVE_TIMER_MAX" :current="current2" @finish="onExtraFinish" @stop="onCountStop"></count-down>
      </view>
        <!--      <button @click="onStop">stop</button>-->
<!--      <button @click="onStop2">stop2</button>-->
    </view>
    <image src="/static/history_background.png" class="history_background"></image>
    <view id="answerHistory">
      <history-block :data="history"/>
    </view>
    <view id="submit">
      <form>
        <view style="display: flex">
          <input placeholder="输入答案" placeholder-style="color: #bac3ab; font-size: 12px" name="myAnswer"  class="input" adjust-position="false" maxlength="20" v-model='inputAnswer' />

<!--            TODO: 是否有效-->
        <view form-type="submit" class="btn" @click='onSubmitAnswer'>发送</view>
        </view>
      </form>
    </view>
<!--    临时测试用-->
<!--    <button @click="onEnd"> end </button>-->
<!--    <button @click="pop">test</button>-->
    <uni-popup ref="popup" type="message">
      <uni-popup-message type="warn" :message="popMessage"/>
<!--      <view id="popup">-->
<!--        test-->
<!--      </view>-->
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
      answerColor: {
        0: '#000000',
        1: '#ad2b29',
        2: '#5e270d',
      },
      playerColor: {
        0: '#535353',
        1: '#fcecbb',
      },
      TURN_TIMER_MAX: 6,
      CUMULATIVE_TIMER_MAX: 20,
      //shot countDown
      active1: true,
      //long countDown
      active2: false,
      current1: 0,
      current2: 0,
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
      subject: "古 梦 雁/长 舟 送 寄 事 神 不 生 西风 多少/1000010011",
      answer:  [],
      history: [],
      inputAnswer: '',
    }
  },
  onLoad() {
    this.registerSocketMessageListener();
  },
  methods:{
    pop(){
      this.$refs.popup.open()
    },
    onFinish(){
      this.info = 'finish'
      // this.$refs.countdown1.active = false
      // this.$refs.countdown2.active = true
      this.active1 = false
      this.active2 = true
    },
    onExtraFinish(){
      this.active2 = false
      console.log('timeout')
      /*this.sendMessage({
        'type': 'timeout'
      })*/
    },
    //change count down side
    changeSide(){
      this.history.push(this.answer)
      this.answer = []

      console.log(this.current1)
      this.current1 = this.TURN_TIMER_MAX
      if(this.side) {
        console.log("side the other" + this.side)
        this.side = 0
        this.current2 = this.myExtraTime
        console.log("myExtraTime" + this.myExtraTime)
      }
      else {
        console.log("side my" + this.side)
        this.side = 1
        this.current2 = this.sideExtraTime
        console.log("sideExtraTime" + this.sideExtraTime)
      }
      this.active1 = true
    },
    onStop(){
      this.active1 = false
      this.current1 = 0
      setTimeout(this.changeSide, 2000)
      // this.$refs.countdown1.active = false
    },
    onStop2(){
      this.active2 = false
      this.current1 = 0
      setTimeout(this.changeSide, 2000)
    },
    onCountStop(val){
      console.log('countStop')
      if(this.side) {
        console.log(val)
        this.sideExtraTime = val
      }
      else{
        console.log(val)
        this.myExtraTime = val
      }
    },
    onEnd(){
      uni.redirectTo({
        'url': '/pages/EndPage/EndPage'
      })
    },
    onSubmitAnswer(e){
      this.sendSocketMessage({
        'type': 'answer',
        'text': this.inputAnswer,
      })
      //计时器停止
      this.active1 = false
      this.active2 = false
    },

    onSocketMessage() {
      const msg = this.popSocketMessage(['game_status', 'game_update']);
      switch (msg.type){
        case 'game_status': {
          console.log('game_status', msg);
          this.mode = msg.mode
          this.subject = msg.subject
          this.history = []
          for(let i = 0; i < msg.history.length; i++){
            let sentence = this.historySentenceParse(msg.history[i])
            this.history.push(sentence)
          }
          break
        }
        case 'game_update':{
          console.log('game_update', msg);
          this.answer = this.historySentenceParse(msg.text)
          switch (this.mode){
            case 'B': {
              let index = parseInt(msg.update)
              this.subject.subject1[index].show = 2
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
          //TODO: 计时器重置
          this.active1 = false
          this.current1 = 0
          setTimeout(this.changeSide, 2000)
          break
        }
        case 'game_end':{
          uni.navigateTo({
            'url': '/pages/EndPage/EndPage'
          })
          break
        }
        case 'invalid_answer':{

        }
      }
    }
  }
}
</script>

<style scoped>
  .circle {
    width: 10%;
    border-radius: 20%;
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
    bottom: 0;
    left: 0;
    width: 100%;
    margin-bottom: 5%;
  }
  .kati {
    font-family: "楷体","楷体_GB2312";
    font-size: 21px;
    padding: 0 1px;
  }

  .btn{
    background-color: #366440;
    font-family: "楷体","楷体_GB2312";
    color: white;
    border-radius: 6px;
    font-size: 14px;
    padding: 6px 10px;
    width: 10%;
    text-align: center;
    vertical-align: center;
    margin-right: 5%;
  }
  .count_down{
    margin: 10px 5px;
  }
</style>
