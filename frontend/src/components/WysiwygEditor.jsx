import ReactQuill from 'react-quill-new';
import 'react-quill-new/dist/quill.snow.css';
import PropTypes from 'prop-types';

const WysiwygEditor = ({
    label,
    name,
    value,
    onChange,
    required = false,
    error,
    helperText,
    className = '',
}) => {
    const handleChange = (content) => {
        // Create a synthetic event object to match other input handlers
        onChange({
            target: {
                name,
                value: content,
            }
        });
    };

    const modules = {
        toolbar: [
            [{ 'header': [1, 2, 3, false] }],
            ['bold', 'italic', 'underline', 'strike'],
            [{ 'list': 'ordered' }, { 'list': 'bullet' }],
            [{ 'color': [] }, { 'background': [] }],
            ['link', 'image'],
            ['clean']
        ],
    };

    return (
        <div className={`w-full ${className}`}>
            {label && (
                <label htmlFor={name} className="text-field-label">
                    {label}
                    {required && <span className="text-red-500 ml-1">*</span>}
                </label>
            )}
            <div className={`wysiwyg-container ${error ? 'border-red-500' : ''}`}>
                <ReactQuill
                    theme="snow"
                    value={value}
                    onChange={handleChange}
                    modules={modules}
                    className="bg-surface-container-high rounded-lg overflow-hidden"
                />
            </div>
            {error && <p className="text-field-error-message">{error}</p>}
            {!error && helperText && (
                <p className="text-field-helper">{helperText}</p>
            )}
            
            <style jsx="true">{`
                .wysiwyg-container .ql-toolbar {
                    border-top-left-radius: 8px;
                    border-top-right-radius: 8px;
                    border-color: rgb(var(--md-sys-color-outline) / 0.5);
                    background-color: rgb(var(--md-sys-color-surface-container-highest));
                }
                .wysiwyg-container .ql-container {
                    border-bottom-left-radius: 8px;
                    border-bottom-right-radius: 8px;
                    border-color: rgb(var(--md-sys-color-outline) / 0.5);
                    min-height: 200px;
                    font-family: inherit;
                    font-size: 0.875rem;
                }
                .dark .wysiwyg-container .ql-snow .ql-stroke {
                    stroke: rgb(var(--md-sys-color-on-surface));
                }
                .dark .wysiwyg-container .ql-snow .ql-fill {
                    fill: rgb(var(--md-sys-color-on-surface));
                }
                .dark .wysiwyg-container .ql-snow .ql-picker {
                    color: rgb(var(--md-sys-color-on-surface));
                }
            `}</style>
        </div>
    );
};

WysiwygEditor.propTypes = {
    label: PropTypes.string,
    name: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    required: PropTypes.bool,
    error: PropTypes.string,
    helperText: PropTypes.string,
    className: PropTypes.string,
};

export default WysiwygEditor;
