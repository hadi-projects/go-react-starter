import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { forgotPassword } from '../api/auth';
import Card from '../components/Card';
import TextField from '../components/TextField';
import Button from '../components/Button';
import toast from 'react-hot-toast';

const ForgotPassword = () => {
    const [email, setEmail] = useState('');

    const forgotPasswordMutation = useMutation({
        mutationFn: forgotPassword,
        onSuccess: (data) => {
            toast.success(data.message || 'Reset link sent! Check your email.');
            setEmail('');
        },
        onError: (error) => {
            toast.error(error.response?.data?.meta?.message || 'Failed to send reset link.');
        },
    });

    const handleSubmit = (e) => {
        e.preventDefault();
        if (!email) {
            toast.error('Email is required');
            return;
        }
        forgotPasswordMutation.mutate(email);
    };

    return (
        <div className="min-h-screen bg-surface-variant flex items-center justify-center p-6">
            <div className="w-full max-w-md">
                <div className="text-center mb-8">
                    <h1 className="text-4xl font-bold text-primary-500 mb-2">Forgot Password</h1>
                    <p className="text-gray-600">Enter your email to receive a reset link</p>
                </div>

                <Card className="p-8">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <TextField
                            label="Email"
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="your@email.com"
                            required
                        />

                        <Button
                            type="submit"
                            fullWidth
                            disabled={forgotPasswordMutation.isPending}
                        >
                            {forgotPasswordMutation.isPending ? 'Sending...' : 'Send Reset Link'}
                        </Button>
                    </form>

                    <div className="mt-6 text-center">
                        <Link to="/login" className="text-primary-500 hover:text-primary-600 font-medium">
                            ← Back to Login
                        </Link>
                    </div>
                </Card>
            </div>
        </div>
    );
};

export default ForgotPassword;
