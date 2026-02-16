import { useState } from 'react';
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
    const [showPassword, setShowPassword] = useState(false);
    const isPassword = type === 'password';
    const inputType = isPassword ? (showPassword ? 'text' : 'password') : type;

    const inputClasses = `
    text-field pr-12
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
            <div className="relative">
                <input
                    id={name}
                    type={inputType}
                    name={name}
                    value={value}
                    onChange={onChange}
                    placeholder={placeholder}
                    required={required}
                    disabled={disabled}
                    className={inputClasses}
                />
                {isPassword && (
                    <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute right-3 top-1/2 -translate-y-1/2 p-2 rounded-full hover:bg-surface-variant text-gray-600 transition-all duration-200"
                        title={showPassword ? 'Hide password' : 'Show password'}
                    >
                        {showPassword ? (
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l18 18" />
                            </svg>
                        ) : (
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                            </svg>
                        )}
                    </button>
                )}
            </div>
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
