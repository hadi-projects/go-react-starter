import { Link } from 'react-router-dom';
import Button from '../components/Button';

const Landing = () => {
    return (
        <div className="min-h-screen bg-gradient-to-br from-primary-500 via-primary-400 to-secondary-200">
            {/* Navigation */}
            <nav className="container mx-auto px-6 py-6">
                <div className="flex justify-between items-center">
                    <h1 className="text-2xl font-bold text-white">Go Starter</h1>
                    <div className="space-x-4">
                        <Link to="/login">
                            <Button variant="outline" className="text-white border-white hover:bg-white/10">
                                Login
                            </Button>
                        </Link>
                        <Link to="/register">
                            <Button variant="secondary">
                                Get Started
                            </Button>
                        </Link>
                    </div>
                </div>
            </nav>

            {/* Hero Section */}
            <div className="container mx-auto px-6 py-20">
                <div className="max-w-4xl mx-auto text-center">
                    <h2 className="text-6xl font-bold text-white mb-6 leading-tight">
                        Build Amazing Apps <br />
                        <span className="text-secondary-100">Faster Than Ever</span>
                    </h2>
                    <p className="text-xl text-white/90 mb-12 max-w-2xl mx-auto">
                        A modern full-stack starter with Go backend and React frontend.
                        Authentication, authorization, and caching built-in.
                    </p>
                    <div className="flex justify-center gap-6">
                        <Link to="/register">
                            <Button className="px-8 py-4 text-lg shadow-md3-3 hover:shadow-md3-4">
                                Get Started Free
                            </Button>
                        </Link>
                        <Button
                            variant="outline"
                            className="px-8 py-4 text-lg text-white border-white hover:bg-white/10"
                        >
                            View Documentation
                        </Button>
                    </div>
                </div>

                {/* Feature Cards */}
                <div className="grid md:grid-cols-3 gap-8 mt-24 max-w-6xl mx-auto">
                    {features.map((feature, index) => (
                        <div
                            key={index}
                            className="bg-white/10 backdrop-blur-md rounded-md3-lg p-8 text-white hover:bg-white/20 transition-all duration-300 transform hover:-translate-y-2"
                        >
                            <div className="text-4xl mb-4">{feature.icon}</div>
                            <h3 className="text-2xl font-semibold mb-3">{feature.title}</h3>
                            <p className="text-white/80">{feature.description}</p>
                        </div>
                    ))}
                </div>
            </div>

            {/* Footer Section */}
            <div className="container mx-auto px-6 py-12 text-center">
                <p className="text-white/70">
                    © 2026 Go Starter. Built with ❤️ using Go and React.
                </p>
            </div>
        </div>
    );
};

const features = [
    {
        icon: '🚀',
        title: 'Fast Performance',
        description: 'Built with Go and Vite for blazing fast performance and development experience.',
    },
    {
        icon: '🔒',
        title: 'Secure by Default',
        description: 'JWT authentication, RBAC, rate limiting, and security best practices included.',
    },
    {
        icon: '⚡',
        title: 'Redis Caching',
        description: 'Built-in Redis caching for optimal performance and reduced database load.',
    },
];

export default Landing;
