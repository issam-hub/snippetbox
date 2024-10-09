/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/templates/**/*.templ","./**/*.html", "./**/*.templ", "./**/*.go"],
  theme: {
    extend: {
      fontFamily:{
        'inter': '"Inter", sans-serif'
      }
    },
  },
  plugins: [],
}

