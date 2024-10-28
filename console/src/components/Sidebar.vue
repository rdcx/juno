<template>
    <div class="h-screen dark:text-white dark:bg-gray-950 border-r pb-6 border-gray-800 px-6 flex flex-col">
        <div class="h-16 pl-2 flex items-center">
            <h1 class="text-3xl dark:text-white"><span class="font-thin">data</span><span class="text-emerald-700 dark:text-emerald-300">juno</span>
            </h1>
        </div>
        <div class="flex flex-1 gap-y-7 flex-col">
            <div class="flex space-y-1  mt-8 font-semibold flex-col">
                <SidebarMenuItem :icon="HomeIcon" :to="{ name: 'home' }" label="Home" />
                <SidebarMenuItem :icon="BanknotesIcon" :to="{ name: 'tokens' }" label="Tokens" />
            </div>
            <div class="flex flex-col">
                <div class="text-sm font-semibold leading-6 text-indigo-200">Data</div>
                <div class="mt-2">
                    <SidebarMenuItem :icon="CircleStackIcon" :to="{ name: 'data' }" label="Stores" />
                </div>
            </div>
            <div class="flex flex-col">
                <div class="text-sm font-semibold leading-6 text-indigo-200">Extraction</div>
                <div class="mt-2">
                    <SidebarMenuItem :icon="BeakerIcon" :to="{ name: 'strategies' }" label="Strategies" />
                    <SidebarMenuItem :icon="ArrowsPointingInIcon" :to="{ name: 'selectors' }" label="Selectors" />
                    <SidebarMenuItem :icon="ViewColumnsIcon" :to="{ name: 'fields' }" label="Fields" />
                    <SidebarMenuItem :icon="FunnelIcon" :to="{ name: 'filters' }" label="Filters" />
                    <SidebarMenuItem :icon="QueueListIcon" :to="{ name: 'jobs' }" label="Jobs" />
                </div>
            </div>
            <div class="flex flex-1 flex-col">
                <div class="text-sm font-semibold leading-6 text-indigo-200">Crawling</div>
                <div class="mt-2">
                    <SidebarMenuItem :icon="BugAntIcon" :to="{ name: 'crawl' }" label="Jobs" />
                </div>
            </div>
            
            <div class="mt-auto">
                <SidebarMenuItem :icon="CogIcon" :to="{ name: 'settings' }" label="Settings" />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ArrowsPointingInIcon, BanknotesIcon, BeakerIcon, BugAntIcon, CircleStackIcon, CubeTransparentIcon, DocumentIcon, FunnelIcon, HomeIcon, QueueListIcon, ServerStackIcon, ViewColumnsIcon } from '@heroicons/vue/24/outline'
import { useRoute } from 'vue-router';
import SidebarMenuItem from './SidebarMenuItem.vue';
import { CogIcon } from '@heroicons/vue/24/outline';
const route = useRoute();
</script>