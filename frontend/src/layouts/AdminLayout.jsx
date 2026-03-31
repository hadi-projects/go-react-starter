import { useEffect, useState, useMemo, useCallback } from "react";
import { useNavigate, Outlet, useLocation, Link } from "react-router-dom";
import Sidebar from "../components/Sidebar";
import Button from "../components/Button";
import { getHealthStatus, getMe } from "../api/admin";
import { logoutApi } from "../api/auth";
import { useTheme } from "../context/ThemeContext";
import { PERMS } from "../utils/permissions";

const AdminLayout = () => {
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const [user, setUser] = useState(null);
  const [cacheStatus, setCacheStatus] = useState("unknown"); // unknown, connected, disconnected
  const [kafkaStatus, setKafkaStatus] = useState("unknown"); // unknown, connected, disconnected
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  // Function to refresh user profile and permissions from server
  const refreshUserData = useCallback(async () => {
    try {
      const response = await getMe();
      if (response.success && response.data) {
        const updatedUser = response.data;
        setUser(updatedUser);
        localStorage.setItem("user", JSON.stringify(updatedUser));
      }
    } catch (error) {
      console.error("Failed to refresh user data:", error);
    }
  }, []);

  // Periodic refresh of user metadata (every 60 seconds)
  useEffect(() => {
    const interval = setInterval(refreshUserData, 60000);
    return () => clearInterval(interval);
  }, [refreshUserData]);

  // Auto-collapse sidebar on small screens
  useEffect(() => {
    const handleResize = () => {
      setSidebarCollapsed(window.innerWidth < 1024);
    };
    handleResize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    const token = localStorage.getItem("token");
    const userData = localStorage.getItem("user");

    if (!token) {
      navigate("/login");
      return;
    }

    if (userData) {
      setUser(JSON.parse(userData));
    }

    const fetchHealthStatus = async () => {
      if (!userData) return;
      const parsedUser = JSON.parse(userData);
      const mask = parsedUser.permissions_mask
        ? BigInt(parsedUser.permissions_mask)
        : 0n;
      if ((mask & PERMS.MANAGE_CACHE) === 0n) return;
      try {
        const response = await getHealthStatus();
        setCacheStatus(response.data.redis);
        setKafkaStatus(response.data.kafka);
      } catch (error) {
        console.error("Failed to fetch health status:", error);
        setCacheStatus("disconnected");
        setKafkaStatus("disconnected");
      }
    };

    fetchHealthStatus();

    // Poll every 30 seconds (only if user has permission)
    if (userData) {
      const parsedUser = JSON.parse(userData);
      const mask = parsedUser.permissions_mask
        ? BigInt(parsedUser.permissions_mask)
        : 0n;
      if ((mask & PERMS.MANAGE_CACHE) !== 0n) {
        const interval = setInterval(fetchHealthStatus, 30000);
        return () => clearInterval(interval);
      }
    }
  }, [navigate]);

  const handleLogout = async () => {
    await logoutApi("manual");
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    navigate("/login");
  };

  // Navigation Configuration
  const navigationSections = [
    {
      label: "Main",
      items: [
        {
          path: "/dashboard",
          label: "Dashboard",
          icon: (
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
              />
            </svg>
          ),
        },
      ],
    },
    {
      label: "Management",
      items: [
        {
          label: "Administrator",
          icon: (
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          ),
          permission: [PERMS.GET_USER, PERMS.GET_ROLE, PERMS.GET_PERMISSION],
          subItems: [
            {
              path: "/admin/users",
              label: "Users",
              permission: PERMS.GET_USER,
              icon: (
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
              ),
            },
            {
              path: "/admin/roles",
              label: "Roles",
              permission: PERMS.GET_ROLE,
              icon: (
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
              ),
            },
            {
              path: "/admin/permissions",
              label: "Permissions",
              permission: PERMS.GET_PERMISSION,
              icon: (
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
                    d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                  />
                </svg>
              ),
            },
          ],
        },
        {
          label: "Logs",
          icon: (
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          ),
          permission: [
            PERMS.GET_AUDIT_LOG,
            PERMS.GET_ALL_LOGS,
            PERMS.GET_AUTH_LOG,
            PERMS.GET_HTTP_LOG,
          ],
          subItems: [
            {
              path: "/admin/logs/audit",
              label: "Audit",
              permission: PERMS.GET_AUDIT_LOG,
              icon: (
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
                    d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                  />
                </svg>
              ),
            },
            {
              path: "/admin/logs/system",
              label: "System",
              permission: PERMS.GET_ALL_LOGS,
              icon: (
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
                    d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
              ),
            },
            {
              path: "/admin/logs/http",
              label: "HTTP",
              permission: PERMS.GET_HTTP_LOG,
              icon: (
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
                    d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                  />
                </svg>
              ),
            },
          ],
        },
        {
          path: "/admin/generator",
          label: "Module Generator",
          permission: PERMS.CREATE_MODULE,
          icon: (
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
              />
            </svg>
          ),
        },
        {
          label: "Produk",
          path: "/admin/produk",
          permission: PERMS.GET_PRODUK,
          icon: (
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
              />
            </svg>
          ),
        },

        // [GENERATOR_INSERT_ADMIN_ITEM]
      ],
    },
  ];

  // Filter navigation based on user permissions
  const filteredNavigation = useMemo(() => {
    if (!user || user.role_id === 1) return navigationSections; // Admin full access (fallback check on role_id)

    const checkPermission = (item) => {
      if (item.permission === undefined) return true;
      if (!user.permissions_mask) return false;
      const mask = BigInt(user.permissions_mask);

      if (typeof item.permission === "bigint")
        return (mask & item.permission) !== 0n;
      if (Array.isArray(item.permission))
        return item.permission.some((p) => (mask & p) !== 0n);
      return false;
    };

    const filterItems = (items) => {
      return items
        .map((item) => {
          if (item.subItems) {
            const filteredSubItems = filterItems(item.subItems);
            return filteredSubItems.length > 0
              ? { ...item, subItems: filteredSubItems }
              : null;
          }
          return checkPermission(item) ? item : null;
        })
        .filter(Boolean);
    };

    return navigationSections
      .map((section) => ({
        ...section,
        items: filterItems(section.items),
      }))
      .filter((section) => section.items.length > 0);
  }, [user]);

  if (!user) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center">
        <div className="text-primary-500 animate-pulse">
          Loading MD3 Expressive...
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen bg-surface overflow-hidden">
      {/* Sidebar Navigation */}
      <Sidebar
        sections={filteredNavigation}
        onLogout={handleLogout}
        collapsed={sidebarCollapsed}
        onToggleCollapse={() => setSidebarCollapsed((v) => !v)}
      />

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0 bg-surface-container-low relative">
        {/* Compact Header */}
        <header className="h-12 flex items-center justify-between px-5 bg-surface-container-low border-b border-outline-variant/30">
          <h2 className="text-sm font-medium text-surface-on">
            Administration
          </h2>
          <div className="flex items-center gap-3">
            {/* Theme Toggle */}
            <button
              onClick={toggleTheme}
              className="p-1.5 rounded-full hover:bg-surface-variant/30 text-surface-on-variant transition-all duration-200"
              title={`Switch to ${theme === "light" ? "dark" : "light"} mode`}
            >
              {theme === "light" ? (
                <svg
                  className="w-5 h-5 transition-transform"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
                  />
                </svg>
              ) : (
                <svg
                  className="w-5 h-5 transition-transform"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 3v1m0 16v1m9-9h1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
                  />
                </svg>
              )}
            </button>

            {/* System Status Indicators - Only visible with manage-cache permission */}
            {(BigInt(user?.permissions_mask || 0) & PERMS.MANAGE_CACHE) !==
              0n && (
              <div className="flex items-center gap-2">
                <div
                  className="flex items-center gap-1.5 px-2 py-1 rounded-full bg-surface-variant/30 text-surface-on-variant text-xs font-medium"
                  title={`Redis: ${cacheStatus}`}
                >
                  <div
                    className={`w-2 h-2 rounded-full ${cacheStatus === "connected" ? "bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.4)]" : "bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.4)]"}`}
                  ></div>
                  <span className="hidden sm:inline">Redis</span>
                </div>
                <div
                  className="flex items-center gap-1.5 px-2 py-1 rounded-full bg-surface-variant/30 text-surface-on-variant text-xs font-medium"
                  title={`Kafka: ${kafkaStatus}`}
                >
                  <div
                    className={`w-2 h-2 rounded-full ${kafkaStatus === "connected" ? "bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.4)]" : "bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.4)]"}`}
                  ></div>
                  <span className="hidden sm:inline">Kafka</span>
                </div>
              </div>
            )}

            <div className="flex flex-col items-end mr-1">
              <span className="text-xs font-medium text-surface-on">
                {user.email}
              </span>
              <span className="text-[10px] text-surface-on-variant font-medium">
                System Administrator
              </span>
            </div>
            <Link
              to="/profile"
              className="w-8 h-8 rounded-full bg-primary-container flex items-center justify-center text-primary-on-container border border-primary/20 cursor-pointer hover:bg-primary-container/80 transition-all text-xs"
              title="My Profile"
            >
              {user.email.charAt(0).toUpperCase()}
            </Link>
          </div>
        </header>

        {/* Scrollable Page Content */}
        <main className="flex-1 overflow-y-auto p-5 custom-scrollbar">
          <div className="max-w-7xl mx-auto animate-fade-in-up">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};

export default AdminLayout;
