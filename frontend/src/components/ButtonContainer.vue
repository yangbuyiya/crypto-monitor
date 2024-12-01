<template>
	<div class="button-container">
		<button class="icon-button" @click="handleClick" aria-label="add">
			<FontAwesomeIcon icon="fa-solid fa-plus" />
		</button>
		<button class="icon-button" @click="hideWindow" aria-label="hide">
			<FontAwesomeIcon icon="fa-solid fa-minus" />
		</button>
		<button class="icon-button" @click="toggleKeepOnTop" aria-label="pin">
			<FontAwesomeIcon
				:icon="
					isKeepOnTop ? 'fa-solid fa-thumbtack' : 'fa-solid fa-thumbtack-slash'
				" />
		</button>
		<button class="icon-button" @click="closeApp" aria-label="close">
			<FontAwesomeIcon icon="fa-solid fa-xmark" />
		</button>
	</div>
</template>

<script setup>
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import { ref } from "vue";
	import {
		WindowSetAlwaysOnTop,
		WindowMinimise,
		Quit,
	} from "../../wailsjs/runtime/runtime";

	const isKeepOnTop = ref(false);

	function handleClick() {
		console.log("Button clicked");
	}

	const hideWindow = () => {
		WindowMinimise();
	};

	const toggleKeepOnTop = () => {
		isKeepOnTop.value = !isKeepOnTop.value;
		WindowSetAlwaysOnTop(isKeepOnTop.value);
	};

	const closeApp = () => {
		// 关闭应用程序，假设 Wails 后端有对应的关闭函数
		Quit();
	};
</script>

<style scoped>
	.button-container {
		display: flex;
		gap: 10px;
		justify-content: center;
		align-items: center;
		margin-top: 12px;
	}

	.icon-button {
		-webkit-app-region: no-drag; /* 确保按钮不可拖动 */
		background: none;
		border: none;
		padding: 0;
		cursor: pointer;
		color: inherit; /* 继承父元素的颜色 */
	}

	.icon-button:focus {
		outline: none; /* 移除聚焦时的外框 */
	}

	.icon-button:hover {
		opacity: 0.8; /* 悬停时略微透明 */
	}

	.icon-button svg {
		width: 16px;
		height: 16px;
	}
</style>
