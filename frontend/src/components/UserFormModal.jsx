import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import Modal from './Modal';
import TextField from './TextField';
import Button from './Button';

const UserFormModal = ({ isOpen, onClose, onSubmit, user, loading = false }) => {
    const isEdit = !!user;
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        role_id: 2,
    });
    const [errors, setErrors] = useState({});

    useEffect(() => {
        if (user) {
            setFormData({
                email: user.email || '',
                password: '',
                role_id: user.role_id || 2,
            });
        } else {
            setFormData({
                email: '',
                password: '',
                role_id: 2,
            });
        }
        setErrors({});
    }, [user, isOpen]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: name === 'role_id' ? parseInt(value) : value,
        }));
        // Clear error for this field
        if (errors[name]) {
            setErrors((prev) => ({ ...prev, [name]: '' }));
        }
    };

    const validate = () => {
        const newErrors = {};

        if (!formData.email.trim()) {
            newErrors.email = 'Email is required';
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
            newErrors.email = 'Invalid email format';
        }

        if (!isEdit && !formData.password) {
            newErrors.password = 'Password is required';
        }

        if (formData.password && formData.password.length < 6) {
            newErrors.password = 'Password must be at least 6 characters';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        if (validate()) {
            // For edit, only send password if it's filled
            const submitData = { ...formData };
            if (isEdit && !submitData.password) {
                delete submitData.password;
            }
            onSubmit(submitData);
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={isEdit ? 'Edit User' : 'Create User'}
            maxWidth="max-w-lg"
        >
            <form onSubmit={handleSubmit}>
                <div className="space-y-4 mb-6">
                    <TextField
                        label="Email"
                        name="email"
                        type="email"
                        value={formData.email}
                        onChange={handleChange}
                        error={errors.email}
                        required
                    />

                    <TextField
                        label="Password"
                        name="password"
                        type="password"
                        value={formData.password}
                        onChange={handleChange}
                        error={errors.password}
                        helperText={isEdit ? 'Leave empty to keep current password' : ''}
                        required={!isEdit}
                    />

                    <div>
                        <label className="text-field-label">
                            Role ID
                        </label>
                        <select
                            name="role_id"
                            value={formData.role_id}
                            onChange={handleChange}
                            className="text-field"
                        >
                            <option value={1}>Admin (1)</option>
                            <option value={2}>User (2)</option>
                        </select>
                    </div>
                </div>

                <div className="flex justify-end gap-3">
                    <Button variant="outline" onClick={onClose} disabled={loading} type="button">
                        Cancel
                    </Button>
                    <Button type="submit" disabled={loading}>
                        {loading ? 'Saving...' : isEdit ? 'Update' : 'Create'}
                    </Button>
                </div>
            </form>
        </Modal>
    );
};

UserFormModal.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onSubmit: PropTypes.func.isRequired,
    user: PropTypes.object,
    loading: PropTypes.bool,
};

export default UserFormModal;
