import PropTypes from 'prop-types';

const Table = ({ columns, data, loading = false, hideEmptyState = false }) => {
    if (loading) {
        return (
            <div className="w-full overflow-x-auto">
                <table className="w-full">
                    <thead className="bg-surface-container-low">
                        <tr>
                            {columns.map((col, index) => (
                                <th key={index} className="px-6 py-4 text-left text-sm font-semibold text-surface-on">
                                    {col.header}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody className="bg-surface-container">
                        {[1, 2, 3, 4, 5].map((row) => (
                            <tr key={row} className="border-b border-outline-variant/30">
                                {columns.map((col, index) => (
                                    <td key={index} className="px-6 py-4">
                                        <div className="h-4 bg-surface-variant/30 rounded animate-pulse"></div>
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        );
    }

    if (!data || data.length === 0) {
        if (hideEmptyState) {
            return null;
        }
        return (
            <div className="w-full text-center py-12">
                <p className="text-surface-on-variant text-lg">No data available</p>
            </div>
        );
    }

    return (
        <div className="w-full overflow-x-auto bg-surface-container rounded-md3-lg border border-outline-variant/30 transition-colors duration-300">
            <table className="w-full">
                <thead className="bg-surface-variant border-b border-outline-variant/30">
                    <tr>
                        {columns.map((col, index) => (
                            <th
                                key={index}
                                className="px-6 py-4 text-left text-sm font-semibold text-surface-on uppercase tracking-wider"
                            >
                                {col.header}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody className="divide-y divide-outline-variant/20">
                    {data.map((row, rowIndex) => (
                        <tr
                            key={row.id || rowIndex}
                            className="hover:bg-primary-container/20 transition-colors duration-200"
                        >
                            {columns.map((col, colIndex) => (
                                <td key={colIndex} className="px-6 py-4 text-sm text-surface-on-variant whitespace-nowrap">
                                    {col.render ? col.render(row) : row[col.accessor]}
                                </td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

Table.propTypes = {
    columns: PropTypes.arrayOf(
        PropTypes.shape({
            header: PropTypes.string.isRequired,
            accessor: PropTypes.string,
            render: PropTypes.func,
        })
    ).isRequired,
    data: PropTypes.array.isRequired,
    loading: PropTypes.bool,
    hideEmptyState: PropTypes.bool,
};

export default Table;
