import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import Card from '../components/Card';
import TextField from '../components/TextField';
import Button from '../components/Button';
import apiClient from '../api/client';

const Login = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        email: '',
        password: '',
    });
    const [errors, setErrors] = useState({});

    const loginMutation = useMutation({
        mutationFn: async (credentials) => {
            const response = await apiClient.post('/auth/login', credentials);
            return response.data;
        },
        onSuccess: (data) => {
            // Save token to localStorage
            localStorage.setItem('token', data.data.token);
            // Redirect to dashboard or home
            navigate('/dashboard');
        },
        onError: (error) => {
            setErrors({
                submit: error.response?.data?.meta?.message || 'Login failed. Please try again.',
            });
        },
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
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
        <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-100 flex items-center justify-center p-6">
            <div className="w-full max-w-md">
                {/* Logo/Header */}
                <div className="text-center mb-8">
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
                                <input type="checkbox" className="mr-2 rounded" />
                                <span className="text-gray-600">Remember me</span>
                            </label>
                            <a href="#" className="text-primary-500 hover:text-primary-600">
                                Forgot password?
                            </a>
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
