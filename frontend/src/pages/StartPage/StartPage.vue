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
