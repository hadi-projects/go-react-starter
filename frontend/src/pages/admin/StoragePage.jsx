import { useState, useEffect, useCallback } from 'react';
import { toast } from 'react-hot-toast';
import Button from '../../components/Button';
import Card from '../../components/Card';
import Table from '../../components/Table';
import Modal from '../../components/Modal';
import Pagination from '../../components/Pagination';
import FileUploadDropzone from '../../components/FileUploadDropzone';
import ShareLinkModal from '../../components/ShareLinkModal';
import usePermission from '../../hooks/usePermission';
import { PERMS } from '../../utils/permissions';
import {
    uploadFile,
    getFiles,
    deleteFile,
    downloadOwnFile,
    getShareLinks,
} from '../../api/storage';

// ── Helpers ───────────────────────────────────────────────────────────────────

const getMimeIcon = (mimeType = '') => {
    if (mimeType.startsWith('image/'))
        return (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
        );
    if (mimeType === 'application/pdf')
        return (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
        );
    if (mimeType.startsWith('video/'))
        return (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M15 10l4.553-2.069A1 1 0 0121 8.87v6.26a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
        );
    if (mimeType.startsWith('audio/'))
        return (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
            </svg>
        );
    if (mimeType.includes('zip') || mimeType.includes('compressed') || mimeType.includes('archive'))
        return (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
            </svg>
        );
    return (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
    );
};

const AccessTypeBadge = ({ type }) => {
    const cfg = {
        one_time:  { label: '1×', color: 'bg-amber-500/10 text-amber-600 dark:text-amber-400' },
        unlimited: { label: '∞',  color: 'bg-green-500/10 text-green-600 dark:text-green-400'  },
        limited:   { label: 'N×', color: 'bg-blue-500/10 text-blue-600 dark:text-blue-400'     },
        timed:     { label: '⏰', color: 'bg-purple-500/10 text-purple-600 dark:text-purple-400' },
    };
    const { label, color } = cfg[type] || cfg.unlimited;
    return (
        <span className={`inline-flex px-1.5 py-0.5 rounded text-[10px] font-bold ${color}`}>
            {label}
        </span>
    );
};

// ── ShareLinks panel (inline, shown below the file row) ───────────────────────

