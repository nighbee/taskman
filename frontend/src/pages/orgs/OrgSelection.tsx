import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { useAuthStore } from '@/store/authStore';
import { useDataStore } from '@/store/dataStore';
import { Building2, Plus, Key, Users } from 'lucide-react';
import { toast } from 'sonner';
import HealthCheck from '@/components/HealthCheck';

const OrgSelection = () => {
  const { user } = useAuthStore();
  const { orgs, createOrg, joinOrg, selectOrg, isLoading, error, loadOrganizations } = useDataStore();
  const navigate = useNavigate();
  
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isJoinOpen, setIsJoinOpen] = useState(false);
  const [newOrgName, setNewOrgName] = useState('');
  const [inviteCode, setInviteCode] = useState('');

  // Defensive programming: ensure orgs is always an array
  const safeOrgs = Array.isArray(orgs) ? orgs : [];
  
  // For now, show all orgs since the backend doesn't include members in the response
  // TODO: Update backend to include members or create a separate endpoint
  const userOrgs = safeOrgs;

  const handleCreateOrg = async () => {
    if (!newOrgName.trim() || !user) return;
    
    try {
      const org = await createOrg(newOrgName);
      console.log('Created org:', org); // Debug log
      
      if (!org.id) {
        toast.error('Organization created but ID not returned');
        return;
      }
      
      toast.success(`Organization "${org.name}" created!`);
      setIsCreateOpen(false);
      setNewOrgName('');
      selectOrg(org.id);
      navigate(`/orgs/${org.id}/projects`);
    } catch (error) {
      console.error('Error creating org:', error); // Debug log
      toast.error('Failed to create organization');
    }
  };

  const handleJoinOrg = async () => {
    if (!inviteCode.trim() || !user) return;
    
    try {
      const success = await joinOrg(inviteCode.toUpperCase());
      
      if (success) {
        toast.success('Successfully joined organization!');
        setIsJoinOpen(false);
        setInviteCode('');
      } else {
        toast.error('Invalid invite code');
      }
    } catch (error) {
      toast.error('Failed to join organization');
    }
  };

  const handleSelectOrg = (orgId: string) => {
    selectOrg(orgId);
    navigate(`/orgs/${orgId}/projects`);
  };

  return (
    <div className="min-h-[calc(100vh-4rem)] px-6 py-12 bg-gradient-subtle">
      <div className="container max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Your Organizations</h1>
          <p className="text-muted-foreground">
            Select an organization to start managing projects
          </p>
          <div className="mt-2">
            <HealthCheck />
          </div>
          {error && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-md">
              <p className="text-red-800 text-sm">{error}</p>
              <button 
                onClick={loadOrganizations}
                className="mt-2 text-red-600 hover:text-red-800 text-sm underline"
              >
                Retry
              </button>
            </div>
          )}
        </div>

        <div className="grid gap-6 mb-6">
          {userOrgs.length === 0 && (
            <Card className="border-dashed">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <Building2 className="h-12 w-12 text-muted-foreground mb-4" />
                <p className="text-muted-foreground text-center mb-6">
                  You haven't joined any organizations yet.
                  <br />
                  Create one or join with an invite code.
                </p>
              </CardContent>
            </Card>
          )}

          {userOrgs.map((org, index) => (
            <Card
              key={org.id || `org-${index}`}
              className="cursor-pointer hover:shadow-lg transition-shadow"
              onClick={() => org.id && handleSelectOrg(org.id)}
            >
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="text-xl">{org.name}</CardTitle>
                    <CardDescription className="mt-2 flex items-center gap-2">
                      <Users className="h-4 w-4" />
                      {org.members ? org.members.length : 'Unknown'} member{org.members && org.members.length !== 1 ? 's' : ''}
                    </CardDescription>
                  </div>
                  <div className="text-xs text-muted-foreground">
                    Code: <span className="font-mono font-semibold">{org.invite_code}</span>
                  </div>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>

        <div className="flex flex-col sm:flex-row gap-4">
          <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
            <DialogTrigger asChild>
              <Button className="flex-1 gap-2">
                <Plus className="h-4 w-4" />
                Create Organization
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create New Organization</DialogTitle>
                <DialogDescription>
                  Start a new workspace for your team
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4 pt-4">
                <div className="space-y-2">
                  <Label htmlFor="orgName">Organization Name</Label>
                  <Input
                    id="orgName"
                    placeholder="Acme Corporation"
                    value={newOrgName}
                    onChange={(e) => setNewOrgName(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleCreateOrg()}
                  />
                </div>
                <Button onClick={handleCreateOrg} className="w-full" disabled={isLoading}>
                  {isLoading ? 'Creating...' : 'Create'}
                </Button>
              </div>
            </DialogContent>
          </Dialog>

          <Dialog open={isJoinOpen} onOpenChange={setIsJoinOpen}>
            <DialogTrigger asChild>
              <Button variant="outline" className="flex-1 gap-2">
                <Key className="h-4 w-4" />
                Join with Code
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Join Organization</DialogTitle>
                <DialogDescription>
                  Enter the invite code provided by your team
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4 pt-4">
                <div className="space-y-2">
                  <Label htmlFor="inviteCode">Invite Code</Label>
                  <Input
                    id="inviteCode"
                    placeholder="ACME99"
                    value={inviteCode}
                    onChange={(e) => setInviteCode(e.target.value.toUpperCase())}
                    onKeyDown={(e) => e.key === 'Enter' && handleJoinOrg()}
                    className="font-mono"
                  />
                </div>
                <Button onClick={handleJoinOrg} className="w-full" disabled={isLoading}>
                  {isLoading ? 'Joining...' : 'Join Organization'}
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </div>
  );
};

export default OrgSelection;
