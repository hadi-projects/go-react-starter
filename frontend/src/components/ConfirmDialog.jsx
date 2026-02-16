import PropTypes from 'prop-types';
import Modal from './Modal';
import Button from './Button';

const ConfirmDialog = ({ isOpen, onClose, onConfirm, title, message, loading = false }) => {
    return (
        <Modal isOpen={isOpen} onClose={onClose} title={title} maxWidth="max-w-sm">
            <div className="mb-6">
                <p className="text-gray-700">{message}</p>
            </div>

            <div className="flex justify-end gap-3">
                <Button variant="outline" onClick={onClose} disabled={loading}>
                    Cancel
                </Button>
                <Button
                    onClick={onConfirm}
                    disabled={loading}
                    className="bg-red-500 hover:bg-red-600 text-white"
                >
                    {loading ? 'Deleting...' : 'Confirm'}
                </Button>
            </div>
        </Modal>
    );
};

ConfirmDialog.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onConfirm: PropTypes.func.isRequired,
    title: PropTypes.string.isRequired,
    message: PropTypes.string.isRequired,
    loading: PropTypes.bool,
};

export default ConfirmDialog;
