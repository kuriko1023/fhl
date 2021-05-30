<template>
  <view>
    <template v-if='!connected'>
      {{ status }}
    </template>
    <template v-else>
      <p>房主：{{ host }}</p>
      <p>房主状态 {{ hostStatus }}</p>
      <button @click="sitDown">坐下</button>
      <button @click="startGame">开始游戏</button>
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
  },
  methods: {
    onSocketMessage() {
      console.log('on message!');
      console.log('peeked message ', this.peekSocketMessage());
    },
    sitDown() {
      uni.sendSocketMessage({
        data: {
          type: 'ready',
        },
      });
    },
    startGame() {
      uni.redirectTo({
        url: "/pages/ChoosePage/ChoosePage"
      })
    }
  }
}
</script>

<style scoped>

</style>
