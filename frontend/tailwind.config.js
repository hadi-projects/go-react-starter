/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // MD3 Purple Theme
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
        },
        secondary: {
          50: '#F5F3FF',
          100: '#E8DEF8',
          200: '#D0BCFF',
          300: '#B89EFF',
          400: '#A080FF',
          500: '#885DFF',
          600: '#7043E6',
          700: '#5929CC',
          800: '#4210B3',
          900: '#2B0099',
        },
        surface: {
          DEFAULT: '#FFFBFE',
          variant: '#E7E0EC',
          dim: '#DED8E1',
          bright: '#FFFFFF',
        },
        outline: {
          DEFAULT: '#79747E',
          variant: '#CAC4D0',
        },
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
