<template>
    <main>
        <div class="flex items-center justify-center h-screen">
            <div class="mx-auto w-[32rem] relative group rounded">
                <div class="relative group h-full w-full p-2">
                    <div
                        class="hidden dark:block absolute -inset-10 blur transition opacity-10 rounded-lg bg-gradient-to-r from-emerald-800 to-emerald-950 rounded-lg dark:from-emerald-600 dark:to-gray-950">
                    </div>
                    <div class="relative dark:bg-gray-950 p-8 rounded-xl">
                        <h1 class="text-3xl dark:text-white"><span class="font-thin">data</span><span
                                class="text-emerald-700 dark:text-emerald-300">juno</span>
                        </h1>
                        <h1 class="dark:text-white text-lg">Welcome back</h1>
                        <form @submit.prevent="login" class="mt-4">
                            <div class="mb-4">
                                <label for="email"
                                    class="block text-sm font-medium text-gray-700 dark:text-gray-200">Email</label>
                                <input type="email" id="email" v-model="email" placeholder="john@example.com"
                                    class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-emerald-500 focus:border-emerald-500 sm:text-sm"
                                    required>
                            </div>

                            <div class="mb-4">
                                <label for="password"
                                    class="block text-sm font-medium text-gray-700 dark:text-gray-200">Password</label>
                                <input type="password" id="password" v-model="password"
                                    class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-emerald-500 focus:border-emerald-500 sm:text-sm"
                                    required>
                            </div>

                            <div v-if="error" class="text-red-500 text-sm mb-4">
                                {{ error }}
                            </div>

                            <div>
                                <button type="submit"
                                    class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-emerald-600 hover:bg-emerald-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500">
                                    Login
                                </button>
                            </div>

                            <div class="mt-4">
                                <p class="text-sm text-gray-600 dark:text-gray-300">Don't have an account? <router-link
                                        to="/register"
                                        class="font-medium text-emerald-600 dark:text-emerald-300 hover:text-emerald-500">Register</router-link>
                                </p>
                            </div>
                        </form>
                    </div>


                </div>
            </div>
        </div>
    </main>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { AuthService } from '@/services';
import type { AuthResponse, LoginPayload } from '@/types/AuthTypes';
import { ProfileService } from '@/services';
import { useRouter } from 'vue-router';
import type { ProfileResponse } from '@/types/ProfileTypes';

const authStore = useAuthStore();
const router = useRouter();

const email = ref('');
const password = ref('');
const error = ref('');

const login = () => {
    const payload: LoginPayload = {
        email: email.value,
        password: password.value
    };
    AuthService.login(payload)
        .then((res: AuthResponse) => {
            console.log(res)
            authStore.setToken(res.token);
            ProfileService.getProfile()
                .then((res: ProfileResponse) => {
                    authStore.setUser(res.user);
                    router.push({ name: 'home' });
                });
        })
        .catch((err: any) => {
            console.log(err)
            error.value = err.response?.data?.message;
        });
}
</script>
