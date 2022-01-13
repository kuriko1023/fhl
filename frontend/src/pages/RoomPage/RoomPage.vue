<template>
  <view>
    <image :src=backgroundImage class="background"></image>
    <!--<image :src="staticRes('room_background.png')" class="background"></image>-->
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
  apiServer,
  wsServer,
  redirect,
  retrieveServerProfile,
  requestLocalProfile,
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
  setup() {
    // Data
    const profileInitialized = ref(false);
    const connected = ref(false);
    const status = ref('获取玩家信息');
    const room = ref('');
    const host = ref('');
    const hostAvatar = ref('');
    const hostStatus = ref('');
    const guest = ref('');
    const guestAvatar = ref('');
    const isHost = ref(false);

    // Methods
    const onSocketMessage = () => {
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
        connected.value = true;
        host.value = msg.host;
        hostAvatar.value = `${apiServer}/avatar/${msg.host_avatar}`;
        hostStatus.value = msg.host_status;  // absent, present, ready
        guest.value = (msg.guest || '');
        guestAvatar.value = `${apiServer}/avatar/${msg.guest_avatar || ''}`;
        G.host = host.value;
        G.guest = guest.value;
        G.hostAvatar = hostAvatar.value;
        G.guestAvatar = guestAvatar.value;
      }
    };
    const sendProfileUpdate = () => {
      sendSocketMessage({
        type: 'profile',
        nickname: G.my.nickname,
        avatar: G.my.avatar,
      })
    };
    const sitDown = () => {
      requestLocalProfile(() => sendSocketMessage({type: 'ready'}));
    };
    const startGame = () => {
      sendSocketMessage({type: 'start_generate'});
      console.log('redirecting');
      // redirect("/pages/ChoosePage/ChoosePage");
    };

    // Initialization
    retrieveServerProfile(() => {
      status.value = '连接房间';

      if (G.myRoom) {
        delete G.myRoom;
        room.value = G.my.id;
      } else {
        const room = Taro.getLaunchOptionsSync().query.room;
        room.value = room;
      }

      const wsUrl = async () =>
        `${wsServer}/channel/${room.value}/${await getLoginCode()}`;
      connectSocket({
        url: wsUrl,
        success: () => {
          registerSocketMessageListener({ onSocketMessage });
        },
        fail: () => {
          status.value = '连接失败';
        },
      });

      isHost.value = G.isHost = (room.value === G.my.id);
    });

    return {
      backgroundImage,
      spinnerImage,

      profileInitialized,
      connected,
      status,
      room,
      host,
      hostAvatar,
      hostStatus,
      guest,
      guestAvatar,
      isHost,

      sitDown,
      startGame,
    };
  },
  onShareAppMessage (res) {
    console.log(this.room);
    return {
      title: '一起来玩飞花令吧',
      path: '/pages/RoomPage/RoomPage?room=' + this.room,
      imageUrl: startBackgroundImage,
    };
  },
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
