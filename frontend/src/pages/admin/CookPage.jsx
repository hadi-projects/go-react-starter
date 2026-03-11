import { useState, useEffect } from 'react';
import { toast } from 'react-hot-toast';
import Button from '../../components/Button';
import Card from '../../components/Card';
import Table from '../../components/Table';
import Modal from '../../components/Modal';
import TextField from '../../components/TextField';
import {
    getAllCooks,
    createCook,
    updateCook,
    deleteCook
} from '../../api/cook';

const CookPage = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);
    const [formData, setFormData] = useState({
        name: '',
    });

    const columns = [
        { header: 'Name', accessor: 'name' },
        { header: 'Created At', accessor: 'created_at', render: (val) => new Date(val).toLocaleString() },
    ];

    const fetchData = async () => {
        setLoading(true);
        try {
            const res = await getAllCooks();
            setData(res.data?.data || []);
        } catch (err) {
            toast.error('Failed to fetch data');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

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
                await updateCook(editingId, formData);
                toast.success('Updated successfully');
            } else {
                await createCook(formData);
                toast.success('Created successfully');
            }
            setIsModalOpen(false);
            fetchData();
        } catch (err) {
            toast.error(err.response?.data?.meta?.message || 'Operation failed');
        }
    };

    const handleDelete = async (id) => {
        if (window.confirm('Are you sure you want to delete this item?')) {
            try {
                await deleteCook(id);
                toast.success('Deleted successfully');
                fetchData();
            } catch (err) {
                toast.error('Failed to delete');
            }
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-surface-on tracking-tight">Cook Management</h1>
                    <p className="text-sm text-surface-on-variant mt-1">Manage your cook instances.</p>
                </div>
                <Button variant="primary" onClick={() => handleOpenModal()}>
                    Add Cook
                </Button>
            </div>

            <Card className="p-0 overflow-hidden">
                <Table
                    columns={columns}
                    data={data}
                    loading={loading}
                    onEdit={handleOpenModal}
                    onDelete={handleDelete}
                />
            </Card>

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={editingId ? 'Edit Cook' : 'Add Cook'}
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

export default CookPage;
