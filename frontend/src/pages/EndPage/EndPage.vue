<template>
  <view>
    <image class="background" :src=gameBackgroundImage />
    <image class="background" :src="staticRes('game_background.png')" />
    <view>
      <view style="text-align: center;">
<!--        <image mode="widthFix" class="result_background" src="/static/result.png"></image>-->
        <image mode="widthFix" class="result"
          :src="(win === 1 ? victoryImage : win === 0 ? tieImage : loseImage)" />
      </view>
      <view style="text-align: center; margin: 10px 0 5px 0; padding: 0 30px">
        <view >
          <image class="circle" style="width: 24px; height: 24px" :src="hostAvatar" mode="widthFix"></image>
          <text style="font-size: 12px; color: #666666">{{ host }}</text>
        </view>
        <image class="vs" :src=vsImage mode="widthFix"></image>
        <view >
          <image class="circle" style="width: 24px; height: 24px" :src="guestAvatar" mode="widthFix"></image>
          <text style="font-size: 12px; color: #666666">{{ guest }}</text>
        </view>
      </view>
      <image :src=historyBackgroundImage class="history_background"></image>
      <image :src="staticRes('history_background.png')" class="history_background"></image>
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
import {
  G,
  staticRes,
  redirect,
  registerSocketMessageListener,
  popSocketMessage,
  historySentenceParse,
  parseSubject,
} from '../../utils';

import gameBackgroundImage from '../../static/game_background_scaled.jpg';
import historyBackgroundImage from '../../static/history_background_scaled.jpg';
import victoryImage from '../../static/victory.png';
import tieImage from '../../static/tie.png';
import loseImage from '../../static/lose.png';
import vsImage from '../../static/vs.png';

import Subject from "../../components/Subject";
import History from "../../components/History";

export default {
  name: "EndPage",
  components:{
    "subject-block": Subject,
    "history-block": History
  },
  data() {
    return{
      gameBackgroundImage,
      historyBackgroundImage,
      victoryImage, tieImage, loseImage,
      vsImage,

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
    this.host = G.host;
    this.guest = G.guest;
    this.hostAvatar = G.hostAvatar;
    this.guestAvatar = G.guestAvatar;
    registerSocketMessageListener({ onSocketMessage: this.onSocketMessage });
  },
  methods:{
    staticRes,
    onBack(){
      redirect('/pages/RoomPage/RoomPage')
    },
    onSocketMessage() {
      const msg = popSocketMessage('end_status');
      if (msg._none) return;

      this.mode = msg.mode;
      this.subject = parseSubject(msg.mode, msg.subject);
      this.history = msg.history.map((item) => historySentenceParse(item));
      this.win = msg.winner;
    },
  }
}

</script>

<style>
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
