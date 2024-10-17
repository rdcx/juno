<template>
    <div>
        <router-link :to="to">
            <div :class="['p-2 rounded-lg', isCurrentPage ? 'dark:bg-gray-800 text-emerald-700 dark:text-gray-300' : 'dark:hover:bg-gray-800']">
                <div class="flex items-center gap-x-3">
                    <component :is="icon" class="w-6 h-6" />
                    <div>{{ label }}</div>
                </div>
            </div>
        </router-link>
    </div>
</template>

<script lang="ts" setup>
import { computed, type Component } from 'vue';
import { useRoute, type RouteLocation } from 'vue-router';

const props = defineProps<{
    icon: Component;
    to: { name: string };
    label: string;
}>();

const route = useRoute();

const isCurrentPage = computed(() => {
    return route.name === props.to.name;
});
</script>