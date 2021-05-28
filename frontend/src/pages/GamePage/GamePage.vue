<template>
  <view>
    <view id="subject">
      <subject-block :mode="mode" :text="subject">
      </subject-block>
    </view>
    <view id="answering">
      <image class="circle" src="/static/picture.png" mode="widthFix"></image>
      <view style="display: flex">
        <view v-for="word in answer" :key="word.value">
          <text :style="{color: answerColor[word.highlight]}">{{word.value}}</text>
        </view>
      </view>
<!--      <progress percent="80" stroke-width="3"/>-->
<!--      <progress percent="100" stroke-width="3"/>-->
      <count-down :active="active1" :color="playerColor[side]" :time="10" :current="current1" @finish="onFinish" ref="countdown1"></count-down>
      <count-down :active="active2" color="#65d4e5" :time="30" :current="current2" @finish="onExtraFinish" @stop="onCountStop"></count-down>
      <button @click="onStop">stop</button>
      <button @click="onStop2">stop2</button>
    </view>
    <view id="answerHistory">
      <history-block :data="history"/>
    </view>
    <view id="submit">
      <form @submit="onSubmitAnswer">
          <input name="myAnswer" adjust-position="false" maxlength="20"/>
        <button form-type="submit" >发送</button>
      </form>
    </view>
<!--    临时测试用-->
    <button @click="onEnd"> end </button>
<!--    <button @click="pop">test</button>-->
<!--    <uni-pop-up ref="popup" type="message">-->
<!--      <uni-pop-up-message type="warning" :message="popMessage"/>-->
<!--    </uni-pop-up>-->
  </view>
</template>

<script>
import uniPopUp from '@dcloudio/uni-ui/lib/uni-popup/uni-popup.vue';
import uniPopUpMessage from '@dcloudio/uni-ui/lib/uni-popup-message/uni-popup-message.vue';
import Subject from "@/components/Subject";
import History from "@/components/History";
import CountDown from "@/components/CountDown";
export default {
  name: "GamePage",
  // onReady: function(){
  //   uni.onSocketMessage( (res)=>{
  //         let message = JSON.parse(res.data)
  //
  //       }
  //
  //   )
  // },
  components: {
    "subject-block": Subject,
    "history-block": History,
    "count-down": CountDown,
    // uniPopUp,
    // uniPopUpMessage,
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
      //shot countDown
      active1: true,
      //long countDown
      active2: false,
      current1: 0,
      current2: 0,
      side: 0,
      myExtraTime: 30,
      sideExtraTime: 30,
      info: 'counting...',
      popMessage: 'warning',
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
      subject: "春花秋月何时了/2",
      answer:  [
        {
          'value':'长',
          'highlight': 1,
        },
        {
          'value': '风',
          'highlight': 0,
        },
        {
          'value': '万',
          'highlight': 0,
        },
        {
          'value': '里',
          'highlight': 0,
        },
        {
          'value': '送',
          'highlight': 0,
        },
        {
          'value': '秋',
          'highlight': 0,
        },
        {
          'value': '雁',
          'highlight': 2,
        },
      ],
      history: [
      [
            {
              'value':'长',
              'highlight': 1,
            },
            {
              'value': '风',
              'highlight': 0,
            },
            {
              'value': '万',
              'highlight': 0,
            },
            {
              'value': '里',
              'highlight': 0,
            },
            {
              'value': '送',
              'highlight': 0,
            },
            {
              'value': '秋',
              'highlight': 0,
            },
            {
              'value': '雁',
              'highlight': 1,
            },
          ]
      ]
    }
  },
  methods:{
    // pop(){
    //   this.$refs.popup.open()
    // },
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
      this.sendMessage({
        'type': 'timeout'
      })
    },
    //change count down side
    changeSide(){
      console.log(this.current1)
      this.current1 = 10
      if(this.side) {
        this.side = 0
        this.current2 = this.myExtraTime
      }
      else {
        this.side = 1
        this.current2 = this.sideExtraTime
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
        this.sideExtraTime = val
      }
      else{
        this.myExtraTime = val
      }
    },
    onEnd(){
      uni.redirectTo({
        'url': '/pages/EndPage/EndPage'
      })
    },
    onSubmitAnswer(e){
      this.sendMessage({
            'type': 'answer',
            'text': e.detail.value.myAnswer
          }
      )
      //计时器停止
      this.active1 = false
      this.active2 = false
    },

    onMessage(msg){
      switch (msg.type){
        case 'game_status': {
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
    width: 20%;
  }
</style>