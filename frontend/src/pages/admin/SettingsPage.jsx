import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import Card from '../../components/Card';
import Button from '../../components/Button';
import TextField from '../../components/TextField';
import FileUploadDropzone from '../../components/FileUploadDropzone';
import { useSettings } from '../../context/SettingsContext';
import { getSettingsByCategory, updateSettings } from '../../api/settings';
import { uploadFile } from '../../api/storage';
import toast from 'react-hot-toast';

const SettingsPage = () => {
    const { category = 'website' } = useParams();
    const queryClient = useQueryClient();
    const { refreshSettings } = useSettings();

    // Fetch settings for the category from URL
    const { data: settingsData, isLoading } = useQuery({
        queryKey: ['settings', category],
        queryFn: () => getSettingsByCategory(category),
    });

    const settings = settingsData?.data || [];

    // Mutation for updating settings
    const updateMutation = useMutation({
        mutationFn: updateSettings,
        onSuccess: () => {
            queryClient.invalidateQueries(['settings', category]);
            refreshSettings();
            toast.success('Settings updated successfully');
        },
        onError: (error) => {
            toast.error(error.response?.data?.meta?.message || 'Failed to update settings');
        },
    });

    const [formState, setFormState] = useState({});
    const [previews, setPreviews] = useState({});

    // Sync form state with fetched settings
    useEffect(() => {
        if (settings.length > 0) {
            const initialForm = {};
            settings.forEach(s => {
                initialForm[s.key] = s.value;
            });
            setFormState(initialForm);
        }
    }, [settings, category]);

    const handleInputChange = (key, value) => {
        setFormState(prev => ({ ...prev, [key]: value }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        const payload = {};
        settings.forEach(s => {
            payload[s.key] = formState[s.key] ?? s.value;
        });
        updateMutation.mutate(payload);
    };

    const handleFileUpload = async (key, file) => {
        // Create local preview immediately
        const localUrl = URL.createObjectURL(file);
        setPreviews(prev => ({ ...prev, [key]: localUrl }));

        const formData = new FormData();
        formData.append('file', file);
        formData.append('description', `System Setting: ${key}`);

        try {
            const res = await uploadFile(formData);
            // res.data is the axios response body { meta: ..., data: { id: ... } }
            const fileId = res.data.data.id;
            handleInputChange(key, String(fileId));
            toast.success('File uploaded. Save settings to apply.');
        } catch (err) {
            console.error('Upload error:', err);
            toast.error('Failed to upload file');
            setPreviews(prev => {
                const newState = { ...prev };
                delete newState[key];
                return newState;
            });
        }
    };

    const getCategoryLabel = () => {
        const labels = {
            website: 'Website Settings',
            smtp: 'Email (SMTP) Settings',
            storage: 'Storage Settings',
            security: 'Security Settings',
            internal: 'Infrastructure Settings',
            advance: 'Advanced Settings'
        };
        return labels[category] || 'Settings';
    };

    if (isLoading) return <div className="p-8 text-center text-surface-on-variant">Loading settings...</div>;

    return (
        <div className="p-4 sm:p-8 max-w-5xl mx-auto">
            <header className="mb-8 font-jakarta">
                <h1 className="text-3xl font-extrabold text-surface-on tracking-tight capitalize">
                    {getCategoryLabel()}
                </h1>
                <p className="text-surface-on-variant mt-1">Configure your application's global behavior and branding.</p>
            </header>

            <div className="flex flex-col gap-8">
                {/* Content Area */}
                <main className="flex-1">
                    <Card className="p-6 sm:p-8">
                        <form onSubmit={handleSubmit} className="space-y-8">
                            <div className="grid grid-cols-1 gap-y-8">
                                {settings.map(setting => (
                                    <div key={setting.key} className="space-y-2">
                                        <div className="flex flex-col">
                                            <span className="text-sm font-bold text-surface-on flex items-center gap-2">
                                                {setting.label}
                                                {setting.description.includes('⚠️') && (
                                                    <span className="text-[10px] bg-warning/20 text-warning px-1.5 py-0.5 rounded uppercase tracking-wider">Restart Required</span>
                                                )}
                                            </span>
                                            <span className="text-xs text-surface-on-variant">{setting.description}</span>
                                        </div>

                                        {setting.field_type === 'file' ? (
                                            <div className="mt-2 flex items-center gap-6">
                                                <div className="w-24 h-24 rounded-lg bg-surface-variant/50 flex items-center justify-center overflow-hidden border border-outline-variant">
                                                    {(previews[setting.key] || formState[setting.key]) ? (
                                                        <img 
                                                           src={previews[setting.key] || `${import.meta.env.VITE_API_URL}/public/storage/${formState[setting.key]}`}
                                                           alt="Preview" 
                                                           className="w-full h-full object-cover"
                                                           onError={(e) => { 
                                                               if (!previews[setting.key]) {
                                                                   e.target.src = 'https://placehold.co/100?text=No+Image'; 
                                                               }
                                                           }}
                                                        />
                                                    ) : (
                                                        <span className="text-2xl opacity-20">🖼️</span>
                                                    )}
                                                </div>
                                                <div className="flex-1">
                                                    <FileUploadDropzone 
                                                        onUpload={(file) => handleFileUpload(setting.key, file)}
                                                    />
                                                </div>
                                            </div>
                                        ) : setting.field_type === 'boolean' ? (
                                            <div className="flex items-center gap-3">
                                                <button
                                                    type="button"
                                                    onClick={() => handleInputChange(setting.key, String(formState[setting.key]) === 'true' ? 'false' : 'true')}
                                                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2
                                                        ${String(formState[setting.key]) === 'true' ? 'bg-primary' : 'bg-surface-variant'}`}
                                                >
                                                    <span
                                                        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform
                                                            ${String(formState[setting.key]) === 'true' ? 'translate-x-6' : 'translate-x-1'}`}
                                                    />
                                                </button>
                                                <span className="text-sm text-surface-on-variant">
                                                    {String(formState[setting.key]) === 'true' ? 'Enabled' : 'Disabled'}
                                                </span>
                                            </div>
                                        ) : (
                                            <TextField
                                                type={setting.field_type === 'password' ? 'password' : setting.field_type === 'number' ? 'number' : 'text'}
                                                value={formState[setting.key] ?? setting.value}
                                                onChange={(e) => handleInputChange(setting.key, e.target.value)}
                                                placeholder={setting.label}
                                                className="max-w-md"
                                            />
                                        )}
                                    </div>
                                ))}
                            </div>

                            <div className="pt-6 border-t border-outline-variant flex justify-end">
                                <Button 
                                    type="submit" 
                                    isLoading={updateMutation.isPending}
                                    className="px-8 shadow-lg shadow-primary/20"
                                >
                                    Save Changes
                                </Button>
                            </div>
                        </form>
                    </Card>
                    
                    {category === 'internal' && (
                        <div className="mt-6 p-4 bg-warning/10 border border-warning/20 rounded-xl flex gap-3 items-start text-warning text-sm">
                            <span className="text-lg">⚠️</span>
                            <p>
                                <strong>Warning:</strong> Changes to Database or Infrastructure settings require a manual restart of the backend service to take effect. Incorrect values may prevent the application from starting correctly.
                            </p>
                        </div>
                    )}
                </main>
            </div>
        </div>
    );
};

export default SettingsPage;
