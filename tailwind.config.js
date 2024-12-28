/** @type {import('tailwindcss').Config} */
export default {
  // corePlugins: {
  //   preflight: false,
  // },
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'midnight': '#0f172a',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
}