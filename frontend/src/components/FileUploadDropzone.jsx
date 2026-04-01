import { useState, useRef } from 'react';
import PropTypes from 'prop-types';

const MAX_SIZE_BYTES = 50 * 1024 * 1024; // 50 MB

const FileUploadDropzone = ({ onUpload, uploading = false, progress = 0 }) => {
    const [isDragging, setIsDragging] = useState(false);
    const [error, setError] = useState('');
    const inputRef = useRef(null);

    const validateAndUpload = (file) => {
        setError('');
        if (!file) return;
        if (file.size > MAX_SIZE_BYTES) {
            setError('File too large. Maximum size is 50 MB.');
            return;
        }
        onUpload(file);
    };

    const handleDrop = (e) => {
        e.preventDefault();
        setIsDragging(false);
        const file = e.dataTransfer.files[0];
        validateAndUpload(file);
    };

    const handleChange = (e) => {
        validateAndUpload(e.target.files[0]);
        e.target.value = '';
    };

    return (
        <div
            onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
            onDragLeave={() => setIsDragging(false)}
            onDrop={handleDrop}
            onClick={() => !uploading && inputRef.current?.click()}
            className={`relative flex flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed p-8 text-center cursor-pointer transition-all duration-200
                ${isDragging
                    ? 'border-primary bg-primary/5 scale-[1.01]'
                    : 'border-outline-variant/50 hover:border-primary/50 hover:bg-surface-variant/10'}
                ${uploading ? 'pointer-events-none opacity-70' : ''}`}
        >
            <input
                ref={inputRef}
                type="file"
                className="hidden"
                onChange={handleChange}
                disabled={uploading}
            />

            {uploading ? (
                <div className="w-full space-y-2">
                    <p className="text-xs font-medium text-surface-on-variant">Uploading…</p>
                    <div className="w-full h-1.5 bg-surface-variant/40 rounded-full overflow-hidden">
                        <div
                            className="h-full bg-primary rounded-full transition-all duration-300"
                            style={{ width: `${progress}%` }}
                        />
                    </div>
                    <p className="text-xs text-surface-on-variant">{progress}%</p>
                </div>
            ) : (
                <>
                    <div className="p-3 rounded-full bg-primary/10 text-primary">
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                        </svg>
                    </div>
                    <div>
                        <p className="text-xs font-semibold text-surface-on">
                            Drop file here or{' '}
                            <span className="text-primary underline">browse</span>
                        </p>
                        <p className="text-[10px] text-surface-on-variant mt-0.5">
                            Maximum file size: 50 MB
                        </p>
                    </div>
                </>
            )}

            {error && (
                <p className="text-xs text-error font-medium mt-1">{error}</p>
            )}
        </div>
    );
};

FileUploadDropzone.propTypes = {
    onUpload: PropTypes.func.isRequired,
    uploading: PropTypes.bool,
    progress: PropTypes.number,
};

export default FileUploadDropzone;
