/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // MD3 Purple Theme - Expressive (Theme-aware via CSS Variables with Opacity Support)
        primary: {
          DEFAULT: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-primary) / ${opacityValue})` : `rgb(var(--md-sys-color-primary))`,
          500: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-primary) / ${opacityValue})` : `rgb(var(--md-sys-color-primary))`,
          container: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-primary-container) / ${opacityValue})` : `rgb(var(--md-sys-color-primary-container))`,
          'on-container': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-primary-container) / ${opacityValue})` : `rgb(var(--md-sys-color-on-primary-container))`,
        },
        secondary: {
          DEFAULT: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-secondary-container) / ${opacityValue})` : `rgb(var(--md-sys-color-secondary-container))`,
          container: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-secondary-container) / ${opacityValue})` : `rgb(var(--md-sys-color-secondary-container))`,
          'on-container': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-secondary-container) / ${opacityValue})` : `rgb(var(--md-sys-color-on-secondary-container))`,
        },
        surface: {
          DEFAULT: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-surface) / ${opacityValue})` : `rgb(var(--md-sys-color-surface))`,
          variant: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-surface-variant) / ${opacityValue})` : `rgb(var(--md-sys-color-surface-variant))`,
          'on': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-surface) / ${opacityValue})` : `rgb(var(--md-sys-color-on-surface))`,
          'on-variant': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-surface-variant) / ${opacityValue})` : `rgb(var(--md-sys-color-on-surface-variant))`,
          'container-low': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-surface-container-low) / ${opacityValue})` : `rgb(var(--md-sys-color-surface-container-low))`,
          'container': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-surface-container) / ${opacityValue})` : `rgb(var(--md-sys-color-surface-container))`,
          'container-high': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-surface-container-high) / ${opacityValue})` : `rgb(var(--md-sys-color-surface-container-high))`,
          'container-highest': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-surface-container-highest) / ${opacityValue})` : `rgb(var(--md-sys-color-surface-container-highest))`,
        },
        'on-primary': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-primary) / ${opacityValue})` : `rgb(var(--md-sys-color-on-primary))`,
        'on-secondary': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-secondary) / ${opacityValue})` : `rgb(var(--md-sys-color-on-secondary))`,
        outline: {
          DEFAULT: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-outline) / ${opacityValue})` : `rgb(var(--md-sys-color-outline))`,
          variant: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-outline-variant) / ${opacityValue})` : `rgb(var(--md-sys-color-outline-variant))`,
        },
        error: {
          DEFAULT: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-error) / ${opacityValue})` : `rgb(var(--md-sys-color-error))`,
          'on': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-error) / ${opacityValue})` : `rgb(var(--md-sys-color-on-error))`,
          container: ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-error-container) / ${opacityValue})` : `rgb(var(--md-sys-color-error-container))`,
          'on-container': ({ opacityValue }) => opacityValue !== undefined ? `rgb(var(--md-sys-color-on-error-container) / ${opacityValue})` : `rgb(var(--md-sys-color-on-error-container))`,
        },
      },
      fontSize: {
        // MD3 Typography Scale - Refined (smaller)
        'xs': ['0.75rem', { lineHeight: '1rem' }],      // 12px
        'sm': ['0.875rem', { lineHeight: '1.25rem' }],  // 14px
        'base': ['0.9375rem', { lineHeight: '1.375rem' }], // 15px
        'lg': ['1rem', { lineHeight: '1.5rem' }],       // 16px
        'xl': ['1.125rem', { lineHeight: '1.75rem' }],  // 18px
        '2xl': ['1.375rem', { lineHeight: '2rem' }],    // 22px
        '3xl': ['1.75rem', { lineHeight: '2.25rem' }],  // 28px
      },
      borderRadius: {
        'md3-sm': '8px',
        'md3': '12px',
        'md3-lg': '16px',
        'md3-xl': '28px',
      },
      boxShadow: {
        'md3-1': '0px 1px 2px 0px rgba(0, 0, 0, 0.3), 0px 1px 3px 1px rgba(0, 0, 0, 0.15)',
        'md3-2': '0px 1px 2px 0px rgba(0, 0, 0, 0.3), 0px 2px 6px 2px rgba(0, 0, 0, 0.15)',
        'md3-3': '0px 4px 8px 3px rgba(0, 0, 0, 0.15), 0px 1px 3px 0px rgba(0, 0, 0, 0.3)',
        'md3-4': '0px 6px 10px 4px rgba(0, 0, 0, 0.15), 0px 2px 3px 0px rgba(0, 0, 0, 0.3)',
        'md3-5': '0px 8px 12px 6px rgba(0, 0, 0, 0.15), 0px 4px 4px 0px rgba(0, 0, 0, 0.3)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
