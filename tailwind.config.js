/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/**/*.{html,js,go}"],
  theme: {
    extend: {
      colors: {
        main: "#e51636",
        black: "#222222",
        white: "#EEEEEE",
        offwhite: "#cccccc",
        darkgray: "#303030",
        gray: "#444444",
        lightgray: "#555555"
      },
      fontFamily: {
        main: "Montserrat",
        second: "Lato",
      },
    },
  },
  plugins: [],
}