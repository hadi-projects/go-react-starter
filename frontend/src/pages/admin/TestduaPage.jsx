import { useState, useEffect } from 'react';
import { toast } from 'react-hot-toast';
import Button from '../../components/Button';
import Card from '../../components/Card';
import Table from '../../components/Table';
import Modal from '../../components/Modal';
import Pagination from '../../components/Pagination';
import TextField from '../../components/TextField';
import usePermission from '../../hooks/usePermission';
import { 
    getAllTestduas, 
    createTestdua, 
    updateTestdua, 
    deleteTestdua 
} from '../../api/testdua';

const TestduaPage = () => {
    const can = usePermission();
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [paginationMeta, setPaginationMeta] = useState({ total_data: 0, total_pages: 1 });
    const [refreshTrigger, setRefreshTrigger] = useState(0);
    const [formData, setFormData] = useState({
        name: '',
    });

    const columns = [
        { header: 'ID', accessor: 'id' },
        { header: 'Name', accessor: 'name' },
        { header: 'Created At', accessor: 'created_at', render: (row) => new Date(row.created_at).toLocaleString() },
    ];

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                const res = await getAllTestduas({ page: currentPage, limit: itemsPerPage });
                setData(res.data?.data || []);
                setPaginationMeta(res.data?.meta || { total_data: 0, total_pages: 1 });
            } catch (err) {
                toast.error('Failed to fetch data');
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [currentPage, itemsPerPage, refreshTrigger]);

    const handleOpenModal = (item = null) => {
        if (item) {
            setEditingId(item.id);
            setFormData({
                name: item.name,
            });
        } else {
            setEditingId(null);
            setFormData({
                name: '',
            });
        }
        setIsModalOpen(true);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (editingId) {
                await updateTestdua(editingId, formData);
                toast.success('Updated successfully');
            } else {
                await createTestdua(formData);
                toast.success('Created successfully');
            }
            setIsModalOpen(false);
            setRefreshTrigger(t => t + 1);
        } catch (err) {
            toast.error(err.response?.data?.meta?.message || 'Operation failed');
        }
    };

    const handleDelete = async (id) => {
        if (window.confirm('Are you sure you want to delete this item?')) {
            try {
                await deleteTestdua(id);
                toast.success('Deleted successfully');
                setRefreshTrigger(t => t + 1);
            } catch (err) {
                toast.error('Failed to delete');
            }
        }
    };

    const tableActions = [
        ...(can('update-testdua') ? [{ label: 'Edit', onClick: handleOpenModal }] : []),
        ...(can('delete-testdua') ? [{ label: 'Delete', onClick: (row) => handleDelete(row.id), className: 'text-error' }] : []),
    ];

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-surface-on tracking-tight">Testdua Management</h1>
                    <p className="text-sm text-surface-on-variant mt-1">Manage your testdua instances.</p>
                </div>
                {can('create-testdua') && (
                    <Button variant="primary" onClick={() => handleOpenModal()}>
                        Add Testdua
                    </Button>
                )}
            </div>

            <Card className="p-0 overflow-hidden">
                <Table 
                    columns={columns} 
                    data={data} 
                    loading={loading}
                    actions={tableActions}
                />
                {!loading && data.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={paginationMeta.total_pages}
                        totalItems={paginationMeta.total_data}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                        onLimitChange={(newLimit) => {
                            setItemsPerPage(newLimit);
                            setCurrentPage(1);
                        }}
                    />
                )}
            </Card>

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={editingId ? 'Edit Testdua' : 'Add Testdua'}
            >
                <form onSubmit={handleSubmit} className="space-y-4 pt-2">
                    <TextField
                        label="Name"
                        name="name"
                        value={formData.name.toString()}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        
                        
                        required
                    />
                    <div className="flex justify-end gap-3 pt-4">
                        <Button type="button" variant="tonal" onClick={() => setIsModalOpen(false)}>
                            Cancel
                        </Button>
                        <Button type="submit" variant="primary">
                            {editingId ? 'Update' : 'Create'}
                        </Button>
                    </div>
                </form>
            </Modal>
        </div>
    );
};

export default TestduaPage;
