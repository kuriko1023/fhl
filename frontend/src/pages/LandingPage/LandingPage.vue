<template>
  <view id="background">
    <image v-if='profileInitialized' class="startButton" @tap="onEnter"
      :src="staticRes('game_start.png')" />
    <image class="background" :src=backgroundImage />
    <image class="background" :src="staticRes('start_background.png')" />
  </view>
</template>

<script>
import { ref, onMounted } from 'vue';
import {
  G,
  staticRes,
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
        redirect("/pages/RoomPage/RoomPage");
      });
    };

    return {
      staticRes,
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
