import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import Table from '../../components/Table';
import Pagination from '../../components/Pagination';
import Card from '../../components/Card';
import Button from '../../components/Button';
import UserFormModal from '../../components/UserFormModal';
import ConfirmDialog from '../../components/ConfirmDialog';
import { getUsers, createUser, updateUser, deleteUser, getRoles } from '../../api/admin';

const Users = () => {
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [searchTerm, setSearchTerm] = useState('');
    const [debouncedSearch, setDebouncedSearch] = useState('');
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [selectedUser, setSelectedUser] = useState(null);

    // Debounce search term
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(searchTerm);
            setCurrentPage(1); // Reset to first page on search
        }, 300);
        return () => clearTimeout(timer);
    }, [searchTerm]);

    const queryClient = useQueryClient();

    const { data, isLoading, error } = useQuery({
        queryKey: ['users', currentPage, itemsPerPage, debouncedSearch],
        queryFn: () => getUsers(currentPage, itemsPerPage, debouncedSearch),
    });

    // Fetch roles for mapping role_id to role name
    const { data: rolesData } = useQuery({
        queryKey: ['roles'],
        queryFn: () => getRoles(1, 100), // Get all roles
    });

    const createMutation = useMutation({
        mutationFn: createUser,
        onSuccess: () => {
            queryClient.invalidateQueries(['users']);
            setIsCreateModalOpen(false);
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, data }) => updateUser(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries(['users']);
            setIsEditModalOpen(false);
            setSelectedUser(null);
        },
    });

    const deleteMutation = useMutation({
        mutationFn: deleteUser,
        onSuccess: () => {
            queryClient.invalidateQueries(['users']);
            setIsDeleteDialogOpen(false);
            setSelectedUser(null);
        },
    });

    const handleCreate = (formData) => {
        createMutation.mutate(formData);
    };

    const handleEdit = (formData) => {
        updateMutation.mutate({ id: selectedUser.id, data: formData });
    };

    const handleDelete = () => {
        deleteMutation.mutate(selectedUser.id);
    };

    const openEditModal = (user) => {
        setSelectedUser(user);
        setIsEditModalOpen(true);
    };

    const openDeleteDialog = (user) => {
        setSelectedUser(user);
        setIsDeleteDialogOpen(true);
    };

    // Create role lookup map
    const rolesMap = {};
    if (rolesData?.data) {
        rolesData.data.forEach(role => {
            rolesMap[role.id] = role.name;
        });
    }

    const columns = [
        { header: 'ID', accessor: 'id' },
        { header: 'Email', accessor: 'email' },
        {
            header: 'Role',
            render: (row) => rolesMap[row.role_id] || `Role ${row.role_id}`
        },
        {
            header: 'Created At',
            render: (row) => new Date(row.created_at).toLocaleDateString(),
        },
        {
            header: 'Actions',
            render: (row) => (
                <div className="flex gap-2">
                    <button
                        onClick={() => openEditModal(row)}
                        className="px-3 py-1 text-sm text-primary border border-primary/50 rounded-md3-sm hover:bg-primary-container/20 transition-colors"
                    >
                        Edit
                    </button>
                    <button
                        onClick={() => openDeleteDialog(row)}
                        className="px-3 py-1 text-sm text-error border border-error/50 rounded-md3-sm hover:bg-error-container/20 transition-colors"
                    >
                        Delete
                    </button>
                </div>
            ),
        },
    ];

    if (error) {
        return (
            <div className="text-center py-12">
                <p className="text-red-500">Error loading users: {error.message}</p>
            </div>
        );
    }

    const users = data?.data || [];
    const meta = data?.meta?.pagination || { total_data: 0, total_pages: 1 };

    return (
        <div>
            <div className="mb-6 flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold text-surface-on">Users Management</h1>
                    <p className="text-surface-on-variant mt-2">Manage user accounts and roles</p>
                </div>
                <Button onClick={() => setIsCreateModalOpen(true)}>
                    Add New User
                </Button>
            </div>

            {/* Search Input */}
            <div className="mb-4">
                <div className="relative">
                    <input
                        type="text"
                        placeholder="Search by email..."
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
                <Table columns={columns} data={users} loading={isLoading} />
                {!isLoading && users.length > 0 && (
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

            {/* Create Modal */}
            <UserFormModal
                isOpen={isCreateModalOpen}
                onClose={() => setIsCreateModalOpen(false)}
                onSubmit={handleCreate}
                loading={createMutation.isPending}
            />

            {/* Edit Modal */}
            <UserFormModal
                isOpen={isEditModalOpen}
                onClose={() => {
                    setIsEditModalOpen(false);
                    setSelectedUser(null);
                }}
                onSubmit={handleEdit}
                user={selectedUser}
                loading={updateMutation.isPending}
            />

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={isDeleteDialogOpen}
                onClose={() => {
                    setIsDeleteDialogOpen(false);
                    setSelectedUser(null);
                }}
                onConfirm={handleDelete}
                title="Delete User"
                message={`Are you sure you want to delete ${selectedUser?.name}? This action cannot be undone.`}
                loading={deleteMutation.isPending}
            />
        </div>
    );
};

export default Users;
