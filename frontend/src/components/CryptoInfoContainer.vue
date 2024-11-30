<template>
	<div class="crypto-info-container">
		<CryptoInfo
			v-for="pair in pairs"
			:key="pair"
			:pair="pair"
			:price="prices[pair] || 'Loading...'"
			:color="colors[pair] || '#000000'" />
	</div>
</template>

<script setup>
	import CryptoInfo from "./CryptoInfo.vue";
	import { ref, onMounted } from "vue";
	import { EventsOn } from "../../wailsjs/runtime/runtime";

	const pairs = ref(["BTC-USDT", "ETH-USDT", "XRP-USDT"]);
	const prices = ref({});
	const colors = ref({});
	const lastAverage = ref({});
	const period = ref(30); // seconds

	onMounted(() => {
		EventsOn("ticker_update", (data) => {
			prices.value[data.pair] = data.price;
			lastAverage.value[data.pair] =
				(parseFloat(lastAverage.value[data.pair] || 0) * (period.value - 1) +
					parseFloat(data.price)) /
				period.value;
			colors.value[data.pair] =
				lastAverage.value[data.pair] > parseFloat(data.price)
					? "#6E9E84"
					: "#C2446B";
		});
	});
</script>

<style scoped>
	.crypto-info-container {
		display: flex;
		flex-direction: column;
		gap: 10px;
		padding: 5px;
	}
</style>
