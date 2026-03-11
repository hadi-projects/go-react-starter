import { useEffect } from 'react';
import Modal from './Modal';
import Button from './Button';

const PermissionFormModal = ({ isOpen, onClose, onSubmit, permission, isLoading }) => {
    useEffect(() => {
        if (isOpen && permission) {
            // Pre-fill form when editing
            document.getElementById('permission-name').value = permission.name || '';
            document.getElementById('permission-description').value = permission.description || '';
        } else if (isOpen) {
            // Clear form when creating
            document.getElementById('permission-name').value = '';
            document.getElementById('permission-description').value = '';
        }
    }, [isOpen, permission]);

    const handleSubmit = (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const data = {
            name: formData.get('name'),
            description: formData.get('description'),
        };
        onSubmit(data);
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={permission ? 'Edit Permission' : 'Create Permission'}
        >
            <form onSubmit={handleSubmit}>
                <div className="space-y-4">
                    <div>
                        <label htmlFor="permission-name" className="text-field-label">
                            Name <span className="text-error">*</span>
                        </label>
                        <input
                            type="text"
                            id="permission-name"
                            name="name"
                            required
                            className="text-field"
                            placeholder="e.g., create-user"
                        />
                    </div>

                    <div>
                        <label htmlFor="permission-description" className="text-field-label">
                            Description
                        </label>
                        <textarea
                            id="permission-description"
                            name="description"
                            rows="3"
                            className="text-field"
                            placeholder="Describe what this permission allows..."
                        />
                    </div>
                </div>

                <div className="mt-6 flex justify-end gap-3">
                    <Button
                        type="button"
                        variant="outline"
                        onClick={onClose}
                        disabled={isLoading}
                    >
                        Cancel
                    </Button>
                    <Button
                        type="submit"
                        disabled={isLoading}
                    >
                        {isLoading ? 'Saving...' : permission ? 'Update' : 'Create'}
                    </Button>
                </div>
            </form>
        </Modal>
    );
};

export default PermissionFormModal;
