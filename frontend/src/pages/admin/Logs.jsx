import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import Table from '../../components/Table';
import Pagination from '../../components/Pagination';
import Card from '../../components/Card';
import Modal from '../../components/Modal';
import logApi from '../../api/log';

const Logs = () => {
    const { type: logType } = useParams();
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage] = useState(10);
    const [selectedLog, setSelectedLog] = useState(null);
    const [isDetailsModalOpen, setIsDetailsModalOpen] = useState(false);

    // Reset page to 1 when log type changes
    useEffect(() => {
        setCurrentPage(1);
    }, [logType]);

    const { data, isLoading, error } = useQuery({
        queryKey: ['logs', logType],
        queryFn: () => logApi.getLogs({ type: logType }),
    });

    const allColumns = [
        {
            header: 'Time',
            accessor: 'time',
            render: (row) => new Date(row.time).toLocaleString()
        },
        {
            header: 'Type',
            accessor: 'type',
            render: (row) => (
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${row.type === 'auth' ? 'bg-blue-100 text-blue-800' :
                    row.type === 'audit' ? 'bg-purple-100 text-purple-800' :
                        'bg-gray-100 text-gray-800'
                    }`}>
                    {row.type.toUpperCase()}
                </span>
            )
        },
        {
            header: 'Level',
            accessor: 'level',
            render: (row) => (
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${row.level === 'info' ? 'bg-green-100 text-green-800' :
                    row.level === 'warn' ? 'bg-orange-100 text-orange-800' :
                        'bg-red-100 text-red-800'
                    }`}>
                    {row.level.toUpperCase()}
                </span>
            )
        },
        { header: 'Action', accessor: 'action' },
        { header: 'Email', accessor: 'email' },
        { header: 'Message', accessor: 'message' },
        {
            header: 'Actions',
            accessor: 'id',
            render: (row) => (
                <button
                    onClick={() => {
                        setSelectedLog(row);
                        setIsDetailsModalOpen(true);
                    }}
                    className="text-primary-600 hover:text-primary-800 font-medium transition-colors"
                >
                    Details
                </button>
            )
        }
    ];

    const columns = logType === 'system'
        ? allColumns.filter(col => ['type', 'level', 'message', 'id'].includes(col.accessor))
        : allColumns;

    if (error) {
        return (
            <div className="text-center py-12">
                <p className="text-red-500">Error loading logs: {error.message}</p>
            </div>
        );
    }

    const allLogs = data?.data || [];
    const totalItems = allLogs.length;
    const totalPages = Math.ceil(totalItems / itemsPerPage);

    // Manual pagination for now as the backend returns all logs
    const paginatedLogs = allLogs.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

    return (
        <div className="animate-fade-in">
            <div className="mb-6">
                <h1 className="text-3xl font-bold text-gray-900 tracking-tight">
                    {logType.charAt(0).toUpperCase() + logType.slice(1)} Logs
                </h1>
                <p className="text-gray-600 mt-2">Monitor {logType} activities and trails</p>
            </div>

            <Card className="p-0 overflow-hidden border border-outline-variant/30 ring-1 ring-gray-200 bg-white/80 backdrop-blur-md">
                <Table columns={columns} data={paginatedLogs} loading={isLoading} />
                {!isLoading && allLogs.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={totalPages}
                        totalItems={totalItems}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                    />
                )}
                {!isLoading && allLogs.length === 0 && (
                    <div className="py-20 text-center">
                        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 text-gray-400 mb-4">
                            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <h3 className="text-lg font-medium text-gray-900">No logs found</h3>
                        <p className="text-gray-500">There are no log entries matching your criteria.</p>
                    </div>
                )}
            </Card>
            <Modal
                isOpen={isDetailsModalOpen}
                onClose={() => setIsDetailsModalOpen(false)}
                title="Log Details"
                maxWidth="max-w-4xl"
            >
                {selectedLog && (
                    <div className="space-y-4">
                        <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                                <p className="text-gray-500">Timestamp</p>
                                <p className="font-medium">{new Date(selectedLog.time).toLocaleString()}</p>
                            </div>
                            <div>
                                <p className="text-gray-500">Log Type</p>
                                <p className="font-medium uppercase">{selectedLog.type}</p>
                            </div>
                            <div>
                                <p className="text-gray-500">Level</p>
                                <p className="font-medium uppercase">{selectedLog.level}</p>
                            </div>
                            <div>
                                <p className="text-gray-500">Request ID</p>
                                <p className="font-mono text-xs">{selectedLog.request_id || '-'}</p>
                            </div>
                        </div>

                        <div>
                            <p className="text-gray-500 text-sm mb-2">Message</p>
                            <div className="p-3 bg-gray-50 rounded border border-gray-200 italic">
                                {selectedLog.message}
                            </div>
                        </div>

                        {selectedLog.details && Object.keys(selectedLog.details).length > 0 && (
                            <div>
                                <p className="text-gray-500 text-sm mb-2">Full Details (JSON)</p>
                                <pre className="p-4 bg-gray-900 text-green-400 rounded-lg overflow-x-auto text-xs font-mono max-h-96">
                                    {JSON.stringify(selectedLog.details, null, 2)}
                                </pre>
                            </div>
                        )}

                        <div className="pt-4 flex justify-end">
                            <button
                                onClick={() => setIsDetailsModalOpen(false)}
                                className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-800 rounded-md transition-colors"
                            >
                                Close
                            </button>
                        </div>
                    </div>
                )}
            </Modal>
        </div>
    );
};

export default Logs;
