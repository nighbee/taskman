import { useState } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useDataStore, type User, type ProjectStatus } from '@/store/dataStore';
import { useAuthStore } from '@/store/authStore';
import { toast } from 'sonner';

interface CreateProjectModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  orgId: string;
  orgMembers: User[];
}

export const CreateProjectModal = ({
  open,
  onOpenChange,
  orgId,
  orgMembers,
}: CreateProjectModalProps) => {
  const { user } = useAuthStore();
  const { createProject } = useDataStore();
  
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [deadline, setDeadline] = useState('');
  const [status, setStatus] = useState<ProjectStatus>('idea');
  const [selectedAssignees, setSelectedAssignees] = useState<string[]>([]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!name.trim()) {
      toast.error('Please enter a project name');
      return;
    }

    if (!user) return;

    try {
      // Convert date to ISO string if provided
      let deadlineISO: string | undefined = undefined;
      if (deadline) {
        // Convert YYYY-MM-DD to ISO datetime string
        deadlineISO = new Date(deadline + 'T23:59:59.000Z').toISOString();
      }

      await createProject(
        orgId,
        name,
        description,
        deadlineISO,
        selectedAssignees,
        user.id,
        status
      );

      toast.success('Project created successfully!');
      onOpenChange(false);
      
      // Reset form
      setName('');
      setDescription('');
      setDeadline('');
      setStatus('idea');
      setSelectedAssignees([]);
    } catch (error) {
      console.error('Error creating project:', error);
      toast.error('Failed to create project');
    }
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
          <DialogTitle>Create New Project</DialogTitle>
          <DialogDescription>
            Add a new project to your organization
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Project Name *</Label>
            <Input
              id="name"
              placeholder="Website Redesign"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              placeholder="Describe the project goals..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="deadline">Deadline</Label>
            <Input
              id="deadline"
              type="date"
              value={deadline}
              onChange={(e) => setDeadline(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="status">Project Stage</Label>
            <Select value={status} onValueChange={(value: ProjectStatus) => setStatus(value)}>
              <SelectTrigger>
                <SelectValue placeholder="Select project stage" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="idea">💡 Idea</SelectItem>
                <SelectItem value="in-progress">🚀 In Progress</SelectItem>
                <SelectItem value="finished">✅ Finished</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Assign Team Members</Label>
            <p className="text-xs text-muted-foreground">
              You will be automatically assigned to this project as the creator.
            </p>
            <div className="space-y-2 max-h-[200px] overflow-y-auto border rounded-md p-3">
              {(() => {
                // Create a list that includes current user + org members, removing duplicates
                const allMembers = [...(orgMembers || [])];
                
                // Add current user if not already in the list
                if (user && !allMembers.find(member => member.id === user.id)) {
                  allMembers.unshift(user);
                }
                
                return allMembers.length > 0 ? (
                  allMembers.map(member => (
                    <div key={member.id} className="flex items-center space-x-2">
                      <Checkbox
                        id={`member-${member.id}`}
                        checked={selectedAssignees.includes(member.id)}
                        onCheckedChange={() => toggleAssignee(member.id)}
                      />
                      <label
                        htmlFor={`member-${member.id}`}
                        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                      >
                        {member.full_name || member.email} ({member.email})
                        {member.id === user?.id && (
                          <span className="ml-1 text-xs text-primary font-medium">(You)</span>
                        )}
                      </label>
                    </div>
                  ))
                ) : (
                  <p className="text-sm text-muted-foreground">No team members available</p>
                );
              })()}
            </div>
          </div>

          <div className="flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Create Project</Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};
