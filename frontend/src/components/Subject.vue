<template>
  <view>
    <view>
      <text v-for="item in data.subject1" :key="item">
        {{item}}
      </text>
    </view>
    <view>
      <text>
        {{data.subject2}}
      </text>
    </view>
  </view>
</template>

<script>
export default {
  name: "Subject",
  props: {
      data: {
        type: Object,
      },
      mode: {
        type: String
      }
  }
}
</script>

<style scoped>

</style>