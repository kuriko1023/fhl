<template>
  <view>
    <image class="background" src="/static/game_background.png" ></image>
    <view>
      <view style="text-align: center;">
      <image mode="widthFix" class="result" :src=" win === 1 ? '/static/victory.png' : win === 0 ? '/static/tie.png' : '/static/lose.png' "></image>
      </view>
      <image src="/static/history_background.png" class="history_background"></image>
    <view  class="info">
    <view>
      <subject-block :mode="mode" :subject="subject" />
    </view>
    <view>
      <history-block :data="history"/>
    </view>
    </view>
<!--    <button @click="onBack" class="btn1">返回</button>-->
  </view>
  </view>
</template>

<script>
import Subject from "@/components/Subject";
import History from "@/components/History";
export default {
name: "EndPage",
  components:{
    "subject-block": Subject,
    "history-block": History
  },
  data() {
    return{
      mode: '',
      subject: {},
      history: [],
      win: 0,
    }
  },
  onLoad() {
    this.registerSocketMessageListener();
  },
  methods:{
    onBack(){
      uni.redirectTo({
        'url': '/pages/RoomPage/RoomPage'
      })
    },
    onSocketMessage() {
      const msg = this.popSocketMessage('end_status');
      if (msg._none) return;

      this.mode = msg.mode;
      this.subject = this.parseSubject(msg.mode, msg.subject);
      this.history = msg.history.map((item) => this.historySentenceParse(item));
      this.win = msg.win;
    },
  }
}

</script>

<style scoped>
.result{
  width: 30%;
  margin-top: 10px;
  margin-bottom: 10px;
}
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}

.btn1{
  background-color: #689a74;
  color: white;
  border-radius: 10px;
  margin-left: 18%;
  margin-right:10%;
  font-size: 14px;
}

.info{
  position: absolute;
  height: 80%;
  width: 88%;
  margin: 0 6% 10px 6%;
  overflow: scroll;
}
.history_background{
  position: absolute;
  height: 80%;
  width: 88%;
  margin: 0 6% 10px 6%;
}
</style>
