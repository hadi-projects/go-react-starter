import { Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import Landing from './pages/Landing'
import Login from './pages/Login'
import Register from './pages/Register'
import ForgotPassword from './pages/ForgotPassword'
import ResetPassword from './pages/ResetPassword'
import Dashboard from './pages/Dashboard'
import AdminLayout from './layouts/AdminLayout'
import Users from './pages/admin/Users'
import Roles from './pages/admin/Roles'
import Permissions from './pages/admin/Permissions'
import Logs from './pages/admin/Logs'
import HttpLogs from './pages/admin/HttpLogs'
import GeneratorPage from './pages/admin/GeneratorPage'
import { ThemeProvider } from './context/ThemeContext'
import ProdukPage from './pages/admin/ProdukPage';
import TwoFAChallengePage from './pages/TwoFAChallengePage';
import ProfilePage from './pages/admin/ProfilePage';

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

        {/* Admin Routes with Sidebar */}
        <Route path="/" element={<AdminLayout />}>
          <Route path="dashboard" element={<Dashboard />} />

          <Route path="admin/users" element={<Users />} />
          <Route path="admin/roles" element={<Roles />} />
          <Route path="admin/permissions" element={<Permissions />} />
          <Route path="admin/logs" element={<Navigate to="/admin/logs/all" replace />} />
          <Route path="admin/logs/http" element={<HttpLogs />} />
          <Route path="admin/logs/:type" element={<Logs />} />
          <Route path="admin/generator" element={<GeneratorPage />} />
          <Route path="admin/produk" element={<ProdukPage />} />
          <Route path="profile" element={<ProfilePage />} />
					// [GENERATOR_INSERT_ROUTE]
        </Route>
      </Routes>
    </ThemeProvider>
  )
}

export default App
