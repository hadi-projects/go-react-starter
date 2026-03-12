import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import Table from '../../components/Table';
import Pagination from '../../components/Pagination';
import Card from '../../components/Card';
import Button from '../../components/Button';
import PermissionFormModal from '../../components/PermissionFormModal';
import ConfirmDialog from '../../components/ConfirmDialog';
import { getPermissions, createPermission, updatePermission, deletePermission, exportPermissions } from '../../api/admin';

const Permissions = () => {
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [searchTerm, setSearchTerm] = useState('');
    const [debouncedSearch, setDebouncedSearch] = useState('');
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [selectedPermission, setSelectedPermission] = useState(null);
    const [isExporting, setIsExporting] = useState(false);

    const queryClient = useQueryClient();

    // Debounce search term
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(searchTerm);
            setCurrentPage(1);
        }, 300);
        return () => clearTimeout(timer);
    }, [searchTerm]);

    const { data, isLoading, error } = useQuery({
        queryKey: ['permissions', currentPage, itemsPerPage],
        queryFn: () => getPermissions(currentPage, itemsPerPage),
    });

    const createMutation = useMutation({
        mutationFn: createPermission,
        onSuccess: () => {
            queryClient.invalidateQueries(['permissions']);
            setIsCreateModalOpen(false);
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, data }) => updatePermission(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries(['permissions']);
            setIsEditModalOpen(false);
            setSelectedPermission(null);
        },
    });

    const deleteMutation = useMutation({
        mutationFn: deletePermission,
        onSuccess: () => {
            queryClient.invalidateQueries(['permissions']);
            setIsDeleteDialogOpen(false);
            setSelectedPermission(null);
        },
    });

    const handleCreate = (data) => {
        createMutation.mutate(data);
    };

    const handleEdit = (data) => {
        updateMutation.mutate({ id: selectedPermission.id, data });
    };

    const handleDelete = () => {
        deleteMutation.mutate(selectedPermission.id);
    };

    const handleExport = async (format) => {
        setIsExporting(true);
        try {
            const response = await exportPermissions(format);
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            const filename = format === 'csv' ? 'permissions.csv' : 'permissions.xlsx';
            link.setAttribute('download', filename);
            document.body.appendChild(link);
            link.click();
            link.remove();
        } catch (err) {
            console.error('Export failed:', err);
        } finally {
            setIsExporting(false);
        }
    };

    const openEditModal = (permission) => {
        setSelectedPermission(permission);
        setIsEditModalOpen(true);
    };

    const openDeleteDialog = (permission) => {
        setSelectedPermission(permission);
        setIsDeleteDialogOpen(true);
    };

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
                <p className="text-red-500">Error loading permissions: {error.message}</p>
            </div>
        );
    }

    const permissions = data?.data || [];
    const meta = data?.meta?.pagination || { total_data: 0, total_pages: 1 };

    // Filter permissions based on search term
    const filteredPermissions = permissions.filter(permission =>
        permission.name.toLowerCase().includes(debouncedSearch.toLowerCase()) ||
        permission.description.toLowerCase().includes(debouncedSearch.toLowerCase())
    );

    return (
        <div>
            <div className="mb-6 flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold text-surface-on">Permissions Management</h1>
                    <p className="text-surface-on-variant mt-2">Manage system permissions and access control</p>
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
                    <Button onClick={() => setIsCreateModalOpen(true)}>
                        Add New Permission
                    </Button>
                </div>
            </div>

            {/* Search Input */}
            <div className="mb-4">
                <div className="relative">
                    <input
                        type="text"
                        placeholder="Search permissions by name or description..."
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

            <Card className="p-0 overflow-hidden">
                <Table columns={columns} data={filteredPermissions} loading={isLoading} actions={[
                    { label: 'Edit', onClick: openEditModal },
                    { label: 'Delete', onClick: openDeleteDialog, className: 'text-error' },
                ]} />
                {!isLoading && permissions.length > 0 && (
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

            <PermissionFormModal
                isOpen={isCreateModalOpen}
                onClose={() => setIsCreateModalOpen(false)}
                onSubmit={handleCreate}
                isLoading={createMutation.isPending}
            />

            <PermissionFormModal
                isOpen={isEditModalOpen}
                onClose={() => {
                    setIsEditModalOpen(false);
                    setSelectedPermission(null);
                }}
                onSubmit={handleEdit}
                permission={selectedPermission}
                isLoading={updateMutation.isPending}
            />

            <ConfirmDialog
                isOpen={isDeleteDialogOpen}
                onClose={() => {
                    setIsDeleteDialogOpen(false);
                    setSelectedPermission(null);
                }}
                onConfirm={handleDelete}
                title="Delete Permission"
                message={`Are you sure you want to delete permission "${selectedPermission?.name}"? This action cannot be undone.`}
                isLoading={deleteMutation.isPending}
            />
        </div>
    );
};

export default Permissions;
