import { createRouter, createWebHistory } from 'vue-router';
import DashboardLayout from '@/components/DashboardLayout.vue';
import AuthView from '@/views/AuthView.vue';
import ManageNodes from '@/views/ManageNodes.vue';
import ManageBalancers from '@/views/ManageBalancers.vue';
import { useUserStore } from '@/stores/user';

// Define routes
const routes = [
  {
    path: '/',
    component: DashboardLayout,
    children: [
      {
        name: 'add-node',
        path: '/add-node',
        component: () => import('@/views/AddNode.vue'),
      },
      { name: 'nodes', path: 'nodes', component: ManageNodes },
      { name: 'balancers', path: 'balancers', component: ManageBalancers },
    ],
    meta: { requiresAuth: true }, // Protected routes
  },

  {
    path: '/login',
    component: AuthView,
  },
];

// Create router
const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Add route guard
router.beforeEach((to, from, next) => {
  const userStore = useUserStore();

  // Check if route requires authentication
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!userStore.token || !userStore.user) {
      // If no token or user, redirect to login
      next('/login');
    } else {
      next(); // Allow access to route
    }
  } else {
    next(); // Route does not require authentication, allow access
  }
});

export default router;
