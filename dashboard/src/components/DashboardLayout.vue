<template>
    <div v-if="isLoggedIn" class="flex h-screen bg-gray-900 text-white">
      <!-- Sidebar -->
      <aside class="w-72 h-full bg-gray-950">
        <div class="items-center justify-center space-y-4 p-8">
          <h1 class="text-2xl text-white mb-4 p-2">Juno</h1>
          <nav>
            <ul>
              <li class="mb-2">
                <router-link class="block p-2 hover:bg-gray-700 rounded-xl" to="/nodes">Nodes</router-link>
              </li>
              <li class="mb-2">
                <router-link class="block p-2 hover:bg-gray-700 rounded-xl" to="/balancers">Balancers</router-link>
              </li>
            </ul>
          </nav>
        </div>
      </aside>
  
      <!-- Main Content -->
      <main class="flex-1 p-6">
        <!-- Top Navigation -->
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-3xl">Node Management</h2>
          <div class="space-x-4">
            <span v-if="email">{{ email.split('@')[0] }}</span>
            <button class="px-4 py-2 bg-gray-700 rounded hover:bg-gray-600" @click="logout">Logout</button>
          </div>
        </div>
  
        <!-- Dashboard Content -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-lg">
          <router-view />
        </div>
      </main>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, computed } from 'vue';
  import { useUserStore } from '@/stores/user';
  import { useRouter } from 'vue-router';
  
  export default defineComponent({
    name: 'DashboardLayout',
    setup() {
      const userStore = useUserStore();
      const router = useRouter();
  
      const isLoggedIn = computed(() => !!userStore.token);
      const email = computed(() => userStore.user?.email);
  
      const logout = () => {
        userStore.logout();
        router.push('/login');
      };
  
      return {
        isLoggedIn,
        email,
        logout,
      };
    },
  });
  </script>
  