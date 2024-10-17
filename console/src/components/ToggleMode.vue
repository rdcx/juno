
<script lang="ts" setup>
import { MoonIcon } from '@heroicons/vue/24/solid'
import { SunIcon } from '@heroicons/vue/24/solid'
import { ref } from 'vue';
import { useDark } from '@vueuse/core';
import { useToggle } from '@vueuse/core';

const isDark = useDark();
const toggleDark = useToggle(isDark);

</script>

<template>
    <div class="fixed bottom-0 right-0 m-4">
        <button @click="toggleDark()" class="rounded bg-indigo-700 p-3 text-white">
            <div class="flex gap-2">
                <MoonIcon v-if="!isDark" class="w-6 h-6" />
                <SunIcon v-if="isDark" class="w-6 h-6" />
                {{ isDark ? 'light' : 'dark' }}
            </div>
        </button> 
    </div>
</template>
