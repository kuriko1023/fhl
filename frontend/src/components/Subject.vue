<template>
  <view>
    <!-- {{JSON.stringify(subject)}} -->
    <template v-if="mode === 'D'">
      <view style="display: flex; flex-wrap:wrap;" class="subjectText">
        <view v-for="item in subject.subject1" :key="item.value">
          <p class="kati word" :style="{color: colorD[item.show]}">{{item.value}}</p>
        </view>
      </view>
    </template>
    <template v-if="mode==='B'">
      <view  class="center">
        <view class="subjectText">
          <view style="display: inline-block;" >
            <view v-for="item in subject.subject1" :key="item.value" style="float: left">
              <p class="kati char" :style="{color: colorB[item.show]}">{{item.value}}</p>
            </view>
          </view>
        </view>
      </view>
    </template>
    <template v-if=" mode==='A' || mode === 'C'">
      <view  class="center">
        <view>
          <view style="display: inline-block">
            <view class="tianzige" v-for="item in subject.subject1" :key="item" style="float: left">
                <text>{{item}}</text>
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
              <p class="kati word" :style="{color: colorD[item.show]}">{{item.value}}</p>
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
    mode: { type: String },
    subject: { type: Object },
  },
  data(){
    return{
      colorB: {
        0: '#000000',
        1: '#ad2b29',
        2: '#828282',
      },
      colorD: {
        0: '#828282',
        1: '#000000',
      },
    }
  },
}
</script>

<style scoped>
.tianzige {
  background-image: url("/static/tianzige.png");
  width: 45px;
  height: 45px;
  background-size: cover;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: space-around;
}
.tianzige text {
  font-family: Kai;
  font-size: 34px;
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
  font-family: Kai;
  font-size: 21px;
  padding: 5px 3px;
}
.char {
  margin: 0 -2px;
}
.word {
  margin: 0 5px;
}

.center {
  text-align: center;
  margin: 12px 0;
}
</style>

<style>

</style>
