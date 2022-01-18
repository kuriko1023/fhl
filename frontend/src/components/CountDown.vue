<template>
  <view>
    <view class="out">
      <view class="in" :style="{width: percent(100 - cur), background: color}"></view>
    </view>
  </view>
</template>

<script>
export default {
name: "CountDown",
  props: {
    active: {
      type: Boolean
    },
    color: {
      type: String
    },
    time: {
      type: Number
    },
    current: {
      type: Number
    },
    update: {
      type: Number
    },
  },
  data(){
  return{
    cur: 0,
    int: -1,
    lastTimestamp: 0,
  }
  },
  watch: {
    active: function(val){
      //console.log('active')
      if(val){
        if (this.int === -1) {
          this.lastTimestamp = Date.now()
          this.int = setInterval(this.intervalFunction, 200)
        }
      }
      else{
        if (this.int !== -1) {
          clearInterval(this.int)
          this.int = -1
        }
       // console.log('a')
        this.$emit('stop', (100 - this.cur) * this.time / 100)
       // console.log('b')
      }
    },
    update: function () {
      const val = this.current;
      this.cur = 100 - (val / this.time) * 100
    }
  },
  methods:{
    percent(num){
      return num + '%'
    },
    intervalFunction(){
        const now = Date.now()
        // 10 = 1000 (ms) / 100 (percent)
        this.cur += (now - this.lastTimestamp) / (10 * this.time);
        this.lastTimestamp = now;
        if(this.cur >= 100){
          clearInterval(this.int)
          this.int = -1;
          this.$emit('finish')
          this.cur = 100
        }
    }
  },
  mounted() {
    if (this.active) {
      if (this.int === -1) {
        this.lastTimestamp = Date.now()
        this.int = setInterval(this.intervalFunction, 200)
      }
    }
  },
  unmounted() {
    if (this.int !== -1) clearInterval(this.int)
  },
}
</script>

<style scoped>
  .out{
    height: 7px;
    border-radius: 2px;
    margin: 3px 0;
  }
  .in{
    height: 7px;
    border-radius: 2px;
    margin: 3px 0;
    transition: width linear 200ms;
  }
</style>
