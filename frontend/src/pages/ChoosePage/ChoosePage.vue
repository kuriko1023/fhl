<template>
  <view>
    <image class="background" src="https://flyhana.starrah.cn/static/game_background.png" ></image>
    <view class="content">
<!--      <uni-row>-->
<!--        <uni-col :span="12">-->
          <view class="modeChoose">
            <radio-group @change="onModeChange" >
              <div><label><radio :disabled="!isHost" class="theme" value="A" :checked="mode==='A'?true:false"/><span class="radio_text">梦笔生花</span></label></div>
              <div><label><radio :disabled="!isHost" class="theme" value="B" :checked="mode==='B'?true:false"/><span class="radio_text">走马观花</span></label></div>
              <div><label><radio :disabled="!isHost" class="theme" value="C" :checked="mode==='C'?true:false"/><span class="radio_text">天女散花</span></label></div>
              <div><label><radio :disabled="!isHost" class="theme" value="D" :checked="mode==='D'?true:false"/><span class="radio_text">雾里看花</span></label></div>
            </radio-group>
          </view>
<!--        </uni-col>-->
<!--        <uni-col :span="6">-->

<!--        </uni-col>-->
<!--          <uni-col :span="6">-->
          <view v-if="mode==='B'||mode==='C'||mode==='D'" class="picker">
            <uni-row>
              <uni-col :span="8">
                <span class="tip" >选择题型：</span>
              </uni-col>
              <uni-col :span="12">
                <picker :disabled="!isHost" :range="range[mode]" @change="onSizeChange" :value="picker" mode="selector" >
                  <view class="picker_btn">
                    {{range[mode][picker]}}
                  </view>
                </picker>
              </uni-col>

            </uni-row>

          </view>
