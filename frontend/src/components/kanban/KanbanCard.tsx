import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Calendar, User, GripVertical, ChevronLeft, ChevronRight } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import type { User as UserType, ProjectStatus } from '@/store/dataStore';

interface KanbanCardProps {
  id: string;
  title: string;
  description?: string;
  deadline?: string;
  assignees?: UserType[];
  isUserInvolved?: boolean;
  currentStatus?: ProjectStatus;
  onStatusChange?: (newStatus: ProjectStatus) => void;
  onClick?: () => void;
}

export const KanbanCard = ({
  id,
  title,
  description,
  deadline,
  assignees = [],
  isUserInvolved = false,
  currentStatus,
  onStatusChange,
  onClick,
}: KanbanCardProps) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const statusOrder: ProjectStatus[] = ['idea', 'in-progress', 'finished'];
  const currentIndex = currentStatus ? statusOrder.indexOf(currentStatus) : -1;
  const canMoveLeft = currentIndex > 0;
  const canMoveRight = currentIndex < statusOrder.length - 1;

  const handleStatusChange = (direction: 'left' | 'right') => {
    if (!currentStatus || !onStatusChange) return;
    
    const newIndex = direction === 'left' ? currentIndex - 1 : currentIndex + 1;
    if (newIndex >= 0 && newIndex < statusOrder.length) {
      onStatusChange(statusOrder[newIndex]);
    }
  };

  return (
    <Card
      ref={setNodeRef}
      style={style}
      className={`cursor-pointer transition-all hover:shadow-md ${
        isUserInvolved ? 'ring-2 ring-primary/50' : ''
      }`}
      onClick={onClick}
    >
      <CardHeader className="p-4 pb-2">
        <div className="flex items-start justify-between gap-2">
          <CardTitle className="text-sm font-semibold line-clamp-2 flex-1">
            {title}
          </CardTitle>
          <div className="flex items-center gap-1">
            {/* Status Control Buttons */}
            {isUserInvolved && onStatusChange && (
              <>
                <Button
                  size="sm"
                  variant="ghost"
                  className="h-6 w-6 p-0"
                  disabled={!canMoveLeft}
                  onClick={(e) => {
                    e.stopPropagation();
                    handleStatusChange('left');
                  }}
                >
                  <ChevronLeft className="h-3 w-3" />
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  className="h-6 w-6 p-0"
                  disabled={!canMoveRight}
                  onClick={(e) => {
                    e.stopPropagation();
                    handleStatusChange('right');
                  }}
                >
                  <ChevronRight className="h-3 w-3" />
                </Button>
              </>
            )}
            <button
              {...attributes}
              {...listeners}
              className="cursor-grab active:cursor-grabbing p-1 hover:bg-muted rounded"
              onClick={(e) => e.stopPropagation()}
            >
              <GripVertical className="h-4 w-4 text-muted-foreground" />
            </button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-4 pt-0">
        {description && (
          <p className="text-xs text-muted-foreground line-clamp-2 mb-3">
            {description}
          </p>
        )}
        
        <div className="flex flex-col gap-2">
          {deadline && (
            <div className="flex items-center gap-2 text-xs text-muted-foreground">
              <Calendar className="h-3 w-3" />
              <span>{new Date(deadline).toLocaleDateString()}</span>
            </div>
          )}
          
          {assignees && assignees.length > 0 && (
            <div className="flex items-center gap-2">
              <User className="h-3 w-3 text-muted-foreground" />
              <div className="flex flex-wrap gap-1">
                {assignees.map(assignee => (
                  <Badge key={assignee.id} variant="secondary" className="text-xs">
                    {assignee.full_name ? assignee.full_name.split(' ')[0] : assignee.email}
                  </Badge>
                ))}
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};
