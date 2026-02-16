import { useEffect, useState, useMemo } from 'react';
import { useNavigate, Outlet, useLocation } from 'react-router-dom';
import Sidebar from '../components/Sidebar';
import Button from '../components/Button';

const AdminLayout = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [user, setUser] = useState(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const userData = localStorage.getItem('user');

        if (!token) {
            navigate('/login');
            return;
        }

        if (userData) {
            setUser(JSON.parse(userData));
        }
    }, [navigate]);

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
    };

    // Navigation Configuration (Original Structure)
    const navigationSections = [
        {
            label: 'Main',
            items: [
                { path: '/dashboard', label: 'Dashboard', icon: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg> },
            ]
        },
        {
            label: 'Management',
            items: [
                {
                    label: 'Administrator',
                    icon: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>,
                    subItems: [
                        { path: '/admin/users', label: 'Users', icon: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /></svg> },
                        { path: '/admin/roles', label: 'Roles', icon: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" /></svg> },
                        { path: '/admin/permissions', label: 'Permissions', icon: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" /></svg> },
                    ]
                }
            ]
        }
    ];

    if (!user) {
        return (
            <div className="min-h-screen bg-surface flex items-center justify-center">
                <div className="text-primary-500 animate-pulse">Loading MD3 Expressive...</div>
            </div>
        );
    }

    return (
        <div className="flex h-screen bg-surface overflow-hidden">
            {/* Sidebar Navigation */}
            <Sidebar sections={navigationSections} onLogout={handleLogout} />

            {/* Main Content Area */}
            <div className="flex-1 flex flex-col min-w-0 bg-surface-container-low relative">
                {/* Compact Header */}
                <header className="h-16 flex items-center justify-between px-8 bg-surface-container-low border-b border-outline-variant/30">
                    <h2 className="text-lg font-medium text-surface-on">Administration</h2>
                    <div className="flex items-center gap-6">
                        <div className="flex flex-col items-end mr-2">
                            <span className="text-sm font-medium text-surface-on">{user.email}</span>
                            <span className="text-[11px] text-surface-on-variant font-medium">System Administrator</span>
                        </div>
                        <div className="w-10 h-10 rounded-full bg-primary-container flex items-center justify-center text-primary-on-container shadow-md3-1 cursor-pointer hover:shadow-md3-2 transition-all">
                            {user.email.charAt(0).toUpperCase()}
                        </div>
                    </div>
                </header>

                {/* Scrollable Page Content */}
                <main className="flex-1 overflow-y-auto p-8 custom-scrollbar">
                    <div className="max-w-7xl mx-auto animate-fade-in-up">
                        <Outlet />
                    </div>
                </main>
            </div>
        </div>
    );
};

export default AdminLayout;
