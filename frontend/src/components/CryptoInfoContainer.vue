<template>
	<div class="crypto-info-container">
		<!-- 滚动区域 -->
		<div class="scroll-container">
			<div class="crypto-info-grid">
				<div 
					v-for="(pair, index) in paginatedPairs" 
					:key="pair"
					class="crypto-info-item"
					:draggable="isDragging === index"
					@mousedown="startLongPress(index)"
					@mouseup="clearLongPress"
					@mouseleave="clearLongPress"
					@dragstart="handleDragStart($event, index)"
					@dragover.prevent
					@dragenter.prevent
					@drop="handleDrop($event, index)"
					@click.stop
					:class="{ 
						'dragging': draggedIndex === index,
						'press-feedback': isPressing === index
					}"
				>
					<CryptoInfo
						:pair="pair"
						:price="prices[pair] || 'Loading...'"
						:color="colors[pair] || '#FFFFFF'"
						:trend="trend[pair] || ''"
						:percentage="percentage[pair] || ''"
						@dblclick="handleDoubleClick(pair)" />
					<div class="icon-placeholder" v-if="showAddPair"></div>
					<FontAwesomeIcon
						icon="fa-solid fa-xmark"
						color="red"
						v-if="showAddPair"
						@click="removePair(pair)"
						aria-label="remove pair"
						style="cursor: pointer" />
				</div>
			</div>
		</div>

		<!-- 分页器 -->
		<div class="pagination" v-if="totalPages > 1">
			<button
				class="page-button"
				@click="currentPage--"
				:disabled="currentPage === 1">
				<FontAwesomeIcon icon="fa-solid fa-chevron-left" />
			</button>
			<span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
			<button
				class="page-button"
				@click="currentPage++"
				:disabled="currentPage === totalPages">
				<FontAwesomeIcon icon="fa-solid fa-chevron-right" />
			</button>
		</div>
	</div>
</template>

<script setup>
	import CryptoInfo from "./CryptoInfo.vue";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	import { ref, onMounted, computed, watch, onUnmounted } from "vue";
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
	const percentage = ref({});
	const currentPage = ref(1);
	const itemsPerPage = 10;
	const draggedIndex = ref(null);
	const isDragging = ref(null);
	const isPressing = ref(null);
	let longPressTimer = null;
	const LONG_PRESS_DURATION = 200; // 长按触发时间（毫秒）

	const totalPages = computed(() => {
		return Math.ceil(pairs.value.length / itemsPerPage);
	});

	const paginatedPairs = computed(() => {
		const start = (currentPage.value - 1) * itemsPerPage;
		const end = start + itemsPerPage;
		return pairs.value.slice(start, end);
	});

	watch(pairs, (newPairs) => {
		const maxPage = Math.ceil(newPairs.length / itemsPerPage);
		if (currentPage.value > maxPage) {
			currentPage.value = Math.max(1, maxPage);
		}
	});

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

	// 开始长按
	const startLongPress = (index) => {
		clearLongPress();
		isPressing.value = index;
		EventsEmit("set_dragging", true);
		
		// 阻止事件冒泡
		event.stopPropagation();
		
		longPressTimer = setTimeout(() => {
			isDragging.value = index;
			const event = new MouseEvent('dragstart', {
				bubbles: true,
				cancelable: true,
				view: window
			});
			const element = document.querySelector(`.crypto-info-item:nth-child(${index + 1})`);
			if (element) {
				element.dispatchEvent(event);
			}
		}, LONG_PRESS_DURATION);
	};

	// 清除长按状态
	const clearLongPress = () => {
		if (longPressTimer) {
			clearTimeout(longPressTimer);
			longPressTimer = null;
		}
		isPressing.value = null;
		if (!draggedIndex.value) {
			isDragging.value = null;
		}
		// 通知后端恢复窗口拖动
		EventsEmit("set_dragging", false);
	};

	// 处理拖动开始
	const handleDragStart = (e, index) => {
		draggedIndex.value = index;
		isDragging.value = index;
		e.dataTransfer.effectAllowed = 'move';
	};

	// 处理放置
	const handleDrop = (e, dropIndex) => {
		e.preventDefault();
		if (draggedIndex.value === null || draggedIndex.value === dropIndex) {
			clearLongPress();
			isDragging.value = null;
			draggedIndex.value = null;
			// 确保恢复窗口拖动
			EventsEmit("set_dragging", false);
			return;
		}

		// 获取当前所有交易对
		const allPairs = [...pairs.value];
		
		// 计算实际的索引（考虑分页）
		const actualDragIndex = (currentPage.value - 1) * itemsPerPage + draggedIndex.value;
		const actualDropIndex = (currentPage.value - 1) * itemsPerPage + dropIndex;
		
		// 交换位置
		const temp = allPairs[actualDragIndex];
		allPairs[actualDragIndex] = allPairs[actualDropIndex];
		allPairs[actualDropIndex] = temp;

		// 更新存储和状态
		localStorage.setItem("cryptoPairs", JSON.stringify(allPairs));
		pairs.value = allPairs;
		EventsEmit("crypto_pairs_changed", allPairs);
		
		clearLongPress();
		isDragging.value = null;
		draggedIndex.value = null;
		// 确保恢复窗口拖动
		EventsEmit("set_dragging", false);
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

			percentage.value[data.pair] = data.percentage;
		});

		EventsOn("crypto_pairs_changed", () => {
			loadPairsFromStorage();
		});
	});

	onUnmounted(() => {
		clearLongPress();
		EventsEmit("set_dragging", false);
	});
