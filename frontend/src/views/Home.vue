<template>
	<ButtonContainer @toggle-add-pair="toggleAddPair" />
	<CryptoInfoContainer :show-add-pair="showAddPair" />
	<div v-if="errorMessages.length" class="error-container">
		<h3>Subscription Error</h3>
		<ul>
			<li v-for="error in errorMessages" :key="error">{{ error }}</li>
		</ul>
	</div>
</template>

<script setup>
	import ButtonContainer from "../components/ButtonContainer.vue";
	import CryptoInfoContainer from "../components/CryptoInfoContainer.vue";
	import { ref, onMounted, onUnmounted, nextTick, watch } from "vue";
	import { WindowSetSize, EventsOn, EventsOff } from "../../wailsjs/runtime";

	const showAddPair = ref(false);
	const homeRef = ref(null);
	const errorMessages = ref([]);

	const toggleAddPair = () => {
		showAddPair.value = !showAddPair.value;
	};

	const handleSubscriptionError = (data) => {
		errorMessages.value.push(`${data.pair}: ${data.error}`);
	};

	onMounted(() => {
		updateWindowHeight();

		EventsOn("ticker_subscription_error", handleSubscriptionError);
	});

	onUnmounted(() => {
		EventsOff("ticker_subscription_error", handleSubscriptionError);
	});
</script>

<style scoped>
	.error-container {
		background-color: #ffe6e6;
		border: 1px solid #ffcccc;
		padding: 10px;
		margin-top: 20px;
		border-radius: 5px;
	}
	.error-container h3 {
		color: #cc0000;
	}
</style>
