import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
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
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { KanbanColumn } from '@/components/kanban/KanbanColumn';
import { KanbanCard } from '@/components/kanban/KanbanCard';
import { useDataStore, type TaskStatus } from '@/store/dataStore';
import { useAuthStore } from '@/store/authStore';
import { ArrowLeft, Calendar, Users } from 'lucide-react';
import { CreateTaskModal } from '@/components/modals/CreateTaskModal';

const statusConfig = {
  'not-started': { title: 'Not Started', color: 'muted' as const },
  'in-progress': { title: 'In Progress', color: 'warning' as const },
  done: { title: 'Done', color: 'success' as const },
};

const ProjectDetail = () => {
  const { orgId, projectId } = useParams<{ orgId: string; projectId: string }>();
  const { user } = useAuthStore();
  const { projects, getTasksByProject, updateTaskStatus } = useDataStore();
  
  const [activeId, setActiveId] = useState<string | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  const project = projects.find(p => p.id === projectId);
  const tasks = projectId ? getTasksByProject(projectId) : [];

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  const isUserAssignee = (project?.assignees && project.assignees.some(a => a.id === user?.id)) || project?.createdBy === user?.id;

  const getTasksByStatus = (status: TaskStatus) => {
    return tasks
      .filter(t => t.status === status)
      .sort((a, b) => {
        const aIsUser = a.createdBy === user?.id || (a.assignees && a.assignees.some(u => u.id === user?.id));
        const bIsUser = b.createdBy === user?.id || (b.assignees && b.assignees.some(u => u.id === user?.id));
        if (aIsUser && !bIsUser) return -1;
        if (!aIsUser && bIsUser) return 1;
        return 0;
      });
  };

  const handleDragStart = (event: DragStartEvent) => {
    if (!isUserAssignee) return;
    setActiveId(event.active.id as string);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveId(null);

    if (!over || !isUserAssignee) return;

    const activeTask = tasks.find(t => t.id === active.id);
    const overStatus = over.id as TaskStatus;

    if (activeTask && activeTask.status !== overStatus) {
      updateTaskStatus(activeTask.id, overStatus);
    }
  };

  const activeTask = tasks.find(t => t.id === activeId);

  if (!project) {
    return <div className="p-6">Project not found</div>;
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] p-6 bg-gradient-subtle">
      <div className="container max-w-7xl mx-auto">
        <div className="mb-6">
          <Button variant="ghost" size="sm" asChild className="mb-4 gap-2">
            <Link to={`/orgs/${orgId}/projects`}>
              <ArrowLeft className="h-4 w-4" />
              Back to Projects
            </Link>
          </Button>

          <div className="bg-card rounded-lg p-6 shadow-md mb-6">
            <div className="flex items-start justify-between mb-4">
              <div className="flex-1">
                <h1 className="text-3xl font-bold mb-2">{project.name}</h1>
                <p className="text-muted-foreground">{project.description}</p>
              </div>
              <Badge
                variant={
                  project.status === 'finished'
                    ? 'default'
                    : project.status === 'in-progress'
                    ? 'secondary'
                    : 'outline'
                }
                className="ml-4"
              >
                {statusConfig[project.status as keyof typeof statusConfig]?.title || project.status}
              </Badge>
            </div>

            <div className="flex flex-wrap gap-4 text-sm">
              {project.deadline && (
                <div className="flex items-center gap-2 text-muted-foreground">
                  <Calendar className="h-4 w-4" />
                  <span>Due: {new Date(project.deadline).toLocaleDateString()}</span>
                </div>
              )}
              
              {project.assignees && project.assignees.length > 0 && (
                <div className="flex items-center gap-2">
                  <Users className="h-4 w-4 text-muted-foreground" />
                  <div className="flex gap-1">
                    {project.assignees.map(assignee => (
                      <Badge key={assignee.id} variant="secondary">
                        {assignee.full_name || assignee.email}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">Tasks</h2>
            {isUserAssignee && (
              <Button onClick={() => setIsCreateModalOpen(true)} size="sm">
                Add Task
              </Button>
            )}
          </div>

          {!isUserAssignee && (
            <div className="bg-warning/10 border border-warning/30 rounded-lg p-4 mb-4">
              <p className="text-sm text-warning-foreground">
                You are viewing this project as an observer. Only assignees can create or move tasks.
              </p>
            </div>
          )}
        </div>

        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          <div className="flex gap-6 overflow-x-auto pb-6">
            {(Object.keys(statusConfig) as TaskStatus[]).map(status => {
              const statusTasks = getTasksByStatus(status);
              const config = statusConfig[status];

              return (
                <KanbanColumn
                  key={status}
                  id={status}
                  title={config.title}
                  color={config.color}
                  items={statusTasks.map(task => {
                    const isUserInvolved =
                      task.createdBy === user?.id ||
                      (task.assignees && task.assignees.some(u => u.id === user?.id));

                    return (
                      <KanbanCard
                        key={task.id}
                        id={task.id}
                        title={task.name}
                        description={task.description}
                        deadline={task.deadline}
                        assignees={task.assignees}
                        isUserInvolved={isUserInvolved}
                      />
                    );
                  })}
                  itemIds={statusTasks.map(t => t.id)}
                  showAddButton={false}
                />
              );
            })}
          </div>

          <DragOverlay>
            {activeTask && (
              <KanbanCard
                id={activeTask.id}
                title={activeTask.name}
                description={activeTask.description}
                deadline={activeTask.deadline}
                assignees={activeTask.assignees}
              />
            )}
          </DragOverlay>
        </DndContext>
      </div>

      {projectId && project && (
        <CreateTaskModal
          open={isCreateModalOpen}
          onOpenChange={setIsCreateModalOpen}
          projectId={projectId}
          projectAssignees={project.assignees || []}
        />
      )}
    </div>
  );
};

export default ProjectDetail;
