<script setup>
	import Home from "./views/Home.vue";
	import { onMounted, onUnmounted } from 'vue';
	import { EventsEmit } from "../wailsjs/runtime/runtime";

	let mouseLeaveTimer;

	onMounted(() => {
		const handleMouseEnter = () => {
			clearTimeout(mouseLeaveTimer);
			EventsEmit("mouse_enter");
		};

		const handleMouseLeave = () => {
			// 添加延迟，避免鼠标快速移动时的闪烁
			mouseLeaveTimer = setTimeout(() => {
				EventsEmit("mouse_leave");
			}, 300);
		};

		document.addEventListener('mouseenter', handleMouseEnter);
		document.addEventListener('mouseleave', handleMouseLeave);

		onUnmounted(() => {
			document.removeEventListener('mouseenter', handleMouseEnter);
			document.removeEventListener('mouseleave', handleMouseLeave);
			clearTimeout(mouseLeaveTimer);
		});
	});
</script>

<template>
	<div id="app">
		<Home />
	</div>
</template>

<style>
	#logo {
		display: block;
		width: 50%;
		height: 50%;
		margin: auto;
		padding: 10% 0 0;
		background-position: center;
		background-repeat: no-repeat;
		background-size: 100% 100%;
		background-origin: content-box;
	}

	#app {
		height: 100vh;
		text-align: center;
		overflow: hidden;
		transition: width 0.3s ease;  /* 添加过渡动画 */
	}
</style>
