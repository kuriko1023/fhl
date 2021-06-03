<template>
  <view>
    <image class="background" src="https://flyhana.starrah.cn/static/game_background.png" ></image>
    <view>
      <view style="text-align: center;">
<!--        <image mode="widthFix" class="result_background" src="/static/result.png"></image>-->
        <image mode="widthFix" class="result" :src=" win === 1 ? 'https://flyhana.starrah.cn/static/victory.png' : win === 0 ? 'https://flyhana.starrah.cn/static/tie.png' : 'https://flyhana.starrah.cn/static/lose.png' "></image>
      </view>
      <view style="text-align: center; margin: 10px 0 5px 0; padding: 0 30px">
        <uni-row>

          <uni-col :span="8">
        <view >
        <image class="circle" :src="hostAvatar" mode="widthFix"></image>
        <p style="font-size: 12px; color: #666666">{{ host }}</p>
        </view>
          </uni-col>
          <uni-col :span="8">
           <image class="vs"  src="https://flyhana.starrah.cn/static/vs.png" mode="widthFix"></image>
<!--            <p style="font-family: 'STKaiti'; font-size: 24px; font-weight: bold">对战</p>-->
          </uni-col>
          <uni-col :span="8">
        <view >
          <image class="circle" :src="guestAvatar" mode="widthFix"></image>
          <p style="font-size: 12px; color: #666666">{{ guest }}</p>
        </view>
          </uni-col>
        </uni-row>
      </view>
      <image src="https://flyhana.starrah.cn/static/history_background.png" class="history_background"></image>
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
      host: '',
      guest: '',
      hostAvatar: '',
      guestAvatar: '',
    }
  },
  onLoad() {
    this.host = getApp().globalData.host;
    this.guest = getApp().globalData.guest;
    this.hostAvatar = getApp().globalData.hostAvatar;
    this.guestAvatar = getApp().globalData.guestAvatar;
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
      this.win = msg.winner;
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
.result_background{
  width: 35%;
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
  height: 65%;
  width: 88%;
  margin: 0 6% 10px 6%;
  overflow: scroll;
}
.history_background{
  position: absolute;
  height: 65%;
  width: 88%;
  margin: 0 6% 10px 6%;
}
.circle {
  width: 35px;
  border-radius: 50%;
}
.vs{
  width: 60px;
}
</style>
