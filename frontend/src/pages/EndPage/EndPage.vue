<template>
  <view>
    <image class="background" src="/static/game_background.png" ></image>
    <view>
      <view style="text-align: center;">
      <text class="result">胜利</text>
      </view>
      <image src="/static/history_background.png" class="history_background"></image>
    <view  class="info">
    <view>
      <subject-block :text="subject" :mode="mode"/>
    </view>
    <view>
      <history-block :data="history"/>
    </view>
    </view>
<!--    <button @click="onBack" class="btn1">返回</button>-->
  </view>
  </view>
</template>

<script>
import Subject from "@/components/Subject";
import History from "@/components/History";
export default {
name: "EndPage",
  components:{
    "subject-block": Subject,
    "history-block": History
  },
  data() {
    return{
      mode: "C",
      subject: "古 梦 雁/长 舟 送 寄 事 神 不 生 西风 多少",
      // mode: "B",
      // subject: "春花秋月何时了/2",
      history: [
        [
          {
            'value':'长',
            'highlight': 1,
          },
          {
            'value': '风',
            'highlight': 0,
          },
          {
            'value': '万',
            'highlight': 0,
          },
          {
            'value': '里',
            'highlight': 0,
          },
          {
            'value': '送',
            'highlight': 0,
          },
          {
            'value': '秋',
            'highlight': 0,
          },
          {
            'value': '雁',
            'highlight': 1,
          },
        ]
      ]
    }
  },
  methods:{
    onBack(){
      uni.redirectTo({
        'url': '/pages/RoomPage/RoomPage'
      })
    }
  }
}

</script>

<style scoped>
.result{
  font-size: 54px;
  font-family:"书体坊向佳红毛笔行书";
  margin-top: 10px;
}
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}

.btn1{
  background-color: #689a74;
  color: white;
  border-radius: 10px;
  margin-left: 18%;
  margin-right:10%;
  font-size: 14px;
}

.info{
  position: absolute;
  height: 100%;
  width: 88%;
  margin: 0 6% 10px 6%;
  overflow: scroll;
}
.history_background{
  position: absolute;
  height: 100%;
  width: 88%;
  margin: 0 6% 10px 6%;
}
</style>