import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import Card from '../components/Card';
import Button from '../components/Button';
import apiClient from '../api/client';

const TwoFAChallengePage = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const tempToken = location.state?.tempToken;
    const [code, setCode] = useState('');
    const [error, setError] = useState('');

    const verifyMutation = useMutation({
        mutationFn: async ({ temp_token, code }) => {
            const res = await apiClient.post('/auth/2fa/verify', { temp_token, code });
            return res.data;
        },
        onSuccess: (data) => {
            localStorage.setItem('token', data.data.access_token);
            if (data.data.refresh_token) {
                localStorage.setItem('refresh_token', data.data.refresh_token);
            }
            localStorage.setItem('user', JSON.stringify(data.data.user));
            navigate('/dashboard');
        },
        onError: (err) => {
            setError(err.response?.data?.meta?.message || 'Invalid code. Please try again.');
        },
    });

    const handleSubmit = (e) => {
        e.preventDefault();
        setError('');
        if (!tempToken) {
            setError('Session expired. Please login again.');
            return;
        }
        verifyMutation.mutate({ temp_token: tempToken, code });
    };

    return (
        <div className="min-h-screen bg-surface-variant flex items-center justify-center p-6">
            <div className="w-full max-w-sm">
                <div className="text-center mb-8">
                    <div className="text-6xl mb-4">🔐</div>
                    <h1 className="text-2xl font-bold text-surface-on mb-2">Two-Factor Authentication</h1>
                    <p className="text-surface-on-variant text-sm">Enter the 6-digit code from your Authenticator app</p>
                </div>
                <Card className="p-8">
                    <form onSubmit={handleSubmit} className="space-y-5">
                        {error && (
                            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md3 text-sm dark:bg-red-900/20 dark:border-red-800 dark:text-red-400">
                                {error}
                            </div>
                        )}
                        <div>
                            <label className="text-field-label">Authentication Code</label>
                            <input
                                type="text"
                                inputMode="numeric"
                                pattern="[0-9]{6}"
                                maxLength={6}
                                value={code}
                                onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))}
                                placeholder="000000"
                                className="text-field text-center text-2xl tracking-[0.5em] font-mono"
                                autoFocus
                            />
                        </div>
                        <Button type="submit" fullWidth disabled={verifyMutation.isPending || code.length !== 6}>
                            {verifyMutation.isPending ? 'Verifying...' : 'Verify Code'}
                        </Button>
                        <div className="flex flex-col items-center space-y-3 pt-2">
                            <button
                                type="button"
                                onClick={() => navigate('/twofa/reset-request', { state: { tempToken } })}
                                className="text-sm text-primary hover:underline transition-all"
                            >
                                Lost access to OTP?
                            </button>
                            <button
                                type="button"
                                onClick={() => navigate('/login')}
                                className="w-full text-sm text-surface-on-variant hover:text-primary transition-colors"
                            >
                                ← Back to Login
                            </button>
                        </div>
                    </form>
                </Card>
            </div>
        </div>
    );
};

export default TwoFAChallengePage;
