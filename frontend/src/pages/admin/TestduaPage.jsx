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
    deleteTestdua,
    exportTestdua
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
    const [isExporting, setIsExporting] = useState(false);
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

    const handleExport = async (format) => {
        setIsExporting(true);
        try {
            const response = await exportTestdua(format);
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            const filename = format === 'csv' ? 'testdua.csv' : 'testdua.xlsx';
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
                <div className="flex gap-2">
                    <div className="flex bg-surface-variant/20 p-1 rounded-lg shrink-0">
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
                    {can('create-testdua') && (
                        <Button variant="primary" onClick={() => handleOpenModal()}>
                            Add Testdua
                        </Button>
                    )}
                </div>
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
