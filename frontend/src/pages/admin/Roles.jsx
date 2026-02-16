import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import Table from '../../components/Table';
import Pagination from '../../components/Pagination';
import Card from '../../components/Card';
import { getRoles } from '../../api/admin';

const Roles = () => {
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage] = useState(10);
    const [searchTerm, setSearchTerm] = useState('');
    const [debouncedSearch, setDebouncedSearch] = useState('');

    // Debounce search term
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(searchTerm);
            setCurrentPage(1);
        }, 300);
        return () => clearTimeout(timer);
    }, [searchTerm]);

    const { data, isLoading, error } = useQuery({
        queryKey: ['roles', currentPage, itemsPerPage],
        queryFn: () => getRoles(currentPage, itemsPerPage),
    });

    const columns = [
        { header: 'ID', accessor: 'id' },
        { header: 'Name', accessor: 'name' },
        { header: 'Description', accessor: 'description' },
        {
            header: 'Created At',
            render: (row) => new Date(row.created_at).toLocaleDateString(),
        },
    ];

    if (error) {
        return (
            <div className="text-center py-12">
                <p className="text-red-500">Error loading roles: {error.message}</p>
            </div>
        );
    }

    const roles = data?.data || [];
    const meta = data?.meta?.pagination || { total_data: 0, total_pages: 1 };

    // Filter roles based on search term
    const filteredRoles = roles.filter(role =>
        role.name.toLowerCase().includes(debouncedSearch.toLowerCase()) ||
        role.description.toLowerCase().includes(debouncedSearch.toLowerCase())
    );

    return (
        <div>
            <div className="mb-6">
                <h1 className="text-3xl font-bold text-gray-900">Roles Management</h1>
                <p className="text-gray-600 mt-2">Manage user roles and access levels</p>
            </div>

            {/* Search Input */}
            <div className="mb-4">
                <div className="relative">
                    <input
                        type="text"
                        placeholder="Search roles by name or description..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    />
                    <svg
                        className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </div>
            </div>

            <Card className="p-0 overflow-hidden">
                <Table columns={columns} data={filteredRoles} loading={isLoading} />
                {!isLoading && roles.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={meta.total_pages}
                        totalItems={meta.total_data}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                    />
                )}
            </Card>
        </div>
    );
};

export default Roles;
