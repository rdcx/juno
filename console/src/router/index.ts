import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth' // Adjust the path if necessary

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
    },
    {
      path: '/data',
      name: 'data',
      component: () => import('../views/DataView.vue'),
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue'),
    },
    {
      path: '/strategies',
      name: 'strategies',
      component: () => import('../views/StrategiesView.vue'),
    },
    {
      path: '/selectors',
      name: 'selectors',
      component: () => import('../views/SelectorsView.vue'),
    },
    {
      path: '/fields',
      name: 'fields',
      component: () => import('../views/FieldsView.vue'),
    },
    {
      path: '/filters',
      name: 'filters',
      component: () => import('../views/FiltersView.vue'),
    },
    {
      path: '/jobs',
      name: 'jobs',
      component: () => import('../views/JobsView.vue'),
    },
    {
      path: '/tokens',
      name: 'tokens',
      component: () => import('../views/TokensView.vue'),
    },
    {
      path: '/crawl',
      name: 'crawl',
      component: () => import('../views/CrawlView.vue'),
    }
  ],
})

// Add a navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // Check if the user is authenticated
  if (!authStore.isAuthenticated && to.name !== 'login' && to.name !== 'register') {
    // If not authenticated, redirect to the login page
    next({ name: 'login' })
  } else {
    // Otherwise, proceed to the requested route
    next()
  }
})

export default router