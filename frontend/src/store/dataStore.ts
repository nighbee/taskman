import { create } from 'zustand';
import { apiClient, type User, type Organization, type Project, type Task } from '@/lib/api';

export type ProjectStatus = 'idea' | 'in-progress' | 'finished';
export type TaskStatus = 'not-started' | 'in-progress' | 'done';

// Re-export types from API for convenience
export type { User, Task };

// Map backend Organization to frontend Org
export interface Org extends Organization {
  members?: User[]; // This will be populated separately
}

// Map backend Project to frontend Project
export interface ProjectWithAssignees extends Project {
  assignees?: User[]; // This will be populated separately
}

// Map backend Task to frontend Task
export interface TaskWithAssignees extends Task {
  assignees?: User[]; // This will be populated separately
}

interface DataState {
  orgs: Org[];
  projects: Project[];
  tasks: Task[];
  selectedOrgId: string | null;
  selectedProjectId: string | null;
  isLoading: boolean;
  error: string | null;
  
  // Data loading
  loadOrganizations: () => Promise<void>;
  loadProjects: (orgId: string) => Promise<void>;
  loadTasks: (orgId: string, projectId: string) => Promise<void>;
  
  // Org actions
  createOrg: (name: string) => Promise<Org>;
  joinOrg: (inviteCode: string) => Promise<boolean>;
  selectOrg: (orgId: string) => void;
  
  // Project actions
  createProject: (orgId: string, name: string, description: string, deadline?: string, assigneeIds?: string[], createdBy?: string, status?: ProjectStatus) => Promise<void>;
  updateProjectStatus: (orgId: string, projectId: string, status: ProjectStatus) => Promise<void>;
  selectProject: (projectId: string) => void;
  
  // Task actions
  createTask: (orgId: string, projectId: string, name: string, description: string, deadline?: string, assigneeIds?: string[]) => Promise<void>;
  updateTaskStatus: (orgId: string, projectId: string, taskId: string, status: TaskStatus) => Promise<void>;
  bulkMoveTasks: (orgId: string, projectId: string, taskIds: string[], status: TaskStatus) => Promise<void>;
  
  // Getters
  getOrgById: (orgId: string) => Org | undefined;
  getProjectsByOrg: (orgId: string) => Project[];
  getTasksByProject: (projectId: string) => Task[];
  
  // Utilities
  clearError: () => void;
}

export const useDataStore = create<DataState>((set, get) => ({
  orgs: [],
  projects: [],
  tasks: [],
  selectedOrgId: null,
  selectedProjectId: null,
  isLoading: false,
  error: null,
  
  // Data loading methods
  loadOrganizations: async () => {
    set({ isLoading: true, error: null });
    try {
      const orgs = await apiClient.getOrganizations();
      set({ orgs, isLoading: false });
    } catch (error) {
      set({ 
        orgs: [], // Reset to empty array on error
        error: error instanceof Error ? error.message : 'Failed to load organizations',
        isLoading: false 
      });
    }
  },
  
  loadProjects: async (orgId: string) => {
    set({ isLoading: true, error: null });
    try {
      const projects = await apiClient.getProjects(orgId);
      set(state => ({ 
        projects: [...state.projects.filter(p => p.org_id !== orgId), ...projects],
        isLoading: false 
      }));
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to load projects',
        isLoading: false 
      });
    }
  },
  
  loadTasks: async (orgId: string, projectId: string) => {
    set({ isLoading: true, error: null });
    try {
      const tasks = await apiClient.getTasks(orgId, projectId);
      set(state => ({ 
        tasks: [...state.tasks.filter(t => t.project_id !== projectId), ...tasks],
        isLoading: false 
      }));
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to load tasks',
        isLoading: false 
      });
    }
  },
  
  // Organization actions
  createOrg: async (name: string) => {
    set({ isLoading: true, error: null });
    try {
      const newOrg = await apiClient.createOrganization({ name });
      set(state => ({ 
        orgs: [...(Array.isArray(state.orgs) ? state.orgs : []), newOrg],
        isLoading: false 
      }));
      return newOrg;
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to create organization',
        isLoading: false 
      });
      throw error;
    }
  },
  
  joinOrg: async (inviteCode: string) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.joinOrganization({ invite_code: inviteCode });
      // Reload organizations to get updated data
      await get().loadOrganizations();
      set({ isLoading: false });
      return true;
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to join organization',
        isLoading: false 
      });
      return false;
    }
  },
  
  selectOrg: (orgId: string) => {
    set({ selectedOrgId: orgId, selectedProjectId: null });
    // Load projects for the selected organization
    get().loadProjects(orgId);
  },
  
  // Project actions
  createProject: async (orgId: string, name: string, description: string, deadline?: string, assigneeIds?: string[], createdBy?: string, status?: ProjectStatus) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.createProject(orgId, {
        name,
        description,
        deadline,
        assignee_ids: assigneeIds,
        created_by: createdBy,
        status: status || 'idea',
      });
      // Reload projects for the organization
      await get().loadProjects(orgId);
      set({ isLoading: false });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to create project',
        isLoading: false 
      });
      throw error;
    }
  },
  
  updateProjectStatus: async (orgId: string, projectId: string, status: ProjectStatus) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.updateProjectStatus(orgId, projectId, status);
      set(state => ({
        projects: state.projects.map(p =>
          p.id === projectId ? { ...p, status } : p
        ),
        isLoading: false
      }));
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to update project status',
        isLoading: false 
      });
      throw error;
    }
  },
  
  selectProject: (projectId: string) => {
    set({ selectedProjectId: projectId });
    const selectedOrgId = get().selectedOrgId;
    if (selectedOrgId) {
      // Load tasks for the selected project
      get().loadTasks(selectedOrgId, projectId);
    }
  },
  
  // Task actions
  createTask: async (orgId: string, projectId: string, name: string, description: string, deadline?: string, assigneeIds?: string[]) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.createTask(orgId, projectId, {
        name,
        description,
        deadline,
        assignee_ids: assigneeIds,
      });
      // Reload tasks for the project
      await get().loadTasks(orgId, projectId);
      set({ isLoading: false });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to create task',
        isLoading: false 
      });
      throw error;
    }
  },
  
  updateTaskStatus: async (orgId: string, projectId: string, taskId: string, status: TaskStatus) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.updateTaskStatus(orgId, projectId, taskId, { status });
      set(state => ({
        tasks: state.tasks.map(t =>
          t.id === taskId ? { ...t, status } : t
        ),
        isLoading: false
      }));
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to update task status',
        isLoading: false 
      });
      throw error;
    }
  },
  
  bulkMoveTasks: async (orgId: string, projectId: string, taskIds: string[], status: TaskStatus) => {
    set({ isLoading: true, error: null });
    try {
      await apiClient.bulkMoveTasks(orgId, projectId, { task_ids: taskIds, status });
      set(state => ({
        tasks: state.tasks.map(t =>
          taskIds.includes(t.id) ? { ...t, status } : t
        ),
        isLoading: false
      }));
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to move tasks',
        isLoading: false 
      });
      throw error;
    }
  },
  
  // Getters
  getOrgById: (orgId: string) => get().orgs.find(o => o.id === orgId),
  
  getProjectsByOrg: (orgId: string) => get().projects.filter(p => p.org_id === orgId),
  
  getTasksByProject: (projectId: string) => get().tasks.filter(t => t.project_id === projectId),
  
  // Utilities
  clearError: () => set({ error: null }),
}));
