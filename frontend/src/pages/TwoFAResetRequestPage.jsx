import { useState } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import Card from '../components/Card';
import Button from '../components/Button';
import apiClient from '../api/client';
import toast from 'react-hot-toast';

const TwoFAResetRequestPage = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const tempToken = location.state?.tempToken;
    const [isSent, setIsSent] = useState(false);

    const resetRequestMutation = useMutation({
        mutationFn: async (temp_token) => {
            const res = await apiClient.post('/auth/2fa/reset-request', { temp_token });
            return res.data;
        },
        onSuccess: (data) => {
            toast.success(data.meta?.message || 'Reset link sent to your email!');
            setIsSent(true);
        },
        onError: (err) => {
            toast.error(err.response?.data?.meta?.message || 'Failed to send reset link.');
        },
    });

    const handleRequest = () => {
        if (!tempToken) {
            toast.error('Session expired. Please login again.');
            navigate('/login');
            return;
        }
        resetRequestMutation.mutate(tempToken);
    };

    if (isSent) {
        return (
            <div className="min-h-screen bg-surface-variant flex items-center justify-center p-6">
                <div className="w-full max-w-md">
                    <Card className="p-8 text-center space-y-6">
                        <div className="text-6xl">📧</div>
                        <h1 className="text-2xl font-bold text-surface-on">Check your email</h1>
                        <p className="text-surface-on-variant px-4">
                            We've sent a link to disable 2FA to your registered email address.
                        </p>
                        <div className="pt-4">
                            <Link 
                                to="/login" 
                                className="text-primary hover:underline font-medium"
                            >
                                ← Back to Login
                            </Link>
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
                    <h1 className="text-2xl font-bold text-surface-on mb-2">Reset 2FA Access</h1>
                    <p className="text-surface-on-variant text-sm px-4">
                        Lost access to your authenticator? We'll send a link to your email to safely disable 2FA.
                    </p>
                </div>
                <Card className="p-8">
                    <div className="space-y-6">
                        <div className="bg-amber-50 border border-amber-200 text-amber-700 px-4 py-3 rounded-md3 text-xs dark:bg-amber-900/20 dark:border-amber-800 dark:text-amber-400">
                           By continuing, you confirm you've lost access to your 2FA device.
                        </div>
                        <Button 
                            onClick={handleRequest} 
                            fullWidth 
                            disabled={resetRequestMutation.isPending}
                        >
                            {resetRequestMutation.isPending ? 'Sending...' : 'Send Reset Link'}
                        </Button>
                        <button
                            type="button"
                            onClick={() => navigate('/2fa-challenge', { state: { tempToken } })}
                            className="w-full text-sm text-surface-on-variant hover:text-primary transition-colors"
                        >
                            ← Back to Challenge
                        </button>
                    </div>
                </Card>
            </div>
        </div>
    );
};

export default TwoFAResetRequestPage;
