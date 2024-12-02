<template>
	<div class="crypto-info">
		<div class="header">
			<img
				:src="cryptoIconURL"
				:alt="pair.split('-')[0].toLowerCase() + ' icon'" />
			<div :style="{ color: color }">{{ pair }}</div>
		</div>
		<div :style="{ color: color }">{{ price }} {{ trend }}</div>
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
		trend: {
			type: String,
			default: "",
		},
	});

	const cryptoIconURL = computed(() => {
		return `https://cdn.jsdelivr.net/gh/vadimmalykhin/binance-icons/crypto/${props.pair
			.split("-")[0]
			.toLowerCase()}.svg`;
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

	.header {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		margin-bottom: 8px;
	}

	.header img {
		width: 16px;
		height: 16px;
	}

	.header div {
		font-size: 16px;
		font-weight: bold;
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
