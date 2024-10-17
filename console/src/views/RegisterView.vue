<template>
    <main>
        <div class="flex items-center justify-center h-screen">
            <div class="mx-auto w-[32rem] relative group rounded">
                <div class="relative group h-full w-full p-2">
                    <div
                        class="hidden dark:block absolute -inset-0.5 blur transition opacity-30 rounded-lg bg-gradient-to-r from-emerald-800 to-emerald-950 rounded-lg dark:from-emerald-600 dark:to-emerald-800">
                    </div>
                    <div class="relative dark:bg-gray-950 p-8 rounded-xl">
                        <h1 class="text-3xl dark:text-white"><span
                                class="text-emerald-700 dark:text-emerald-300">juno</span>
                        </h1>
                        <h1 class="dark:text-white text-lg">console.register</h1>
                        <form @submit.prevent="register" class="mt-4">
                            <div class="mb-4">
                                <label for="name"
                                    class="block text-sm font-medium text-gray-700 dark:text-gray-200">Name</label>
                                <input type="name" id="name" v-model="name"
                                    class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-emerald-500 focus:border-emerald-500 sm:text-sm"
                                    required>
                            </div>

                            <div class="mb-4">
                                <label for="email"
                                    class="block text-sm font-medium text-gray-700 dark:text-gray-200">Email</label>
                                <input type="email" id="email" v-model="email"
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
                                    Register
                                </button>
                            </div>

                            <div class="mt-4">
                                <p class="text-sm text-gray-600 dark:text-gray-300">Already have an account?
                                    <router-link to="/login"
                                        class="font-medium text-emerald-600 dark:text-emerald-300 hover:text-emerald-500">Login</router-link>
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
import { AuthService } from '@/api';
import type { AuthResponse, LoginPayload, RegisterPayload } from '@/types/AuthTypes';
import { ProfileService } from '@/api';
import { useRouter } from 'vue-router';
import type { ProfileResponse } from '@/types/ProfileTypes';
import { useNotificationStore } from '@/stores/notification';

const authStore = useAuthStore();
const router = useRouter();

const name = ref('');
const email = ref('');
const password = ref('');
const error = ref('');

const register = () => {
    const payload: RegisterPayload = {
        name: name.value,
        email: email.value,
        password: password.value
    };
    AuthService.register(payload)
        .then((res: AuthResponse) => {
            if (res.status == "success") {
                AuthService.login({ email: email.value, password: password.value })
                    .then((res: AuthResponse) => {
                        authStore.setToken(res.token);
                        ProfileService.getProfile()
                            .then((res: ProfileResponse) => {
                                authStore.setUser(res.user);
                                router.push({ name: 'home' });

                                useNotificationStore().createNotification(
                                    'Account created',
                                    'success'
                                );
                            });
                    });
                return
            }
            error.value = res.message;
        })
        .catch((err: any) => {
            error.value = err.response?.data?.message;
        });
}
</script>
