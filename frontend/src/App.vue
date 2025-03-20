<script setup>
	import Home from "./views/Home.vue";
	import { onMounted, onUnmounted } from 'vue';
	import { EventsEmit } from "../wailsjs/runtime/runtime";

	onMounted(() => {
		const handleMouseEnter = () => {
			EventsEmit("mouse_enter");
		};

		const handleMouseLeave = () => {
			EventsEmit("mouse_leave");
		};

		document.addEventListener('mouseenter', handleMouseEnter);
		document.addEventListener('mouseleave', handleMouseLeave);

		onUnmounted(() => {
			document.removeEventListener('mouseenter', handleMouseEnter);
			document.removeEventListener('mouseleave', handleMouseLeave);
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
	}
</style>
