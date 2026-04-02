import { createContext, useContext, useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { getSettingsByCategory } from '../api/settings';

const SettingsContext = createContext();

export const SettingsProvider = ({ children }) => {
    const [settings, setSettings] = useState({
        app_name: import.meta.env.VITE_APP_NAME || 'Go-React Starter',
        logo: null,
        favicon: null,
    });

    const refreshSettings = async () => {
        try {
            const res = await getSettingsByCategory('website');
            const websiteSettings = res.data.data;

            const newSettings = { ...settings };
            websiteSettings.forEach(s => {
                if (s.key === 'app_name' && s.value) newSettings.app_name = s.value;
                if (s.key === 'logo' && s.value) newSettings.logo = s.value;
                if (s.key === 'favicon' && s.value) newSettings.favicon = s.value;
            });

            setSettings(newSettings);

            // Update document title and favicon
            if (newSettings.app_name) {
                // We'll let individual pages handle specific titles, 
                // but we can set a base one here if needed.
            }

            if (newSettings.favicon) {
                const link = document.querySelector("link[rel*='icon']") || document.createElement('link');
                link.type = 'image/x-icon';
                link.rel = 'shortcut icon';
                link.href = `${import.meta.env.VITE_API_URL}/public/share/${newSettings.favicon}`;
                document.getElementsByTagName('head')[0].appendChild(link);
            }
        } catch (err) {
            console.error('Failed to fetch settings:', err);
        }
    };

    useEffect(() => {
        refreshSettings();
    }, []);

    return (
        <SettingsContext.Provider value={{ ...settings, refreshSettings }}>
            {children}
        </SettingsContext.Provider>
    );
};

SettingsProvider.propTypes = {
    children: PropTypes.node.isRequired,
};

export const useSettings = () => {
    const context = useContext(SettingsContext);
    if (!context) {
        throw new Error('useSettings must be used within a SettingsProvider');
    }
    return context;
};
