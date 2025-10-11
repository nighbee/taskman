// API client for connecting to the TaskMan backend

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Types matching the backend models
export interface User {
  id: string;
  email: string;
  full_name: string;
  created_at: string;
  updated_at: string;
}

export interface Organization {
  id: string;
  name: string;
  invite_code: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface Project {
  id: string;
  org_id: string;
  name: string;
  description: string;
  status: 'idea' | 'in-progress' | 'finished';
  deadline?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
  task_count?: number;
}

export interface Task {
  id: string;
  project_id: string;
  name: string;
  description: string;
  status: 'not-started' | 'in-progress' | 'done';
  deadline?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface CreateOrgRequest {
  name: string;
}

export interface JoinOrgRequest {
  invite_code: string;
}

export interface CreateProjectRequest {
  name: string;
  description: string;
  deadline?: string;
  assignee_ids?: string[];
}

export interface CreateTaskRequest {
  name: string;
  description: string;
  deadline?: string;
  assignee_ids?: string[];
}

export interface UpdateTaskStatusRequest {
  status: 'not-started' | 'in-progress' | 'done';
}

export interface BulkMoveTasksRequest {
  task_ids: string[];
  status: 'not-started' | 'in-progress' | 'done';
}

// API client class
class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    // Load token from localStorage on initialization
    this.token = localStorage.getItem('auth_token');
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
    }

    return response.json();
  }

  // Auth endpoints
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    this.setToken(response.token);
    return response;
  }

  async register(email: string, password: string, fullName: string): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, full_name: fullName }),
    });
    this.setToken(response.token);
    return response;
  }

  // Organization endpoints
  async getOrganizations(): Promise<Organization[]> {
    return this.request<Organization[]>('/organizations');
  }

  async createOrganization(data: CreateOrgRequest): Promise<Organization> {
    return this.request<Organization>('/organizations', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async joinOrganization(data: JoinOrgRequest): Promise<Organization> {
    return this.request<Organization>('/organizations/join', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Project endpoints
  async getProjects(orgId: string): Promise<Project[]> {
    return this.request<Project[]>(`/organizations/${orgId}/projects`);
  }

  async createProject(orgId: string, data: CreateProjectRequest): Promise<Project> {
    return this.request<Project>(`/organizations/${orgId}/projects`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateProjectStatus(
    orgId: string,
    projectId: string,
    status: 'idea' | 'in-progress' | 'finished'
  ): Promise<Project> {
    return this.request<Project>(`/organizations/${orgId}/projects/${projectId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status }),
    });
  }

  // Task endpoints
  async getTasks(orgId: string, projectId: string): Promise<Task[]> {
    return this.request<Task[]>(`/organizations/${orgId}/projects/${projectId}/tasks`);
  }

  async createTask(
    orgId: string,
    projectId: string,
    data: CreateTaskRequest
  ): Promise<Task> {
    return this.request<Task>(`/organizations/${orgId}/projects/${projectId}/tasks`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateTaskStatus(
    orgId: string,
    projectId: string,
    taskId: string,
    data: UpdateTaskStatusRequest
  ): Promise<Task> {
    return this.request<Task>(`/organizations/${orgId}/projects/${projectId}/tasks/${taskId}/move`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async bulkMoveTasks(
    orgId: string,
    projectId: string,
    data: BulkMoveTasksRequest
  ): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/organizations/${orgId}/projects/${projectId}/tasks/bulk-move`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  // Utility methods
  setToken(token: string | null) {
    this.token = token;
    if (token) {
      localStorage.setItem('auth_token', token);
    } else {
      localStorage.removeItem('auth_token');
    }
  }

  clearToken() {
    this.setToken(null);
  }

  isAuthenticated(): boolean {
    return !!this.token;
  }
}

// Export singleton instance
export const apiClient = new ApiClient(API_BASE_URL);
