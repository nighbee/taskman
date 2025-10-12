import { useState } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';
import { useDataStore, type User } from '@/store/dataStore';
import { useAuthStore } from '@/store/authStore';
import { toast } from 'sonner';

interface CreateTaskModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  projectId: string;
  projectAssignees: User[];
}

export const CreateTaskModal = ({
  open,
  onOpenChange,
  projectId,
  projectAssignees,
}: CreateTaskModalProps) => {
  const { user } = useAuthStore();
  const { createTask } = useDataStore();
  
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [deadline, setDeadline] = useState('');
  const [selectedAssignees, setSelectedAssignees] = useState<string[]>([]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!name.trim()) {
      toast.error('Please enter a task name');
      return;
    }

    if (!user) return;

    createTask(
      projectId,
      name,
      description,
      deadline,
      selectedAssignees,
      user.id
    );

    toast.success('Task created successfully!');
    onOpenChange(false);
    
    // Reset form
    setName('');
    setDescription('');
    setDeadline('');
    setSelectedAssignees([]);
  };

  const toggleAssignee = (userId: string) => {
    setSelectedAssignees(prev =>
      prev.includes(userId)
        ? prev.filter(id => id !== userId)
        : [...prev, userId]
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Create New Task</DialogTitle>
          <DialogDescription>
            Add a new task to this project
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="taskName">Task Name *</Label>
            <Input
              id="taskName"
              placeholder="Design homepage mockups"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="taskDescription">Description</Label>
            <Textarea
              id="taskDescription"
              placeholder="Describe what needs to be done..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="taskDeadline">Deadline</Label>
            <Input
              id="taskDeadline"
              type="date"
              value={deadline}
              onChange={(e) => setDeadline(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label>Assign to (Project Members Only)</Label>
            <div className="space-y-2 max-h-[150px] overflow-y-auto border rounded-md p-3">
              {projectAssignees && projectAssignees.length > 0 ? (
                projectAssignees.map(member => (
                  <div key={member.id} className="flex items-center space-x-2">
                    <Checkbox
                      id={`task-member-${member.id}`}
                      checked={selectedAssignees.includes(member.id)}
                      onCheckedChange={() => toggleAssignee(member.id)}
                    />
                    <label
                      htmlFor={`task-member-${member.id}`}
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                    >
                      {member.full_name || member.email}
                    </label>
                  </div>
                ))
              ) : (
                <p className="text-sm text-muted-foreground">No project members available</p>
              )}
            </div>
          </div>

          <div className="flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Create Task</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
