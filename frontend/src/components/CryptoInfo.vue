<template>
	<div class="crypto-info">
		<div :style="{ color: color }">{{ pair }}</div>
		<div :style="{ color: color }">{{ price }}</div>
	</div>
</template>

<script setup>
	import { computed } from "vue";

	const props = defineProps({
		pair: {
			type: String,
			required: true,
		},
		price: {
			type: String,
			required: true,
		},
		color: {
			type: String,
			default: "#000000",
		},
	});

	const color = computed(() => {
		return props.color || "#000000";
	});
</script>

<style scoped>
	.crypto-info {
		background-color: #13191d;
		border-radius: 8px;
		padding: 10px;
		margin: 5px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		min-width: 100px;
		text-align: center;
	}

	.crypto-info div:first-child {
		font-size: 16px;
		font-weight: bold;
		margin-bottom: 8px;
	}

	.crypto-info div:last-child {
		font-size: 16px;
		font-weight: bold;
	}

	.crypto-info:hover {
		transform: translateY(-2px);
		transition: transform 0.2s ease;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.15);
	}
	.crypto-info {
		user-select: none;
	}
</style>
