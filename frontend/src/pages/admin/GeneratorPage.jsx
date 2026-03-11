import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import Button from '../../components/Button';
import { generateModule } from '../../api/admin';

const GeneratorPage = () => {
    const navigate = useNavigate();
    const [config, setConfig] = useState({
        module_name: '',
        table_name: '',
        audit_log: true,
        fields: [
            { name: 'name', type: 'string', binding: 'required', searchable: true, unique: false }
        ]
    });

    const [loading, setLoading] = useState(false);

    const handleFieldChange = (index, field, value) => {
        const newFields = [...config.fields];
        newFields[index][field] = value;
        setConfig({ ...config, fields: newFields });
    };

    const addField = () => {
        setConfig({
            ...config,
            fields: [...config.fields, { name: '', type: 'string', binding: '', searchable: false, unique: false }]
        });
    };

    const removeField = (index) => {
        const newFields = config.fields.filter((_, i) => i !== index);
        setConfig({ ...config, fields: newFields });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            await generateModule(config);
            toast.success('Module generated successfully! Please restart the backend if needed.');

            // Redirect to the new module's page
            // Convert to lowercase for the URL path
            const routePath = `/admin/${config.module_name.toLowerCase()}`;
            navigate(routePath);

            // Note: A full page reload might be necessary if the frontend needs to fetch new routes,
            // but we'll try client-side navigation first. If backend needs restart, this will fail gracefully.
        } catch (error) {
            toast.error(error.response?.data?.meta?.message || 'Failed to generate module');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="space-y-8">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold text-surface-on tracking-tight">Module Generator</h1>
                    <p className="text-surface-on-variant mt-2">Scaffold a new backend module instantly.</p>
                </div>
            </div>

            <form onSubmit={handleSubmit} className="bg-surface-container rounded-[28px] p-8 border border-outline-variant/30 space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-surface-on px-1">Module Name (e.g. Product)</label>
                        <input
                            type="text"
                            value={config.module_name}
                            onChange={(e) => setConfig({ ...config, module_name: e.target.value })}
                            className="w-full bg-surface-container-high border border-outline rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary/50 transition-all text-surface-on"
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-surface-on px-1">Table Name (e.g. products)</label>
                        <input
                            type="text"
                            value={config.table_name}
                            onChange={(e) => setConfig({ ...config, table_name: e.target.value })}
                            className="w-full bg-surface-container-high border border-outline rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all"
                            required
                        />
                    </div>
                </div>

                <div className="flex items-center gap-2 px-1">
                    <input
                        type="checkbox"
                        id="audit_log"
                        checked={config.audit_log}
                        onChange={(e) => setConfig({ ...config, audit_log: e.target.checked })}
                        className="w-5 h-5 rounded border-outline text-primary focus:ring-primary"
                    />
                    <label htmlFor="audit_log" className="text-sm font-medium text-surface-on">Enable Audit Logging</label>
                </div>

                <div className="space-y-4">
                    <div className="flex justify-between items-center px-1">
                        <h3 className="text-lg font-medium text-surface-on">Fields</h3>
                        <Button type="button" variant="tonal" size="sm" onClick={addField}>
                            Add Field
                        </Button>
                    </div>

                    <div className="space-y-4">
                        {config.fields.map((field, index) => (
                            <div key={index} className="grid grid-cols-1 md:grid-cols-6 gap-4 p-4 bg-surface-container-low rounded-2xl border border-outline-variant/20 relative animate-fade-in">
                                <div className="md:col-span-2 space-y-1">
                                    <label className="text-[11px] font-bold uppercase text-surface-on-variant px-1">Name</label>
                                    <input
                                        type="text"
                                        value={field.name}
                                        onChange={(e) => handleFieldChange(index, 'name', e.target.value)}
                                        className="w-full bg-surface-container border border-outline rounded-lg px-3 py-2 text-sm focus:outline-none text-surface-on"
                                        placeholder="field_name"
                                        required
                                    />
                                </div>
                                <div className="space-y-1">
                                    <label className="text-[11px] font-bold uppercase text-surface-on-variant px-1">Type</label>
                                    <select
                                        value={field.type}
                                        onChange={(e) => handleFieldChange(index, 'type', e.target.value)}
                                        className="w-full bg-surface-container border border-outline rounded-lg px-3 py-2 text-sm focus:outline-none text-surface-on"
                                    >
                                        <option value="string">String</option>
                                        <option value="int">Integer</option>
                                        <option value="float">Float</option>
                                        <option value="boolean">Boolean</option>
                                        <option value="datetime">Datetime</option>
                                        <option value="wysiwyg">Rich Text</option>
                                    </select>
                                </div>
                                <div className="space-y-1">
                                    <label className="text-[11px] font-bold uppercase text-surface-on-variant px-1">Binding</label>
                                    <input
                                        type="text"
                                        value={field.binding}
                                        onChange={(e) => handleFieldChange(index, 'binding', e.target.value)}
                                        className="w-full bg-surface-container border border-outline rounded-lg px-3 py-2 text-sm focus:outline-none text-surface-on"
                                        placeholder="required"
                                    />
                                </div>
                                <div className="flex items-end justify-center pb-2 gap-4">
                                    <label className="flex items-center gap-2 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            checked={field.searchable}
                                            onChange={(e) => handleFieldChange(index, 'searchable', e.target.checked)}
                                            className="w-4 h-4 rounded border-outline text-primary"
                                        />
                                        <span className="text-[11px] font-medium text-surface-on">Search</span>
                                    </label>
                                    <label className="flex items-center gap-2 cursor-pointer">
                                        <input
                                            type="checkbox"
                                            checked={field.unique}
                                            onChange={(e) => handleFieldChange(index, 'unique', e.target.checked)}
                                            className="w-4 h-4 rounded border-outline text-error"
                                        />
                                        <span className="text-[11px] font-medium text-surface-on">Unique</span>
                                    </label>
                                </div>
                                <div className="flex items-end justify-end pb-1">
                                    <button
                                        type="button"
                                        onClick={() => removeField(index)}
                                        className="p-2 text-error hover:bg-error-container/20 rounded-full transition-colors"
                                        disabled={config.fields.length === 1}
                                    >
                                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="pt-6 flex justify-end">
                    <Button type="submit" variant="primary" size="lg" disabled={loading}>
                        {loading ? 'Generating...' : 'Generate Module'}
                    </Button>
                </div>
            </form>
        </div>
    );
};

export default GeneratorPage;