</script>

<style scoped>
	.crypto-info-container {
		display: flex;
		flex-direction: column;
		height: calc(100vh - 20px); /* 减去padding */
		padding: 10px;
		gap: 10px;
	}

	.scroll-container {
		flex: 1;
		overflow-y: auto;
		overflow-x: hidden;
		/* 自定义滚动条样式 */
		scrollbar-width: thin;
		scrollbar-color: rgba(255, 255, 255, 0.3) transparent;
	}

	/* Webkit浏览器的滚动条样式 */
	.scroll-container::-webkit-scrollbar {
		width: 6px;
	}

	.scroll-container::-webkit-scrollbar-track {
		background: transparent;
	}

	.scroll-container::-webkit-scrollbar-thumb {
		background-color: rgba(255, 255, 255, 0.3);
		border-radius: 3px;
	}

	.crypto-info-grid {
		display: flex;
		flex-direction: column;
		gap: 8px;
		padding: 5px;
	}

	.crypto-info-item {
		display: flex;
		align-items: center;
		justify-content: center;
		user-select: none;
		transition: transform 0.2s ease, opacity 0.2s ease, background-color 0.2s ease;
		-webkit-app-region: no-drag;
		padding: 8px;
		border-radius: 4px;
	}

	.crypto-info-item.press-feedback {
		background-color: rgba(255, 255, 255, 0.1);
		transform: scale(0.98);
	}

	.crypto-info-item.dragging {
		opacity: 0.5;
		transform: scale(0.95);
		background-color: rgba(255, 255, 255, 0.15);
		cursor: move;
	}

	.crypto-info-item:not(.dragging) {
		cursor: default;
	}

	.pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 10px;
		padding: 10px 0;
		background-color: rgba(0, 0, 0, 0.2);
		border-radius: 4px;
		margin-top: auto; /* 确保分页器始终在底部 */
	}

	.page-button {
		background: none;
		border: 1px solid rgba(255, 255, 255, 0.3);
		border-radius: 4px;
		padding: 4px 8px;
		cursor: pointer;
		color: inherit;
		transition: background-color 0.2s ease;
	}

	.page-button:hover:not(:disabled) {
		background-color: rgba(255, 255, 255, 0.1);
	}

	.page-button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.page-info {
		font-size: 14px;
		min-width: 50px;
		text-align: center;
	}
</style>
