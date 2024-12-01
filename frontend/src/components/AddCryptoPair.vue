<template>
	<div class="add-crypto-pair">
		<input
			v-model="newPair"
			@keyup.enter="addPair"
			placeholder="BTC-USDT"
			:class="{ error: isError }"
			aria-label="add pair" />
		<button @click="addPair" aria-label="add pair">
			<FontAwesomeIcon icon="fa-solid fa-check" />
		</button>
		<div v-if="isError" class="error-message">
			Check your pair format. e.g. BTC-USDT
		</div>
	</div>
</template>

<script setup>
	import { ref } from "vue";
	import { EventsEmit } from "../../wailsjs/runtime/runtime";
	import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
	const newPair = ref("");
	const isError = ref(false);

	const addPair = () => {
		if (newPair.value.trim() === "") {
			isError.value = true;
			return;
		}

		const pair = newPair.value.trim().toUpperCase();
		// check pair format
		if (!pair.match(/^[A-Z]+-[A-Z]+$/)) {
			isError.value = true;
			return;
		}
		const storedPairs = JSON.parse(localStorage.getItem("cryptoPairs")) || [];

		if (!storedPairs.includes(pair)) {
			storedPairs.push(pair);
			localStorage.setItem("cryptoPairs", JSON.stringify(storedPairs));
			EventsEmit("crypto_pairs_changed", [pair]);
		}

		isError.value = false;
		newPair.value = "";
	};
</script>

<style scoped>
	.add-crypto-pair {
		display: flex;
		gap: 5px;
		justify-content: center;
		align-items: center;
		margin-top: 20px;
		flex-wrap: wrap;
	}

	.add-crypto-pair input {
		padding: 8px 12px;
		border: 1px solid #ccc;
		border-radius: 6px;
		width: 45%;
	}

	.add-crypto-pair input.error {
		border-color: #ff4444;
	}

	.add-crypto-pair button {
		background: none;
		border: none;
		cursor: pointer;
		color: inherit;
	}

	.error-message {
		color: #ff4444;
		font-size: 12px;
		width: 100%;
		text-align: center;
		margin-top: 5px;
	}
</style>
