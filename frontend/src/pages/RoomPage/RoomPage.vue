<template>
  <view>
    <template v-if='!connected'>
      {{ status }}
    </template>
    <template v-else>
      <p>房主：{{ host }}</p>
      <p>房主状态：{{ hostStatus }}</p>
      <p>客人：{{ guest }}</p>
      <button @click="sitDown">坐下</button>
      <template v-if='hostStatus === "ready" && guest !== ""'>
        <button @click="startGame">开始游戏</button>
      </template>
    </template>
  </view>
</template>

<script>
export default {
  name: "RoomPage",
  data() {
    return {
      connected: false,
      status: '- 状态 -',

      host: '',
      hostStatus: '',
      guest: '',
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

</style>
