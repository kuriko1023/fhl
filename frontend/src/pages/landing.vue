<template>
  <view id="background">
    <image v-if='profileInitialized' class="startButton" @click="enterPressed" :src="staticRes('game_start.png')"></image>
    <image class="background" src="/static/start_background_scaled.jpg" ></image>
    <image class="background" :src="staticRes('start_background.png')" ></image>
    <text class="kai info">
      版本 0.2.1 (2024 年 3 月 27 日)
    </text>
  </view>
</template>

<script>
export default {
name: "LandingPage",
  data: () => ({
    profileInitialized: false,
  }),
  onLoad() {
    this.retrieveServerProfile(() => this.profileInitialized = true)
    if (window.location.search.startsWith('?room=')) {
      uni.navigateTo({
        url: "/pages/room"
      })
    }
  },
  onShow() {
    this.closeSocket()
  },
  methods: {
    enterPressed() {
      this.requestLocalProfile(() => this.enterRoom())
    },
    enterRoom() {
      getApp().globalData.myRoom = true;
      uni.navigateTo({
        url: "/pages/room"
      })
    }
  }
}
</script>

<style scoped>
.startButton{
  position: absolute;
  width: 43%;
  height: 8%;
  z-index: 1;
  top: 87%;
  left: 55%;
}
.info {
  color: #422;
  font-size: 12px;
  position: absolute;
  left: 6px;
  bottom: 6px;
}
</style>

<style>
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}
</style>
