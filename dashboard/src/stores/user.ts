import { defineStore } from 'pinia';
import axios from 'axios';

interface UserState {
  token: string | null;
  user: { email: string } | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem('token'),
    user: JSON.parse(localStorage.getItem('user') || 'null'),
  }),

  actions: {
    async login(email: string, password: string) {
      try {
        const response = await axios.post('http://localhost:8080/auth/token', { email, password });
        const { token } = response.data;

        this.token = token;

        localStorage.setItem('token', token);
        
        const profileResponse = await axios.get('http://localhost:8080/profile', {
            headers: { Authorization: token },
        });
        
        this.user = profileResponse.data.user;
        localStorage.setItem('user', JSON.stringify(this.user));
        
        return this.user;

      } catch (error) {
        console.error('Login failed:', error);
        throw error;
      }
    },

    async signup(email: string, password: string) {
      try {
        const response = await axios.post('http://localhost:8080/users', { email, password });
        return response.data;
      } catch (error) {
        console.error('Signup failed:', error);
        throw error;
      }
    },

    logout() {
      this.token = null;
      this.user = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    },
  },
});
