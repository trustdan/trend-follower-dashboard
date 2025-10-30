// Theme store for managing dark/light mode
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type Theme = 'light' | 'dark';

// Initialize from localStorage or default to 'light'
function getInitialTheme(): Theme {
	if (!browser) return 'light';

	const stored = localStorage.getItem('theme');
	if (stored === 'dark' || stored === 'light') {
		return stored;
	}

	// Check system preference
	if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
		return 'dark';
	}

	return 'light';
}

function createThemeStore() {
	const { subscribe, set, update } = writable<Theme>(getInitialTheme());

	return {
		subscribe,
		toggle: () => update(theme => {
			const newTheme = theme === 'light' ? 'dark' : 'light';

			if (browser) {
				localStorage.setItem('theme', newTheme);

				// Update document class for theme
				if (newTheme === 'dark') {
					document.documentElement.classList.add('dark');
				} else {
					document.documentElement.classList.remove('dark');
				}
			}

			return newTheme;
		}),
		setTheme: (theme: Theme) => {
			if (browser) {
				localStorage.setItem('theme', theme);

				if (theme === 'dark') {
					document.documentElement.classList.add('dark');
				} else {
					document.documentElement.classList.remove('dark');
				}
			}

			set(theme);
		},
		initialize: () => {
			if (browser) {
				const theme = getInitialTheme();
				if (theme === 'dark') {
					document.documentElement.classList.add('dark');
				}
			}
		}
	};
}

export const theme = createThemeStore();
