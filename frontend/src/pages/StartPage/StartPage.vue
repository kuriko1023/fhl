<template>
  <view id="background">
    <button plain="true" @click="onEnter" class="startButton">
      <image src="/static/game_start.png"></image>
    </button>
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
.background{
  position: absolute;
  height: 100%;
  width: 100%;
}
.startButton{
  position: relative;
  width: 10%;
  height: 10%;
}
</style>