<template>
  <view>
    <view class="out">
      <view class="in" :style="{width: curPercent, background: color}"></view>
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
    total: {
      type: Number
    },
    start: {
      type: Number
    },
    end: {
      type: Number
    },
  },
  data(){
    return {
      curPercent: '100%',
      timerId: -1,
    }
  },
  methods:{
    intervalFunction(){
      if (this.active) {
        const t = Date.now() / 1000
        const avail = this.end - this.start
        const remain = avail - Math.max(0, Math.min(avail, t - this.start))
        this.curPercent = (remain / this.total * 100) + '%'
      }
    }
  },
  mounted() {
    this.timerId = setInterval(this.intervalFunction, 200)
  },
  unmounted() {
    if (this.timerId !== -1) clearInterval(this.timerId)
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
    transition: width 200ms ease;
  }
</style>