<!--        </uni-col>-->
<!--      </uni-row>-->
      <view >
        <template v-if="isSubject" class="rule">
          <text class="tip" style="margin: 8px 10px 0 10px">题目内容</text>
          <subject-block :mode="mode" :subject="subject"></subject-block>
        </template>
        <template v-else >
          <text class="tip" style="margin: 8px 10px 0 10px">题目规则</text>
          <view class="content_background">
            <text>
              {{rule[mode]}}
            </text>
          </view>
        </template>
      </view>
      <view class="bottom">
        <uni-row v-if='isHost'>
          <uni-col :span="12">
          <view>
            <button @click="generate" class="btn1" :disabled="isGenReqPending">{{isSubject?'换一换':'生成题目'}}</button>
          </view>
          </uni-col>
          <uni-col :span="12">
          <view>
            <button class="btn2" @click="onConfirm" :disabled="!isSubject">
              确定
            </button>
          </view>
          </uni-col>
        </uni-row>
        <uni-row v-else>
          <view>
            <button class="btn-full">
              请等待房主选题
            </button>
          </view>
        </uni-row>
      </view>
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
      isHost: false,
      mode: 'A',
      picker: 0,
      range: {
        'B': ['五字', '六字', '七字', '八字', '九字'],
        'C': ['一词-十词', '三词-十六词'],
        'D': ['五词', '六词', '七词', '八词', '九词', '十词']
      },
      rangeValue: {
        'B': [5, 6, 7, 8, 9],
        'C': [1, 3],
        'D': [5, 6, 7, 8, 9, 10],
      },
      size: 0,
      isSubject: false,
      subject: {},
      rule: {
        'A': "题目为一个字或词，玩家需轮流说出带有该字（词）的诗句。\n" +
            "如：\n" +
            "　给出的题目为【霜】，玩家回答【空里流霜不觉飞】即符合规则。",
        'B': "题目为一句诗句，其中的字按顺序依次作为关键字，玩家轮流各自说出包含当前关键字的诗句。\n" +
            "如：\n" +
            "　给出的题目为【春花秋月何时了】，则两个玩家需各自说出含有春、花、秋、月、何、时、了的诗句。",
        'C': "题目为一组固定字词，与一组可消去字词，玩家轮流从固定字词与可消去字词中各选择一个，说出同时含有两者的诗句。每个消去词只能被选择一次。\n" +
            "如：\n" +
            "　给出的题目为固定字词：【山、天】与可消去字词：【水、青、红、明月、松】，玩家回答【斜阳却在青山外】后，即可消去【青】。",
        'D': "题目为两组可消去字词。玩家轮流从两组字词中各选择一个，说出同时含有两者的诗句。所有词都只能被选择一次。\n" +
            "如：\n" +
            "　给出的题目为【三、何】【二十、月】，玩家回答【江月何年初照人】，即可消去【何】【月】。",
      },

      isGenReqPending: false,
    }
  },
  onLoad() {
    this.registerSocketMessageListener();
    this.isHost = getApp().globalData.isHost;
  },
  methods:{
    sendChoice(){
      this.sendSocketMessage({
        'type': 'set_mode',
        'mode': this.mode,
        'size': (this.mode === 'A' ? 0 : this.rangeValue[this.mode][this.picker]),
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
      this.sendSocketMessage({
        type: 'start_game',
      });
    },
    generate(){
      this.isGenReqPending = true;
      this.sendSocketMessage({
        'type': 'generate',
        'mode': this.mode,
        'size': (this.mode === 'A' ? 0 : this.rangeValue[this.mode][this.picker]),
      })
    },
    onSocketMessage() {
      if (this.tryPeekSocketMessage('game_status')) {
        uni.redirectTo({
          url: '/pages/GamePage/GamePage'
        });
        return;
      }

      const msg = this.popSocketMessage('generated');
      if (msg._none) return;
      //更新题目，与选择器绑定
      this.isGenReqPending = false
      this.mode = msg.mode
      if(this.mode !== 'A') {
        for(let i = 0; i < this.range[this.mode].length; i++){
          if(this.rangeValue[this.mode][i] == msg.size){
            this.picker = i
          }
        }
        this.size = this.range[this.mode][this.picker]
      }
      if(msg.subject !== null){
        this.subject = this.parseSubject(msg.mode, msg.subject)
        if (this.mode === 'B')
          this.subject.subject1[0].show = 0;
      }
      //显示题目
      this.isSubject = (msg.subject !== null)
    }
  }
}
</script>

<style scoped>
  .content{
    position: absolute;
    top: 0;
    left: 0;
    padding: 20px;
  }
  .background{
    position: absolute;
    height: 100%;
    width: 100%;
    z-index: -1;
  }
  .radio_text{
    font-size: 14px;
    color: #444444
  }
  .modeChoose{
    margin: 0 10px 10px 10px;
  }
  .modeChoose div {
    display: inline-block;
    width: 50%;
    text-align: center;
  }
  .picker_btn {
    background-color: white;
    padding: 3px 0 3px 10px;
    font-size: 13px !important;
  }
  .btn2{
    background-color: #366440;
    color: white;
    border-radius: 10px;
    margin-left: 10%;
    margin-right:18%;
    font-size: 14px;
  }
  .btn1{
    background-color: #689a74;
    color: white;
    border-radius: 10px;
    margin-left: 18%;
    margin-right:10%;
    font-size: 14px;
    transition: background-color 0.05s ease;
  }
  .btn1[disabled] {
    background-color: #cdc;
    color: white;
  }
  .btn-full{
    background-color: #689a74;
    color: white;
    border-radius: 10px;
    margin-left: 9%;
    margin-right:9%;
    font-size: 14px;
  }
  .tip{
    font-size: 14px;
    color: #366440;
    font-weight: bold;
  }
  .content_background {
    background-color: #fdf8ed;
    border: 1.3px solid #975f5b;
    border-radius: 10px;
    padding: 7px 10px;
    horiz-align: center;
    min-height: 100px;
    margin: 5px 10px;
    font-size: 14px;
    line-height: 1.6;
  }
  .bottom {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    margin-bottom: 5%;
  }
  .picker {
    margin: 10px;
  }
</style>

