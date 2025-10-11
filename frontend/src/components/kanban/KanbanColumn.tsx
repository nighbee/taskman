import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface KanbanColumnProps {
  id: string;
  title: string;
  color: 'muted' | 'warning' | 'success';
  items: React.ReactNode[];
  itemIds: string[];
  onAddClick?: () => void;
  showAddButton?: boolean;
}

const colorClasses = {
  muted: 'bg-muted text-muted-foreground',
  warning: 'bg-warning/20 text-warning-foreground',
  success: 'bg-success/20 text-success-foreground',
};

export const KanbanColumn = ({
  id,
  title,
  color,
  items,
  itemIds,
  onAddClick,
  showAddButton = true,
}: KanbanColumnProps) => {
  const { setNodeRef, isOver } = useDroppable({ id });

  return (
    <div className="flex flex-col h-full min-w-[320px] bg-muted/30 rounded-lg p-4">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <h3 className="font-semibold text-sm uppercase tracking-wide">
            {title}
          </h3>
          <Badge variant="outline" className={colorClasses[color]}>
            {items.length}
          </Badge>
        </div>
        {showAddButton && onAddClick && (
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={onAddClick}
          >
            <Plus className="h-4 w-4" />
          </Button>
        )}
      </div>

      <div
        ref={setNodeRef}
        className={`flex-1 space-y-3 overflow-y-auto pr-2 min-h-[200px] transition-colors ${
          isOver ? 'bg-primary/5 rounded-md' : ''
        }`}
      >
        <SortableContext items={itemIds} strategy={verticalListSortingStrategy}>
          {items}
        </SortableContext>
      </div>
    </div>
  );
};
