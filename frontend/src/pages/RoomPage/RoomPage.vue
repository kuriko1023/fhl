<template>
  <view>
    <image src="https://flyhana.starrah.cn/static/room.png" class="background"></image>
    <view style="width: 100%; height: 100%; padding-top: 20%">
      <template v-if='!connected'>
        {{ status }}
      </template>
      <view v-else class="center">
        <view style="margin: 10px 0">
        {{ room }}<br><br>
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center;">
            <image class="circle" src="https://flyhana.starrah.cn/static/picture.png" mode="widthFix"></image>
            <p style="font-size: 12px; color: #666666">{{ host }}</p>
          </view>
          <p style="position: relative; margin-left: 100px; margin-top: 15px;" class="status">{{ hostStatus === 'ready' ? '已准备' : hostStatus === 'present' ? '在线' : '离线' }}</p>
        </view>
        </view>
        <view>
        <view style="display: inline-block">
          <view style="float: left; width: 100px; text-align: center;">
            <image class="circle" :src="'https://flyhana.starrah.cn/static/' + (guest !== '' ? 'picture1.jpg' : 'grey_avatar_132.jpg')" mode="widthFix"></image>
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
              <button @click="startGame" :disabled='!(hostStatus === "ready" && guest !== "")' class="btn2">开始游戏</button>
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
      connected: false,
      status: '- 状态 -',

      room: '自己的房间',

      host: '',
      hostStatus: '',
      guest: '',
    };
  },
  onLoad() {
    uni.showShareMenu({
      path: '/ABCDEFG',
    });

    const firstLoad = getApp().globalData.firstLoad;
    if (!firstLoad) {
      getApp().globalData.firstLoad = true;
      const room = uni.getEnterOptionsSync().query.room;
      this.room = room;
    }

    uni.connectSocket({
      url: 'wss://flyhana.starrah.cn/channel/my/!kuriko1023',
      success: () => {
        // setTimeout(() => this.connected = true, 1000)
        this.connected = true;
        this.registerSocketMessageListener();
      },
      fail: () => {
        this.status = '连接失败';
      },
    });
    getApp().globalData.isHost = true;
  },
  onShareAppMessage (res) {
    return {
      title: '分享标题',
      path: '/pages/RoomPage/RoomPage?room=别人的房间',
      imageUrl: 'https://flyhana.starrah.cn/static/tianzige.png',
    };
  },
  methods: {
    onSocketMessage() {
      const msg = this.popSocketMessage(['room_status', 'start_generate']);
      if (msg.type === 'room_status') {
        this.host = msg.host;
        this.hostStatus = msg.host_status;  // absent, present, ready
        this.guest = (msg.guest || '');
      } else if (msg.type === 'start_generate') {
        uni.redirectTo({
          url: "/pages/ChoosePage/ChoosePage"
        })
      }
    },
    sitDown() {
      this.sendSocketMessage({type: 'ready'});
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
