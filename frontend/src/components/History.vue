<template>
  <view>
    <ul>
      <li v-for="item in data" :key="item">
        <view style="display: flex">
          <view v-for="word in item" :key="word.value">
            <text :style="{color: wordColor[word.highlight]}">{{word.value}}</text>
          </view>
        </view>
      </li>
    </ul>
  </view>
</template>

<script>
export default {
name: "History",
  props: {
    data: {
      type: Array
    }
  },
  data(){
    return{
      wordColor: {
        0: '#000000',
        1: '#ad2b29',
      }
    }
  }
}
</script>

<style scoped>

</style>