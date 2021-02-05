<template>
  <view>
    <view id="subject">
      <subject-block mode="mode" :data="subject">
      </subject-block>
    </view>
    <view id="answering">
      <image class="circle" src="/static/picture.png" mode="widthFix"></image>
      <view>
        <text>{{answer}}</text>
      </view>
      <progress percent="80" stroke-width="3"/>
      <progress percent="100" stroke-width="3"/>
    </view>
    <view id="answerHistory">
      <history-block :data="history"/>
    </view>
    <view id="submit">
      <form @submit="onSubmit">
        <input name="myAnswer" maxlength="20"/>
        <button form-type="submit">发送</button>
      </form>
    </view>
<!--    临时测试用-->
    <button @click="onEnd"> end </button>
  </view>
</template>

<script>
import Subject from "@/components/Subject";
import History from "@/components/History";
export default {
name: "GamePage",
  components: {
    "subject-block": Subject,
    "history-block": History,
  },
  data() {
    return {
      mode: "C",
      subject: {
        "subject1": ["雁", "古", "梦"],
        "subject2": "长 舟 送 寄 事 神 不 生 西风 多少",
      },
      answer: "长风万里送秋雁",
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
    onSubmit(e){
      console.log(e.details.value.myAnswer)
    },
    onEnd(){
      uni.redirectTo({
        'url': '/pages/EndPage/EndPage'
      })
    }
  }
}
</script>

<style scoped>
  .circle {
    width: 20%;
  }
</style>