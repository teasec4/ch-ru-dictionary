/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  theme: {
    extend: {
      colors: {
        brand: {
          primary: "#3B82F6",
          secondary: "#8B5CF6",
        },
        ui: {
          background: "#509296",
          surface: "#FFFFFF",
        },
        text: {
          primary: "#111827",
          secondary: "#6B7280",
        },
      },
    },
  },
  plugins: [],
};