import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useTheme } from '../context/ThemeContext';
import { getPublicFileInfo, downloadViaShareLink } from '../api/storage';

// ── Mime type icon ────────────────────────────────────────────────────────────

const MimeIcon = ({ mimeType = '', size = 'lg' }) => {
    const cls = size === 'lg' ? 'w-12 h-12' : 'w-5 h-5';

    if (mimeType.startsWith('image/'))
        return (
            <svg className={cls} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
        );
    if (mimeType === 'application/pdf')
        return (
            <svg className={cls} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
        );
    if (mimeType.startsWith('video/'))
        return (
            <svg className={cls} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M15 10l4.553-2.069A1 1 0 0121 8.87v6.26a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
        );
    if (mimeType.startsWith('audio/'))
        return (
            <svg className={cls} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
            </svg>
        );
    if (mimeType.includes('zip') || mimeType.includes('compressed') || mimeType.includes('archive'))
        return (
            <svg className={cls} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                    d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
            </svg>
        );
    return (
        <svg className={cls} fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
    );
};

// ── Access type badge ─────────────────────────────────────────────────────────

const AccessTypeBadge = ({ type }) => {
    const cfg = {
        one_time:  { label: 'One-time view',   color: 'bg-amber-500/10 text-amber-600 dark:text-amber-400',  icon: '👁' },
        unlimited: { label: 'Always available', color: 'bg-green-500/10 text-green-600 dark:text-green-400',  icon: '∞' },
        limited:   { label: 'Limited views',    color: 'bg-blue-500/10 text-blue-600 dark:text-blue-400',     icon: '🔢' },
        timed:     { label: 'Time-limited',      color: 'bg-purple-500/10 text-purple-600 dark:text-purple-400', icon: '⏰' },
    };
    const { label, color, icon } = cfg[type] || cfg.unlimited;
    return (
        <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium ${color}`}>
            <span>{icon}</span>
            {label}
        </span>
    );
};

// ── Skeleton loader ───────────────────────────────────────────────────────────

const Skeleton = () => (
    <div className="animate-pulse space-y-4">
        <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-xl bg-surface-variant/40" />
            <div className="flex-1 space-y-2">
                <div className="h-4 bg-surface-variant/40 rounded w-2/3" />
                <div className="h-3 bg-surface-variant/30 rounded w-1/3" />
            </div>
        </div>
        <div className="h-px bg-outline-variant/30" />
        <div className="space-y-2">
            <div className="h-3 bg-surface-variant/30 rounded w-1/2" />
            <div className="h-3 bg-surface-variant/30 rounded w-1/3" />
            <div className="h-3 bg-surface-variant/30 rounded w-2/5" />
        </div>
        <div className="h-10 bg-surface-variant/30 rounded-xl" />
    </div>
);

// ── Main page ─────────────────────────────────────────────────────────────────

const SharePage = () => {
    const { token } = useParams();
    const { theme, toggleTheme } = useTheme();

    // page state: 'loading' | 'requires_password' | 'ready' | 'error'
    const [pageState, setPageState] = useState('loading');
    const [fileInfo, setFileInfo] = useState(null);
    const [errorType, setErrorType] = useState(''); // 'not_found' | 'forbidden' | 'generic'

    // password flow
    const [password, setPassword] = useState('');
    const [passwordError, setPasswordError] = useState('');
    const [checkingPassword, setCheckingPassword] = useState(false);
    const [passwordVerified, setPasswordVerified] = useState(false);

    // download
    const [downloading, setDownloading] = useState(false);
    const [downloadError, setDownloadError] = useState('');
    const [downloadDone, setDownloadDone] = useState(false);

    useEffect(() => {
        fetchInfo();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [token]);

    const fetchInfo = async () => {
        setPageState('loading');
        setFileInfo(null);
        setErrorType('');
        setPasswordVerified(false);
        setDownloadDone(false);
        try {
            const res = await getPublicFileInfo(token);
            const info = res.data?.data;
            setFileInfo(info);
            if (info?.requires_password) {
                setPageState('requires_password');
            } else {
                setPageState('ready');
            }
        } catch (err) {
            const status = err.response?.status;
            if (status === 404) setErrorType('not_found');
            else if (status === 403) setErrorType('forbidden');
            else setErrorType('generic');
            setPageState('error');
        }
    };

    // Verify password by attempting a metadata fetch — we re-check on the
    // download endpoint because the info endpoint doesn't validate the password.
    // Instead we do a small "probe" download attempt with the password.
    const handleVerifyPassword = async (e) => {
        e.preventDefault();
        setPasswordError('');
        setCheckingPassword(true);
        try {
            // Attempt a HEAD-like approach: try the download with ?check=1.
            // Since the backend doesn't support HEAD, we try the info endpoint
            // and then validate on actual download. For UX, we accept the password
            // and let any error surface on actual download.
            // A better alternative: attempt download but abort immediately on success.
            // Here we do a download probe (small content) and immediately revoke.
            const res = await downloadViaShareLink(token, password);
            if (res.status === 200) {
                // Password OK — revoke the blob URL immediately (don't trigger download yet)
                setPasswordVerified(true);
                setPageState('ready');
            }
        } catch (err) {
            const status = err.response?.status;
            if (status === 401) {
                setPasswordError('Incorrect password. Please try again.');
            } else if (status === 403) {
                setErrorType('forbidden');
                setPageState('error');
            } else if (status === 404) {
                setErrorType('not_found');
                setPageState('error');
            } else {
                setPasswordError('Unable to verify password. Please try again.');
            }
        } finally {
            setCheckingPassword(false);
        }
    };

    const handleDownload = async () => {
        setDownloading(true);
        setDownloadError('');
        try {
            const res = await downloadViaShareLink(token, passwordVerified ? password : '');
            const disposition = res.headers?.['content-disposition'] || '';
            const filenameMatch = disposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
            const filename = filenameMatch
                ? filenameMatch[1].replace(/['"]/g, '')
                : fileInfo?.original_name || 'download';

            const url = window.URL.createObjectURL(new Blob([res.data], { type: fileInfo?.mime_type }));

            if (fileInfo?.allow_download) {
                // Force download
                const a = document.createElement('a');
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                a.remove();
            } else {
                // Open inline in new tab (view-only)
                window.open(url, '_blank');
            }

            window.URL.revokeObjectURL(url);
            setDownloadDone(true);

            // For one_time links, show consumed notice after a short delay
            if (fileInfo?.access_type === 'one_time') {
                setTimeout(() => fetchInfo(), 1500);
            }
        } catch (err) {
            const status = err.response?.status;
            if (status === 403) {
                setDownloadError('This link is no longer accessible. It may have been consumed or revoked.');
            } else if (status === 401) {
                setDownloadError('Password is required or incorrect.');
                setPageState('requires_password');
                setPasswordVerified(false);
            } else {
                setDownloadError('Download failed. Please try again.');
            }
        } finally {
            setDownloading(false);
        }
    };

    // ── Render helpers ────────────────────────────────────────────────────────

    const renderError = () => {
        const configs = {
            not_found: {
                icon: (
                    <svg className="w-14 h-14 text-surface-on-variant/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                            d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                    </svg>
                ),
                title: 'Link Not Found',
                message: 'This share link does not exist or has already expired.',
            },
            forbidden: {
                icon: (
                    <svg className="w-14 h-14 text-surface-on-variant/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                            d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                ),
                title: 'Access Denied',
                message: 'This link is no longer accessible. It may have reached its view limit, expired, or been manually revoked.',
            },
            generic: {
                icon: (
                    <svg className="w-14 h-14 text-surface-on-variant/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                ),
                title: 'Something Went Wrong',
                message: 'We could not load this file. Please check the link and try again.',
            },
        };
        const cfg = configs[errorType] || configs.generic;
        return (
            <div className="text-center py-6 space-y-3">
                <div className="flex justify-center">{cfg.icon}</div>
                <h2 className="text-base font-bold text-surface-on">{cfg.title}</h2>
                <p className="text-xs text-surface-on-variant max-w-xs mx-auto">{cfg.message}</p>
                <div className="pt-2">
                    <Link
                        to="/"
                        className="text-xs text-primary hover:underline font-medium"
                    >
                        ← Go to homepage
                    </Link>
                </div>
            </div>
        );
    };

    const renderPasswordForm = () => (
        <div className="space-y-5">
            {/* File preview (name + size, no detailed info) */}
            <div className="flex items-center gap-3 p-3 rounded-xl bg-surface-container border border-outline-variant/30">
                <div className="p-2.5 rounded-lg bg-surface-variant/30 text-surface-on-variant">
                    <MimeIcon mimeType={fileInfo?.mime_type} size="sm" />
                </div>
                <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-surface-on truncate">{fileInfo?.original_name}</p>
                    <p className="text-xs text-surface-on-variant">{fileInfo?.size_human}</p>
                </div>
            </div>

            <div className="flex items-center gap-2.5 p-3 rounded-lg bg-amber-500/10 border border-amber-500/20">
                <svg className="w-4 h-4 text-amber-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                        d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                <p className="text-xs font-medium text-surface-on">This file is password protected.</p>
            </div>

            <form onSubmit={handleVerifyPassword} className="space-y-3">
                <div>
                    <label className="text-field-label">Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => { setPassword(e.target.value); setPasswordError(''); }}
                        placeholder="Enter password to access this file"
                        className="text-field"
                        autoFocus
                        required
                    />
                    {passwordError && (
                        <p className="text-field-error-message">{passwordError}</p>
                    )}
                </div>
                <button
                    type="submit"
                    disabled={checkingPassword || !password}
                    className="w-full flex items-center justify-center gap-2 py-2.5 px-4 rounded-full bg-primary text-on-primary text-sm font-semibold transition-all hover:brightness-110 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {checkingPassword ? (
                        <>
                            <span className="w-4 h-4 border-2 border-on-primary border-t-transparent rounded-full animate-spin" />
                            Verifying…
                        </>
                    ) : (
                        <>
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                    d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z" />
                            </svg>
                            Unlock File
                        </>
                    )}
                </button>
            </form>
        </div>
    );

    const renderFileInfo = () => (
        <div className="space-y-5">
            {/* File card */}
            <div className="flex items-center gap-4 p-4 rounded-xl bg-surface-container border border-outline-variant/30">
                <div className="p-3 rounded-xl bg-primary/10 text-primary flex-shrink-0">
                    <MimeIcon mimeType={fileInfo?.mime_type} size="lg" />
                </div>
                <div className="flex-1 min-w-0">
                    <h2 className="text-sm font-bold text-surface-on truncate" title={fileInfo?.original_name}>
                        {fileInfo?.original_name}
                    </h2>
                    <p className="text-xs text-surface-on-variant mt-0.5">
                        {fileInfo?.size_human}
                        <span className="mx-1.5 opacity-40">·</span>
                        <span className="font-mono">{fileInfo?.mime_type}</span>
                    </p>
                </div>
            </div>

            {/* Meta info */}
            <div className="space-y-2">
                {fileInfo?.label && (
                    <div className="flex items-center justify-between text-xs">
                        <span className="text-surface-on-variant">Label</span>
                        <span className="font-medium text-surface-on">{fileInfo.label}</span>
                    </div>
                )}

                <div className="flex items-center justify-between text-xs">
                    <span className="text-surface-on-variant">Access type</span>
                    <AccessTypeBadge type={fileInfo?.access_type} />
                </div>

                {fileInfo?.access_type === 'limited' && fileInfo?.max_views != null && (
                    <div className="flex items-center justify-between text-xs">
                        <span className="text-surface-on-variant">Views remaining</span>
                        <span className="font-semibold text-surface-on">
                            {Math.max(0, fileInfo.max_views - fileInfo.view_count)} / {fileInfo.max_views}
                        </span>
                    </div>
                )}

                {fileInfo?.access_type === 'timed' && fileInfo?.expires_at && (
                    <div className="flex items-center justify-between text-xs">
                        <span className="text-surface-on-variant">Expires</span>
                        <span className="font-medium text-surface-on">
                            {new Date(fileInfo.expires_at).toLocaleString()}
                        </span>
                    </div>
                )}

                <div className="flex items-center justify-between text-xs">
                    <span className="text-surface-on-variant">Download</span>
                    <span className={`font-medium ${fileInfo?.allow_download ? 'text-green-500' : 'text-amber-500'}`}>
                        {fileInfo?.allow_download ? 'Allowed' : 'View only'}
                    </span>
                </div>

                {fileInfo?.requires_password && (
                    <div className="flex items-center justify-between text-xs">
                        <span className="text-surface-on-variant">Password</span>
                        <span className="font-medium text-primary">Protected</span>
                    </div>
                )}
            </div>

            <div className="h-px bg-outline-variant/30" />

            {/* One-time warning */}
            {fileInfo?.access_type === 'one_time' && !downloadDone && (
                <div className="flex items-start gap-2.5 p-3 rounded-lg bg-amber-500/10 border border-amber-500/20 text-xs">
                    <svg className="w-4 h-4 text-amber-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                    <p className="text-surface-on">
                        <span className="font-bold">One-time link: </span>
                        This link will be permanently deactivated after you access the file.
                    </p>
                </div>
            )}

            {/* Download success state */}
            {downloadDone && fileInfo?.access_type !== 'one_time' && (
                <div className="flex items-center gap-2.5 p-3 rounded-lg bg-green-500/10 border border-green-500/20 text-xs">
                    <svg className="w-4 h-4 text-green-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <p className="text-surface-on font-medium">
                        {fileInfo?.allow_download ? 'Download started!' : 'File opened in a new tab.'}
                    </p>
                </div>
            )}

            {/* Download error */}
            {downloadError && (
                <div className="flex items-start gap-2.5 p-3 rounded-lg bg-error/10 border border-error/20 text-xs">
                    <svg className="w-4 h-4 text-error flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                            d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    <p className="text-error">{downloadError}</p>
                </div>
            )}

            {/* CTA button */}
            <button
                onClick={handleDownload}
                disabled={downloading}
                className="w-full flex items-center justify-center gap-2 py-3 px-4 rounded-full bg-primary text-on-primary text-sm font-semibold transition-all hover:brightness-110 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
            >
                {downloading ? (
                    <>
                        <span className="w-4 h-4 border-2 border-on-primary border-t-transparent rounded-full animate-spin" />
                        Preparing…
                    </>
                ) : fileInfo?.allow_download ? (
                    <>
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                        </svg>
                        Download File
                    </>
                ) : (
                    <>
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                        </svg>
                        View File
                    </>
                )}
            </button>

            {fileInfo?.allow_download === false && (
                <p className="text-center text-[10px] text-surface-on-variant">
                    The sender has disabled downloads. The file will open in your browser.
                </p>
            )}
        </div>
    );

    // ── Page layout ───────────────────────────────────────────────────────────

    return (
        <div className="min-h-screen bg-surface flex flex-col">
            {/* Minimal top bar */}
            <header className="h-12 flex items-center justify-between px-5 border-b border-outline-variant/30 bg-surface-container-low flex-shrink-0">
                <Link to="/" className="flex items-center gap-2 group">
                    <div className="w-6 h-6 rounded-md bg-primary flex items-center justify-center">
                        <svg className="w-3.5 h-3.5 text-on-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5}
                                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                        </svg>
                    </div>
                    <span className="text-sm font-semibold text-surface-on group-hover:text-primary transition-colors">
                        File Share
                    </span>
                </Link>

                <button
                    onClick={toggleTheme}
                    className="p-1.5 rounded-full hover:bg-surface-variant/30 text-surface-on-variant transition-all"
                    title={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
                >
                    {theme === 'light' ? (
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                        </svg>
                    ) : (
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M12 3v1m0 16v1m9-9h1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                        </svg>
                    )}
                </button>
            </header>

            {/* Content */}
            <main className="flex-1 flex items-center justify-center p-4">
                <div className="w-full max-w-sm">
                    {/* Card */}
                    <div className="bg-surface-container rounded-2xl border border-outline-variant/30 shadow-lg overflow-hidden">
                        {/* Card header */}
                        <div className="px-5 py-4 border-b border-outline-variant/20 bg-surface-container-low">
                            <p className="text-[10px] font-bold text-surface-on-variant uppercase tracking-widest">
                                Shared File
                            </p>
                        </div>

                        {/* Card body */}
                        <div className="p-5">
                            {pageState === 'loading' && <Skeleton />}
                            {pageState === 'error' && renderError()}
                            {pageState === 'requires_password' && renderPasswordForm()}
                            {pageState === 'ready' && renderFileInfo()}
                        </div>
                    </div>

                    {/* Footer note */}
                    <p className="text-center text-[10px] text-surface-on-variant/50 mt-4">
                        Powered by File Share · Secure file sharing
                    </p>
                </div>
            </main>
        </div>
    );
};

export default SharePage;
