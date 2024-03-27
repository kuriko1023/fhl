<script>
  import GamePage from "@/pages/game";
  import { staticRes } from 'utils';
	export default {
		onLaunch: function() {
            // 加载字体
            const loadFont = () => uni.loadFontFace({
              global: true,
              family: 'Kai',
              source: `url("${staticRes('FZKTJW.woff2')}")`,
              success: (x) => { console.log('font loaded!', x) },
              fail: (x) => {
                console.log('font load failed!', x)
                setTimeout(loadFont, 500)
              },
            })
            setTimeout(loadFont, 100)

		},
		onShow: function() {
		},
		onHide: function() {
		}
	}
</script>

<style>
uni-page-wrapper {
  overflow: hidden;
}
uni-page {
  min-height: 100%;
  margin: 0 auto;
  width: 100%;
  max-width: 60vh;
  position: relative;
}

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

  .kai {
    font-family: Kai;
  }
</style>
