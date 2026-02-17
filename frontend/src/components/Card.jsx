import PropTypes from 'prop-types';

const Card = ({
    children,
    className = '',
    elevated = true,
}) => {
    const classes = `
    card
    ${className}
  `.trim().replace(/\s+/g, ' ');

    return (
        <div className={classes}>
            {children}
        </div>
    );
};

Card.propTypes = {
    children: PropTypes.node.isRequired,
    className: PropTypes.string,
    elevated: PropTypes.bool,
};

export default Card;
