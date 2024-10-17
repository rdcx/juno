/** @type {import('tailwindcss').Config} */
export default {
  content: [],
  darkMode: 'class',
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      maxWidth: {
        '9xl': '1440px',
        '10xl': '1920px',
      },
    },
  },
  plugins: [],
}

