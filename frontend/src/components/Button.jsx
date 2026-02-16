import PropTypes from 'prop-types';

const Button = ({
    children,
    variant = 'primary',
    type = 'button',
    onClick,
    disabled = false,
    className = '',
    fullWidth = false,
}) => {
    const baseStyles = 'btn';

    const variantStyles = {
        primary: 'btn-primary',
        secondary: 'btn-secondary',
        outline: 'btn-outline',
    };

    const classes = `
    ${baseStyles}
    ${variantStyles[variant]}
    ${fullWidth ? 'w-full' : ''}
    ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
    ${className}
  `.trim().replace(/\s+/g, ' ');

    return (
        <button
            type={type}
            onClick={onClick}
            disabled={disabled}
            className={classes}
        >
            {children}
        </button>
    );
};

Button.propTypes = {
    children: PropTypes.node.isRequired,
    variant: PropTypes.oneOf(['primary', 'secondary', 'outline']),
    type: PropTypes.oneOf(['button', 'submit', 'reset']),
    onClick: PropTypes.func,
    disabled: PropTypes.bool,
    className: PropTypes.string,
    fullWidth: PropTypes.bool,
};

export default Button;
