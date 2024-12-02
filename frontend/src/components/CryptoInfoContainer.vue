<template>
	<div class="crypto-info-container">
		<div class="crypto-info-item" v-for="pair in pairs" :key="pair">
			<div class="icon-placeholder" v-if="showAddPair"></div>
			<CryptoInfo
				:pair="pair"
				:price="prices[pair] || 'Loading...'"
				:color="colors[pair] || '#FFFFFF'"
				:trend="trend[pair] || ''"
				@dblclick="handleDoubleClick(pair)" />
			<FontAwesomeIcon
				icon="fa-solid fa-xmark"
				color="red"
				v-if="showAddPair"
				@click="removePair(pair)"
				aria-label="remove pair"
				style="cursor: pointer" />
		</div>
	</div>
</template>

<script setup>
	import CryptoInfo from "./CryptoInfo.vue";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import { ref, onMounted } from "vue";
	import {
		EventsOn,
		EventsEmit,
		BrowserOpenURL,
	} from "../../wailsjs/runtime/runtime";
	import { defineProps } from "vue";

	const props = defineProps({
		showAddPair: {
			type: Boolean,
			default: false,
		},
	});

	const pairs = ref([]);
	const prices = ref({});
	const colors = ref({});
	const trend = ref({});
	const lastAverage = ref({});
	const period = ref(60);

	const loadPairsFromStorage = () => {
		const storedPairs = JSON.parse(localStorage.getItem("cryptoPairs")) || [];
		pairs.value = storedPairs;
		EventsEmit("listen_crypto_pairs", storedPairs); // 通知后端监听这些币种
	};

	const handleDoubleClick = (pair) => {
		const formattedPair = pair.toLowerCase();
		console.log(formattedPair);
		BrowserOpenURL(`https://www.okx.com/trade-spot/${formattedPair}`);
	};

	const removePair = (pair) => {
		const storedPairs = JSON.parse(localStorage.getItem("cryptoPairs")) || [];
		const updatedPairs = storedPairs.filter((p) => p !== pair);
		localStorage.setItem("cryptoPairs", JSON.stringify(updatedPairs));
		pairs.value = updatedPairs;
		EventsEmit("crypto_pairs_changed", updatedPairs);
	};

	onMounted(() => {
		loadPairsFromStorage();

		EventsEmit("crypto_pairs_changed", pairs.value); // trigger backend to listen pairs

		EventsOn("ticker_update", (data) => {
			if (!pairs.value.includes(data.pair)) return;

			prices.value[data.pair] = data.price;
			lastAverage.value[data.pair] =
				(parseFloat(lastAverage.value[data.pair] || 0) * (period.value - 1) +
					parseFloat(data.price)) /
				period.value;

			// diff between current price and average price
			const currentPrice = parseFloat(data.price);
			const avgPrice = lastAverage.value[data.pair];
			const diffPercent =
				Math.abs((currentPrice - avgPrice) / avgPrice) * 100 * 50;

			const brightness = Math.min(Math.max(diffPercent, 40), 100); // 将最低亮度提高到40

			// adjust color by price
			const hue = currentPrice > avgPrice ? 150 : 0; // 150 -> green, 0 -> red
			colors.value[data.pair] = `hsl(${hue}, 50%, ${brightness}%)`;

			trend.value[data.pair] =
				currentPrice > avgPrice ? "↑" : currentPrice < avgPrice ? "↓" : "";
		});

		EventsOn("crypto_pairs_changed", () => {
			loadPairsFromStorage();
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

	.crypto-info-item {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.icon-placeholder {
		width: 16px;
	}
</style>
