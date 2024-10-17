<template>
    <div class="h-screen dark:text-white dark:bg-gray-950 border-r pb-6 border-gray-800 px-6 flex flex-col">
        <div class="h-16 pl-2 flex items-center">
            <h1 class="text-3xl dark:text-white"><span class="text-emerald-700 dark:text-emerald-300">juno</span>
            </h1>
        </div>
        <div class="flex-1 space-y-1 -mx-2 mt-8 font-semibold">
            <SidebarMenuItem :icon="HomeIcon" :to="{ name: 'home' }" label="Home" />
            <SidebarMenuItem :icon="CircleStackIcon" :to="{ name: 'data' }" label="Data" />
            <SidebarMenuItem :icon="CubeTransparentIcon" :to="{ name: 'queries' }" label="Queries" />
        </div>

        <div class="mt-auto">
            <SidebarMenuItem :icon="CogIcon" :to="{ name: 'settings' }" label="Settings" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { CircleStackIcon, CubeTransparentIcon, DocumentIcon, HomeIcon, ServerStackIcon } from '@heroicons/vue/24/outline'
import { useRoute } from 'vue-router';
import SidebarMenuItem from './SidebarMenuItem.vue';
import { CogIcon } from '@heroicons/vue/24/outline';
const route = useRoute();
</script>