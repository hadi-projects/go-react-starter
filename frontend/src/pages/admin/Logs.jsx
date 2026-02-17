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
            header: 'Level',
            accessor: 'level',
            render: (row) => (
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${row.level === 'info' ? 'bg-primary-container/20 text-primary' :
                    row.level === 'warn' ? 'bg-secondary-container/20 text-secondary' :
                        'bg-error-container/20 text-error'
                    }`}>
                    {row.level.toUpperCase()}
                </span>
            )
        },
        {
            header: 'Action',
            accessor: 'action',
            render: (row) => (
                <div className="truncate max-w-xs whitespace-nowrap overflow-hidden">
                    {row.action}
                </div>
            )
        },
        {
            header: 'Email',
            accessor: 'email',
            render: (row) => (
                <div className="truncate max-w-xs whitespace-nowrap overflow-hidden">
                    {row.email}
                </div>
            )
        },
        {
            header: 'Message',
            accessor: 'message',
            render: (row) => (
                <div className="truncate max-w-xs whitespace-nowrap overflow-hidden">
                    {row.message}
                </div>
            )
        },
        {
            header: 'Time',
            accessor: 'time',
            render: (row) => {
                const date = new Date(row.time);
                const day = String(date.getDate()).padStart(2, '0');
                const month = String(date.getMonth() + 1).padStart(2, '0');
                const year = String(date.getFullYear()).slice(-2);
                const hours = String(date.getHours()).padStart(2, '0');
                const minutes = String(date.getMinutes()).padStart(2, '0');
                const seconds = String(date.getSeconds()).padStart(2, '0');
                return (
                    <div className="whitespace-nowrap">
                        {`${day}/${month}/${year} ${hours}:${minutes}:${seconds}`}
                    </div>
                );
            }
        },
        {
            header: 'Actions',
            accessor: 'id',
            render: (row) => (
                <button
                    onClick={() => {
                        setSelectedLog(row);
                        setIsDetailsModalOpen(true);
                    }}
                    className="text-primary hover:bg-primary-container/20 px-2 py-1 rounded transition-colors font-medium"
                >
                    Details
                </button>
            )
        }
    ];

    const columns = logType === 'system'
        ? allColumns.filter(col => ['level', 'message', 'time', 'id'].includes(col.accessor))
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
                <h1 className="text-3xl font-bold text-surface-on tracking-tight">
                    {logType.charAt(0).toUpperCase() + logType.slice(1)} Logs
                </h1>
                <p className="text-surface-on-variant mt-2">Monitor {logType} activities and trails</p>
            </div>

            <Card className="p-0 overflow-hidden border border-outline-variant/30 bg-surface-container">
                <Table columns={columns} data={paginatedLogs} loading={isLoading} hideEmptyState={true} />
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
                        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-surface-variant/20 text-surface-on-variant mb-4">
                            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <h3 className="text-lg font-medium text-surface-on">No logs found</h3>
                        <p className="text-surface-on-variant">There are no log entries matching your criteria.</p>
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
                                <p className="text-surface-on-variant text-xs opacity-70">Timestamp</p>
                                <p className="font-medium text-surface-on">{new Date(selectedLog.time).toLocaleString()}</p>
                            </div>
                            <div>
                                <p className="text-surface-on-variant text-xs opacity-70">Log Type</p>
                                <p className="font-medium text-surface-on uppercase">{selectedLog.type}</p>
                            </div>
                            <div>
                                <p className="text-surface-on-variant text-xs opacity-70">Level</p>
                                <p className="font-medium text-surface-on uppercase">{selectedLog.level}</p>
                            </div>
                            <div>
                                <p className="text-surface-on-variant text-xs opacity-70">Request ID</p>
                                <p className="font-mono text-xs text-surface-on">{selectedLog.request_id || '-'}</p>
                            </div>
                        </div>

                        <div>
                            <p className="text-surface-on-variant text-sm mb-2 opacity-70">Message</p>
                            <div className="p-3 bg-surface-variant/10 rounded border border-outline-variant/30 italic text-surface-on">
                                {selectedLog.message}
                            </div>
                        </div>

                        {selectedLog.details && Object.keys(selectedLog.details).length > 0 && (
                            <div>
                                <p className="text-surface-on-variant text-sm mb-2 opacity-70">Full Details (JSON)</p>
                                <pre className="p-4 bg-gray-900 dark:bg-black text-green-400 rounded-lg overflow-x-auto text-xs font-mono max-h-96 border border-outline-variant/30">
                                    {JSON.stringify(selectedLog.details, null, 2)}
                                </pre>
                            </div>
                        )}

                        <div className="pt-4 flex justify-end">
                            <button
                                onClick={() => setIsDetailsModalOpen(false)}
                                className="px-4 py-2 bg-surface-variant/20 hover:bg-surface-variant/30 text-surface-on rounded-md3 transition-colors"
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
