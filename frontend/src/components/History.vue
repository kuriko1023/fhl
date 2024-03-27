<template>
  <view class="history">
    <ul style="padding: 10px">
      <li v-for="item in data" :key="item">
        <view  class="center">
          <view style='line-height: 0; margin-bottom: 0.5ex'>
            <view style="display: inline-block; line-height: 0.5">
              <view v-for="word in item" :key="word.word" style="float: left">
                <p class="kati" :style="{color: wordColor[word.highlight]}">{{word.word}}</p>
              </view>
            </view>
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
        2: '#5e270d',
      }
    }
  }
}
</script>

<style scoped>
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}

.kati {
  font-family: Kai;
  font-size: 21px;
  padding: 5px 1px;
}

.center {
  text-align: center;
}
</style>
