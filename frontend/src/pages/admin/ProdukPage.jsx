import { useState, useEffect } from 'react';
import { toast } from 'react-hot-toast';
import Button from '../../components/Button';
import Card from '../../components/Card';
import Table from '../../components/Table';
import Modal from '../../components/Modal';
import TextField from '../../components/TextField';
import {
    getAllProduks,
    createProduk,
    updateProduk,
    deleteProduk
} from '../../api/produk';

const ProdukPage = () => {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);
    const [formData, setFormData] = useState({
        name: '',
        harga: 0,
    });

    const columns = [
        { header: 'Name', accessor: 'name' },
        { header: 'Harga', accessor: 'harga' },
        { header: 'Created At', accessor: 'created_at', render: (row) => new Date(row.created_at).toLocaleString() },
    ];

    const fetchData = async () => {
        setLoading(true);
        try {
            const res = await getAllProduks();
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
                harga: item.harga,
            });
        } else {
            setEditingId(null);
            setFormData({
                name: '',
                harga: 0,
            });
        }
        setIsModalOpen(true);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (editingId) {
                await updateProduk(editingId, formData);
                toast.success('Updated successfully');
            } else {
                await createProduk(formData);
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
                await deleteProduk(id);
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
                    <h1 className="text-2xl font-bold text-surface-on tracking-tight">Produk Management</h1>
                    <p className="text-sm text-surface-on-variant mt-1">Manage your produk instances.</p>
                </div>
                <Button variant="primary" onClick={() => handleOpenModal()}>
                    Add Produk
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
                title={editingId ? 'Edit Produk' : 'Add Produk'}
            >
                <form onSubmit={handleSubmit} className="space-y-4 pt-2">
                    <TextField
                        label="Name"
                        name="name"
                        value={formData.name.toString()}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}


                        required
                    />
                    <TextField
                        label="Harga"
                        name="harga"
                        value={formData.harga.toString()}
                        onChange={(e) => setFormData({ ...formData, harga: parseInt(e.target.value) || 0 })}

                        type="number"
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

export default ProdukPage;
