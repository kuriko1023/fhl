<template>
  <view>
    <image class="background" src="https://flyhana.starrah.cn/static/game_background.png" ></image>
    <view class="content">
<!--      <uni-row>-->
<!--        <uni-col :span="12">-->
          <view class="modeChoose">
            <radio-group @change="onModeChange" >
              <label><radio class="theme" value="A" :checked="mode==='A'?true:false"/><span class="radio_text">单字飞花</span></label>
              <label><radio class="theme" value="B" :checked="mode==='B'?true:false"/><span class="radio_text">多字飞花</span></label>
              <label><radio class="theme" value="C" :checked="mode==='C'?true:false"/><span class="radio_text">超级飞花</span></label>
              <label><radio class="theme" value="D" :checked="mode==='D'?true:false"/><span class="radio_text">谜之飞花</span></label>
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
                <picker :range="range[mode]" @change="onSizeChange" :value="picker" mode="selector" >
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
        <uni-row>
          <uni-col :span="12">
          <view>
            <button @click="generate" class="btn1">{{isSubject?'换一换':'生成题目'}}</button>
          </view>
          </uni-col>
          <uni-col :span="12">
          <view>
            <button class="btn2" @click="onConfirm">
              确定
            </button>
          </view>
          </uni-col>
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
        'A': "题目为一个字或词，玩家需轮流说出带有该字（词）的合法诗句。",
        'B': "题目为一句诗句，从该诗句中选取若干个字，玩家依次说出包含当前关键字的合法诗句。  如：\n" +
            "给出的题目诗句为 “春花秋月何时了” , 则玩家依次说出含有春、花、秋、月、何、时的诗句。",
        'C': "题目生成一组固定字词，与一组可消去字词，玩家轮流从固定字词与可消去字词中各选择一个，说出同时含有两者的诗句，即可消去一个字词。\n" +
            "  如：\n" +
            "    给出的题目为 固定字词：【山 天】 可消去字词：【水、青、红、明月、松】， 玩家回答 【斜阳却在青山外】后，即可消去字词 青。",
        'D': "题目生成两组可消去字词。玩家轮流从两组可消去字词中各选择一个，说出同时含有两者的诗句，即可分别消去字词\n" +
            "  如：\n" +
            "  给出的题目为 【三、何】【二十、月】，玩家回答【江月何年初照人】，即可消去字词何， 月。",
      }
    }
  },
  onLoad() {
    this.registerSocketMessageListener();
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
      this.sendSocketMessage({
        'type': 'generate',
        'mode': this.mode,
        'size': (this.mode === 'A' ? 0 : this.rangeValue[this.mode][this.picker]),
      })
    },
    onSocketMessage() {
      if (this.peekSocketMessage().type === 'game_status') {
        uni.navigateTo({
          url: '/pages/GamePage/GamePage'
        });
        return;
      }

      const msg = this.popSocketMessage('generated');
      if (msg._none) return;
      //更新题目，与选择器绑定
      this.mode = msg.mode
      if(this.mode !== 'A') {
        for(let i = 0; i < this.range[this.mode].length; i++){
          if(this.range[this.mode][i] == msg.size){
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
      if(!this.isSubject){
        this.isSubject = true
      }
    }
  }
}
</script>

<style scoped>
  .content{
    position: absolute;
    top: 0;
    left: 0;
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
    /*width: 50%;*/
    margin: 0 10px 0 10px;
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
  }
  .tip{
    font-size: 14px;
    color: #366440;
    font-weight: bold;
  }
  .content {
    margin: 20px;
  }
  .content_background {
    background-color: #fdf8ed;
    border: 1.3px solid #975f5b;
    border-radius: 10px;
    padding: 10px;
    horiz-align: center;
    position: absolute;
    min-height: 100px;
    margin: 5px 10px;
    font-size: 14px;
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

