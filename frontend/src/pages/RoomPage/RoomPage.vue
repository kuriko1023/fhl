<template>
  <view>
    <image src="https://flyhana.starrah.cn/static/room.png" class="background"></image>
    <view style="width: 100%; height: 100%; padding-top: 20%">
      <template v-if='!connected'>
        {{ status }}
      </template>
      <view v-else class="center">
        <view style="margin: 10px 0">
        room = [{{ room.substr(0, 8) }}]<br><br>
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center;">
            <image class="circle" :src="hostAvatar" mode="widthFix"></image>
            <p style="font-size: 12px; color: #666666">{{ host }}</p>
          </view>
          <p style="position: relative; margin-left: 100px; margin-top: 15px;" class="status">{{ hostStatus === 'ready' ? '已准备' : hostStatus === 'present' ? '在线' : '离线' }}</p>
        </view>
        </view>
        <view>
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center;">
            <image class="circle" :src="guestAvatar" mode="widthFix"></image>
            <p style="font-size: 12px; color: #666666">{{ guest !== '' ? guest : '客人' }}</p>
          </view>
          <p style="position: relative; margin-left: 100px; margin-top: 15px;" class="status">{{ guest !== '' ? '已准备' : '未进入' }}</p>
        </view>
        </view>
<!--        <p>客人：{{ guest }}</p>-->
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
        const room = uni.getEnterOptionsSync().query.room;
        this.room = room;
      }

      uni.login({success: (res) => this.connectSocket({
        // url: 'wss://flyhana.starrah.cn/channel/my/!kuriko1023',
        url: 'wss://flyhana.starrah.cn/channel/' + this.room + '/' + res.code,
        success: () => {
          // setTimeout(() => this.connected = true, 1000)
          this.connected = true;
          this.registerSocketMessageListener();
        },
        fail: () => {
          this.status = '连接失败';
        },
      })});

      this.isHost = getApp().globalData.isHost = (this.room === getApp().globalData.my.id);
    });
  },
  onUnload() {
    this.closeSocket();
  },
  onShareAppMessage (res) {
    return {
      title: '分享标题',
      path: '/pages/RoomPage/RoomPage?room=' + this.room,
      imageUrl: 'https://flyhana.starrah.cn/static/tianzige.png',
    };
  },
  methods: {
    onSocketMessage() {
      const msg = this.popSocketMessage(['room_status', 'start_generate']);
      if (msg.type === 'room_status') {
        if (getApp().globalData.my.create) {
          // 需要发送个人信息
          delete getApp().globalData.my.create
          this.sendProfileUpdate()
          return; // 下次收到消息时更新
        }
        this.host = msg.host;
        this.hostAvatar = 'https://flyhana.starrah.cn/avatar/' + msg.host_avatar;
        this.hostStatus = msg.host_status;  // absent, present, ready
        this.guest = (msg.guest || '');
        this.guestAvatar = 'https://flyhana.starrah.cn/avatar/' + (msg.guest_avatar || '');
      } else if (msg.type === 'start_generate') {
        uni.redirectTo({
          url: "/pages/ChoosePage/ChoosePage"
        })
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
        url: "/pages/ChoosePage/ChoosePage"
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
</style>
