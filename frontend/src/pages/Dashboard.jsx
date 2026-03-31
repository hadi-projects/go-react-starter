import { Link } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import { toast } from "react-hot-toast";
import Button from "../components/Button";
import Card from "../components/Card";
import { clearCache } from "../api/admin";
import { getDashboardStats } from "../api/statistics";

const Dashboard = () => {
  const clearCacheMutation = useMutation({
    mutationFn: clearCache,
    onSuccess: () => {
      toast.success("Cache cleared successfully!");
    },
    onError: (error) => {
      toast.error(
        error.response?.data?.meta?.message || "Failed to clear cache",
      );
    },
  });

  const handleClearCache = () => {
    if (
      window.confirm(
        "Are you sure you want to clear all cache? This action cannot be undone.",
      )
    ) {
      clearCacheMutation.mutate();
    }
  };

  // Fetch real statistics from API
  const { data: statsData, isLoading } = useQuery({
    queryKey: ["dashboard-stats"],
    queryFn: getDashboardStats,
  });

  // Statistics data with smooth colors - refined for both light and dark modes
  const stats = [
    {
      id: 1,
      title: "Total Users",
      value: statsData?.data?.total_users?.toString() || "0",
      icon: (
        <svg
          className="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
          />
        </svg>
      ),
      gradient:
        "from-blue-50 to-indigo-50 dark:from-blue-900/10 dark:to-indigo-900/10",
      iconBg:
        "bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-800/20 dark:to-indigo-800/20",
      iconColor: "text-blue-600 dark:text-blue-400",
    },
    {
      id: 2,
      title: "Total Roles",
      value: statsData?.data?.total_roles?.toString() || "0",
      icon: (
        <svg
          className="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
          />
        </svg>
      ),
      gradient:
        "from-purple-50 to-pink-50 dark:from-purple-900/10 dark:to-pink-900/10",
      iconBg:
        "bg-gradient-to-br from-purple-100 to-pink-100 dark:from-purple-800/20 dark:to-pink-800/20",
      iconColor: "text-purple-600 dark:text-purple-400",
    },
    {
      id: 3,
      title: "Permissions",
      value: statsData?.data?.total_permissions?.toString() || "0",
      icon: (
        <svg
          className="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
          />
        </svg>
      ),
      gradient:
        "from-amber-50 to-orange-50 dark:from-amber-900/10 dark:to-orange-900/10",
      iconBg:
        "bg-gradient-to-br from-amber-100 to-orange-100 dark:from-amber-800/20 dark:to-orange-800/20",
      iconColor: "text-amber-600 dark:text-amber-400",
    },
  ];

  const recentActivities = [
    {
      id: 1,
      action: "User john@example.com registered",
      time: "2 minutes ago",
    },
    { id: 2, action: 'Role "Manager" created', time: "15 minutes ago" },
    { id: 3, action: 'Permission "edit-user" updated', time: "1 hour ago" },
    { id: 4, action: "User jane@example.com logged in", time: "2 hours ago" },
    { id: 5, action: "System backup completed", time: "3 hours ago" },
  ];

  return (
    <div className="max-w-7xl mx-auto">
      {/* Page Header */}
      <div className="mb-5">
        <h1 className="text-xl font-bold text-surface-on mb-1">Dashboard</h1>
        <p className="text-xs text-surface-on-variant">
          Welcome back! Here's an overview of your system.
        </p>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-5">
        {isLoading
          ? // Loading skeleton
            [1, 2, 3].map((i) => (
              <div
                key={i}
                className="bg-surface-container rounded-xl p-4 border border-outline-variant/20"
              >
                <div className="animate-pulse">
                  <div className="h-4 bg-surface-variant rounded w-1/3 mb-4"></div>
                  <div className="h-12 bg-surface-variant rounded w-2/3"></div>
                </div>
              </div>
            ))
          : stats.map((stat) => (
              <div
                key={stat.id}
                className={`relative overflow-hidden rounded-xl p-4 bg-gradient-to-br ${stat.gradient}
                                      border border-outline-variant/30 transition-all duration-300 hover:border-outline-variant/60
                                      hover:-translate-y-1`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <p className="text-xs font-medium text-surface-on-variant mb-2">
                      {stat.title}
                    </p>
                    <h3 className="text-3xl font-bold text-surface-on tracking-tight">
                      {stat.value}
                    </h3>
                  </div>
                  <div
                    className={`${stat.iconBg} ${stat.iconColor} p-2.5 rounded-lg`}
                  >
                    {stat.icon}
                  </div>
                </div>
                {/* Decorative element */}
                <div className="absolute -bottom-2 -right-2 w-16 h-16 bg-white/20 rounded-full blur-2xl"></div>
              </div>
            ))}
      </div>

      {/* Two Column Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Recent Activity */}
        <div className="lg:col-span-2 bg-surface-container rounded-xl p-4 border border-outline-variant/20">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-sm font-bold text-surface-on">
              Recent Activity
            </h2>
            <span className="text-xs font-medium text-surface-on-variant bg-surface-variant/30 px-2.5 py-1 rounded-full">
              Last 24 hours
            </span>
          </div>
          <div className="space-y-1">
            {recentActivities.map((activity) => (
              <div
                key={activity.id}
                className="flex items-start gap-3 p-2 rounded-lg hover:bg-surface-variant/20
                                         transition-colors duration-200 group"
              >
                <div className="w-1.5 h-1.5 bg-primary rounded-full mt-1.5 flex-shrink-0 group-hover:scale-125 transition-transform"></div>
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-surface-on font-medium">
                    {activity.action}
                  </p>
                  <p className="text-[10px] text-surface-on-variant mt-0.5">
                    {activity.time}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-surface-container rounded-xl p-4 border border-outline-variant/20">
          <h2 className="text-sm font-bold text-surface-on mb-3">
            Quick Actions
          </h2>
          <div className="space-y-2">
            <Link to="/admin/users">
              <button
                className="w-full flex items-center gap-2.5 p-2.5 rounded-lg bg-primary/5
                                             hover:bg-primary/10 text-primary font-medium transition-all duration-200
                                             hover:scale-[1.02] group"
              >
                <div className="p-1.5 bg-primary/10 rounded-md group-hover:bg-primary/20 transition-colors">
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                    />
                  </svg>
                </div>
                <span className="text-xs">Manage Users</span>
              </button>
            </Link>

            <Link to="/admin/roles">
              <button
                className="w-full flex items-center gap-2.5 p-2.5 rounded-lg bg-surface-variant/20
                                             hover:bg-surface-variant/40 text-surface-on font-medium transition-all duration-200
                                             hover:scale-[1.02] group"
              >
                <div className="p-1.5 bg-surface-variant/30 rounded-md group-hover:bg-surface-variant/50 transition-colors">
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                    />
                  </svg>
                </div>
                <span className="text-xs">Manage Roles</span>
              </button>
            </Link>

            <Link to="/admin/permissions">
              <button
                className="w-full flex items-center gap-2.5 p-2.5 rounded-lg bg-surface-variant/20
                                             hover:bg-surface-variant/40 text-surface-on font-medium transition-all duration-200
                                             hover:scale-[1.02] group"
              >
                <div className="p-1.5 bg-surface-variant/30 rounded-md group-hover:bg-surface-variant/50 transition-colors">
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                    />
                  </svg>
                </div>
                <span className="text-xs">Manage Permissions</span>
              </button>
            </Link>

            {/* Clear Cache Button */}
            <div className="pt-2 mt-2 border-t border-outline-variant/30">
              <button
                onClick={handleClearCache}
                disabled={clearCacheMutation.isPending}
                className="w-full flex items-center gap-2.5 p-2.5 rounded-lg bg-error/5 hover:bg-error/10
                                         text-error font-medium transition-all duration-200 hover:scale-[1.02] group
                                         disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <div className="p-1.5 bg-error/10 rounded-md group-hover:bg-error/20 transition-colors">
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                  </svg>
                </div>
                <span className="text-xs">
                  {clearCacheMutation.isPending ? "Clearing..." : "Clear Cache"}
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
