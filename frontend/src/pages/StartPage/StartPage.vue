<template>
  <view id="background">
    <image v-if='profileInitialized' class="startButton" @click="onEnter" src="https://flyhana.starrah.cn/static/game_start.png"></image>
    <image class="background" src="https://flyhana.starrah.cn/static/start_background.png" ></image>
  </view>
</template>

<script>
export default {
name: "StartPage",
  data: () => ({
    profileInitialized: false,
  }),
  onLoad() {
    this.retrieveServerProfile(() => this.profileInitialized = true)
  },
  methods: {
    onEnter(){
      this.requestLocalProfile(() => this.enterRoom())
    },
    enterRoom() {
      getApp().globalData.myRoom = true;
      uni.navigateTo({
        url: "/pages/RoomPage/RoomPage"
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
</style>

<style>
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}
</style>
