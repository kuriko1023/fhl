<template>
  <view>
    <ul>
      <li v-for="item in data" :key="item.value">
        <text>
          {{item.value}}
        </text>
      </li>
    </ul>
  </view>
</template>

<script>
export default {
name: "History",
  props: {
    data: {
      type: Array
    }
  }
}
</script>

<style scoped>

</style>