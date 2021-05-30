<template>
  <view>
    <template v-if="mode === 'D'">
      <view style="display: flex; flex-wrap:wrap;" class="subjectText">
        <view v-for="item in subject.subject1" :key="item.value">
          <text class="kati" :style="{color: colorD[item.show]}">{{item.value}}</text>
        </view>
      </view>
    </template>
    <template v-if="mode==='B'">
      <view  class="center">
        <view class="subjectText">
          <view style="display: inline-block;" >
            <view v-for="item in subject.subject1" :key="item.value" style="float: left">
              <text class="kati" :style="{color: colorB[item.show]}">{{item.value}}</text>
            </view>
          </view>
        </view>
      </view>
    </template>
    <template v-if=" mode==='A' || mode === 'C'">
      <view  class="center">
        <view>
          <view style="display: inline-block">
            <view v-for="item in subject.subject1" :key="item" style="float: left">
                <text class="tianzige">{{item}}</text>
            </view>
          </view>
        </view>
      </view>
    </template>
    <template v-if=" mode=== 'C' || mode==='D'">
      <view  class="center">
        <view class="subjectText">
          <view style="display: inline-block;" >
            <view v-for="item in subject.subject2" :key="item.value" style="float: left">
              <text class="kati" :style="{color: colorD[item.show]}">{{item.value}}</text>
            </view>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<script>
export default {
  name: "Subject",
  props: {
      text: {
        type: String,
      },
      mode: {
        type: String
      }
  },
  data(){
    return{
      colorB: {
        0: '#000000',
        1: '#ad2b29'
      },
      colorD: {
        0: '#828282',
        1: '#000000',
      },
      subject: {}
    }
  },
  watch: {
    text: function () { this.refresh() },
    mode: function () { this.refresh() },
  },
  methods: { refresh() {
    console.log(this.text)
    this.subject = {}
    switch (this.mode) {
      case 'A': {
        this.subject.subject1 = []
        this.subject.subject1.push(this.text)
        break
      }
      case 'B': {
        // 1 要用的字   2 当前的字
        let parse = this.text.split('/')
        this.subject.subject1 = []
        for (let i = 0; i < parse[0].length; i++) {
          let tmpObject = {}
          tmpObject.value = parse[0][i]
          if(parseInt(parse[1]) === i){
            tmpObject.show = 1
          }
          else  tmpObject.show = 0
          this.subject.subject1.push(tmpObject)
        }
        break
      }
      case 'C': {
        // 1 未使用  0 已使用
        let parse = this.text.split('/')
        this.subject.subject1 = parse[0].split(' ')
        let tmpArray = parse[1].split(' ')
        this.subject.subject2 = []
        for (let i = 0; i < tmpArray.length; i++) {
          let tmpObject = {}
          tmpObject.value = tmpArray[i]
          if(parse[2] !== undefined) {
            tmpObject.show = parseInt(parse[2][i])
            console.log("gamePage")
          }
          else {
            tmpObject.show = 1
            console.log("endPage")
          }
          this.subject.subject2.push(tmpObject)
        }
        break
      }
      case 'D': {
        // 1 未使用  0 已使用
        let parse = this.text.split('/')
        let tmpArray1 = parse[0].split(' ')
        let tmpArray2 = parse[1].split(' ')
        this.subject.subject1 = []
        this.subject.subject2 = []
        for (let i = 0; i < tmpArray1.length; i++) {
          let tmpObject1 = {}
          tmpObject1.value = tmpArray1[i]
          if(parse[2] !== undefined) {
            tmpObject1.show = parseInt(parse[2][i])
          }
          else tmpObject1.show = 1
          this.subject.subject1.push(tmpObject1)
          let tmpObject2 = {}
          tmpObject2.value = tmpArray2[i]
          if(parse[3] !== undefined) {
            tmpObject2.show = parseInt(parse[3][i])
          }
          else tmpObject2.show = 1
          this.subject.subject2.push(tmpObject2)
        }
        break
      }
    }
  } },
}
</script>

<style scoped>
.tianzige {
  font-family: "楷体","楷体_GB2312";
  font-size: 34px;
  background-size: cover;
  background-image: url("/static/tianzige.png");
  margin: 0 7px 0 7px;
}
.subjectText {
  display: flex;
  background-color: #fdf8ed;
  border: 1.3px solid #975f5b;
  border-radius: 10px;
  padding: 6px 15px;
  margin: 5px 10px;
  /*horiz-align: center;*/
  /*position: absolute;*/
  /*margin: 10px;*/
}

.kati {
  font-family: "楷体","楷体_GB2312";
  font-size: 21px;
  padding: 5px 3px;
}

.center {
  text-align: center;
  margin: 12px 0;
}
</style>

<style>

</style>
