import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import Modal from './Modal';
import TextField from './TextField';
import Button from './Button';
import { getPermissions } from '../api/admin';

const RoleFormModal = ({ isOpen, onClose, onSubmit, role, loading = false }) => {
    const isEdit = !!role;
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        permission_ids: [],
    });
    const [allPermissions, setAllPermissions] = useState([]);
    const [fetchingPermissions, setFetchingPermissions] = useState(false);
    const [errors, setErrors] = useState({});

    useEffect(() => {
        const fetchPermissions = async () => {
            setFetchingPermissions(true);
            try {
                // Fetch a large enough limit to get all permissions for the checkboxes
                const data = await getPermissions(1, 100);
                setAllPermissions(data.data || []);
            } catch (error) {
                console.error('Failed to fetch permissions:', error);
            } finally {
                setFetchingPermissions(false);
            }
        };

        if (isOpen) {
            fetchPermissions();
            if (role) {
                setFormData({
                    name: role.name || '',
                    description: role.description || '',
                    permission_ids: role.permissions?.map(p => p.id) || [],
                });
            } else {
                setFormData({
                    name: '',
                    description: '',
                    permission_ids: [],
                });
            }
            setErrors({});
        }
    }, [role, isOpen]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
        if (errors[name]) {
            setErrors((prev) => ({ ...prev, [name]: '' }));
        }
    };

    const handlePermissionToggle = (permissionId) => {
        setFormData((prev) => {
            const currentIds = [...prev.permission_ids];
            const index = currentIds.indexOf(permissionId);
            if (index > -1) {
                currentIds.splice(index, 1);
            } else {
                currentIds.push(permissionId);
            }
            return { ...prev, permission_ids: currentIds };
        });
    };

    const validate = () => {
        const newErrors = {};
        if (!formData.name.trim()) {
            newErrors.name = 'Name is required';
        }
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        if (validate()) {
            onSubmit(formData);
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={isEdit ? 'Edit Role' : 'Create Role'}
            maxWidth="max-w-2xl"
        >
            <form onSubmit={handleSubmit}>
                <div className="space-y-4 mb-6">
                    <TextField
                        label="Role Name"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        error={errors.name}
                        placeholder="e.g. Moderator"
                        required
                    />

                    <TextField
                        label="Description"
                        name="description"
                        value={formData.description}
                        onChange={handleChange}
                        placeholder="Describe what this role can do"
                    />

                    <div>
                        <label className="text-field-label mb-3">
                            Permissions
                            {fetchingPermissions && <span className="ml-2 text-xs text-surface-on-variant">(Loading...)</span>}
                        </label>
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 bg-surface-variant p-4 rounded-md3 max-h-60 overflow-y-auto">
                            {allPermissions.map((permission) => (
                                <label
                                    key={permission.id}
                                    className="flex items-center space-x-3 cursor-pointer group"
                                >
                                    <div className="relative flex items-center">
                                        <input
                                            type="checkbox"
                                            className="peer h-5 w-5 cursor-pointer appearance-none rounded border border-outline transition-all checked:bg-primary-500 checked:border-primary-500"
                                            checked={formData.permission_ids.includes(permission.id)}
                                            onChange={() => handlePermissionToggle(permission.id)}
                                        />
                                        <span className="absolute text-white opacity-0 peer-checked:opacity-100 top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 pointer-events-none">
                                            <svg className="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="4">
                                                <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                                            </svg>
                                        </span>
                                    </div>
                                    <span className="text-sm text-surface-on group-hover:text-primary transition-colors">
                                        {permission.name}
                                    </span>
                                </label>
                            ))}
                            {allPermissions.length === 0 && !fetchingPermissions && (
                                <p className="text-sm text-surface-on-variant col-span-full italic">No permissions available</p>
                            )}
                        </div>
                    </div>
                </div>

                <div className="flex justify-end gap-3">
                    <Button variant="outline" onClick={onClose} disabled={loading} type="button">
                        Cancel
                    </Button>
                    <Button type="submit" disabled={loading || fetchingPermissions}>
                        {loading ? 'Saving...' : isEdit ? 'Update Role' : 'Create Role'}
                    </Button>
                </div>
            </form>
        </Modal>
    );
};

RoleFormModal.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onSubmit: PropTypes.func.isRequired,
    role: PropTypes.object,
    loading: PropTypes.bool,
};

export default RoleFormModal;
