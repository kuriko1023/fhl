<template>
  <view>
    <view>
      <subject-block :data="subject" mode="C"/>
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
      subject: {
        "subject1": ["雁", "古", "梦"],
        "subject2": "长 舟 送 寄 事 神 不 生 西风 多少",
      },
      history: [
        {
          "value": "长风万里送秋雁"
        },
        {
          "value": "千古兴亡多少事"
        }
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