import PropTypes from 'prop-types';

const Table = ({ columns, data, loading = false }) => {
    if (loading) {
        return (
            <div className="w-full overflow-x-auto">
                <table className="w-full">
                    <thead className="bg-surface-variant">
                        <tr>
                            {columns.map((col, index) => (
                                <th key={index} className="px-6 py-4 text-left text-sm font-semibold text-gray-900">
                                    {col.header}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody>
                        {[1, 2, 3, 4, 5].map((row) => (
                            <tr key={row} className="border-b border-outline-variant">
                                {columns.map((col, index) => (
                                    <td key={index} className="px-6 py-4">
                                        <div className="h-4 bg-gray-200 rounded animate-pulse"></div>
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
        return (
            <div className="w-full text-center py-12">
                <p className="text-gray-500 text-lg">No data available</p>
            </div>
        );
    }

    return (
        <div className="w-full overflow-x-auto bg-white rounded-md3-lg border border-outline-variant/30">
            <table className="w-full">
                <thead className="bg-surface-variant">
                    <tr>
                        {columns.map((col, index) => (
                            <th
                                key={index}
                                className="px-6 py-4 text-left text-sm font-semibold text-gray-900 border-b border-outline-variant"
                            >
                                {col.header}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {data.map((row, rowIndex) => (
                        <tr
                            key={row.id || rowIndex}
                            className="border-b border-outline-variant hover:bg-primary-50 transition-colors duration-150"
                        >
                            {columns.map((col, colIndex) => (
                                <td key={colIndex} className="px-6 py-4 text-sm text-gray-700">
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
};

export default Table;
