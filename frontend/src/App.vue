<script>
  import GamePage from "@/pages/GamePage/GamePage";
	export default {
		onLaunch: function() {
			console.log('App Launch')
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
</style>
