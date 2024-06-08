/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "markata/templates/**/*.{html,js}",
    "markout/**/*.{html,js}",
    "static/**/*.{html,js}",
  ],
  plugins: [require("@tailwindcss/typography")],
  theme: {
    extend: {
      boxShadow: {
        xlc: "0 0 60px 15px rgba(0, 0, 0, 0.3)",
      },
      saturate: {
        25: ".25",
        75: ".75",
      },
      screens: {
        print: { raw: "print" },
      },
      width: {
        128: "32rem",
      },
      boxShadow: {
        xlc: "0 0 60px 15px rgba(0, 0, 0, 0.3)",
        lgc: "0 0 20px 0px rgba(0, 0, 0, 0.3)",
      },
    },
  },
};
