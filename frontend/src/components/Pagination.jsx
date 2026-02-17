import PropTypes from 'prop-types';
import Button from './Button';

const Pagination = ({ currentPage, totalPages, onPageChange, totalItems, itemsPerPage, onLimitChange }) => {
    const startItem = (currentPage - 1) * itemsPerPage + 1;
    const endItem = Math.min(currentPage * itemsPerPage, totalItems);

    const getPageNumbers = () => {
        const pages = [];
        const maxVisible = 5;

        if (totalPages <= maxVisible) {
            for (let i = 1; i <= totalPages; i++) {
                pages.push(i);
            }
        } else {
            if (currentPage <= 3) {
                for (let i = 1; i <= 4; i++) pages.push(i);
                pages.push('...');
                pages.push(totalPages);
            } else if (currentPage >= totalPages - 2) {
                pages.push(1);
                pages.push('...');
                for (let i = totalPages - 3; i <= totalPages; i++) pages.push(i);
            } else {
                pages.push(1);
                pages.push('...');
                for (let i = currentPage - 1; i <= currentPage + 1; i++) pages.push(i);
                pages.push('...');
                pages.push(totalPages);
            }
        }

        return pages;
    };

    return (
        <div className="flex flex-col sm:flex-row items-center justify-between px-6 py-4 bg-surface-container border-t border-outline-variant/30 gap-4 transition-colors duration-300">
            <div className="flex items-center gap-4">
                <div className="text-sm text-surface-on-variant">
                    Showing <span className="font-semibold text-surface-on">{startItem}</span> to <span className="font-semibold text-surface-on">{endItem}</span> of <span className="font-semibold text-surface-on">{totalItems}</span> results
                </div>

                {onLimitChange && (
                    <div className="flex items-center gap-2">
                        <label htmlFor="limit" className="text-sm text-surface-on-variant opacity-70">Show:</label>
                        <select
                            id="limit"
                            value={itemsPerPage}
                            onChange={(e) => onLimitChange(Number(e.target.value))}
                            className="text-sm border border-outline rounded p-1 focus:ring-2 focus:ring-primary/50 focus:outline-none bg-surface-container-high text-surface-on transition-all shadow-sm"
                        >
                            {[10, 20, 50, 100].map(limit => (
                                <option key={limit} value={limit}>{limit}</option>
                            ))}
                        </select>
                    </div>
                )}
            </div>

            <div className="flex items-center gap-2">
                <Button
                    variant="outline"
                    onClick={() => onPageChange(currentPage - 1)}
                    disabled={currentPage === 1}
                    className="px-3 py-2 text-sm h-9 flex items-center shadow-sm border-outline-variant/50 hover:bg-surface-variant/20"
                >
                    <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                    Previous
                </Button>

                <div className="flex items-center gap-1">
                    {getPageNumbers().map((page, index) => (
                        page === '...' ? (
                            <span key={`ellipsis-${index}`} className="px-2 text-surface-on-variant opacity-50">
                                ...
                            </span>
                        ) : (
                            <button
                                key={page}
                                onClick={() => onPageChange(page)}
                                className={`w-9 h-9 text-sm font-medium rounded transition-all duration-200 ${currentPage === page
                                    ? 'bg-primary text-on-primary shadow-md transform scale-105'
                                    : 'text-surface-on-variant hover:bg-primary-container/20 hover:text-primary border border-transparent'
                                    }`}
                            >
                                {page}
                            </button>
                        )
                    ))}
                </div>

                <Button
                    variant="outline"
                    onClick={() => onPageChange(currentPage + 1)}
                    disabled={currentPage === totalPages}
                    className="px-3 py-2 text-sm h-9 flex items-center shadow-sm border-outline-variant/50 hover:bg-surface-variant/20"
                >
                    Next
                    <svg className="w-4 h-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                </Button>
            </div>
        </div>
    );
};

Pagination.propTypes = {
    currentPage: PropTypes.number.isRequired,
    totalPages: PropTypes.number.isRequired,
    onPageChange: PropTypes.func.isRequired,
    totalItems: PropTypes.number.isRequired,
    itemsPerPage: PropTypes.number.isRequired,
    onLimitChange: PropTypes.func,
};

export default Pagination;
