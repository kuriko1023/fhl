<template>
  <view>
    <image :src=backgroundImage class="background"></image>
    <image :src="staticRes('room_background.png')" class="background"></image>
    <view style="width: 100%; height: 100%; padding-top: 20%">
      <view v-if='!connected' class='center'
        style='width: 80%; height: 60px; position: relative'>
        <view class='vertical-center' style='width: 32px; height: 32px; left: 30%'>
          <!-- Dave Gandy, https://www.flaticon.com/authors/dave-gandy -->
          <image :src=spinnerImage
            class='spinning'
            style='width: 100%; height: 100%' />
        </view>
        <text style='left: 50%' class='vertical-center'>{{ status }}</text>
      </view>
      <view v-else class="center">
        <view style="margin: 10px 0">
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center; position: relative;">
            <text class='badge'>房主</text>
            <image class="circle" :src="hostAvatar" mode="widthFix"></image>
            <text style="font-size: 12px; color: #666666; height: 18px">{{ host }}</text>
          </view>
          <text style="position: relative; margin-left: 100px; margin-top: 15px;" class="status">{{ hostStatus === 'ready' ? '已准备' : hostStatus === 'present' ? '在线' : '离线' }}</text>
        </view>
        </view>
        <view>
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center; position: relative;">
            <text class='badge'>客人</text>
            <image class="circle" :src="guestAvatar" mode="widthFix"></image>
            <text style="font-size: 12px; color: #666666; height: 18px">{{ guest !== '' ? guest : '客人' }}</text>
          </view>
          <text style="position: relative; margin-left: 100px; margin-top: 15px;" class="status">{{ guest !== '' ? '已准备' : '未进入' }}</text>
        </view>
        </view>
<!--        <p>客人：{{ guest }}</text>-->
      </view>
      <view v-if="isHost && connected" style='text-align: center; font-size: 14px; margin-top: 4ex'>
        <view>点击右上角「…」按钮</view>
        <view>邀请好友加入房间</view>
      </view>
      <view class="bottom">
        <view style='width: 100%; height: 20px'>
          <button @tap="sitDown" class="btn1">坐下</button>
        </view>
        <view style='width: 100%; height: 20px'>
          <button @tap="startGame" :disabled='!(hostStatus === "ready" && guest !== "" && isHost)' class="btn2">开始游戏</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { ref } from 'vue';
import {
  G,
  staticRes,
  apiServer,
  wsServer,
  redirect,
  retrieveServerProfile,
  requestLocalProfile,
  enterQueryParam,
  getLoginCode,
  connectSocket,
  registerSocketMessageListener,
  sendSocketMessage,
  peekSocketMessage,
  tryPeekSocketMessage,
  tryPopSocketMessage,
  popSocketMessage,
} from '../../utils';

import backgroundImage from '../../static/room_background_scaled.jpg';
import spinnerImage from '../../static/spinner-of-dots.png';
import startBackgroundImage from '../../static/start_background_scaled.jpg';

export default {
  name: "RoomPage",
  data() {
    return {
      backgroundImage,
      spinnerImage,
      startBackgroundImage,

      profileInitialized: false,
      connected: false,
      status: '获取玩家信息',

      room: '',

      host: '',
      hostAvatar: '',
      hostStatus: '',
      guest: '',
      guestAvatar: '',

      isHost: false,
    };
  },
  onLoad() {
    retrieveServerProfile(() => {
      this.status = '连接房间';

      if (G.myRoom) {
        delete G.myRoom;
        this.room = G.my.id;
      } else {
        const room = uni.getEnterOptionsSync().query.room;
        this.room = room;
      }

      const urlPromise = async () =>
        `${wsServer}/channel/${this.room}/${await getLoginCode()}`;
      connectSocket({
        url: urlPromise,
        success: () => {
          registerSocketMessageListener({ onSocketMessage: this.onSocketMessage });
        },
        fail: () => {
          this.status = '连接失败';
        },
      });

      this.isHost = G.isHost = (this.room === G.my.id);
    });
  },
  onShareAppMessage (res) {
    return {
      title: '一起来玩飞花令吧',
      path: '/pages/RoomPage/RoomPage?room=' + this.room,
      imageUrl: startBackgroundImage,
    };
  },
  methods: {
    staticRes,
    onSocketMessage() {
      if (tryPeekSocketMessage('generated')) {
        redirect("/pages/ChoosePage/ChoosePage")
        return
      } else if (tryPeekSocketMessage('game_status')) {
        redirect("/pages/GamePage/GamePage")
        return
      }

      const msg = popSocketMessage('room_status');
      if (msg.type === 'room_status') {
        if (G.my.create) {
          // 需要发送个人信息
          delete G.my.create
          this.sendProfileUpdate()
          return; // 下次收到消息时更新
        }
        this.connected = true;
        this.host = msg.host;
        this.hostAvatar = `${apiServer}/avatar/${msg.host_avatar}`;
        this.hostStatus = msg.host_status;  // absent, present, ready
        this.guest = (msg.guest || '');
        this.guestAvatar = `${apiServer}/avatar/${msg.guest_avatar || ''}`;
        G.host = this.host;
        G.guest = this.guest;
        G.hostAvatar = this.hostAvatar;
        G.guestAvatar = this.guestAvatar;
      }
    },
    sendProfileUpdate() {
      sendSocketMessage({
        type: 'profile',
        nickname: G.my.nickname,
        avatar: G.my.avatar,
      })
    },
    sitDown() {
      requestLocalProfile(() => sendSocketMessage({type: 'ready'}));
    },
    startGame() {
      sendSocketMessage({type: 'start_generate'});
      redirect("/pages/ChoosePage/ChoosePage")
    }
  }
}
</script>

<style>
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}
.status{
  font-size: 26px;
  font-family: 华文行楷;
  font-weight: bold;
  color: #49443d;
}
.circle {
  width: 45px;
  height: 45px;
  border-radius: 50%;
}
.center {
  /*position: absolute;*/
  width: fit-content;
  /*left: 20%;*/
  margin: auto;
}
.vertical-center {
  margin: 0;
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
}
.badge {
  position: absolute;
  left: -1.5em;
  top: -0.5ex;
  font-size: 14px;
  background-color: #84765e;
  color: white;
  padding: 2px 5px;
  border-radius: 4px;
}

.bottom {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  margin-bottom: 8%;
}

.btn2{
  background-color: #84765e;
  color: white;
  border-radius: 10px;
  margin-left: 10%;
  margin-right:18%;
  font-size: 14px;
}
.btn1{
  background-color: #a49b8c;
  color: white;
  border-radius: 10px;
  margin-left: 18%;
  margin-right:10%;
  font-size: 14px;
}

@keyframes spinning {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.spinning {
  animation: spinning 2s linear infinite;
}
</style>
