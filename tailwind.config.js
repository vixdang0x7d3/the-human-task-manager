/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './internal/http/templates/**/*.templ',
    './internal/http/assets/**/*.js'
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

