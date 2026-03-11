import { useState, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import Table from '../../components/Table';
import Pagination from '../../components/Pagination';
import Card from '../../components/Card';
import Modal from '../../components/Modal';
import logApi from '../../api/log';

const HttpLogs = () => {
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [selectedLog, setSelectedLog] = useState(null);
    const [isDetailsModalOpen, setIsDetailsModalOpen] = useState(false);
    const [activeTab, setActiveTab] = useState('request'); // 'request' | 'response'

    // Filters
    const [methodFilter, setMethodFilter] = useState('');
    const [pathFilter, setPathFilter] = useState('');
    const [statusFilter, setStatusFilter] = useState('');

    const queryParams = useMemo(() => {
        const params = {
            page: currentPage,
            limit: itemsPerPage,
        };
        if (methodFilter) params.method = methodFilter;
        if (pathFilter) params.path = pathFilter;
        if (statusFilter) params.status_code = statusFilter;
        return params;
    }, [currentPage, itemsPerPage, methodFilter, pathFilter, statusFilter]);

    const { data, isLoading, error } = useQuery({
        queryKey: ['http-logs', queryParams],
        queryFn: () => logApi.getHttpLogs(queryParams),
        refetchInterval: 30000, 
    });

    const columns = [
        {
            header: 'Req ID',
            accessor: 'request_id',
            render: (row) => (
                <div className="font-mono text-xs text-surface-on-variant truncate max-w-[80px]" title={row.request_id}>
                    {row.request_id?.split('-')[0] || '-'}
                </div>
            )
        },
        {
            header: 'Method',
            accessor: 'method',
            render: (row) => {
                const colors = {
                    GET: 'bg-blue-500/10 text-blue-600 dark:text-blue-400',
                    POST: 'bg-green-500/10 text-green-600 dark:text-green-400',
                    PUT: 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400',
                    DELETE: 'bg-red-500/10 text-red-600 dark:text-red-400',
                    PATCH: 'bg-orange-500/10 text-orange-600 dark:text-orange-400'
                };
                const colorClass = colors[row.method] || 'bg-surface-variant/30 text-surface-on';
                return (
                    <span className={`px-2 py-1 rounded text-xs font-bold ${colorClass}`}>
                        {row.method}
                    </span>
                );
            }
        },
        {
            header: 'Status',
            accessor: 'status_code',
            render: (row) => {
                const status = row.status_code;
                let colorClass = 'bg-surface-variant/30 text-surface-on';
                if (status >= 200 && status < 300) colorClass = 'bg-green-500/10 text-green-600 dark:text-green-400';
                else if (status >= 300 && status < 400) colorClass = 'bg-blue-500/10 text-blue-600 dark:text-blue-400';
                else if (status >= 400 && status < 500) colorClass = 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400';
                else if (status >= 500) colorClass = 'bg-red-500/10 text-red-600 dark:text-red-400';
                
                return (
                    <span className={`px-2 py-1 rounded text-xs font-bold ${colorClass}`}>
                        {status}
                    </span>
                );
            }
        },
        {
            header: 'Path',
            accessor: 'path',
            render: (row) => (
                <div className="truncate max-w-[200px]" title={row.path}>
                    {row.path}
                </div>
            )
        },
        {
            header: 'User Email',
            accessor: 'user_email',
            render: (row) => (
                <div className="truncate max-w-[150px] text-sm text-surface-on-variant" title={row.user_email || '-'}>
                    {row.user_email || '-'}
                </div>
            )
        },
        {
            header: 'Latency',
            accessor: 'latency',
            render: (row) => {
                const ms = row.latency;
                const color = ms > 1000 ? 'text-red-500' : ms > 500 ? 'text-yellow-500' : 'text-surface-on-variant';
                return <span className={`font-mono text-sm ${color}`}>{ms}ms</span>;
            }
        },
        {
            header: 'Time',
            accessor: 'created_at',
            render: (row) => {
                const date = new Date(row.created_at);
                return (
                    <div className="whitespace-nowrap text-sm">
                        {date.toLocaleString()}
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
                        setActiveTab('request');
                        setIsDetailsModalOpen(true);
                    }}
                    className="text-primary hover:bg-primary-container/20 px-2 py-1 rounded transition-colors font-medium text-sm"
                >
                    Detail
                </button>
            )
        }
    ];

    const parseJSON = (str) => {
        if (!str) return null;
        try {
            return JSON.parse(str);
        } catch (e) {
            return str;
        }
    };

    const renderJsonBlock = (data) => {
        if (!data) return <div className="text-surface-on-variant italic">Empty</div>;
        const parsed = parseJSON(data);
        if (typeof parsed === 'string') {
            return <div className="whitespace-pre-wrap text-surface-on break-all">{parsed}</div>;
        }
        return (
            <pre className="p-4 bg-gray-900 dark:bg-black text-green-400 rounded-lg overflow-auto text-xs font-mono border border-outline-variant/30 relative">
                <button 
                    className="absolute top-2 right-2 p-1.5 bg-white/10 hover:bg-white/20 rounded text-white transition-colors"
                    onClick={(e) => {
                        e.stopPropagation();
                        navigator.clipboard.writeText(JSON.stringify(parsed, null, 2));
                    }}
                    title="Copy to clipboard"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
                </button>
                {JSON.stringify(parsed, null, 2)}
            </pre>
        );
    };

    if (error) {
        return (
            <div className="text-center py-12">
                <p className="text-red-500">Error loading HTTP logs: {error.message}</p>
            </div>
        );
    }

    const logsList = data?.data || [];
    const meta = data?.meta?.pagination || { total_data: 0, total_pages: 0 };

    return (
        <div className="animate-fade-in">
            <div className="mb-6">
                <h1 className="text-3xl font-bold text-surface-on tracking-tight">
                    HTTP Logs
                </h1>
                <p className="text-surface-on-variant mt-2">Monitor incoming HTTP requests and responses</p>
            </div>
            
            <Card className="mb-6 p-4 flex flex-wrap gap-4 items-end bg-surface border border-outline-variant/30">
                <div className="flex-1 min-w-[200px]">
                    <label className="block text-sm font-medium text-surface-on-variant mb-1">Path</label>
                    <input 
                        type="text" 
                        placeholder="/api/v1/..." 
                        value={pathFilter}
                        onChange={(e) => setPathFilter(e.target.value)}
                        className="w-full px-3 py-2 bg-surface hover:bg-surface-variant/30 border border-outline-variant rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary text-surface-on transition-colors"
                    />
                </div>
                <div className="w-[120px]">
                    <label className="block text-sm font-medium text-surface-on-variant mb-1">Method</label>
                    <select 
                        value={methodFilter}
                        onChange={(e) => setMethodFilter(e.target.value)}
                        className="w-full px-3 py-2 bg-surface border border-outline-variant rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-surface-on appearance-none"
                    >
                        <option value="">All</option>
                        <option value="GET">GET</option>
                        <option value="POST">POST</option>
                        <option value="PUT">PUT</option>
                        <option value="DELETE">DELETE</option>
                        <option value="PATCH">PATCH</option>
                    </select>
                </div>
                <div className="w-[120px]">
                    <label className="block text-sm font-medium text-surface-on-variant mb-1">Status</label>
                    <input 
                        type="number" 
                        placeholder="e.g. 200" 
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value)}
                        className="w-full px-3 py-2 bg-surface border border-outline-variant rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-surface-on"
                        min="100" max="599"
                    />
                </div>
                <button 
                    onClick={() => {
                        setPathFilter('');
                        setMethodFilter('');
                        setStatusFilter('');
                        setCurrentPage(1);
                    }}
                    className="px-4 py-2 bg-surface-variant/30 hover:bg-surface-variant/50 text-surface-on rounded-md transition-colors"
                >
                    Clear Filters
                </button>
            </Card>

            <Card className="p-0 overflow-hidden border border-outline-variant/30 bg-surface-container">
                <Table columns={columns} data={logsList} loading={isLoading} hideEmptyState={true} />
                {!isLoading && logsList.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={meta.total_pages}
                        totalItems={meta.total_data}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                        onLimitChange={(newLimit) => {
                            setItemsPerPage(newLimit);
                            setCurrentPage(1);
                        }}
                    />
                )}
                {!isLoading && logsList.length === 0 && (
                    <div className="py-20 text-center">
                        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-surface-variant/20 text-surface-on-variant mb-4">
                            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                            </svg>
                        </div>
                        <h3 className="text-lg font-medium text-surface-on">No logs found</h3>
                        <p className="text-surface-on-variant">There are no HTTP logs matching your criteria.</p>
                    </div>
                )}
            </Card>
            
            <Modal
                isOpen={isDetailsModalOpen}
                onClose={() => setIsDetailsModalOpen(false)}
                title="HTTP Request / Response Detail"
                maxWidth="max-w-5xl"
            >
                {selectedLog && (
                    <div className="flex flex-col h-[70vh]">
                        {/* Summary Header */}
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 p-4 bg-surface-variant/10 rounded-lg border border-outline-variant/30 mb-4 shrink-0">
                            <div>
                                <p className="text-xs text-surface-on-variant mb-1">Method & Path</p>
                                <div className="flex items-center gap-2">
                                    <span className="font-bold text-sm text-surface-on">{selectedLog.method}</span>
                                    <span className="font-mono text-sm text-surface-on truncate" title={selectedLog.path}>{selectedLog.path}</span>
                                </div>
                            </div>
                            <div>
                                <p className="text-xs text-surface-on-variant mb-1">Status & Latency</p>
                                <div className="flex items-center gap-2">
                                    <span className={`font-bold text-sm ${selectedLog.status_code >= 400 ? 'text-red-500' : 'text-green-500'}`}>
                                        {selectedLog.status_code}
                                    </span>
                                    <span className="text-surface-on-variant">•</span>
                                    <span className="font-mono text-sm text-surface-on">{selectedLog.latency}ms</span>
                                </div>
                            </div>
                            <div>
                                <p className="text-xs text-surface-on-variant mb-1">User / IP</p>
                                <div className="flex items-center gap-2">
                                    <span className="text-sm font-medium text-surface-on truncate block" title={selectedLog.user_email || 'Anonymous'}>
                                        {selectedLog.user_email || 'Anonymous'}
                                    </span>
                                    <span className="text-surface-on-variant">•</span>
                                    <span className="text-sm font-mono text-surface-on">{selectedLog.client_ip}</span>
                                </div>
                            </div>
                            <div>
                                <p className="text-xs text-surface-on-variant mb-1">Timestamp</p>
                                <p className="text-sm text-surface-on">{new Date(selectedLog.created_at).toLocaleString()}</p>
                            </div>
                            <div className="col-span-2 md:col-span-4 border-t border-outline-variant/30 pt-2 mt-2">
                                <p className="text-xs text-surface-on-variant mb-1">Req ID / Agent</p>
                                <p className="text-xs font-mono text-surface-on truncate" title={`${selectedLog.request_id} | ${selectedLog.user_agent}`}>
                                    {selectedLog.request_id} | {selectedLog.user_agent}
                                </p>
                            </div>
                            {selectedLog.middleware_trace && (
                                <div className="col-span-2 md:col-span-4 border-t border-outline-variant/30 pt-2">
                                    <p className="text-xs text-surface-on-variant mb-2">Middleware Trace</p>
                                    <div className="flex flex-wrap gap-2">
                                        {selectedLog.middleware_trace.split(' -> ').map((trace, index) => (
                                            <div key={index} className="flex items-center gap-2">
                                                <span className="px-2 py-0.5 rounded bg-primary/10 text-primary text-[10px] font-bold">
                                                    {trace}
                                                </span>
                                                {index < selectedLog.middleware_trace.split(' -> ').length - 1 && (
                                                    <svg className="w-3 h-3 text-surface-on-variant" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M9 5l7 7-7 7" />
                                                    </svg>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Tabs */}
                        <div className="flex border-b border-outline-variant/30 mb-4 shrink-0">
                            <button
                                className={`px-6 py-3 font-medium text-sm transition-colors border-b-2 ${
                                    activeTab === 'request' 
                                        ? 'border-primary text-primary bg-primary/5' 
                                        : 'border-transparent text-surface-on-variant hover:text-surface-on hover:bg-surface-variant/10'
                                }`}
                                onClick={() => setActiveTab('request')}
                            >
                                Request Data
                            </button>
                            <button
                                className={`px-6 py-3 font-medium text-sm transition-colors border-b-2 flex items-center gap-2 ${
                                    activeTab === 'response' 
                                        ? 'border-primary text-primary bg-primary/5' 
                                        : 'border-transparent text-surface-on-variant hover:text-surface-on hover:bg-surface-variant/10'
                                }`}
                                onClick={() => setActiveTab('response')}
                            >
                                Response Data
                                {selectedLog.status_code >= 400 && (
                                    <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse"></span>
                                )}
                            </button>
                        </div>

                        {/* Tab Content */}
                        <div className="flex-1 overflow-y-auto pr-2 custom-scrollbar">
                            {activeTab === 'request' && (
                                <div className="space-y-6">
                                    <div>
                                        <h3 className="text-sm font-semibold text-surface-on mb-3 flex items-center gap-2">
                                            <svg className="w-4 h-4 text-surface-on-variant" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16m-7 6h7" /></svg>
                                            Request Headers
                                        </h3>
                                        {renderJsonBlock(selectedLog.request_headers)}
                                    </div>
                                    <div>
                                        <h3 className="text-sm font-semibold text-surface-on mb-3 flex items-center gap-2">
                                            <svg className="w-4 h-4 text-surface-on-variant" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" /></svg>
                                            Request Body
                                        </h3>
                                        {renderJsonBlock(selectedLog.request_body)}
                                    </div>
                                </div>
                            )}

                            {activeTab === 'response' && (
                                <div className="space-y-6">
                                    <div>
                                        <h3 className="text-sm font-semibold text-surface-on mb-3 flex items-center gap-2">
                                            <svg className="w-4 h-4 text-surface-on-variant" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16m-7 6h7" /></svg>
                                            Response Headers
                                        </h3>
                                        {renderJsonBlock(selectedLog.response_headers)}
                                    </div>
                                    <div>
                                        <h3 className="text-sm font-semibold text-surface-on mb-3 flex items-center gap-2">
                                            <svg className="w-4 h-4 text-surface-on-variant" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" /></svg>
                                            Response Body
                                        </h3>
                                        {renderJsonBlock(selectedLog.response_body)}
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </Modal>
        </div>
    );
};

export default HttpLogs;
