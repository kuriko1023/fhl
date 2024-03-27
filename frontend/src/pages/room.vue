<!-- http://localhost:8080/?uid=u1#/pages/room -->
<template>
  <view>
    <image src="/static/room_background_scaled.jpg" class="background"></image>
    <image :src="staticRes('room_background.png')" class="background"></image>
    <view style="width: 100%; height: 100%; padding-top: 20%">
      <view v-if='!connected' class='center'
        style='width: 80%; height: 60px; position: relative'>
        <view class='vertical-center' style='width: 32px; height: 32px; left: 30%'>
          <!-- Dave Gandy, https://www.flaticon.com/authors/dave-gandy -->
          <image src='/static/spinner-of-dots.png'
            class='spinning'
            style='width: 100%; height: 100%' />
        </view>
        <p style='left: 50%' class='vertical-center'>{{ status }}</p>
      </view>
      <view v-else class="center">
        <view style="margin: 10px 0">
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center; position: relative;">
            <p class='badge'>房主</p>
            <image class="circle" :src="hostAvatar" mode="widthFix"></image>
            <p style="font-size: 12px; color: #666666; height: 18px">{{ host }}</p>
          </view>
          <p style="position: relative; margin-left: 100px; margin-top: 15px;" class="status">{{ hostStatus === 'ready' ? '已准备' : hostStatus === 'present' ? '在线' : '离线' }}</p>
        </view>
        </view>
        <view>
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center; position: relative;">
            <p class='badge'>客人</p>
            <image class="circle" :src="guestAvatar" mode="widthFix"></image>
            <p style="font-size: 12px; color: #666666; height: 18px">{{ guest !== '' ? guest : '客人' }}</p>
          </view>
          <p style="position: relative; margin-left: 100px; margin-top: 15px;" class="status">{{ guest !== '' ? '已准备' : '未进入' }}</p>
        </view>
        </view>
<!--        <p>客人：{{ guest }}</p>-->
      </view>
      <view v-if="isHost && connected" style='text-align: center; font-size: 14px; margin-top: 4ex; line-height: 1.6'>
        点击右上角「…」按钮<br>邀请好友加入房间
      </view>
      <view class="bottom">
        <uni-row>
          <uni-col :span="12">
            <view>
              <button @click="sitDown" class="btn1">坐下</button>
            </view>
          </uni-col>
          <uni-col :span="12">
            <view>
              <button @click="startGame" :disabled='!(hostStatus === "ready" && guest !== "" && isHost)' class="btn2">开始游戏</button>
            </view>
          </uni-col>
        </uni-row>
      </view>
    </view>
  </view>
</template>

<script>
import { apiServer, wsServer } from 'utils';
export default {
  name: "RoomPage",
  data() {
    return {
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
    this.retrieveServerProfile(() => {
      this.status = '连接房间';

      if (getApp().globalData.myRoom) {
        delete getApp().globalData.myRoom;
        this.room = getApp().globalData.my.id;
      } else {
        let room = 'my';
        if (uni.getSystemInfoSync().uniPlatform === 'mp-weixin') {
          room = uni.getEnterOptionsSync().query.room;
        } else {
          if (window.location.search.startsWith('?room='))
            room = window.location.search.substring('?room='.length);
        }
        this.room = room;
      }

      const urlPromise = () => new Promise((resolve, reject) => {
        this.adaptedLogin({
          success: (code) => resolve(
            `${wsServer}/channel/${this.room}/${code}`,
          ),
          fail: () => reject(),
        });
      });
      this.connectSocket({
        url: urlPromise,
        success: () => {
          this.registerSocketMessageListener();
        },
        fail: () => {
          this.status = '连接失败';
        },
      });

      this.isHost = getApp().globalData.isHost = (this.room === getApp().globalData.my.id);
    });
  },
  onShareAppMessage (res) {
    return {
      title: '一起来玩飞花令吧',
      path: '/pages/RoomPage/RoomPage?room=' + this.room,
      imageUrl: '/static/start_background.jpg',
    };
  },
  methods: {
    onSocketMessage() {
      if (this.tryPeekSocketMessage('generated') ||
          this.tryPeekSocketMessage('generate_wait')) {
        uni.redirectTo({
          url: "/pages/subject"
        })
        return
      } else if (this.tryPeekSocketMessage('game_status')) {
        uni.redirectTo({
          url: "/pages/game"
        })
        return
      }

      const msg = this.popSocketMessage('room_status');
      if (msg.type === 'room_status') {
        if (getApp().globalData.my.create) {
          // 需要发送个人信息
          delete getApp().globalData.my.create
          this.sendProfileUpdate()
          return; // 下次收到消息时更新
        }
        this.connected = true;
        this.host = msg.host;
        this.hostAvatar = `${apiServer}/avatar/${msg.host_avatar}`;
        this.hostStatus = msg.host_status;  // absent, present, ready
        this.guest = (msg.guest || '');
        this.guestAvatar = `${apiServer}/avatar/${msg.guest_avatar || ''}`;
        getApp().globalData.host = this.host;
        getApp().globalData.guest = this.guest;
        getApp().globalData.hostAvatar = this.hostAvatar;
        getApp().globalData.guestAvatar = this.guestAvatar;
      }
    },
    sendProfileUpdate() {
      this.sendSocketMessage({
        type: 'profile',
        nickname: getApp().globalData.my.nickname,
        avatar: getApp().globalData.my.avatar,
      })
    },
    sitDown() {
      this.requestLocalProfile(() => this.sendSocketMessage({type: 'ready'}));
    },
    startGame() {
      this.sendSocketMessage({type: 'start_generate'});
      uni.redirectTo({
        url: "/pages/subject"
      })
    }
  }
}
</script>

<style scoped>
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}
.status{
  font-size: 26px;
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
  bottom: 8vh;
  left: calc(50% - 30vh);
  width: 100%;
  max-width: 60vh;
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
