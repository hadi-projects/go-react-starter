import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import Card from '../components/Card';
import Button from '../components/Button';
import apiClient from '../api/client';
import toast from 'react-hot-toast';

const TwoFAResetConfirmPage = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const token = searchParams.get('token');
    const [isConfirmed, setIsConfirmed] = useState(false);

    const confirmMutation = useMutation({
        mutationFn: async (token) => {
            const res = await apiClient.post('/auth/2fa/reset-confirm', { token });
            return res.data;
        },
        onSuccess: (data) => {
            toast.success(data.meta?.message || '2FA disabled successfully!');
            setIsConfirmed(true);
        },
        onError: (err) => {
            toast.error(err.response?.data?.meta?.message || 'Failed to disable 2FA.');
        },
    });

    const handleConfirm = () => {
        if (!token) {
            toast.error('Token is missing.');
            return;
        }
        confirmMutation.mutate(token);
    };

    if (isConfirmed) {
        return (
            <div className="min-h-screen bg-surface-variant flex items-center justify-center p-6">
                <div className="w-full max-w-md">
                    <Card className="p-8 text-center space-y-6">
                        <div className="text-6xl text-center">✅</div>
                        <h1 className="text-2xl font-bold text-surface-on">2FA Disabled Successfully</h1>
                        <p className="text-surface-on-variant px-4">
                            Your account is now accessible without 2FA. You can re-enable it later in your profile settings.
                        </p>
                        <div className="pt-4">
                            <Button 
                                onClick={() => navigate('/login')} 
                                fullWidth
                            >
                                Continue to Login
                            </Button>
                        </div>
                    </Card>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-surface-variant flex items-center justify-center p-6">
            <div className="w-full max-w-sm">
                <div className="text-center mb-8">
                    <div className="text-6xl mb-4">🛡️</div>
                    <h1 className="text-2xl font-bold text-surface-on mb-2">Disable 2FA</h1>
                    <p className="text-surface-on-variant text-sm px-4">
                        Confirm you want to disable Two-Factor Authentication for your account.
                    </p>
                </div>
                <Card className="p-8">
                    <div className="space-y-6">
                        <Button 
                            onClick={handleConfirm} 
                            fullWidth 
                            disabled={confirmMutation.isPending || !token}
                        >
                            {confirmMutation.isPending ? 'Processing...' : 'Confirm & Disable 2FA'}
                        </Button>
                        <div className="text-center">
                            <Link 
                                to="/login" 
                                className="text-sm text-surface-on-variant hover:text-primary transition-colors font-medium"
                            >
                                ← Back to Login
                            </Link>
                        </div>
                    </div>
                </Card>
            </div>
        </div>
    );
};

export default TwoFAResetConfirmPage;
