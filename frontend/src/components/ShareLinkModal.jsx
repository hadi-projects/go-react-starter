import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import Modal from './Modal';
import Button from './Button';
import TextField from './TextField';
import { toast } from 'react-hot-toast';
import { createShareLink, updateShareLink, revokeShareLink, getShareLinkLogs } from '../api/storage';

const ACCESS_TYPES = [
    {
        value: 'unlimited',
        label: 'Always accessible',
        desc: 'Link stays active until you manually revoke it.',
        icon: (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064" />
            </svg>
        ),
    },
    {
        value: 'one_time',
        label: 'One-time view',
        desc: 'Link is automatically deactivated after the first access.',
        icon: (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
        ),
    },
    {
        value: 'limited',
        label: 'Limited views',
        desc: 'Link deactivates after a set number of accesses.',
        icon: (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
            </svg>
        ),
    },
    {
        value: 'timed',
        label: 'Time-based expiry',
        desc: 'Link expires at a specific date and time.',
        icon: (
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        ),
    },
];

const ShareLinkModal = ({ isOpen, onClose, fileId, existingLink, onRefresh }) => {
    const isEdit = !!existingLink;

    const [form, setForm] = useState({
        label: '',
        access_type: 'unlimited',
        max_views: 10,
        expires_at: '',
        password: '',
        allow_download: true,
    });

    const [loading, setLoading] = useState(false);
    const [created, setCreated] = useState(null);
    const [tab, setTab] = useState('form');
    const [logs, setLogs] = useState([]);
    const [logsLoading, setLogsLoading] = useState(false);
    const [copied, setCopied] = useState(false);

    useEffect(() => {
        if (!isOpen) return;
        if (existingLink) {
            setForm({
                label: existingLink.label || '',
                access_type: existingLink.access_type || 'unlimited',
                max_views: existingLink.max_views || 10,
                expires_at: existingLink.expires_at
                    ? new Date(existingLink.expires_at).toISOString().slice(0, 16)
                    : '',
                password: '',
                allow_download: existingLink.allow_download ?? true,
            });
        } else {
            setForm({
                label: '',
                access_type: 'unlimited',
                max_views: 10,
                expires_at: '',
                password: '',
                allow_download: true,
            });
            setCreated(null);
        }
        setTab('form');
        setLogs([]);
    }, [existingLink, isOpen]);

    const fetchLogs = async () => {
        if (!existingLink) return;
        setLogsLoading(true);
        try {
            const res = await getShareLinkLogs(existingLink.id);
            setLogs(res.data?.data || []);
        } catch {
            toast.error('Failed to load access logs');
        } finally {
            setLogsLoading(false);
        }
    };

    const handleTabChange = (t) => {
        setTab(t);
        if (t === 'logs') fetchLogs();
    };

    const buildPayload = () => {
        const payload = {
            label: form.label,
            access_type: form.access_type,
            allow_download: form.allow_download,
        };
        if (form.password) payload.password = form.password;
        if (form.access_type === 'limited') payload.max_views = Number(form.max_views);
        if (form.access_type === 'timed' && form.expires_at)
            payload.expires_at = new Date(form.expires_at).toISOString();
        return payload;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            if (isEdit) {
                const payload = buildPayload();
                // allow clearing password (empty string = remove)
                if (form.password !== undefined) payload.password = form.password;
                await updateShareLink(existingLink.id, payload);
                toast.success('Share link updated');
                onRefresh();
                onClose();
            } else {
                const res = await createShareLink(fileId, buildPayload());
                const newLink = res.data?.data;
                setCreated(newLink);
                toast.success('Share link created!');
                onRefresh();
            }
        } catch (err) {
            toast.error(err.response?.data?.meta?.message || 'Operation failed');
        } finally {
            setLoading(false);
        }
    };

    const handleRevoke = async () => {
        if (!window.confirm('Are you sure you want to revoke this share link? It will no longer be accessible.')) return;
        try {
            await revokeShareLink(existingLink.id);
            toast.success('Share link revoked');
            onRefresh();
            onClose();
        } catch {
            toast.error('Failed to revoke link');
        }
    };

    const handleCopy = (url) => {
        navigator.clipboard.writeText(url);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={isEdit ? 'Edit Share Link' : 'Create Share Link'}
            maxWidth="max-w-lg"
        >
            {/* ── Created success state ─────────────────────────────── */}
            {created && !isEdit && (
                <div className="space-y-4 pt-2">
                    <div className="flex items-center gap-2.5 p-3 bg-green-500/10 rounded-lg border border-green-500/20">
                        <svg className="w-5 h-5 text-green-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                        <p className="text-xs font-semibold text-surface-on">Share link created successfully!</p>
                    </div>

                    <div>
                        <p className="text-xs font-medium text-surface-on-variant mb-1.5">Share URL</p>
                        <div className="flex gap-2">
                            <input
                                readOnly
                                value={created.share_url}
                                className="text-field flex-1 text-xs"
                            />
                            <button
                                onClick={() => handleCopy(created.share_url)}
                                className={`px-3 py-1.5 rounded-lg text-xs font-semibold whitespace-nowrap transition-all
                                    ${copied
                                        ? 'bg-green-500 text-white'
                                        : 'bg-primary text-on-primary hover:brightness-110'}`}
                            >
                                {copied ? '✓ Copied!' : 'Copy'}
                            </button>
                        </div>
                    </div>

                    <div className="rounded-lg border border-outline-variant/30 p-3 space-y-1.5 text-xs">
                        <div className="flex justify-between">
                            <span className="text-surface-on-variant">Type</span>
                            <span className="font-medium text-surface-on capitalize">
                                {ACCESS_TYPES.find(a => a.value === created.access_type)?.label}
                            </span>
                        </div>
                        {created.max_views != null && (
                            <div className="flex justify-between">
                                <span className="text-surface-on-variant">Max views</span>
                                <span className="font-medium text-surface-on">{created.max_views}</span>
                            </div>
                        )}
                        {created.expires_at && (
                            <div className="flex justify-between">
                                <span className="text-surface-on-variant">Expires</span>
                                <span className="font-medium text-surface-on">
                                    {new Date(created.expires_at).toLocaleString()}
                                </span>
                            </div>
                        )}
                        <div className="flex justify-between">
                            <span className="text-surface-on-variant">Download</span>
                            <span className={`font-medium ${created.allow_download ? 'text-green-500' : 'text-amber-500'}`}>
                                {created.allow_download ? 'Allowed' : 'View only'}
                            </span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-surface-on-variant">Password</span>
                            <span className={`font-medium ${created.has_password ? 'text-primary' : 'text-surface-on-variant'}`}>
                                {created.has_password ? 'Protected' : 'None'}
                            </span>
                        </div>
                    </div>

                    <div className="flex justify-end gap-2 pt-1">
                        <Button variant="tonal" onClick={() => setCreated(null)}>
                            Create Another
                        </Button>
                        <Button variant="primary" onClick={onClose}>
                            Done
                        </Button>
                    </div>
                </div>
            )}

            {/* ── Form / edit view ──────────────────────────────────── */}
            {(!created || isEdit) && (
                <div className="pt-1">
                    {/* Tabs — only in edit mode */}
                    {isEdit && (
                        <div className="flex border-b border-outline-variant/30 mb-4 -mx-1">
                            {[
                                { key: 'form', label: 'Settings' },
                                { key: 'logs', label: 'Access Logs' },
                            ].map(({ key, label }) => (
                                <button
                                    key={key}
                                    onClick={() => handleTabChange(key)}
                                    className={`px-4 py-2 text-xs font-semibold border-b-2 transition-colors -mb-px
                                        ${tab === key
                                            ? 'border-primary text-primary'
                                            : 'border-transparent text-surface-on-variant hover:text-surface-on'}`}
                                >
                                    {label}
                                </button>
                            ))}
                        </div>
                    )}

                    {/* ── Settings tab ────────────────────────────── */}
                    {tab === 'form' && (
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <TextField
                                label="Label (optional)"
                                name="label"
                                value={form.label}
                                onChange={(e) => setForm({ ...form, label: e.target.value })}
                                placeholder="e.g. Client A access"
                            />

                            {/* Access type selector */}
                            <div>
                                <p className="text-field-label">Access Type</p>
                                <div className="grid grid-cols-1 gap-2">
                                    {ACCESS_TYPES.map((at) => (
                                        <label
                                            key={at.value}
                                            className={`flex items-start gap-3 p-2.5 rounded-lg border cursor-pointer transition-all
                                                ${form.access_type === at.value
                                                    ? 'border-primary bg-primary/5'
                                                    : 'border-outline-variant/40 hover:border-outline-variant bg-surface-container'}`}
                                        >
                                            <input
                                                type="radio"
                                                name="access_type"
                                                value={at.value}
                                                checked={form.access_type === at.value}
                                                onChange={(e) => setForm({ ...form, access_type: e.target.value })}
                                                className="mt-0.5 accent-primary flex-shrink-0"
                                            />
                                            <div className={`mt-0.5 flex-shrink-0 ${form.access_type === at.value ? 'text-primary' : 'text-surface-on-variant'}`}>
                                                {at.icon}
                                            </div>
                                            <div className="min-w-0">
                                                <p className="text-xs font-semibold text-surface-on">{at.label}</p>
                                                <p className="text-[10px] text-surface-on-variant mt-0.5">{at.desc}</p>
                                            </div>
                                        </label>
                                    ))}
                                </div>
                            </div>

                            {/* Conditional: max views */}
                            {form.access_type === 'limited' && (
                                <TextField
                                    label="Max Views"
                                    name="max_views"
                                    type="number"
                                    value={String(form.max_views)}
                                    onChange={(e) => setForm({ ...form, max_views: e.target.value })}
                                    required
                                />
                            )}

                            {/* Conditional: expires at */}
                            {form.access_type === 'timed' && (
                                <div>
                                    <label className="text-field-label">Expires At</label>
                                    <input
                                        type="datetime-local"
                                        value={form.expires_at}
                                        onChange={(e) => setForm({ ...form, expires_at: e.target.value })}
                                        className="text-field"
                                        required
                                    />
                                </div>
                            )}

                            {/* Advanced options */}
                            <div className="border border-outline-variant/30 rounded-lg p-3 space-y-3 bg-surface-container">
                                <p className="text-[10px] font-bold text-surface-on-variant uppercase tracking-widest">
                                    Advanced Options
                                </p>

                                <label className="flex items-start gap-2.5 cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={form.allow_download}
                                        onChange={(e) => setForm({ ...form, allow_download: e.target.checked })}
                                        className="mt-0.5 accent-primary flex-shrink-0"
                                    />
                                    <div>
                                        <p className="text-xs font-semibold text-surface-on">Allow Download</p>
                                        <p className="text-[10px] text-surface-on-variant">
                                            If unchecked, recipients can only view the file inline (images, PDFs).
                                        </p>
                                    </div>
                                </label>

                                <TextField
                                    label={isEdit
                                        ? 'New Password (leave blank to keep current, empty string to remove)'
                                        : 'Password (optional)'}
                                    name="password"
                                    type="password"
                                    value={form.password}
                                    onChange={(e) => setForm({ ...form, password: e.target.value })}
                                    placeholder="Leave blank for no password"
                                />
                            </div>

                            {/* Actions */}
                            <div className="flex items-center pt-2 gap-3">
                                {isEdit && (
                                    <button
                                        type="button"
                                        onClick={handleRevoke}
                                        className="text-xs text-error hover:underline font-semibold"
                                    >
                                        Revoke Link
                                    </button>
                                )}
                                <div className="flex gap-2 ml-auto">
                                    <Button type="button" variant="tonal" onClick={onClose}>
                                        Cancel
                                    </Button>
                                    <Button type="submit" variant="primary" disabled={loading}>
                                        {loading
                                            ? 'Saving…'
                                            : isEdit
                                                ? 'Update Link'
                                                : 'Create Link'}
                                    </Button>
                                </div>
                            </div>
                        </form>
                    )}

                    {/* ── Access logs tab ──────────────────────────── */}
                    {tab === 'logs' && (
                        <div className="space-y-2">
                            {logsLoading ? (
                                <div className="text-center py-10">
                                    <div className="inline-block w-5 h-5 border-2 border-primary border-t-transparent rounded-full animate-spin" />
                                    <p className="text-xs text-surface-on-variant mt-2">Loading logs…</p>
                                </div>
                            ) : logs.length === 0 ? (
                                <div className="text-center py-10">
                                    <svg className="w-8 h-8 mx-auto text-surface-on-variant/30 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5}
                                            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                                    </svg>
                                    <p className="text-xs text-surface-on-variant">No accesses recorded yet.</p>
                                </div>
                            ) : (
                                <div className="space-y-1 max-h-72 overflow-y-auto custom-scrollbar pr-1">
                                    {logs.map((log) => (
                                        <div
                                            key={log.id}
                                            className="flex items-start gap-2.5 p-2.5 rounded-lg bg-surface-container border border-outline-variant/20"
                                        >
                                            <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                                            <div className="flex-1 min-w-0">
                                                <p className="text-[10px] font-mono text-surface-on truncate">
                                                    {log.ip_address || 'Unknown IP'}
                                                </p>
                                                <p className="text-[10px] text-surface-on-variant truncate">
                                                    {log.user_agent || 'Unknown agent'}
                                                </p>
                                            </div>
                                            <p className="text-[10px] text-surface-on-variant whitespace-nowrap flex-shrink-0">
                                                {new Date(log.accessed_at).toLocaleString()}
                                            </p>
                                        </div>
                                    ))}
                                </div>
                            )}

                            <div className="flex justify-end pt-2">
                                <Button type="button" variant="tonal" onClick={onClose}>Close</Button>
                            </div>
                        </div>
                    )}
                </div>
            )}
        </Modal>
    );
};

ShareLinkModal.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    fileId: PropTypes.number,
    existingLink: PropTypes.object,
    onRefresh: PropTypes.func.isRequired,
};

export default ShareLinkModal;
