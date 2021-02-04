<template>
	<view class="content">
		<image class="logo" src="/static/logo.png"></image>
		<view>
			<text class="title">{{title}}</text>
			<br>
			<text>User ID: {{uid}}</text>
		</view>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				title: 'Hello World!',
				uid: 'unknown',
			}
		},
		onLoad() {
			this.login()
		},
		methods: {
			login() {
				this.uid = 'loading'
				uni.login({
					provider: 'weixin',
					success: (ret) => {
						this.uid = ret.code
						console.log(this.uid)
					},
					fail: () => {
						this.uid = 'fail'
					},
				})
			}
		}
	}
</script>

<style>
	.content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.logo {
		height: 200rpx;
		width: 200rpx;
		margin: 200rpx auto 50rpx auto;
	}

	.text-area {
		display: flex;
		justify-content: center;
	}

	.title {
		font-size: 36rpx;
		color: #8f8f94;
	}
</style>
