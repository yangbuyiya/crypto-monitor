<template>
    <div class="settings-container">
        <h3>默认交易对设置</h3>
        <div class="pairs-container">
            <div v-for="(pair, index) in defaultPairs" :key="index" class="pair-item">
                <input 
                    v-model="defaultPairs[index]" 
                    :class="{ error: !isValidPair(pair) }"
                    @change="validateAndSave"
                />
                <button @click="removePair(index)" class="remove-btn">
                    <FontAwesomeIcon icon="fa-solid fa-trash" />
                </button>
            </div>
            <div class="add-new">
                <input 
                    v-model="newPair"
                    placeholder="添加新交易对"
                    :class="{ error: newPairError }"
                />
                <button @click="addNewPair" class="add-btn">
                    <FontAwesomeIcon icon="fa-solid fa-plus" />
                </button>
            </div>
        </div>
        <div v-if="error" class="error-message">
            请使用正确的格式，如：BTC-USDT
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { EventsEmit } from "../../wailsjs/runtime/runtime";

const defaultPairs = ref([]);
const newPair = ref('');
const error = ref(false);
const newPairError = ref(false);

onMounted(() => {
    const stored = localStorage.getItem('defaultCryptoPairs');
    if (stored) {
        defaultPairs.value = JSON.parse(stored);
    } else {
        defaultPairs.value = ["BTC-USDT", "ETH-USDT", "SOL-USDT"];
        savePairs();
    }
});

const isValidPair = (pair) => {
    return pair.match(/^[A-Z]+-[A-Z]+$/);
};

const validateAndSave = () => {
    error.value = false;
    const allValid = defaultPairs.value.every(pair => isValidPair(pair));
    if (allValid) {
        savePairs();
    } else {
        error.value = true;
    }
};

const savePairs = () => {
    localStorage.setItem('defaultCryptoPairs', JSON.stringify(defaultPairs.value));
    // 更新当前显示的交易对
    const currentPairs = JSON.parse(localStorage.getItem('cryptoPairs')) || [];
    const newPairs = [...new Set([...defaultPairs.value, ...currentPairs])];
    localStorage.setItem('cryptoPairs', JSON.stringify(newPairs));
    EventsEmit("crypto_pairs_changed", newPairs);
};

const removePair = (index) => {
    defaultPairs.value.splice(index, 1);
    savePairs();
};

const addNewPair = () => {
    newPairError.value = false;
    const pair = newPair.value.trim().toUpperCase();
    
    if (!isValidPair(pair)) {
        newPairError.value = true;
        return;
    }
    
    if (!defaultPairs.value.includes(pair)) {
        defaultPairs.value.push(pair);
        savePairs();
        newPair.value = '';
    }
};
</script>

<style scoped>
.settings-container {
    padding: 15px;
    max-width: 300px;
    margin: 0 auto;
}

.pairs-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 15px;
}

.pair-item, .add-new {
    display: flex;
    gap: 8px;
    align-items: center;
}

input {
    padding: 6px 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    flex: 1;
}

input.error {
    border-color: #ff4444;
}

button {
    background: none;
    border: none;
    cursor: pointer;
    color: inherit;
    padding: 4px 8px;
}

.remove-btn:hover {
    color: #ff4444;
}

.error-message {
    color: #ff4444;
    font-size: 12px;
    margin-top: 8px;
    text-align: center;
}

h3 {
    margin: 0 0 15px 0;
    text-align: center;
}
</style> 