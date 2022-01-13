<template>
  <view id="background">
    <image v-if='profileInitialized' class="startButton" @tap="onEnter" src="http://123.57.21.143:8000/game_start.png"></image>
    <image class="background" :src=backgroundImage />
    <image class="background" src="http://123.57.21.143:8000/start_background.png" ></image>
  </view>
</template>

<script>
import { ref, onMounted } from 'vue';
import {
  G,
  redirect,
  retrieveServerProfile,
  requestLocalProfile,
  closeSocket,
} from '../../utils';

import backgroundImage from '../../static/start_background_scaled.jpg';

export default {
  name: "StartPage",
  setup() {
    const profileInitialized = ref(false);

    onMounted(() => {
      retrieveServerProfile(() => profileInitialized.value = true);
      closeSocket();
    });

    const onEnter = () => {
      requestLocalProfile(() => {
        G.myRoom = true;
        console.log('entering room', G.my);
        // redirect("/pages/RoomPage/RoomPage");
      });
    };

    return {
      backgroundImage,

      profileInitialized,

      onEnter,
    };
  },
}
</script>

<style>
.startButton{
  position: absolute;
  width: 43%;
  height: 8%;
  z-index: 1;
  top: 87%;
  left: 55%;
}
.background{
  position: absolute;
  height: 100%;
  width: 100%;
  z-index: -1;
}
</style>
