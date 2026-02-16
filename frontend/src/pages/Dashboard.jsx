import { Link } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { toast } from 'react-hot-toast';
import Button from '../components/Button';
import Card from '../components/Card';
import { clearCache } from '../api/admin';

const Dashboard = () => {
    const clearCacheMutation = useMutation({
        mutationFn: clearCache,
        onSuccess: () => {
            toast.success('Cache cleared successfully!');
        },
        onError: (error) => {
            toast.error(error.response?.data?.meta?.message || 'Failed to clear cache');
        },
    });

    const handleClearCache = () => {
        if (window.confirm('Are you sure you want to clear all cache? This action cannot be undone.')) {
            clearCacheMutation.mutate();
        }
    };

    // Dummy statistics data
    const stats = [
        {
            id: 1,
            title: 'Total Users',
            value: '254',
            change: '+12%',
            trend: 'up',
            icon: (
                <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
            ),
            color: 'bg-blue-500'
        },
        {
            id: 2,
            title: 'Total Roles',
            value: '8',
            change: '+2',
            trend: 'up',
            icon: (
                <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
            ),
            color: 'bg-green-500'
        },
        {
            id: 3,
            title: 'Permissions',
            value: '42',
            change: '+5',
            trend: 'up',
            icon: (
                <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                </svg>
            ),
            color: 'bg-purple-500'
        },
        {
            id: 4,
            title: 'Active Sessions',
            value: '89',
            change: '-3%',
            trend: 'down',
            icon: (
                <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
            ),
            color: 'bg-orange-500'
        }
    ];

    const recentActivities = [
        { id: 1, action: 'User john@example.com registered', time: '2 minutes ago' },
        { id: 2, action: 'Role "Manager" created', time: '15 minutes ago' },
        { id: 3, action: 'Permission "edit-user" updated', time: '1 hour ago' },
        { id: 4, action: 'User jane@example.com logged in', time: '2 hours ago' },
        { id: 5, action: 'System backup completed', time: '3 hours ago' }
    ];

    return (
        <div>
            {/* Page Header */}
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
                <p className="text-gray-600 mt-2">Welcome back! Here's what's happening today.</p>
            </div>

            {/* Statistics Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                {stats.map((stat) => (
                    <Card key={stat.id}>
                        <div className="flex items-start justify-between">
                            <div className="flex-1">
                                <p className="text-sm text-gray-600 mb-1">{stat.title}</p>
                                <h3 className="text-3xl font-bold text-gray-900 mb-2">{stat.value}</h3>
                                <div className="flex items-center gap-1">
                                    <span className={`text-sm font-medium ${stat.trend === 'up' ? 'text-green-600' : 'text-red-600'
                                        }`}>
                                        {stat.change}
                                    </span>
                                    <span className="text-xs text-gray-500">from last month</span>
                                </div>
                            </div>
                            <div className={`${stat.color} p-3 rounded-xl text-white`}>
                                {stat.icon}
                            </div>
                        </div>
                    </Card>
                ))}
            </div>

            {/* Two Column Layout */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                {/* Recent Activity */}
                <Card className="lg:col-span-2">
                    <h2 className="text-xl font-bold text-gray-900 mb-4">Recent Activity</h2>
                    <div className="space-y-4">
                        {recentActivities.map((activity) => (
                            <div key={activity.id} className="flex items-start gap-3 pb-4 border-b border-gray-100 last:border-0">
                                <div className="w-2 h-2 bg-primary-500 rounded-full mt-2"></div>
                                <div className="flex-1">
                                    <p className="text-gray-900">{activity.action}</p>
                                    <p className="text-sm text-gray-500 mt-1">{activity.time}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                </Card>

                {/* Quick Actions */}
                <Card>
                    <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
                    <div className="space-y-3">
                        <Link to="/admin/users">
                            <Button fullWidth className="justify-start">
                                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                                </svg>
                                Manage Users
                            </Button>
                        </Link>
                        <Link to="/admin/roles">
                            <Button variant="secondary" fullWidth className="justify-start">
                                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                                </svg>
                                Manage Roles
                            </Button>
                        </Link>
                        <Link to="/admin/permissions">
                            <Button variant="outline" fullWidth className="justify-start">
                                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                                </svg>
                                Manage Permissions
                            </Button>
                        </Link>

                        {/* Clear Cache Button */}
                        <div className="pt-3 border-t border-gray-200">
                            <Button
                                variant="danger"
                                fullWidth
                                className="justify-start"
                                onClick={handleClearCache}
                                disabled={clearCacheMutation.isPending}
                            >
                                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                </svg>
                                {clearCacheMutation.isPending ? 'Clearing...' : 'Clear Cache'}
                            </Button>
                        </div>
                    </div>
                </Card>
            </div>
        </div>
    );
};

export default Dashboard;
