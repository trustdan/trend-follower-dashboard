<script lang="ts">
	import { theme } from '$lib/stores/theme';
	import { logger } from '$lib/utils/logger';
	import { Sun, Moon, Settings } from 'lucide-svelte';

	let currentTheme: 'light' | 'dark';
	theme.subscribe(value => (currentTheme = value));

	function toggleTheme() {
		theme.toggle();
		logger.themeChange(currentTheme === 'light' ? 'dark' : 'light');
	}
</script>

<header class="bg-[var(--bg-secondary)] border-b border-[var(--border-color)] px-6 py-4">
	<div class="flex items-center justify-between">
		<!-- App Title -->
		<div class="flex items-center gap-3">
			<h1 class="text-2xl font-bold text-[var(--text-primary)]">TF-Engine</h1>
			<span class="text-sm text-[var(--text-secondary)] font-medium">
				Trend Following Dashboard
			</span>
		</div>

		<!-- Right side actions -->
		<div class="flex items-center gap-4">
			<!-- Theme Toggle Button -->
			<button
				on:click={toggleTheme}
				class="p-2 rounded-lg hover:bg-[var(--bg-tertiary)] transition-colors"
				aria-label="Toggle theme"
				title={currentTheme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
			>
				{#if currentTheme === 'light'}
					<Moon size={20} class="text-[var(--text-secondary)]" />
				{:else}
					<Sun size={20} class="text-[var(--text-secondary)]" />
				{/if}
			</button>

			<!-- Settings Button -->
			<button
				class="p-2 rounded-lg hover:bg-[var(--bg-tertiary)] transition-colors"
				aria-label="Settings"
				title="Settings"
			>
				<Settings size={20} class="text-[var(--text-secondary)]" />
			</button>
		</div>
	</div>
</header>
