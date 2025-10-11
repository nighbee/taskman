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
  
  initializeAuth: () => {
    // Check if user is already authenticated via token
    if (apiClient.isAuthenticated()) {
      // We could decode the JWT token here to get user info
      // For now, we'll just set authenticated to true
      // In a real app, you might want to verify the token with the backend
      set({ isAuthenticated: true });
    }
  },
}));
