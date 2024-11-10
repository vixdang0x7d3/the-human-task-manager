/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './internal/template/**/*.templ',
    './static/**/*.js'
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require("daisyui")
  ],

  daisyui: {
    themes: [ "light", "dark", "retro", "coffee"]
  }
}

