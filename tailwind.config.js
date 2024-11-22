/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './internal/http/templates/**/*.templ',
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

