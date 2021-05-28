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
  },
  data(){
  return{
    cur: 0,
    int: 0,
  }
  },
  watch: {
    active: function(val){
      console.log('active')
      if(val){
        this.int = setInterval(this.intervalFunction, this.time)
      }
      else{
        clearInterval(this.int)
        console.log('a')
        this.$emit('stop', (100 - this.cur) * this.time / 100)
        console.log('b')
      }
    },
    current: function(val){
      if(val !== 0) {
        console.log('current')
        this.cur = 100 - (val / this.time) * 100
      }
    }
  },
  methods:{
    percent(num){
      return num + '%'
    },
    intervalFunction(){
        this.cur = this.cur + 0.1
        if(this.cur >= 100){
          clearInterval(this.int)
          this.$emit('finish')
          this.cur = 100
        }
    }
  },
  mounted() {
    if (this.active) {
      this.int = setInterval(this.intervalFunction, this.time)
    }
  }
}
</script>

<style scoped>
  .out{
    height: 10px;
  }
  .in{
    height: 10px;
  }
</style>