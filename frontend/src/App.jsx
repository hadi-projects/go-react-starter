import { Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import Landing from "./pages/Landing";
import Login from "./pages/Login";
import Register from "./pages/Register";
import ForgotPassword from "./pages/ForgotPassword";
import ResetPassword from "./pages/ResetPassword";
import Dashboard from "./pages/Dashboard";
import AdminLayout from "./layouts/AdminLayout";
import Users from "./pages/admin/Users";
import Roles from "./pages/admin/Roles";
import Permissions from "./pages/admin/Permissions";
import Logs from "./pages/admin/Logs";
import HttpLogs from "./pages/admin/HttpLogs";
import GeneratorPage from "./pages/admin/GeneratorPage";
import { ThemeProvider } from "./context/ThemeContext";
import ProdukPage from "./pages/admin/ProdukPage";
import TwoFAChallengePage from "./pages/TwoFAChallengePage";
import ProfilePage from "./pages/admin/ProfilePage";
import StoragePage from "./pages/admin/StoragePage";
import SharePage from "./pages/SharePage";
import SettingsPage from "./pages/admin/SettingsPage";
import TwoFAResetRequestPage from "./pages/TwoFAResetRequestPage";
import TwoFAResetConfirmPage from "./pages/TwoFAResetConfirmPage";

import ApiKeyPage from "./pages/admin/ApiKeyPage";

// [GENERATOR_INSERT_IMPORT]

function App() {
  return (
    <ThemeProvider>
      <Toaster position="top-right" />
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route path="/2fa-challenge" element={<TwoFAChallengePage />} />
        <Route path="/twofa/reset-request" element={<TwoFAResetRequestPage />} />
        <Route path="/twofa/reset-confirm" element={<TwoFAResetConfirmPage />} />

        {/* Admin Routes with Sidebar */}
        <Route path="/" element={<AdminLayout />}>
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="profile" element={<ProfilePage />} />
          
          <Route path="admin">
            <Route path="apikeys" element={<ApiKeyPage />} />
            <Route path="users" element={<Users />} />
            <Route path="roles" element={<Roles />} />
            <Route path="permissions" element={<Permissions />} />
            <Route path="logs" element={<Navigate to="/admin/logs/all" replace />} />
            <Route path="logs/http" element={<HttpLogs />} />
            <Route path="logs/:type" element={<Logs />} />
            <Route path="generator" element={<GeneratorPage />} />
            <Route path="produk" element={<ProdukPage />} />
            <Route path="storage" element={<StoragePage />} />
            <Route path="settings" element={<Navigate to="/admin/settings/website" replace />} />
            <Route path="settings/:category" element={<SettingsPage />} />
          </Route>
          
          // [GENERATOR_INSERT_ROUTE]
        </Route>

        {/* Public share page — outside AdminLayout, no auth required */}
        <Route path="/share/:token" element={<SharePage />} />
      </Routes>
    </ThemeProvider>
  );
}

export default App;
