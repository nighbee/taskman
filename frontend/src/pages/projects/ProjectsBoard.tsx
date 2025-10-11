import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  DndContext,
  DragOverlay,
  closestCorners,
  PointerSensor,
  useSensor,
  useSensors,
  DragStartEvent,
  DragEndEvent,
} from '@dnd-kit/core';
import { arrayMove } from '@dnd-kit/sortable';
import { Button } from '@/components/ui/button';
import { KanbanColumn } from '@/components/kanban/KanbanColumn';
import { KanbanCard } from '@/components/kanban/KanbanCard';
import { useDataStore, type ProjectStatus } from '@/store/dataStore';
import { useAuthStore } from '@/store/authStore';
import { Plus } from 'lucide-react';
import { CreateProjectModal } from '@/components/modals/CreateProjectModal';

const statusConfig = {
  idea: { title: 'Idea', color: 'muted' as const },
  'in-progress': { title: 'In Progress', color: 'warning' as const },
  finished: { title: 'Finished', color: 'success' as const },
};

const ProjectsBoard = () => {
  const { orgId } = useParams<{ orgId: string }>();
  const { user } = useAuthStore();
  const { getProjectsByOrg, updateProjectStatus, selectProject, getOrgById } = useDataStore();
  const navigate = useNavigate();
  
  const [activeId, setActiveId] = useState<string | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const org = orgId ? getOrgById(orgId) : null;
  const projects = orgId ? getProjectsByOrg(orgId) : [];

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  const getProjectsByStatus = (status: ProjectStatus) => {
    return projects
      .filter(p => p.status === status)
      .sort((a, b) => {
        // User's projects first
        const aIsUser = a.createdBy === user?.id || a.assignees.some(u => u.id === user?.id);
        const bIsUser = b.createdBy === user?.id || b.assignees.some(u => u.id === user?.id);
        if (aIsUser && !bIsUser) return -1;
        if (!aIsUser && bIsUser) return 1;
        return 0;
      });
  };

  const handleDragStart = (event: DragStartEvent) => {
    setActiveId(event.active.id as string);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveId(null);

    if (!over) return;

    const activeProject = projects.find(p => p.id === active.id);
    const overStatus = over.id as ProjectStatus;

    if (activeProject && activeProject.status !== overStatus) {
      updateProjectStatus(activeProject.id, overStatus);
    }
  };

  const handleProjectClick = (projectId: string) => {
    selectProject(projectId);
    navigate(`/orgs/${orgId}/projects/${projectId}`);
  };

  const activeProject = projects.find(p => p.id === activeId);

  if (!org) {
    return <div className="p-6">Organization not found</div>;
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] p-6 bg-gradient-subtle">
      <div className="container max-w-7xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-3xl font-bold mb-2">{org.name} - Projects</h1>
            <p className="text-muted-foreground">
              Manage your projects across different stages
            </p>
          </div>
          <Button onClick={() => setIsCreateModalOpen(true)} className="gap-2">
            <Plus className="h-4 w-4" />
            New Project
          </Button>
        </div>

        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          <div className="flex gap-6 overflow-x-auto pb-6">
            {(Object.keys(statusConfig) as ProjectStatus[]).map(status => {
              const statusProjects = getProjectsByStatus(status);
              const config = statusConfig[status];

              return (
                <KanbanColumn
                  key={status}
                  id={status}
                  title={config.title}
                  color={config.color}
                  items={statusProjects.map(project => {
                    const isUserInvolved =
                      project.createdBy === user?.id ||
                      project.assignees.some(u => u.id === user?.id);

                    return (
                      <KanbanCard
                        key={project.id}
                        id={project.id}
                        title={project.name}
                        description={project.description}
                        deadline={project.deadline}
                        assignees={project.assignees}
                        isUserInvolved={isUserInvolved}
                        onClick={() => handleProjectClick(project.id)}
                      />
                    );
                  })}
                  itemIds={statusProjects.map(p => p.id)}
                  showAddButton={false}
                />
              );
            })}
          </div>

          <DragOverlay>
            {activeProject && (
              <KanbanCard
                id={activeProject.id}
                title={activeProject.name}
                description={activeProject.description}
                deadline={activeProject.deadline}
                assignees={activeProject.assignees}
              />
            )}
          </DragOverlay>
        </DndContext>
      </div>

      {orgId && (
        <CreateProjectModal
          open={isCreateModalOpen}
          onOpenChange={setIsCreateModalOpen}
          orgId={orgId}
          orgMembers={org.members}
        />
      )}
    </div>
  );
};

export default ProjectsBoard;