const ShareLinksPanel = ({ file, onCreateNew, onEdit }) => {
    const [links, setLinks] = useState([]);
    const [loading, setLoading] = useState(true);

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const res = await getShareLinks(file.id);
            setLinks(res.data?.data || []);
        } catch {
            // silently fail
        } finally {
            setLoading(false);
        }
    }, [file.id]);

    useEffect(() => { load(); }, [load]);

    const handleCopy = (url) => {
        navigator.clipboard.writeText(url);
        toast.success('Link copied!');
    };

    return (
        <div className="px-4 py-3 bg-surface-container-low border-t border-outline-variant/20 space-y-2">
            <div className="flex items-center justify-between">
                <p className="text-[10px] font-bold text-surface-on-variant uppercase tracking-widest">
                    Share Links
                </p>
                <button
                    onClick={onCreateNew}
                    className="flex items-center gap-1 text-[10px] font-semibold text-primary hover:underline"
                >
                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    New Link
                </button>
            </div>

            {loading && (
                <div className="flex items-center gap-2 py-2">
                    <div className="w-3 h-3 border border-primary border-t-transparent rounded-full animate-spin" />
                    <p className="text-[10px] text-surface-on-variant">Loading…</p>
                </div>
            )}

            {!loading && links.length === 0 && (
                <p className="text-[10px] text-surface-on-variant py-1">
                    No share links yet.{' '}
                    <button onClick={onCreateNew} className="text-primary hover:underline font-medium">
                        Create one
                    </button>
                </p>
            )}

            {!loading && links.map((link) => (
                <div
                    key={link.id}
                    className={`flex items-center gap-2 p-2 rounded-lg border text-[10px] transition-colors
                        ${link.is_active
                            ? 'border-outline-variant/30 bg-surface-container'
                            : 'border-outline-variant/10 bg-surface-variant/10 opacity-50'}`}
                >
                    <AccessTypeBadge type={link.access_type} />
                    <span className="flex-1 truncate font-mono text-surface-on">{link.share_url}</span>
                    <span className="text-surface-on-variant whitespace-nowrap">
                        {link.view_count}{link.max_views ? `/${link.max_views}` : ''} views
                    </span>
                    {link.has_password && (
                        <svg className="w-3 h-3 text-primary flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                    )}
                    {link.is_active && (
                        <button onClick={() => handleCopy(link.share_url)}
                            className="text-primary hover:underline font-semibold whitespace-nowrap">
                            Copy
                        </button>
                    )}
                    <button onClick={() => onEdit(link, load)}
                        className="text-surface-on-variant hover:text-surface-on font-semibold whitespace-nowrap">
                        Edit
                    </button>
                </div>
            ))}
        </div>
    );
};

// ── Main page ─────────────────────────────────────────────────────────────────

const StoragePage = () => {
    const can = usePermission();

    // Data state
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [paginationMeta, setPaginationMeta] = useState({ total_data: 0, total_pages: 1 });
    const [refreshTrigger, setRefreshTrigger] = useState(0);

    // Pagination & search
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [searchTerm, setSearchTerm] = useState('');
    const [debouncedSearch, setDebouncedSearch] = useState('');

    // Upload state
    const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
    const [uploadDescription, setUploadDescription] = useState('');
    const [uploading, setUploading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState(0);
    const [pendingFile, setPendingFile] = useState(null);

    // Share link state
    const [isShareModalOpen, setIsShareModalOpen] = useState(false);
    const [selectedFileForShare, setSelectedFileForShare] = useState(null);
    const [editingShareLink, setEditingShareLink] = useState(null);
    const [shareLinkRefresh, setShareLinkRefresh] = useState(() => () => {});

    // Inline share panel
    const [expandedFileId, setExpandedFileId] = useState(null);

    // Debounce search
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearch(searchTerm);
            setCurrentPage(1);
        }, 300);
        return () => clearTimeout(timer);
    }, [searchTerm]);

    // Fetch files
    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                const res = await getFiles({
                    page: currentPage,
                    limit: itemsPerPage,
                    search: debouncedSearch || undefined,
                });
                setData(res.data?.data?.data || []);
                setPaginationMeta(
                    res.data?.data?.meta || { total_data: 0, total_pages: 1 }
                );
            } catch {
                toast.error('Failed to fetch files');
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [currentPage, itemsPerPage, debouncedSearch, refreshTrigger]);

    const refresh = () => setRefreshTrigger((t) => t + 1);

    // ── Upload handlers ──────────────────────────────────────────────────────

    const handleFilePicked = (file) => {
        setPendingFile(file);
    };

    const handleUploadSubmit = async (e) => {
        e.preventDefault();
        if (!pendingFile) {
            toast.error('Please select a file first.');
            return;
        }
        setUploading(true);
        setUploadProgress(0);
        try {
            const fd = new FormData();
            fd.append('file', pendingFile);
            fd.append('description', uploadDescription);
            await uploadFile(fd, setUploadProgress);
            toast.success(`"${pendingFile.name}" uploaded successfully`);
            setIsUploadModalOpen(false);
            setPendingFile(null);
            setUploadDescription('');
            setUploadProgress(0);
            refresh();
        } catch (err) {
            toast.error(err.response?.data?.meta?.message || 'Upload failed');
        } finally {
            setUploading(false);
        }
    };

    const handleCloseUploadModal = () => {
        if (uploading) return;
        setIsUploadModalOpen(false);
        setPendingFile(null);
        setUploadDescription('');
        setUploadProgress(0);
    };

    // ── Delete ───────────────────────────────────────────────────────────────

    const handleDelete = async (row) => {
        if (!window.confirm(`Delete "${row.original_name}"? This cannot be undone.`)) return;
        try {
            await deleteFile(row.id);
            toast.success('File deleted');
            if (expandedFileId === row.id) setExpandedFileId(null);
            refresh();
        } catch {
            toast.error('Failed to delete file');
        }
    };

    // ── Download ─────────────────────────────────────────────────────────────

    const handleDownload = async (row) => {
        try {
            const res = await downloadOwnFile(row.id);
            const url = window.URL.createObjectURL(new Blob([res.data]));
            const a = document.createElement('a');
            a.href = url;
            a.download = row.original_name;
            document.body.appendChild(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);
        } catch {
            toast.error('Failed to download file');
        }
    };

    // ── Share ────────────────────────────────────────────────────────────────

    const openShareModal = (file, existingLink = null, refreshFn = () => {}) => {
        setSelectedFileForShare(file);
        setEditingShareLink(existingLink);
        setShareLinkRefresh(() => refreshFn);
        setIsShareModalOpen(true);
    };

    const handleShareModalClose = () => {
        setIsShareModalOpen(false);
        setSelectedFileForShare(null);
        setEditingShareLink(null);
        refresh();
        shareLinkRefresh();
    };

    const toggleExpand = (id) => {
        setExpandedFileId((prev) => (prev === id ? null : id));
    };

    // ── Table columns ────────────────────────────────────────────────────────

    const columns = [
        {
            header: '',
            accessor: 'expand',
            render: (row) => (
                <button
                    onClick={() => toggleExpand(row.id)}
                    className="p-1 rounded text-surface-on-variant hover:text-surface-on transition-transform"
                    title={expandedFileId === row.id ? 'Collapse' : 'View share links'}
                >
                    <svg
                        className={`w-3.5 h-3.5 transition-transform duration-200 ${expandedFileId === row.id ? 'rotate-90' : ''}`}
                        fill="none" stroke="currentColor" viewBox="0 0 24 24"
                    >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                </button>
            ),
        },
        {
            header: 'Name',
            accessor: 'original_name',
            render: (row) => (
                <div className="flex items-center gap-2 max-w-xs">
                    <span className="text-surface-on-variant flex-shrink-0">
                        {getMimeIcon(row.mime_type)}
                    </span>
                    <span className="truncate font-medium text-surface-on" title={row.original_name}>
                        {row.original_name}
                    </span>
                </div>
            ),
        },
        {
            header: 'Size',
            accessor: 'size_human',
            render: (row) => (
                <span className="text-surface-on-variant">{row.size_human}</span>
            ),
        },
        {
            header: 'Type',
            accessor: 'mime_type',
            render: (row) => (
                <span className="font-mono text-[10px] text-surface-on-variant bg-surface-variant/30 px-1.5 py-0.5 rounded">
                    {row.mime_type?.split('/')[1] || row.mime_type}
                </span>
            ),
        },
        {
            header: 'Share Links',
            accessor: 'share_count',
            render: (row) => (
                <button
                    onClick={() => toggleExpand(row.id)}
                    className="flex items-center gap-1 group"
                >
                    <svg className="w-3.5 h-3.5 text-surface-on-variant group-hover:text-primary transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                    </svg>
                    <span className={`text-xs font-medium ${row.share_count > 0 ? 'text-primary' : 'text-surface-on-variant'}`}>
                        {row.share_count}
                    </span>
                </button>
            ),
        },
        {
            header: 'Uploaded',
            accessor: 'created_at',
            render: (row) => (
                <span className="text-surface-on-variant">
                    {new Date(row.created_at).toLocaleDateString()}
                </span>
            ),
        },
    ];

    const tableActions = [
        ...(can(PERMS.GET_FILE)
            ? [{
                label: 'Download',
                onClick: handleDownload,
                icon: (
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                ),
            }]
            : []),
        ...(can(PERMS.SHARE_FILE)
            ? [{
                label: 'Share',
                onClick: (row) => openShareModal(row),
                icon: (
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                    </svg>
                ),
            }]
            : []),
        ...(can(PERMS.DELETE_FILE)
            ? [{
                label: 'Delete',
                onClick: handleDelete,
                className: 'text-error',
                icon: (
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                ),
            }]
            : []),
    ];

    // ── Render ────────────────────────────────────────────────────────────────

    return (
        <div className="space-y-5">
            {/* Header */}
            <div className="flex justify-between items-start">
                <div>
                    <h1 className="text-xl font-bold text-surface-on tracking-tight">Storage</h1>
                    <p className="text-xs text-surface-on-variant mt-0.5">
                        {paginationMeta.total_data} file{paginationMeta.total_data !== 1 ? 's' : ''} stored
                    </p>
                </div>
                {can(PERMS.UPLOAD_FILE) && (
                    <Button variant="primary" onClick={() => setIsUploadModalOpen(true)}>
                        <svg className="w-4 h-4 mr-1.5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                        </svg>
                        Upload File
                    </Button>
                )}
            </div>

            {/* Search */}
            <div className="relative">
                <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-surface-on-variant pointer-events-none"
                    fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <input
                    type="text"
                    placeholder="Search by file name or description…"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="text-field pl-9 text-xs"
                />
                {searchTerm && (
                    <button
                        onClick={() => setSearchTerm('')}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-surface-on-variant hover:text-surface-on"
                    >
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                )}
            </div>

            {/* Table */}
            <Card className="p-0 overflow-hidden">
                {/* Custom table rendering to support expandable rows */}
                <div className="w-full overflow-x-auto bg-surface-container border border-outline-variant/30 transition-colors duration-300">
                    <table className="w-full">
                        <thead className="bg-surface-variant border-b border-outline-variant/30">
                            <tr>
                                {[...columns, { header: 'Actions', _isActions: true }].map((col, i) => (
                                    <th
                                        key={i}
                                        className={`px-4 py-2.5 text-left text-xs font-semibold text-surface-on uppercase tracking-wider
                                            ${col._isActions ? 'text-center' : ''}`}
                                    >
                                        {col.header}
                                    </th>
                                ))}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-outline-variant/20">
                            {loading
                                ? [1, 2, 3, 4, 5].map((i) => (
                                    <tr key={i} className="border-b border-outline-variant/30">
                                        {[...columns, {}].map((_, j) => (
                                            <td key={j} className="px-4 py-2.5">
                                                <div className="h-3.5 bg-surface-variant/30 rounded animate-pulse" />
                                            </td>
                                        ))}
                                    </tr>
                                ))
                                : data.length === 0
                                    ? (
                                        <tr>
                                            <td colSpan={columns.length + 1} className="text-center py-12">
                                                <svg className="w-10 h-10 mx-auto text-surface-on-variant/20 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                                                        d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                                                </svg>
                                                <p className="text-sm text-surface-on-variant">
                                                    {debouncedSearch ? 'No files match your search.' : 'No files uploaded yet.'}
                                                </p>
                                                {!debouncedSearch && can(PERMS.UPLOAD_FILE) && (
                                                    <button
                                                        onClick={() => setIsUploadModalOpen(true)}
                                                        className="mt-2 text-xs text-primary hover:underline font-medium"
                                                    >
                                                        Upload your first file
                                                    </button>
                                                )}
                                            </td>
                                        </tr>
                                    )
                                    : data.map((row) => (
                                        <>
                                            <tr
                                                key={row.id}
                                                className={`hover:bg-primary-container/20 transition-colors duration-200
                                                    ${expandedFileId === row.id ? 'bg-primary-container/10' : ''}`}
                                            >
                                                {columns.map((col, colIdx) => (
                                                    <td key={colIdx} className="px-4 py-2 text-xs text-surface-on-variant whitespace-nowrap">
                                                        {col.render ? col.render(row) : row[col.accessor]}
                                                    </td>
                                                ))}
                                                <td className="px-4 py-2 text-xs whitespace-nowrap">
                                                    <div className="flex justify-center gap-1">
                                                        {tableActions.map((action, i) => (
                                                            <button
                                                                key={i}
                                                                onClick={() => action.onClick(row)}
                                                                title={action.label}
                                                                className={`p-1.5 rounded-full hover:bg-surface-variant/40 transition-colors ${action.className || 'text-surface-on-variant'}`}
                                                            >
                                                                <span className="w-4 h-4 block">{action.icon}</span>
                                                            </button>
                                                        ))}
                                                    </div>
                                                </td>
                                            </tr>
                                            {expandedFileId === row.id && (
                                                <tr key={`${row.id}-expand`}>
                                                    <td colSpan={columns.length + 1} className="p-0">
                                                        <ShareLinksPanel
                                                            file={row}
                                                            onCreateNew={() => openShareModal(row)}
                                                            onEdit={(link, reloadFn) => openShareModal(row, link, reloadFn)}
                                                        />
                                                    </td>
                                                </tr>
                                            )}
                                        </>
                                    ))}
                        </tbody>
                    </table>
                </div>

                {!loading && data.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={paginationMeta.total_pages}
                        totalItems={paginationMeta.total_data}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                        onLimitChange={(limit) => { setItemsPerPage(limit); setCurrentPage(1); }}
                    />
                )}
            </Card>

            {/* ── Upload modal ─────────────────────────────────────── */}
            <Modal
                isOpen={isUploadModalOpen}
                onClose={handleCloseUploadModal}
                title="Upload File"
                maxWidth="max-w-md"
            >
                <form onSubmit={handleUploadSubmit} className="space-y-4 pt-2">
                    <FileUploadDropzone
                        onUpload={handleFilePicked}
                        uploading={uploading}
                        progress={uploadProgress}
                    />

                    {pendingFile && !uploading && (
                        <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-primary/5 border border-primary/20">
                            <span className="text-primary flex-shrink-0">
                                {getMimeIcon(pendingFile.type)}
                            </span>
                            <div className="flex-1 min-w-0">
                                <p className="text-xs font-semibold text-surface-on truncate">{pendingFile.name}</p>
                                <p className="text-[10px] text-surface-on-variant">
                                    {(pendingFile.size / (1024 * 1024)).toFixed(2)} MB
                                </p>
                            </div>
                            <button
                                type="button"
                                onClick={() => setPendingFile(null)}
                                className="text-surface-on-variant hover:text-surface-on"
                            >
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            </button>
                        </div>
                    )}

                    <div>
                        <label className="text-field-label">Description (optional)</label>
                        <textarea
                            value={uploadDescription}
                            onChange={(e) => setUploadDescription(e.target.value)}
                            placeholder="Add a description for this file…"
                            rows={2}
                            className="text-field resize-none text-xs"
                            disabled={uploading}
                        />
                    </div>

                    <div className="flex justify-end gap-2 pt-1">
                        <Button type="button" variant="tonal" onClick={handleCloseUploadModal} disabled={uploading}>
                            Cancel
                        </Button>
                        <Button type="submit" variant="primary" disabled={!pendingFile || uploading}>
                            {uploading ? `Uploading ${uploadProgress}%…` : 'Upload'}
                        </Button>
                    </div>
                </form>
            </Modal>

            {/* ── Share link modal ─────────────────────────────────── */}
            <ShareLinkModal
                isOpen={isShareModalOpen}
                onClose={handleShareModalClose}
                fileId={selectedFileForShare?.id}
                existingLink={editingShareLink}
                onRefresh={() => {
                    refresh();
                    shareLinkRefresh();
                }}
            />
        </div>
    );
};

export default StoragePage;
