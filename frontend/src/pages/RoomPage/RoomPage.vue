<template>
  <view>
    <image src="https://flyhana.starrah.cn/static/room.png" class="background"></image>
    <view style="width: 100%; height: 100%; padding-top: 20%">
      <template v-if='!connected'>
        {{ status }}
      </template>
      <view v-else class="center">
        <view style="margin: 10px 0">
        <view style="display: inline-block">
          <view style="float:left">
            <image class="circle" src="https://flyhana.starrah.cn/static/picture.png" mode="widthFix"></image>
            <p style="font-size: 12px; color: #666666">{{ host }}</p>
          </view>
          <p style="float:left; margin-top: 15px;" class="status">{{ hostStatus }}</p>
        </view>
        </view>
        <view>
        <view style="display: inline-block">
          <view style="float:left">
            <image class="circle" src="https://flyhana.starrah.cn/static/picture1.jpg" mode="widthFix"></image>
            <p style="font-size: 12px; color: #666666">{{ guest }}</p>
          </view>
          <p style="float:left; margin-top: 15px;" class="status">{{ guestStatus }}</p>
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
              <button @click="startGame" :disabled='hostStatus === "ready" && guest !== ""' class="btn2">开始游戏</button>
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

      host: 'kuriko',
      hostStatus: '已准备',
      guestStatus: "已准备",
      guest: 'pisces',
    };
  },
  onLoad() {
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
  margin-right: 12px;
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
