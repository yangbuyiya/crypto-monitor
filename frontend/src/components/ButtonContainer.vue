<template>
	<div class="button-container">
		<button class="icon-button" @click="handleEditPair" aria-label="add">
			<FontAwesomeIcon icon="fa-solid fa-plus" />
		</button>
		<button class="icon-button" @click="handleSettings" aria-label="settings">
			<FontAwesomeIcon icon="fa-solid fa-gear" />
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
	<AddCryptoPair v-if="showAddPair" />
	<Settings v-if="showSettings" />
</template>

<script setup>
	import AddCryptoPair from "./AddCryptoPair.vue";
	import Settings from "./Settings.vue";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import { ref } from "vue";
	import {
		WindowSetAlwaysOnTop,
		WindowMinimise,
		Quit,
	} from "../../wailsjs/runtime/runtime";
	import { defineEmits } from "vue";

	const isKeepOnTop = ref(false);
	const showAddPair = ref(false);
	const showSettings = ref(false);

	// 定义事件发射
	const emit = defineEmits(["toggle-add-pair"]);

	function handleEditPair() {
		showAddPair.value = !showAddPair.value;
		showSettings.value = false;
		emit("toggle-add-pair");
	}

	function handleSettings() {
		showSettings.value = !showSettings.value;
		showAddPair.value = false;
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
