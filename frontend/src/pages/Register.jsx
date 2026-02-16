import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import Card from '../components/Card';
import TextField from '../components/TextField';
import Button from '../components/Button';
import apiClient from '../api/client';

const Register = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        confirmPassword: '',
        roleId: 2, // Default to regular user role
    });
    const [errors, setErrors] = useState({});

    const registerMutation = useMutation({
        mutationFn: async (userData) => {
            // Remove confirmPassword before sending to API
            const { confirmPassword, ...dataToSend } = userData;
            const response = await apiClient.post('/auth/register', {
                ...dataToSend,
                role_id: userData.roleId,
            });
            return response.data;
        },
        onSuccess: () => {
            // Redirect to login after successful registration
            navigate('/login');
        },
        onError: (error) => {
            setErrors({
                submit: error.response?.data?.meta?.message || 'Registration failed. Please try again.',
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

    const validateForm = () => {
        const newErrors = {};

        if (!formData.name || formData.name.trim().length < 2) {
            newErrors.name = 'Name must be at least 2 characters';
        }

        if (!formData.email) {
            newErrors.email = 'Email is required';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
            newErrors.email = 'Please enter a valid email';
        }

        if (!formData.password) {
            newErrors.password = 'Password is required';
        } else if (formData.password.length < 6) {
            newErrors.password = 'Password must be at least 6 characters';
        }

        if (formData.password !== formData.confirmPassword) {
            newErrors.confirmPassword = 'Passwords do not match';
        }

        return newErrors;
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        setErrors({});

        const validationErrors = validateForm();
        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors);
            return;
        }

        registerMutation.mutate(formData);
    };

    const getPasswordStrength = (password) => {
        if (!password) return { strength: 0, label: '', color: '' };

        let strength = 0;
        if (password.length >= 6) strength++;
        if (password.length >= 10) strength++;
        if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++;
        if (/\d/.test(password)) strength++;
        if (/[^a-zA-Z0-9]/.test(password)) strength++;

        const labels = ['', 'Weak', 'Fair', 'Good', 'Strong', 'Very Strong'];
        const colors = ['', 'bg-red-500', 'bg-orange-500', 'bg-yellow-500', 'bg-green-500', 'bg-green-600'];

        return { strength, label: labels[strength], color: colors[strength] };
    };

    const passwordStrength = getPasswordStrength(formData.password);

    return (
        <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-100 flex items-center justify-center p-6">
            <div className="w-full max-w-md">
                {/* Logo/Header */}
                <div className="text-center mb-8">
                    <h1 className="text-4xl font-bold text-primary-500 mb-2">Get Started</h1>
                    <p className="text-gray-600">Create your account</p>
                </div>

                {/* Register Card */}
                <Card className="p-8">
                    <form onSubmit={handleSubmit} className="space-y-5">
                        {errors.submit && (
                            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md3">
                                {errors.submit}
                            </div>
                        )}

                        <TextField
                            label="Full Name"
                            type="text"
                            name="name"
                            value={formData.name}
                            onChange={handleChange}
                            placeholder="John Doe"
                            error={errors.name}
                            required
                        />

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

                        <div>
                            <TextField
                                label="Password"
                                type="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                placeholder="Create a strong password"
                                error={errors.password}
                                required
                            />
                            {formData.password && (
                                <div className="mt-2">
                                    <div className="flex gap-1 mb-1">
                                        {[1, 2, 3, 4, 5].map((level) => (
                                            <div
                                                key={level}
                                                className={`h-1 flex-1 rounded ${level <= passwordStrength.strength
                                                        ? passwordStrength.color
                                                        : 'bg-gray-200'
                                                    }`}
                                            />
                                        ))}
                                    </div>
                                    <p className="text-xs text-gray-600">
                                        Strength: {passwordStrength.label}
                                    </p>
                                </div>
                            )}
                        </div>

                        <TextField
                            label="Confirm Password"
                            type="password"
                            name="confirmPassword"
                            value={formData.confirmPassword}
                            onChange={handleChange}
                            placeholder="Confirm your password"
                            error={errors.confirmPassword}
                            required
                        />

                        <div className="pt-2">
                            <label className="flex items-start">
                                <input type="checkbox" className="mr-2 mt-1 rounded" required />
                                <span className="text-sm text-gray-600">
                                    I agree to the{' '}
                                    <a href="#" className="text-primary-500 hover:text-primary-600">
                                        Terms of Service
                                    </a>{' '}
                                    and{' '}
                                    <a href="#" className="text-primary-500 hover:text-primary-600">
                                        Privacy Policy
                                    </a>
                                </span>
                            </label>
                        </div>

                        <Button
                            type="submit"
                            fullWidth
                            disabled={registerMutation.isPending}
                        >
                            {registerMutation.isPending ? 'Creating Account...' : 'Create Account'}
                        </Button>
                    </form>

                    <div className="mt-6 text-center">
                        <p className="text-gray-600">
                            Already have an account?{' '}
                            <Link to="/login" className="text-primary-500 hover:text-primary-600 font-medium">
                                Sign in
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

export default Register;
