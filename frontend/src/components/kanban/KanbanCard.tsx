import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Calendar, User, GripVertical } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { User as UserType } from '@/store/dataStore';

interface KanbanCardProps {
  id: string;
  title: string;
  description?: string;
  deadline?: string;
  assignees?: UserType[];
  isUserInvolved?: boolean;
  onClick?: () => void;
}

export const KanbanCard = ({
  id,
  title,
  description,
  deadline,
  assignees = [],
  isUserInvolved = false,
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
          <CardTitle className="text-sm font-semibold line-clamp-2">
            {title}
          </CardTitle>
          <button
            {...attributes}
            {...listeners}
            className="cursor-grab active:cursor-grabbing p-1 hover:bg-muted rounded"
            onClick={(e) => e.stopPropagation()}
          >
            <GripVertical className="h-4 w-4 text-muted-foreground" />
          </button>
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
          
          {assignees.length > 0 && (
            <div className="flex items-center gap-2">
              <User className="h-3 w-3 text-muted-foreground" />
              <div className="flex flex-wrap gap-1">
                {assignees.map(assignee => (
                  <Badge key={assignee.id} variant="secondary" className="text-xs">
                    {assignee.fullName.split(' ')[0]}
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
