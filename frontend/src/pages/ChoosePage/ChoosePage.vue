<template>
  <view>
    <view class="modeChoose">
      <radio-group @change="onModeChange">
        <radio value="A"/>
        <radio value="B"/>
        <radio value="C"/>
        <radio value="D"/>
      </radio-group>
    </view>
    <view>
      <button @click="onConfirm">
        确定
      </button>
    </view>
  </view>
</template>

<script>
export default {
name: "ChoosePage",
  methods:{
    onModeChange(e){
      console.log(e.detail.value)
    },
    onConfirm(){
      uni.navigateTo({
        'url': '/pages/GamePage/GamePage'
      })
    }
  }
}
</script>

<style scoped>

</style>