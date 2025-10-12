import { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api';

const HealthCheck = () => {
  const [status, setStatus] = useState<'checking' | 'healthy' | 'error'>('checking');
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const checkHealth = async () => {
      try {
        // Try to make a simple request to see if the backend is running
        const response = await fetch('http://localhost:8080/health');
        if (response.ok) {
          setStatus('healthy');
        } else {
          setStatus('error');
          setError(`HTTP ${response.status}: ${response.statusText}`);
        }
      } catch (err) {
        setStatus('error');
        setError(err instanceof Error ? err.message : 'Unknown error');
      }
    };

    checkHealth();
  }, []);

  if (status === 'checking') {
    return <div className="text-sm text-yellow-600">Checking backend connection...</div>;
  }

  if (status === 'error') {
    return (
      <div className="text-sm text-red-600">
        Backend connection failed: {error}
      </div>
    );
  }

  return <div className="text-sm text-green-600">Backend connected ✓</div>;
};

export default HealthCheck;
