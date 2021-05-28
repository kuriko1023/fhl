<template>
  <view>
    <view class="modeChoose">
      <radio-group @change="onModeChange">
        <label><radio value="A" :checked="mode==='A'?true:false"/>单字飞花</label>
        <label><radio value="B" :checked="mode==='B'?true:false"/>多字飞花</label>
        <label><radio value="C" :checked="mode==='C'?true:false"/>超级飞花</label>
        <label><radio value="D" :checked="mode==='D'?true:false"/>谜之飞花</label>
      </radio-group>
    </view>
    <view v-if="mode==='B'||mode==='C'||mode==='D'">
      <picker :range="range[mode]" @change="onSizeChange" :value="picker">
        <view>
          {{range[mode][picker]}}
        </view>
      </picker>
    </view>
    <view>
      <button @click="generate">{{isSubject?'换一换':'生成题目'}}</button>
    </view>
    <view v-if="isSubject">
      <subject-block :mode="mode" :text="subject"></subject-block>
    </view>
    <view v-else>
      <text>
        {{rule[mode]}}
      </text>
    </view>
    <view>
      <button @click="onConfirm">
        确定
      </button>
    </view>
  </view>
</template>

<script>
import Subject from "@/components/Subject";
export default {
name: "ChoosePage",
  components:{
    "subject-block": Subject,
  },
  data(){
    return{
      mode: 'A',
      picker: 0,
      range1: ['5', '6', '7', '8', '9'],
      range: {
        'B': ['5', '6', '7', '8', '9'],
        'C': ['1-10', '3-16'],
        'D': ['5', '6', '7', '8', '9', '10']
      },
      size: 0,
      isSubject: false,
      subject: '古 梦 雁/长 舟 送 寄 事 神 不 生 西风 多少/1000010011',
      rule: {
        'A': "不会吧不会吧不会有人不会玩单字飞花令吧（xxx",
        'B': "不会吧不会吧不会有人不会玩多字飞花令吧（xxx",
        'C': "不会吧不会吧不会有人不会玩超级飞花令吧（xxx",
        'D': "不会吧不会吧不会有人不会玩谜之飞花令吧（xxx",
      }
    }
  },
  methods:{
    sendChoice(){
      this.sendMessage({
        'type': 'set_mode',
        'mode': this.mode,
        'size': (this.size)[0],
      })
    },
    onModeChange(e){
      this.mode = e.detail.value
      this.picker = 0
      if(this.mode !== 'A') {
        this.size = this.range[this.mode][this.picker]
      }
      this.isSubject = false
      this.sendChoice()
    },
    onSizeChange(e){
      this.picker = e.detail.value
      this.size = this.range[this.mode][this.picker]
      this.isSubject = false
      this.sendChoice()
    },
    onConfirm(){
      uni.navigateTo({
        'url': '/pages/GamePage/GamePage'
      })
    },
    generate(){
      //测试用
      if(!this.isSubject){
        this.isSubject = true
      }
      this.sendMessage({
        'type': 'generate',
        'mode': this.mode,
        'size': (this.size)[0],
      })
    },
    onMessage(msg){
      //TODO:更新题目，与选择器绑定
      if(msg.type === 'generated'){
        this.mode = msg.mode
        this.size = msg.size
        for(let i = 0; i < this.range[this.mode].length; i++){
          if(this.range[this.mode][i] == this.size){
            this.picker = i
          }
        }
        if(msg.subject !== null){
          this.subject = msg.subject
        }
        //显示题目
        if(!this.isSubject){
          this.isSubject = true
        }
      }
    }
  }
}
</script>

<style scoped>

</style>