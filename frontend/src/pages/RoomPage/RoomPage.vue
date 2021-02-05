<template>
  <view>
    <button @click="sitDown"> 座位1 </button>
    <button @click="sitDown"> 座位2 </button>
    <button @click="onStart"> 开始游戏 </button>
  </view>
</template>

<script>
export default {
name: "RoomPage",
  methods: {
    sitDown(){

    },
    onStart(){
        uni.redirectTo({
          url: "/pages/ChoosePage/ChoosePage"
        })
      }
  }
}
</script>

<style scoped>

</style>