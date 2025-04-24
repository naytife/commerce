import type { Config } from "tailwindcss";
import flowbitePlugin from 'flowbite/plugin'

const config: Config = {
	plugins: [flowbitePlugin],
	content: ["./src/**/*.{html,js,svelte,ts}",
			"./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}"],
	darkMode: 'class',
	theme: {
		container: {
			center: true,
			padding: "2rem",
			screens: {
				"2xl": "1400px"
			}
		},
		extend: {
			colors: {
			  primary: { 50: '#FFF5F2', 100: '#FFF1EE', 200: '#FFE4DE', 300: '#FFD5CC', 400: '#FFBCAD', 500: '#FE795D', 600: '#EF562F', 700: '#EB4F27', 800: '#CC4522', 900: '#A5371B'},
			}
		  }
	},
};

export default config;
