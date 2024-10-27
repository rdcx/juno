<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import StrategyService from '@/services/strategy/StrategyService';
import {FieldService, FilterService, SelectorService } from '@/services';
import type { Strategy } from '@/types/StrategyTypes';
import type { Selector } from '@/types/SelectorTypes';
import type { Field } from '@/types/FieldTypes';
import type { Filter } from '@/types/FilterTypes';

const props = defineProps<{ strategy: Strategy | null }>();
const emit = defineEmits(['close']);

const strategy = ref(props.strategy);
const availableSelectors = ref<Array<Selector>>([]);
const availableFields = ref<Array<Field>>([]);
const availableFilters = ref<Array<Filter>>([]);
const selectedSelectorId = ref('');
const selectedFieldId = ref('');
const selectedFilterId = ref('');

watch(() => props.strategy, (newStrategy) => {
    strategy.value = newStrategy;
});

const fetchAvailableEntities = async () => {
    availableSelectors.value = (await SelectorService.list()).selectors;
    availableFields.value = (await FieldService.list()).fields;
    availableFilters.value = (await FilterService.list()).filters;
};

const addSelector = async () => {
    await StrategyService.AddSelector({ strategy_id: strategy.value?.id || '', selector_id: selectedSelectorId.value });
    fetchUpdatedStrategy();
};

const removeSelector = async (selector_id: string) => {
    await StrategyService.RemoveSelector({ strategy_id: strategy.value?.id || '', selector_id });
    fetchUpdatedStrategy();
};

const addField = async () => {
    await StrategyService.AddField({ strategy_id: strategy.value?.id || '', field_id: selectedFieldId.value });
    fetchUpdatedStrategy();
};

const removeField = async (field_id: string) => {
    await StrategyService.RemoveField({ strategy_id: strategy.value?.id || '', field_id });
    fetchUpdatedStrategy();
};

const addFilter = async () => {
    await StrategyService.AddFilter({ strategy_id: strategy.value?.id || '', filter_id: selectedFilterId.value });
    fetchUpdatedStrategy();
};

const removeFilter = async (filter_id: string) => {
    await StrategyService.RemoveFilter({ strategy_id: strategy.value?.id || '', filter_id });
    fetchUpdatedStrategy();
};

const fetchUpdatedStrategy = async () => {
    if (!strategy.value) return;
    const response = await StrategyService.list();
    strategy.value = response.strategies.find(s => s.id === strategy.value?.id) || null;
};

onMounted(() => {
    fetchAvailableEntities();
});

</script>

<template>
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-4xl w-full overflow-auto">
            <h3 class="text-xl font-semibold dark:text-gray-300 mb-4">Edit Strategy: {{ strategy?.name }}</h3>

            <div class="grid grid-cols-3 gap-6">
                <!-- Selectors Management -->
                <div class="col-span-1">
                    <h4 class="font-medium dark:text-gray-300">Selectors</h4>
                    <ul class="pl-4 mb-4">
                        <li v-for="selector in strategy?.selectors || []" :key="selector.id" class="flex justify-between items-center">
                            <span class="dark:text-gray-300">{{ selector.name }}</span>
                            <button @click="removeSelector(selector.id)" class="text-red-500 hover:underline">Remove</button>
                        </li>
                    </ul>
                    <select v-model="selectedSelectorId" class="border p-2 rounded w-full mb-2">
                        <option value="" disabled>Select a Selector</option>
                        <option v-for="selector in availableSelectors" :key="selector.id" :value="selector.id">
                            {{ selector.name }}
                        </option>
                    </select>
                    <button @click="addSelector" class="bg-blue-500 text-white rounded p-2 w-full">Add Selector</button>
                </div>

                <!-- Fields Management -->
                <div class="col-span-1">
                    <h4 class="font-medium dark:text-gray-300">Fields</h4>
                    <ul class="pl-4 mb-4">
                        <li v-for="field in strategy?.fields || []" :key="field.id" class="flex justify-between items-center">
                            <span class="dark:text-gray-300">{{ field.name }}</span>
                            <button @click="removeField(field.id)" class="text-red-500 hover:underline">Remove</button>
                        </li>
                    </ul>
                    <select v-model="selectedFieldId" class="border p-2 rounded w-full mb-2">
                        <option value="" disabled>Select a Field</option>
                        <option v-for="field in availableFields" :key="field.id" :value="field.id">
                            {{ field.name }}
                        </option>
                    </select>
                    <button @click="addField" class="bg-blue-500 text-white rounded p-2 w-full">Add Field</button>
                </div>

                <!-- Filters Management -->
                <div class="col-span-1">
                    <h4 class="font-medium dark:text-gray-300">Filters</h4>
                    <ul class="pl-4 mb-4">
                        <li v-for="filter in strategy?.filters || []" :key="filter.id" class="flex justify-between items-center">
                            <span class="dark:text-gray-300">{{ filter.name }}</span>
                            <button @click="removeFilter(filter.id)" class="text-red-500 hover:underline">Remove</button>
                        </li>
                    </ul>
                    <select v-model="selectedFilterId" class="border p-2 rounded w-full mb-2">
                        <option value="" disabled>Select a Filter</option>
                        <option v-for="filter in availableFilters" :key="filter.id" :value="filter.id">
                            {{ filter.name }}
                        </option>
                    </select>
                    <button @click="addFilter" class="bg-blue-500 text-white rounded p-2 w-full">Add Filter</button>
                </div>
            </div>

            <button @click="$emit('close')" class="bg-gray-500 text-white rounded p-2 w-full mt-6">Close</button>
        </div>
    </div>
</template>
