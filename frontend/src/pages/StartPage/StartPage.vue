<template>
  <view id="background">
    <image class="startButton" @click="onEnter" src="/static/game_start.png"></image>
    <image class="background" src="/static/start_background.png" ></image>
  </view>
</template>

<script>
export default {
name: "StartPage",
  methods: {
    onEnter(){
      console.log("1")
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