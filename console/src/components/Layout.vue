<script setup lang="ts">
import { ChevronDownIcon } from '@heroicons/vue/24/solid';
import Sidebar from './Sidebar.vue';
import { useAuthStore } from '@/stores/auth';
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const authStore = useAuthStore();
const user = computed(() => authStore.user);

import {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
} from '@headlessui/vue'
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline';

</script>

<template>
    <div class="flex">
        <div class="w-72">
            <Sidebar />
        </div>
        <div class="flex-grow">
            <div class="w-full dark:bg-gray-950 h-16 border-b border-gray-800 flex justify-between items-center">
                <div class="dark:text-gray-300 px-8">
                    <form class="relative flex flex-1" action="#" method="GET">
                        <label for="search-field" class="sr-only">Search</label>
                        <MagnifyingGlassIcon
                            class="pointer-events-none absolute inset-y-0 left-0 h-full w-5 text-gray-400"
                            aria-hidden="true" />
                        <input id="search-field"
                            class="block h-full w-full text-lg border-0 py-0 pl-8 pr-0 dark:bg-gray-950 dark:text-gray-300 focus:border-0 text-gray-900 placeholder:text-gray-400 focus:ring-0 focus:outline-none"
                            placeholder="Search..." type="search" name="search" />
                    </form>
                </div>
                <div class="ml-auto px-4 ">
                    <!-- Profile dropdown -->
                    <Menu as="div" class="relative">
                        <MenuButton class="-m-1.5 flex items-center p-1.5">
                            <span class="sr-only">Open user menu</span>
                            <img class="h-8 w-8 rounded-full bg-gray-50"
                                src="https://pbs.twimg.com/profile_images/1826933648399925248/E_klgWfw_400x400.jpg"
                                alt="" />
                            <span class="hidden lg:flex lg:items-center">
                                <span class="ml-4 text-sm font-semibold leading-6 text-gray-900 dark:text-gray-300"
                                    aria-hidden="true">{{ user.name }}</span>
                                <ChevronDownIcon class="ml-2 h-5 w-5 text-gray-400" aria-hidden="true" />
                            </span>
                        </MenuButton>
                        <transition enter-active-class="transition ease-out duration-100"
                            enter-from-class="transform opacity-0 scale-95"
                            enter-to-class="transform opacity-100 scale-100"
                            leave-active-class="transition ease-in duration-75"
                            leave-from-class="transform opacity-100 scale-100"
                            leave-to-class="transform opacity-0 scale-95">
                            <MenuItems
                                class="absolute dark:bg-gray-800 right-0 z-10 mt-2.5 w-32 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 focus:outline-none">
                                <MenuItem v-for="item in [
                                    { name: 'Profile', action: () => router.push({ name: 'profile' }) },
                                    { name: 'Settings', action: () => router.push({ name: 'settings' }) },
                                    { name: 'Sign out', action: () => authStore.logout() }
                                ]" :key="item.name" v-slot="{ active }" class="w-full text-left">
                                <button @click="item.action"
                                    :class="[active ? 'bg-gray-50 dark:bg-gray-700' : '', 'block px-3 py-1 text-sm leading-6 text-gray-900 dark:text-gray-300']">{{
                                        item.name }}
                                </button>
                                </MenuItem>
                            </MenuItems>
                        </transition>
                    </Menu>
                </div>
            </div>
            <div class="p-4">
                <slot />
            </div>
        </div>
    </div>
</template>