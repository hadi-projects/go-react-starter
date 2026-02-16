/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // MD3 Purple Theme - Expressive
        primary: {
          50: '#F3E5F5',
          100: '#E1BEE7',
          200: '#CE93D8',
          300: '#BA68C8',
          400: '#AB47BC',
          500: '#6750A4', // Main primary
          600: '#5846A0',
          700: '#4A378B',
          800: '#3C2976',
          900: '#2E1A61',
          container: '#EADDFF',
          'on-container': '#21005D',
        },
        secondary: {
          50: '#F5F3FF',
          100: '#E8DEF8',
          200: '#D0BCFF',
          300: '#B89EFF',
          400: '#A080FF',
          500: '#625B71',
          600: '#4A4458',
          700: '#332D41',
          800: '#1D192B',
          900: '#000000',
          container: '#E8DEF8',
          'on-container': '#1D192B',
        },
        surface: {
          DEFAULT: '#FEF7FF',
          variant: '#E7E0EC',
          dim: '#DED8E1',
          bright: '#FFFFFF',
          'on': '#1D1B20',
          'on-variant': '#49454F',
          'container-lowest': '#FFFFFF',
          'container-low': '#F7F2FA',
          'container': '#F3EDF7',
          'container-high': '#ECE6F0',
          'container-highest': '#E6E0E9',
        },
        'on-primary': '#FFFFFF',
        'on-secondary': '#FFFFFF',
        outline: {
          DEFAULT: '#79747E',
          variant: '#CAC4D0',
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
