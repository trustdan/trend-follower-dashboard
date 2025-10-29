/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			// Custom colors for day/night mode
			colors: {
				// Banner gradients
				'banner-red-start': '#DC2626',
				'banner-red-end': '#991B1B',
				'banner-yellow-start': '#F59E0B',
				'banner-yellow-end': '#FBBF24',
				'banner-green-start': '#10B981',
				'banner-green-end': '#059669',
			},
		},
	},
	plugins: [],
	// Enable dark mode via class
	darkMode: 'class',
};
