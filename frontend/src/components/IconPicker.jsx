import React, { useState } from 'react';
import { HERO_ICONS } from './heroIcons';

const IconPicker = ({ selectedIcon, onSelect }) => {
    const [searchTerm, setSearchTerm] = useState('');

    const filteredIcons = HERO_ICONS.filter(icon => 
        icon.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div className="space-y-4 p-4 bg-surface-container-high rounded-2xl border border-outline-variant/30">
            <div className="flex justify-between items-center">
                <label className="text-sm font-medium text-surface-on">Sidebar Icon</label>
                <input
                    type="text"
                    placeholder="Search icons..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="text-xs bg-surface-container border border-outline rounded-lg px-3 py-1.5 focus:outline-none"
                />
            </div>
            
            <div className="grid grid-cols-6 sm:grid-cols-8 md:grid-cols-10 gap-2 max-h-48 overflow-y-auto p-1 custom-scrollbar">
                {filteredIcons.map((icon) => (
                    <button
                        key={icon.name}
                        onClick={() => onSelect(icon.path)}
                        type="button"
                        className={`p-2 rounded-xl transition-all flex items-center justify-center border-2 ${
                            selectedIcon === icon.path 
                            ? 'border-primary bg-primary/10 text-primary' 
                            : 'border-transparent hover:bg-surface-variant/30 text-surface-on-variant'
                        }`}
                        title={icon.name}
                    >
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d={icon.path} />
                        </svg>
                    </button>
                ))}
            </div>
        </div>
    );
};

export default IconPicker;
