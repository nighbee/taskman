import { create } from 'zustand';
import { apiClient, type User } from '@/lib/api';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<void>;
  logout: () => void;
  clearError: () => void;
  initializeAuth: () => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
  
  login: async (email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.login(email, password);
      set({ 
        user: response.user, 
        isAuthenticated: true, 
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Login failed',
        isLoading: false 
      });
      throw error;
    }
  },
  
  register: async (email: string, password: string, fullName: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.register(email, password, fullName);
      set({ 
        user: response.user, 
        isAuthenticated: true, 
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Registration failed',
        isLoading: false 
      });
      throw error;
    }
  },
  
  logout: () => {
    apiClient.clearToken();
    set({ user: null, isAuthenticated: false, error: null });
  },
  
  clearError: () => set({ error: null }),
  
  initializeAuth: async () => {
    // Check if user is already authenticated via token
    if (apiClient.isAuthenticated()) {
      try {
        // Verify token by making a request to get user info
        // This will fail if token is expired or invalid
        const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';
        const response = await fetch(`${API_BASE_URL}/auth/me`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
            'Content-Type': 'application/json',
          },
        });
        
        if (response.ok) {
          const userData = await response.json();
          set({ 
            user: userData.user, 
            isAuthenticated: true 
          });
        } else {
          // Token is invalid, clear it
          apiClient.clearToken();
          set({ user: null, isAuthenticated: false });
        }
      } catch (error) {
        // Network error or token invalid, clear it
        apiClient.clearToken();
        set({ user: null, isAuthenticated: false });
      }
    }
  },
}));
