import { Routes, Route } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import Landing from './pages/Landing'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import AdminLayout from './layouts/AdminLayout'
import Users from './pages/admin/Users'
import Roles from './pages/admin/Roles'
import Permissions from './pages/admin/Permissions'

function App() {
  return (
    <>
      <Toaster position="top-right" />
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        {/* Admin Routes with Sidebar */}
        <Route path="/" element={<AdminLayout />}>
          <Route path="dashboard" element={<Dashboard />} />

          <Route path="admin/users" element={<Users />} />
          <Route path="admin/roles" element={<Roles />} />
          <Route path="admin/permissions" element={<Permissions />} />
        </Route>
      </Routes>
    </>
  )
}

export default App
