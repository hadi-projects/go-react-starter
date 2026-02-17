import { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import PropTypes from 'prop-types';

const Sidebar = ({ sections = [], title = "Admin Panel", onLogout }) => {
    const location = useLocation();
    const [expandedSections, setExpandedSections] = useState({});

    const isActive = (path) => location.pathname === path;
    const isChildActive = (item) => {
        if (!item.subItems) return isActive(item.path);
        return item.subItems.some(sub => isActive(sub.path));
    };

    // Auto-expand section if a child is active
    useEffect(() => {
        const newExpanded = { ...expandedSections };
        let changed = false;

        sections.forEach(section => {
            section.items.forEach(item => {
                if (item.subItems && isChildActive(item) && !expandedSections[item.label]) {
                    newExpanded[item.label] = true;
                    changed = true;
                }
            });
        });

        if (changed) {
            setExpandedSections(newExpanded);
        }
    }, [location.pathname, sections]);

    const toggleSection = (label) => {
        setExpandedSections(prev => ({
            ...prev,
            [label]: !prev[label]
        }));
    };

    return (
        <aside className="nav-drawer flex-shrink-0 animate-fade-in">
            {/* Header */}
            <div className="h-16 flex items-center px-6">
                <h2 className="text-xl font-semibold text-surface-on">{title}</h2>
            </div>

            {/* Navigation Drawer Content */}
            <nav className="flex-1 overflow-y-auto py-2 custom-scrollbar">
                {sections.map((section, idx) => (
                    <div key={idx} className="mb-4">
                        {section.label && (
                            <h3 className="px-6 py-3 text-xs font-semibold text-surface-on-variant uppercase tracking-wider">
                                {section.label}
                            </h3>
                        )}
                        <ul className="space-y-1 px-3">
                            {section.items.map((item) => {
                                const active = isChildActive(item);
                                const hasSubItems = !!item.subItems;
                                const isExpanded = expandedSections[item.label];

                                return (
                                    <li key={item.label}>
                                        {hasSubItems ? (
                                            <div>
                                                <button
                                                    onClick={() => toggleSection(item.label)}
                                                    className={`w-full flex items-center justify-between px-4 py-3 rounded-full transition-all duration-300 group ${active
                                                        ? 'bg-secondary-container/50 text-secondary-on-container'
                                                        : 'text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on'
                                                        }`}
                                                >
                                                    <div className="flex items-center gap-3">
                                                        <div className="transition-transform duration-300 group-hover:scale-110">
                                                            {item.icon}
                                                        </div>
                                                        <span className="font-medium text-sm">{item.label}</span>
                                                    </div>
                                                    <svg
                                                        className={`w-4 h-4 transition-transform duration-300 ${isExpanded ? 'rotate-180' : ''}`}
                                                        fill="none"
                                                        stroke="currentColor"
                                                        viewBox="0 0 24 24"
                                                    >
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                                                    </svg>
                                                </button>

                                                {/* Sub Items */}
                                                <div className={`overflow-hidden transition-all duration-300 ease-in-out ${isExpanded ? 'max-h-64 mt-1' : 'max-h-0'}`}>
                                                    <ul className="space-y-1 ml-4 pl-4 border-l border-outline-variant/30">
                                                        {item.subItems.map((sub) => {
                                                            const subActive = isActive(sub.path);
                                                            return (
                                                                <li key={sub.path}>
                                                                    <Link
                                                                        to={sub.path}
                                                                        className={`flex items-center gap-3 px-4 py-2.5 rounded-full transition-all duration-300 group ${subActive
                                                                            ? 'bg-secondary-container text-secondary-on-container shadow-sm'
                                                                            : 'text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on'
                                                                            }`}
                                                                    >
                                                                        <div className="transition-transform duration-300 group-hover:scale-110">
                                                                            {sub.icon}
                                                                        </div>
                                                                        <span className="font-medium text-sm">{sub.label}</span>
                                                                    </Link>
                                                                </li>
                                                            );
                                                        })}
                                                    </ul>
                                                </div>
                                            </div>
                                        ) : (
                                            <Link
                                                to={item.path}
                                                className={`flex items-center gap-3 px-4 py-3 rounded-full transition-all duration-300 group ${active
                                                    ? 'bg-secondary-container text-secondary-on-container'
                                                    : 'text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on'
                                                    }`}
                                            >
                                                <div className="transition-transform duration-300 group-hover:scale-110">
                                                    {item.icon}
                                                </div>
                                                <span className="font-medium text-sm">{item.label}</span>
                                            </Link>
                                        )}
                                    </li>
                                );
                            })}
                        </ul>
                    </div>
                ))}
            </nav>

            {/* Bottom Actions */}
            <div className="p-4 border-t border-outline-variant/30">
                <button
                    onClick={onLogout}
                    className="flex items-center gap-3 w-full px-4 py-3 rounded-full text-surface-on-variant hover:bg-surface-container-high hover:text-surface-on transition-all duration-300 group"
                >
                    <div className="transition-transform duration-300 group-hover:scale-110">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                    </div>
                    <span className="font-medium text-sm">Logout</span>
                </button>
            </div>
        </aside>
    );
};

Sidebar.propTypes = {
    sections: PropTypes.arrayOf(PropTypes.shape({
        label: PropTypes.string,
        items: PropTypes.arrayOf(PropTypes.shape({
            path: PropTypes.string,
            label: PropTypes.string.isRequired,
            icon: PropTypes.node.isRequired,
            subItems: PropTypes.arrayOf(PropTypes.shape({
                path: PropTypes.string.isRequired,
                label: PropTypes.string.isRequired,
                icon: PropTypes.node.isRequired,
            }))
        })).isRequired,
    })),
    title: PropTypes.string,
    onLogout: PropTypes.func.isRequired,
};

export default Sidebar;
