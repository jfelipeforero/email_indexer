/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html","./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
            "primary":"black"
      },
      backgroundImage: {
        'my-component': "./src/assets/teepublic_enron-logo-t-shirt-defunct-finance-company-corporate-humor-teepublic_1647802506.large.png"
      },
    },
    container: {
      padding: "2rem",
      center: true, 
    },
    screens: {
      sm: "640px",
      md: "768px",
    },
  },
  plugins: [],
}

