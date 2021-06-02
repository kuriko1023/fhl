<template>
  <view id="background">
    <image class="startButton" @click="onEnter" src="https://flyhana.starrah.cn/static/game_start.png"></image>
    <image class="background" src="https://flyhana.starrah.cn/static/start_background.png" ></image>
  </view>
</template>

<script>
export default {
name: "StartPage",
  methods: {
    onEnter(){
      console.log("1")

      // 发布时使用 getUserProfile
      // uni.getUserProfile({
      // 测试时使用 getUserInfo 会返回「微信用户」和灰色头像
      uni.getUserInfo({
        desc: '用于向其他玩家展示头像和昵称',
        success: (res) => {
          console.log(res.userInfo.nickName, res.userInfo.avatarUrl)
        },
        fail: () => {
          console.log('getUserProfile() failed')
        },
      })

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
