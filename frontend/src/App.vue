<script>
  import GamePage from "@/pages/GamePage/GamePage";
  import { staticRes } from 'utils';
	export default {
		onLaunch: function() {
			console.log('App Launch')

            // 加载字体
            setTimeout(() => uni.loadFontFace({
              family: 'Kai',
              source: `url("${staticRes('FZKTJW.woff2')}")`,
              success: (x) => { console.log('font loaded!', x) },
              fail: (x) => { console.log('font load failed!', x) },
            }), 100)

		},
		onShow: function() {
		  // function connect() {
      //   let int = setInterval(function () {
      //     uni.connectSocket({
      //       url: 'ws://66.42.69.75:2310/channel/my/!kuriko',
      //       success: () => {
      //         console.log('connect success')
      //       }
      //     })
      //     uni.onSocketOpen(function (res) {
      //       clearInterval(int)
      //     })
      //   }, 100)
      // }
      // connect()
      // uni.onSocketClose(function(res) {
      //       connect()
      //     }
      // )
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
</script>

<style>
	/*每个页面公共css */
  radio::before,
  checkbox::before {
    font-size: 10px !important;
    color: #666666 !important;
  }

  radio .wx-radio-input,
  checkbox .wx-checkbox-input,
  radio .uni-radio-input,
  checkbox .uni-checkbox-input {
    margin: 5px;
    width: 12px !important;
    height: 12px !important;
  }

  switch.theme[checked] .wx-switch-input.wx-switch-input-checked,
  checkbox.theme[checked] .wx-checkbox-input,
  radio.theme[checked] .wx-radio-input,
  switch.theme.checked .uni-switch-input.uni-switch-input-checked,
  checkbox.theme.checked .uni-checkbox-input,
  radio.theme.checked .uni-radio-input{
    background-color: #366440 !important;
    border-color: #366440 !important;
    color: #366440 !important;
    font-size: 12px !important;
  }

  .uni-picker-container.uni-picker-header {
    background-color: #eaf6ed !important;
  }
  .uni-picker-container.uni-picker-action.uni-picker-action-confirm{
    color: #366440 !important;
  }

  uni-picker .uni-picker-action.uni-picker-action-confirm
  {color: #366440 !important;}
</style>
