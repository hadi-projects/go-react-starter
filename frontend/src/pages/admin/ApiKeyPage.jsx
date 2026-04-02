import { useState, useEffect } from 'react';
import { toast } from 'react-hot-toast';
import Button from '../../components/Button';
import Card from '../../components/Card';
import Table from '../../components/Table';
import Modal from '../../components/Modal';
import Pagination from '../../components/Pagination';
import usePermission from '../../hooks/usePermission';
import { PERMS } from '../../utils/permissions';
import { getApiKeys, createApiKey, deleteApiKey } from '../../api/api_key';
import { getRoles } from '../../api/admin';

const ApiKeyPage = () => {
    const can = usePermission();
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [paginationMeta, setPaginationMeta] = useState({ total_data: 0, total_pages: 1 });
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [searchTerm, setSearchTerm] = useState('');
    const [refreshTrigger, setRefreshTrigger] = useState(0);

    // Create Modal state
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
    const [roles, setRoles] = useState([]);
    const [newKeyData, setNewKeyData] = useState({
        name: '',
        type: 'sk_tp',
        role_id: '',
        expires_in_days: 0,
        allowed_ips: ''
    });
    const [createdKey, setCreatedKey] = useState(null);

    // Fetch API Keys
    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                const res = await getApiKeys({
                    page: currentPage,
                    limit: itemsPerPage,
                    search: searchTerm || undefined
                });
                setData(res.data?.data?.data || []);
                setPaginationMeta(res.data?.data?.meta || { total_data: 0, total_pages: 1 });
            } catch (err) {
                toast.error('Failed to fetch API Keys');
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [currentPage, itemsPerPage, searchTerm, refreshTrigger]);

    // Fetch Roles for selection
    useEffect(() => {
        const fetchRoles = async () => {
            try {
                const res = await getRoles(1, 100, '', 'api');
                setRoles(res.data?.data || []);
            } catch (err) {
                console.error('Failed to fetch roles');
            }
        };
        if (isCreateModalOpen) fetchRoles();
    }, [isCreateModalOpen]);

    const handleCreate = async (e) => {
        e.preventDefault();
        try {
            const res = await createApiKey({
                ...newKeyData,
                role_id: parseInt(newKeyData.role_id),
                expires_in_days: parseInt(newKeyData.expires_in_days)
            });
            setCreatedKey(res.data.data);
            toast.success('API Key generated successfully');
            setRefreshTrigger(prev => prev + 1);
        } catch (err) {
            toast.error(err.response?.data?.meta?.message || 'Failed to generate API Key');
        }
    };

    const handleDelete = async (id) => {
        if (!window.confirm('Are you sure you want to revoke this API Key? Any application using it will lose access.')) return;
        try {
            await deleteApiKey(id);
            toast.success('API Key revoked');
            setRefreshTrigger(prev => prev + 1);
        } catch (err) {
            toast.error('Failed to revoke API Key');
        }
    };

    const copyToClipboard = (text) => {
        navigator.clipboard.writeText(text);
        toast.success('Key copied to clipboard');
    };

    const columns = [
        { header: 'Name', accessor: 'name', render: (row) => <span className="font-medium text-surface-on">{row.name}</span> },
        { header: 'Key Prefix', accessor: 'prefix', render: (row) => <code className="text-[11px] bg-surface-variant/30 px-1.5 py-0.5 rounded">{row.prefix}</code> },
        { header: 'Role', accessor: 'role_name', render: (row) => <span className="text-xs px-2 py-0.5 rounded-full bg-primary/10 text-primary font-medium">{row.role_name}</span> },
        { header: 'Last Used', accessor: 'last_used_at', render: (row) => <span className="text-surface-on-variant">{row.last_used_at ? new Date(row.last_used_at).toLocaleString() : 'Never'}</span> },
        { header: 'Allowed IPs', accessor: 'allowed_ips', render: (row) => <span className="text-[10px] bg-surface-variant/20 px-1.5 py-0.5 rounded text-surface-on-variant">{row.allowed_ips || 'Global (Any)'}</span> },
        { header: 'Expires', accessor: 'expires_at', render: (row) => <span className="text-surface-on-variant">{row.expires_at ? new Date(row.expires_at).toLocaleDateString() : 'Never'}</span> },
    ];

    const actions = [
        ...(can(PERMS.DELETE_API_KEY) ? [{
            label: 'Revoke',
            onClick: (row) => handleDelete(row.id),
            className: 'text-error',
            icon: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
        }] : [])
    ];

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-xl font-bold text-surface-on tracking-tight">API Keys</h1>
                    <p className="text-xs text-surface-on-variant mt-0.5">Manage external application access</p>
                </div>
                {can(PERMS.CREATE_API_KEY) && (
                    <Button variant="primary" onClick={() => { setIsCreateModalOpen(true); setCreatedKey(null); }}>
                        Generate Key
                    </Button>
                )}
            </div>

            <Card className="p-0 overflow-hidden">
                <Table
                    columns={columns}
                    data={data}
                    loading={loading}
                    actions={actions}
                />
                {!loading && data.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={paginationMeta.total_pages}
                        totalItems={paginationMeta.total_data}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                        onLimitChange={setItemsPerPage}
                    />
                )}
            </Card>

            <Modal
                isOpen={isCreateModalOpen}
                onClose={() => setIsCreateModalOpen(false)}
                title={createdKey ? "Key Generated" : "Generate New API Key"}
                maxWidth="max-w-md"
            >
                {createdKey ? (
                    <div className="space-y-4">
                        <div className="p-4 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-700 dark:text-amber-400 text-xs">
                            <strong>Warning:</strong> This key will only be shown <b>once</b>. Please copy it and store it securely.
                        </div>
                        <div className="space-y-1.5">
                            <label className="text-[10px] font-bold uppercase tracking-wider text-surface-on-variant">Your API Key</label>
                            <div className="flex gap-2">
                                <code className="flex-1 p-3 rounded-lg bg-surface-container font-mono text-sm break-all">
                                    {createdKey.raw_key}
                                </code>
                                <button
                                    onClick={() => copyToClipboard(createdKey.raw_key)}
                                    className="p-3 rounded-lg bg-primary text-on-primary hover:brightness-110"
                                >
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>
                                </button>
                            </div>
                        </div>
                        <Button className="w-full" onClick={() => setIsCreateModalOpen(false)}>I have saved the key</Button>
                    </div>
                ) : (
                    <form onSubmit={handleCreate} className="space-y-4">
                        <div>
                            <label className="text-field-label">Key Name</label>
                            <input
                                type="text"
                                className="text-field"
                                placeholder="e.g. Mobile App, Staging"
                                value={newKeyData.name}
                                onChange={e => setNewKeyData(prev => ({ ...prev, name: e.target.value }))}
                                required
                            />
                        </div>
                        <div>
                            <label className="text-field-label">Key Type</label>
                            <select
                                className="text-field"
                                value={newKeyData.type}
                                onChange={e => setNewKeyData(prev => ({ ...prev, type: e.target.value }))}
                            >
                                <option value="sk_tp">Private Key (sk_tp_...)</option>
                                <option value="uuid">Standard UUID</option>
                            </select>
                        </div>
                        <div>
                            <label className="text-field-label">Assigned Role</label>
                            <select
                                className="text-field"
                                value={newKeyData.role_id}
                                onChange={e => setNewKeyData(prev => ({ ...prev, role_id: e.target.value }))}
                                required
                            >
                                <option value="">Select Role</option>
                                {roles.map(r => (
                                    <option key={r.id} value={r.id}>{r.name}</option>
                                ))}
                            </select>
                            <p className="text-[10px] text-surface-on-variant mt-1.5 px-1">This key will have all permissions associated with this role.</p>
                        </div>
                        <div>
                            <label className="text-field-label">Expires in (Days)</label>
                            <input
                                type="number"
                                className="text-field"
                                placeholder="0 for never"
                                value={newKeyData.expires_in_days}
                                onChange={e => setNewKeyData(prev => ({ ...prev, expires_in_days: e.target.value }))}
                            />
                        </div>
                        <div>
                            <label className="text-field-label">Allowed IPs</label>
                            <input 
                                type="text"
                                className="text-field"
                                placeholder="e.g. 1.2.3.4, 5.6.7.8 (optional)"
                                value={newKeyData.allowed_ips}
                                onChange={e => setNewKeyData(prev => ({ ...prev, allowed_ips: e.target.value }))}
                            />
                            <p className="text-[10px] text-surface-on-variant mt-1.5 px-1">Comma-separated list of IP addresses. Leave empty for access from any IP.</p>
                        </div>
                        <div className="flex justify-end gap-3 pt-2">
                            <Button type="button" variant="tonal" onClick={() => setIsCreateModalOpen(false)}>Cancel</Button>
                            <Button type="submit" variant="primary">Generate Key</Button>
                        </div>
                    </form>
                )}
            </Modal>
        </div>
    );
};

export default ApiKeyPage;
