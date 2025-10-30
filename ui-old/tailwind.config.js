/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				// Banner gradient colors (from overview-plan.md)
				'banner-red-start': '#DC2626',
				'banner-red-end': '#991B1B',
				'banner-yellow-start': '#F59E0B',
				'banner-yellow-end': '#FBBF24',
				'banner-green-start': '#10B981',
				'banner-green-end': '#059669',
				'banner-blue-start': '#3B82F6',
				'banner-blue-end': '#1D4ED8',
				'banner-purple-start': '#8B5CF6',
				'banner-purple-end': '#6D28D9',
			},
			spacing: {
				// 8px base spacing system
				'1': '4px',
				'2': '8px',
				'3': '12px',
				'4': '16px',
				'5': '24px',
				'6': '32px',
				'8': '48px',
				'10': '64px',
			},
			fontSize: {
				'xs': '12px',
				'sm': '14px',
				'base': '16px',
				'lg': '18px',
				'xl': '20px',
				'2xl': '24px',
				'3xl': '30px',
				'4xl': '36px', // Banner text size
			},
		},
	},
	plugins: [],
};
