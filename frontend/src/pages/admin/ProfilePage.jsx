import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import { getMe } from '../../api/admin';
import apiClient from '../../api/client';
import TwoFASettingsCard from '../../components/TwoFASettingsCard';

const ProfilePage = () => {
    const [passwordForm, setPasswordForm] = useState({ current_password: '', new_password: '', confirm_password: '' });
    const [passwordErrors, setPasswordErrors] = useState({});

    const { data: meData, isLoading, refetch } = useQuery({
        queryKey: ['me'],
        queryFn: getMe,
    });

    const user = meData?.data;

    const changePasswordMutation = useMutation({
        mutationFn: async ({ current_password, new_password }) => {
            const stored = JSON.parse(localStorage.getItem('user') || '{}');
            await apiClient.put(`/users/${stored.id}`, { password: new_password });
        },
        onSuccess: () => {
            toast.success('Password changed successfully!');
            setPasswordForm({ current_password: '', new_password: '', confirm_password: '' });
            setPasswordErrors({});
        },
        onError: (err) => {
            toast.error(err.response?.data?.meta?.message || 'Failed to change password');
        },
    });

    const handlePasswordChange = (e) => {
        const { name, value } = e.target;
        setPasswordForm(prev => ({ ...prev, [name]: value }));
        if (passwordErrors[name]) setPasswordErrors(prev => ({ ...prev, [name]: '' }));
    };

    const handlePasswordSubmit = (e) => {
        e.preventDefault();
        const errors = {};
        if (!passwordForm.current_password) errors.current_password = 'Required';
        if (!passwordForm.new_password) errors.new_password = 'Required';
        if (passwordForm.new_password.length < 6) errors.new_password = 'Minimum 6 characters';
        if (passwordForm.new_password !== passwordForm.confirm_password) errors.confirm_password = 'Passwords do not match';
        if (Object.keys(errors).length > 0) { setPasswordErrors(errors); return; }
        changePasswordMutation.mutate(passwordForm);
    };

    const roleLabel = user?.role || 'User';
    const initial = (user?.email || '?').charAt(0).toUpperCase();

    return (
        <div className="max-w-3xl mx-auto">
            {/* Page Header */}
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-surface-on mb-1">My Profile</h1>
                <p className="text-surface-on-variant">Manage your account settings and security</p>
            </div>

            {isLoading ? (
                <div className="space-y-4">
                    {[1, 2, 3].map(i => (
                        <div key={i} className="h-24 bg-surface-container rounded-xl animate-pulse" />
                    ))}
                </div>
            ) : (
                <div className="space-y-6">
                    {/* Identity Card */}
                    <div className="bg-surface-container border border-outline-variant/30 rounded-xl p-6 flex items-center gap-6">
                        <div className="w-20 h-20 rounded-full bg-gradient-to-br from-primary to-primary/70 flex items-center justify-center text-3xl font-bold text-white flex-shrink-0 shadow-lg">
                            {initial}
                        </div>
                        <div className="flex-1 min-w-0">
                            <h2 className="text-xl font-bold text-surface-on truncate">{user?.name || user?.email}</h2>
                            <p className="text-surface-on-variant text-sm mt-0.5">{user?.email}</p>
                            <div className="flex flex-wrap gap-2 mt-3">
                                <span className="px-3 py-1 text-xs font-semibold bg-primary/10 text-primary rounded-full capitalize">
                                    {roleLabel}
                                </span>
                                <span className={`px-3 py-1 text-xs font-semibold rounded-full ${
                                    user?.status === 'active' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' :
                                    user?.status === 'freezed' ? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400' :
                                    'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
                                }`}>
                                    {user?.status ? user.status.charAt(0).toUpperCase() + user.status.slice(1) : 'Active'}
                                </span>
                                {user?.two_fa_enabled && (
                                    <span className="px-3 py-1 text-xs font-semibold bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 rounded-full">
                                        🔐 2FA Active
                                    </span>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Account Info */}
                    <div className="bg-surface-container border border-outline-variant/30 rounded-xl p-6">
                        <h3 className="text-base font-semibold text-surface-on mb-4">Account Information</h3>
                        <dl className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            {[
                                { label: 'User ID', value: `#${user?.id}` },
                                { label: 'Email', value: user?.email },
                                { label: 'Role', value: roleLabel },
                                { label: 'Status', value: user?.status ? user.status.charAt(0).toUpperCase() + user.status.slice(1) : 'Active' },
                            ].map(({ label, value }) => (
                                <div key={label} className="bg-surface-variant/20 rounded-lg px-4 py-3">
                                    <dt className="text-xs font-medium text-surface-on-variant mb-0.5">{label}</dt>
                                    <dd className="text-sm font-semibold text-surface-on">{value}</dd>
                                </div>
                            ))}
                        </dl>
                    </div>

                    {/* Change Password */}
                    <div className="bg-surface-container border border-outline-variant/30 rounded-xl p-6">
                        <h3 className="text-base font-semibold text-surface-on mb-4">Change Password</h3>
                        <form onSubmit={handlePasswordSubmit} className="space-y-4">
                            {[
                                { name: 'current_password', label: 'Current Password' },
                                { name: 'new_password', label: 'New Password' },
                                { name: 'confirm_password', label: 'Confirm New Password' },
                            ].map(({ name, label }) => (
                                <div key={name}>
                                    <label className="text-field-label">{label}</label>
                                    <input
                                        type="password"
                                        name={name}
                                        value={passwordForm[name]}
                                        onChange={handlePasswordChange}
                                        className={`text-field ${passwordErrors[name] ? 'border-error' : ''}`}
                                    />
                                    {passwordErrors[name] && (
                                        <p className="text-xs text-error mt-1">{passwordErrors[name]}</p>
                                    )}
                                </div>
                            ))}
                            <div className="pt-1">
                                <button
                                    type="submit"
                                    disabled={changePasswordMutation.isPending}
                                    className="px-5 py-2.5 bg-primary text-white text-sm font-semibold rounded-xl hover:bg-primary/90 transition-colors disabled:opacity-50"
                                >
                                    {changePasswordMutation.isPending ? 'Saving...' : 'Update Password'}
                                </button>
                            </div>
                        </form>
                    </div>

                    {/* 2FA Settings */}
                    <TwoFASettingsCard user={user} refetch={refetch} />
                </div>
            )}
        </div>
    );
};

export default ProfilePage;
