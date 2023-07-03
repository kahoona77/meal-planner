import colors from "tailwindcss/colors.js";

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "../tmpl/**/*.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#086699",
        secondary: colors.yellow,
        neutral: colors.gray,
      },
    },
  },
  plugins: [],
}

