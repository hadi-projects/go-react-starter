import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import Table from '../../components/Table';
import Pagination from '../../components/Pagination';
import Card from '../../components/Card';
import Button from '../../components/Button';
import RoleFormModal from '../../components/RoleFormModal';
import ConfirmDialog from '../../components/ConfirmDialog';
import { getRoles, createRole, updateRole, deleteRole, exportRoles } from '../../api/admin';

const Roles = () => {
    const queryClient = useQueryClient();
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [searchTerm, setSearchTerm] = useState('');
    const [debouncedSearch, setDebouncedSearch] = useState('');

    // Modal state
    const [isFormModalOpen, setIsFormModalOpen] = useState(false);
    const [isDeleteOpen, setIsDeleteOpen] = useState(false);
    const [selectedRole, setSelectedRole] = useState(null);
    const [isExporting, setIsExporting] = useState(false);

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

    const createMutation = useMutation({
        mutationFn: createRole,
        onSuccess: () => {
            queryClient.invalidateQueries(['roles']);
            setIsFormModalOpen(false);
            toast.success('Role created successfully');
        },
        onError: (err) => {
            toast.error(err.response?.data?.meta?.message || 'Failed to create role');
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, data }) => updateRole(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries(['roles']);
            setIsFormModalOpen(false);
            toast.success('Role updated successfully');
        },
        onError: (err) => {
            toast.error(err.response?.data?.meta?.message || 'Failed to update role');
        },
    });

    const deleteMutation = useMutation({
        mutationFn: deleteRole,
        onSuccess: () => {
            queryClient.invalidateQueries(['roles']);
            setIsDeleteOpen(false);
            toast.success('Role deleted successfully');
        },
        onError: (err) => {
            toast.error(err.response?.data?.meta?.message || 'Failed to delete role');
        },
    });

    const handleCreateRole = () => {
        setSelectedRole(null);
        setIsFormModalOpen(true);
    };

    const handleEditRole = (role) => {
        setSelectedRole(role);
        setIsFormModalOpen(true);
    };

    const handleDeleteRole = (role) => {
        setSelectedRole(role);
        setIsDeleteOpen(true);
    };

    const handleFormSubmit = (formData) => {
        if (selectedRole) {
            updateMutation.mutate({ id: selectedRole.id, data: formData });
        } else {
            createMutation.mutate(formData);
        }
    };

    const handleExport = async (format) => {
        setIsExporting(true);
        try {
            const response = await exportRoles(format);
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            const filename = format === 'csv' ? 'roles.csv' : 'roles.xlsx';
            link.setAttribute('download', filename);
            document.body.appendChild(link);
            link.click();
            link.remove();
        } catch (err) {
            console.error('Export failed:', err);
            toast.error('Export failed');
        } finally {
            setIsExporting(false);
        }
    };

    const columns = [
        { header: 'ID', accessor: 'id' },
        { header: 'Name', accessor: 'name' },
        { header: 'Description', accessor: 'description' },
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
        (role.description?.toLowerCase().includes(debouncedSearch.toLowerCase()) || false)
    );

    return (
        <div>
            <div className="flex flex-col md:flex-row md:items-center justify-between mb-8 gap-4">
                <div>
                    <h1 className="text-2xl font-bold text-surface-on">Roles Management</h1>
                    <p className="text-surface-on-variant mt-1">Manage user roles and their associated permissions</p>
                </div>
                <div className="flex gap-2">
                    <div className="flex bg-surface-variant/20 p-1 rounded-lg">
                        <button
                            onClick={() => handleExport('excel')}
                            className="px-3 py-1.5 text-xs font-semibold hover:bg-surface-variant/30 rounded-md transition-all flex items-center gap-1.5 text-surface-on disabled:opacity-50"
                            disabled={isExporting}
                        >
                            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                            Excel
                        </button>
                        <button
                            onClick={() => handleExport('csv')}
                            className="px-3 py-1.5 text-xs font-semibold hover:bg-surface-variant/30 rounded-md transition-all flex items-center gap-1.5 text-surface-on disabled:opacity-50"
                            disabled={isExporting}
                        >
                            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                            CSV
                        </button>
                    </div>
                    <Button onClick={handleCreateRole}>
                        Create Role
                    </Button>
                </div>
            </div>

            {/* Search Input */}
            <div className="mb-6">
                <div className="relative max-w-md">
                    <input
                        type="text"
                        placeholder="Search roles..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="text-field"
                    />
                    <svg
                        className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-surface-on-variant"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </div>
            </div>

            <Card className="p-0 overflow-hidden border border-outline-variant">
                <Table columns={columns} data={filteredRoles} loading={isLoading} actions={[
                    { label: 'Edit', onClick: handleEditRole },
                    { label: 'Delete', onClick: handleDeleteRole, className: 'text-error' },
                ]} />
                {!isLoading && roles.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={meta.total_pages}
                        totalItems={meta.total_data}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                        onLimitChange={(newLimit) => {
                            setItemsPerPage(newLimit);
                            setCurrentPage(1);
                        }}
                    />
                )}
            </Card>

            {/* Form Modal */}
            <RoleFormModal
                isOpen={isFormModalOpen}
                onClose={() => setIsFormModalOpen(false)}
                onSubmit={handleFormSubmit}
                role={selectedRole}
                loading={createMutation.isPending || updateMutation.isPending}
            />

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={isDeleteOpen}
                onClose={() => setIsDeleteOpen(false)}
                onConfirm={() => deleteMutation.mutate(selectedRole?.id)}
                title="Delete Role"
                message={`Are you sure you want to delete the role "${selectedRole?.name}"? This action cannot be undone.`}
                loading={deleteMutation.isPending}
            />
        </div>
    );
};

export default Roles;
