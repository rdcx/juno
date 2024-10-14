<template>
    <div class="flex items-center justify-center h-screen bg-black text-white">
      <div class="p-8 bg-stone-900 rounded-xl shadow-lg w-96">
        <h2 class="text-3xl mb-4">{{ isSignup ? 'Sign Up' : 'Login' }}</h2>
        <form @submit.prevent="handleSubmit">
          <div class="mb-4">
            <label class="block text-white mb-1">Email</label>
            <input v-model="email" type="email" class="w-full p-2 rounded bg-stone-800 text-white" required />
          </div>
          <div class="mb-4">
            <label class="block text-white mb-1">Password</label>
            <input v-model="password" type="password" class="w-full p-2 rounded bg-stone-800 text-white" required />
          </div>
          <button type="submit" class="w-full bg-stone-700 p-2 rounded hover:bg-stone-600">
            {{ isSignup ? 'Sign Up' : 'Login' }}
          </button>
        </form>
  
        <p class="mt-4">
          <span v-if="isSignup">Already have an account?</span>
          <span v-else>Don't have an account?</span>
          <a href="#" @click.prevent="toggleMode" class="ml-2 text-stone-500 underline">
            {{ isSignup ? 'Login' : 'Sign Up' }}
          </a>
        </p>
      </div>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref } from 'vue';
  import { useUserStore } from '@/stores/user';
  import { useRouter } from 'vue-router';
  
  export default defineComponent({
    name: 'AuthView',
    setup() {
      const email = ref('');
      const password = ref('');
      const isSignup = ref(false);
      const userStore = useUserStore();
      const router = useRouter();
  
      const handleSubmit = async () => {
        try {
          if (isSignup.value) {
            await userStore.signup(email.value, password.value);
            isSignup.value = false; // Switch to login after signup
          } else {
            await userStore.login(email.value, password.value);
            router.push('/'); // Redirect to the dashboard after login
          }
        } catch (error) {
          alert('Failed to authenticate');
        }
      };
  
      const toggleMode = () => {
        isSignup.value = !isSignup.value;
      };
  
      return {
        email,
        password,
        isSignup,
        handleSubmit,
        toggleMode,
      };
    },
  });
  </script>
  