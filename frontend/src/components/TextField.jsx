import PropTypes from 'prop-types';

const TextField = ({
    label,
    type = 'text',
    name,
    value,
    onChange,
    placeholder,
    error,
    helperText,
    required = false,
    disabled = false,
    className = '',
}) => {
    const inputClasses = `
    text-field
    ${error ? 'text-field-error' : ''}
    ${className}
  `.trim().replace(/\s+/g, ' ');

    return (
        <div className="w-full">
            {label && (
                <label htmlFor={name} className="text-field-label">
                    {label}
                    {required && <span className="text-red-500 ml-1">*</span>}
                </label>
            )}
            <input
                id={name}
                type={type}
                name={name}
                value={value}
                onChange={onChange}
                placeholder={placeholder}
                required={required}
                disabled={disabled}
                className={inputClasses}
            />
            {error && <p className="text-field-error-message">{error}</p>}
            {!error && helperText && (
                <p className="text-field-helper">{helperText}</p>
            )}
        </div>
    );
};

TextField.propTypes = {
    label: PropTypes.string,
    type: PropTypes.string,
    name: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    placeholder: PropTypes.string,
    error: PropTypes.string,
    helperText: PropTypes.string,
    required: PropTypes.bool,
    disabled: PropTypes.bool,
    className: PropTypes.string,
};

export default TextField;
