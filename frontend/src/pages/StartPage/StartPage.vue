<template>
  <view>
    <button plain="true" @click="onEnter">开始</button>
<!--      <image src="/static/game_start.png"></image>-->
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

</style>