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

	// max difference ratio (50%)
	const maxDiffRatio = 0.5;

	const calculateColor = (difference, avgPrice) => {
		const ratio =
			Math.min(Math.abs(difference / avgPrice), maxDiffRatio) / maxDiffRatio; // 0 到 1
		let r, g, b;
		let lightness = 40 + ratio * 30; // lightness from 40% to 80%

		if (difference > 0) {
			// white -> green
			r = Math.round(255 * (1 - ratio));
			g = 255;
			b = Math.round(255 * (1 - ratio));
		} else if (difference < 0) {
			// red -> white
			r = 255;
			g = Math.round(255 * (1 - ratio));
			b = Math.round(255 * (1 - ratio));
		} else {
			// equal -> white
			r = 255;
			g = 255;
			b = 255;
		}

		// convert to HSL to adjust lightness
		const rgbToHsl = (r, g, b) => {
			r /= 255;
			g /= 255;
			b /= 255;
			const max = Math.max(r, g, b);
			const min = Math.min(r, g, b);
			let h,
				s,
				l = (max + min) / 2;

			if (max === min) {
				h = s = 0; // achromatic
			} else {
				const d = max - min;
				s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
				switch (max) {
					case r:
						h = (g - b) / d + (g < b ? 6 : 0);
						break;
					case g:
						h = (b - r) / d + 2;
						break;
					case b:
						h = (r - g) / d + 4;
						break;
				}
				h /= 6;
			}
			return [h * 360, s * 100, l * 100];
		};

		[, , lightness] = rgbToHsl(r, g, b);
		// limit lightness between 50% and 80%
		lightness = Math.min(Math.max(lightness, 50), 80);

		return `hsl(${difference > 0 ? 120 : 0}, 100%, ${lightness}%)`;
	};

	const loadPairsFromStorage = () => {
		const storedPairs = JSON.parse(localStorage.getItem("cryptoPairs")) || [];
		pairs.value = storedPairs;
		EventsEmit("listen_crypto_pairs", storedPairs);
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
				(parseFloat(lastAverage.value[data.pair] || period.value) *
					(period.value - 1) +
					parseFloat(data.price)) /
				period.value;

			// diff between current price and average price
			const currentPrice = parseFloat(data.price);
			const avgPrice = lastAverage.value[data.pair];
			const difference = currentPrice - avgPrice;

			colors.value[data.pair] = calculateColor(difference, avgPrice);

			if (difference > 0) {
				trend.value[data.pair] = "↑";
			} else if (difference < 0) {
				trend.value[data.pair] = "↓";
			} else {
				trend.value[data.pair] = "";
			}
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
