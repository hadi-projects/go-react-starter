import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import Card from '../components/Card';
import TextField from '../components/TextField';
import Button from '../components/Button';
import { useSettings } from '../context/SettingsContext';
import apiClient from '../api/client';

const Login = () => {
    const navigate = useNavigate();
    const { logo } = useSettings();
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        remember_me: false,
    });
    const [errors, setErrors] = useState({});

    const loginMutation = useMutation({
        mutationFn: async (credentials) => {
            const response = await apiClient.post('/auth/login', credentials);
            return response.data;
        },
        onSuccess: (data) => {
            // Check if 2FA is required
            if (data.data.requires_2fa) {
                // Navigate to 2FA challenge with temp token
                navigate('/2fa-challenge', { state: { tempToken: data.data.temp_token } });
                return;
            }
            // Normal login: Save tokens and redirect
            localStorage.setItem('token', data.data.access_token);
            if (data.data.refresh_token) {
                localStorage.setItem('refresh_token', data.data.refresh_token);
            }
            localStorage.setItem('user', JSON.stringify(data.data.user));
            navigate('/dashboard');
        },
        onError: (error) => {
            setErrors({
                submit: error.response?.data?.meta?.message || 'Invalid email or password',
            });
        },
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({ 
            ...prev, 
            [name]: type === 'checkbox' ? checked : value 
        }));
        // Clear field error on change
        if (errors[name]) {
            setErrors(prev => ({ ...prev, [name]: '' }));
        }
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        setErrors({});

        // Basic validation
        const newErrors = {};
        if (!formData.email) {
            newErrors.email = 'Email is required';
        }
        if (!formData.password) {
            newErrors.password = 'Password is required';
        }

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        loginMutation.mutate(formData);
    };

    return (
        <div className="min-h-screen bg-surface-variant flex items-center justify-center p-6">
            <div className="w-full max-w-md">
                {/* Logo/Header */}
                <div className="text-center mb-8 flex flex-col items-center">
                    {logo && (
                        <div className="w-16 h-16 rounded-2xl border border-outline-variant/30 overflow-hidden bg-surface-container-high shadow-lg p-2.5 mb-6">
                            <img 
                                src={`${import.meta.env.VITE_API_URL}/public/storage/${logo}`} 
                                alt="Logo"
                                className="w-full h-full object-contain"
                                onError={(e) => { e.target.style.display = 'none'; }}
                            />
                        </div>
                    )}
                    <h1 className="text-4xl font-bold text-primary-500 mb-2">Welcome Back</h1>
                    <p className="text-gray-600">Sign in to your account</p>
                </div>

                {/* Login Card */}
                <Card className="p-8">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        {errors.submit && (
                            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md3">
                                {errors.submit}
                            </div>
                        )}

                        <TextField
                            label="Email"
                            type="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                            placeholder="your@email.com"
                            error={errors.email}
                            required
                        />

                        <TextField
                            label="Password"
                            type="password"
                            name="password"
                            value={formData.password}
                            onChange={handleChange}
                            placeholder="Enter your password"
                            error={errors.password}
                            required
                        />

                        <div className="flex items-center justify-between text-sm">
                            <label className="flex items-center">
                                <input 
                                    type="checkbox" 
                                    name="remember_me"
                                    checked={formData.remember_me}
                                    onChange={handleChange}
                                    className="mr-2 rounded" 
                                />
                                <span className="text-gray-600">Remember me</span>
                            </label>
                            <Link to="/forgot-password" className="text-primary-500 hover:text-primary-600">
                                Forgot password?
                            </Link>
                        </div>

                        <Button
                            type="submit"
                            fullWidth
                            disabled={loginMutation.isPending}
                        >
                            {loginMutation.isPending ? 'Signing in...' : 'Sign In'}
                        </Button>
                    </form>

                    <div className="mt-6 text-center">
                        <p className="text-gray-600">
                            Don't have an account?{' '}
                            <Link to="/register" className="text-primary-500 hover:text-primary-600 font-medium">
                                Sign up
                            </Link>
                        </p>
                    </div>
                </Card>

                {/* Back to Home */}
                <div className="mt-6 text-center">
                    <Link to="/" className="text-gray-600 hover:text-primary-500">
                        ← Back to Home
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default Login;
