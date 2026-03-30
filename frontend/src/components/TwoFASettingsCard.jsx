import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import PropTypes from 'prop-types';
import Button from './Button';
import apiClient from '../api/client';

const TwoFASettingsCard = ({ user }) => {
    const queryClient = useQueryClient();
    const isEnabled = user?.two_fa_enabled ?? false;
    const [code, setCode] = useState('');
    const [enrollData, setEnrollData] = useState(null);
    const [stage, setStage] = useState('idle'); // idle | enrolling | confirming | disabling
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    const enrollMutation = useMutation({
        mutationFn: () => apiClient.post('/auth/2fa/enroll'),
        onSuccess: (res) => {
            setEnrollData(res.data.data);
            setStage('confirming');
            setError('');
        },
        onError: (err) => setError(err.response?.data?.meta?.message || 'Enrollment failed'),
    });

    const confirmMutation = useMutation({
        mutationFn: (code) => apiClient.post('/auth/2fa/confirm', { code }),
        onSuccess: () => {
            setMessage('2FA has been enabled successfully!');
            setStage('idle');
            setEnrollData(null);
            setCode('');
            queryClient.invalidateQueries({ queryKey: ['me'] });
            // Update localStorage user
            const stored = JSON.parse(localStorage.getItem('user') || '{}');
            localStorage.setItem('user', JSON.stringify({ ...stored, two_fa_enabled: true }));
        },
        onError: (err) => setError(err.response?.data?.meta?.message || 'Invalid code'),
    });

    const disableMutation = useMutation({
        mutationFn: (code) => apiClient.delete('/auth/2fa/disable', { data: { code } }),
        onSuccess: () => {
            setMessage('2FA has been disabled.');
            setStage('idle');
            setCode('');
            queryClient.invalidateQueries({ queryKey: ['me'] });
            const stored = JSON.parse(localStorage.getItem('user') || '{}');
            localStorage.setItem('user', JSON.stringify({ ...stored, two_fa_enabled: false }));
        },
        onError: (err) => setError(err.response?.data?.meta?.message || 'Invalid code'),
    });

    const handleConfirm = (e) => {
        e.preventDefault();
        setError('');
        confirmMutation.mutate(code);
    };

    const handleDisable = (e) => {
        e.preventDefault();
        setError('');
        disableMutation.mutate(code);
    };

    return (
        <div className="border border-outline-variant/30 rounded-xl p-6 bg-surface-container">
            <div className="flex items-center justify-between mb-4">
                <div>
                    <h3 className="text-base font-semibold text-surface-on">Two-Factor Authentication</h3>
                    <p className="text-sm text-surface-on-variant mt-0.5">
                        {isEnabled ? 'Your account is protected with HOTP 2FA.' : 'Add an extra layer of security to your account.'}
                    </p>
                </div>
                <span className={`px-2.5 py-1 text-xs font-semibold rounded-full ${isEnabled ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' : 'bg-surface-variant text-surface-on-variant'}`}>
                    {isEnabled ? '● Enabled' : '○ Disabled'}
                </span>
            </div>

            {message && (
                <div className="mb-4 text-sm text-green-700 bg-green-50 dark:bg-green-900/20 dark:text-green-400 border border-green-200 dark:border-green-800 px-3 py-2 rounded-lg">
                    {message}
                </div>
            )}
            {error && (
                <div className="mb-4 text-sm text-red-700 bg-red-50 dark:bg-red-900/20 dark:text-red-400 border border-red-200 dark:border-red-800 px-3 py-2 rounded-lg">
                    {error}
                </div>
            )}

            {stage === 'idle' && !isEnabled && (
                <Button onClick={() => { enrollMutation.mutate(); setStage('enrolling'); }} disabled={enrollMutation.isPending}>
                    {enrollMutation.isPending ? 'Generating...' : 'Enable 2FA'}
                </Button>
            )}

            {stage === 'idle' && isEnabled && (
                <Button variant="outline" onClick={() => { setStage('disabling'); setError(''); }}>
                    Disable 2FA
                </Button>
            )}

            {stage === 'confirming' && enrollData && (
                <div className="space-y-4">
                    <p className="text-sm text-surface-on-variant">
                        1. Scan this QR code with your Authenticator app (e.g. Google Authenticator, Authy).
                    </p>
                    <div className="bg-white p-3 inline-block rounded-lg border border-outline-variant/30">
                        <img
                            src={`https://api.qrserver.com/v1/create-qr-code/?data=${encodeURIComponent(enrollData.qr_url)}&size=180x180`}
                            alt="2FA QR Code"
                            className="w-44 h-44"
                        />
                    </div>
                    <p className="text-xs text-surface-on-variant">Can't scan? Use secret: <code className="bg-surface-variant px-1.5 py-0.5 rounded font-mono text-xs">{enrollData.secret}</code></p>
                    <p className="text-sm text-surface-on-variant">2. Enter the 6-digit code from your app to confirm:</p>
                    <form onSubmit={handleConfirm} className="flex gap-3 items-center">
                        <input
                            type="text"
                            inputMode="numeric"
                            maxLength={6}
                            value={code}
                            onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))}
                            placeholder="000000"
                            className="text-field text-center tracking-widest font-mono w-36"
                        />
                        <Button type="submit" disabled={confirmMutation.isPending || code.length !== 6}>
                            {confirmMutation.isPending ? 'Confirming...' : 'Confirm'}
                        </Button>
                        <button type="button" className="text-sm text-surface-on-variant hover:text-error" onClick={() => { setStage('idle'); setCode(''); setEnrollData(null); }}>
                            Cancel
                        </button>
                    </form>
                </div>
            )}

            {stage === 'disabling' && (
                <div className="space-y-3">
                    <p className="text-sm text-surface-on-variant">Enter your current 6-digit code to disable 2FA:</p>
                    <form onSubmit={handleDisable} className="flex gap-3 items-center">
                        <input
                            type="text"
                            inputMode="numeric"
                            maxLength={6}
                            value={code}
                            onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))}
                            placeholder="000000"
                            className="text-field text-center tracking-widest font-mono w-36"
                        />
                        <Button type="submit" disabled={disableMutation.isPending || code.length !== 6} className="btn-primary bg-red-600 hover:bg-red-700">
                            {disableMutation.isPending ? 'Disabling...' : 'Disable'}
                        </Button>
                        <button type="button" className="text-sm text-surface-on-variant hover:text-surface-on" onClick={() => { setStage('idle'); setCode(''); }}>
                            Cancel
                        </button>
                    </form>
                </div>
            )}
        </div>
    );
};

TwoFASettingsCard.propTypes = {
    user: PropTypes.object,
};

export default TwoFASettingsCard;
