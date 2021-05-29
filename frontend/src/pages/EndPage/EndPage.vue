<template>
  <view>
    <view>
      <subject-block :text="subject" :mode="mode"/>
    </view>
    <view>
      <history-block :data="history"/>
    </view>
    <button @click="onBack">返回</button>
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

</style>