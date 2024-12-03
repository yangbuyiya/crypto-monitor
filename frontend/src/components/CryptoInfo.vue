<template>
	<div class="crypto-info">
		<div class="header">
			<img
				:src="cryptoIconURL"
				:alt="pair.split('-')[0].toLowerCase() + ' icon'" />
			<div>
				{{ pair.split("-")[0] }}
			</div>
			<div
				class="percentage"
				:style="{ color: percentage.startsWith('+') ? '#99FF99' : '#FF9999' }">
				{{ percentage }}
			</div>
		</div>
		<div class="price-info">
			<div class="price" :style="{ color: color }">{{ price }} {{ trend }}</div>
		</div>
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
		percentage: {
			type: String,
			default: "0.00%",
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

	.percentage {
		font-size: 14px !important;
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
