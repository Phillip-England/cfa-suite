/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/**/*.{html,js,go}"],
  theme: {
    extend: {
      colors: {
        main: "#e51636",
        mainZero: "#822b24",
        black: "#222222",
        white: "#EEEEEE",
        offwhite: "#cccccc",
        darkgray: "#303030",
        gray: "#3d3d3d",
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