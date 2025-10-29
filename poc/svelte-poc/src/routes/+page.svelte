<script lang="ts">
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	// Theme store
	const isDark = writable(false);

	// Settings data
	let settings = {
		equity: 0,
		riskPct: 0,
		portfolioCap: 0,
		bucketCap: 0,
	};

	let loading = false;
	let error = '';

	// Load theme preference from localStorage
	onMount(() => {
		const savedTheme = localStorage.getItem('theme');
		if (savedTheme === 'dark') {
			isDark.set(true);
			document.documentElement.classList.add('dark');
		}

		// Fetch settings on mount
		fetchSettings();
	});

	// Toggle theme
	function toggleTheme() {
		isDark.update(v => {
			const newValue = !v;
			if (newValue) {
				document.documentElement.classList.add('dark');
				localStorage.setItem('theme', 'dark');
			} else {
				document.documentElement.classList.remove('dark');
				localStorage.setItem('theme', 'light');
			}
			return newValue;
		});
	}

	// Fetch settings from API
	async function fetchSettings() {
		loading = true;
		error = '';

		try {
			const response = await fetch('/api/settings');
			if (!response.ok) {
				throw new Error(`HTTP ${response.status}: ${response.statusText}`);
			}
			const data = await response.json();
			settings = data;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to fetch settings';
			console.error('Error fetching settings:', err);
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)] p-8">
	<!-- Header with Theme Toggle -->
	<header class="flex justify-between items-center mb-8">
		<h1 class="text-3xl font-bold">TF-Engine Settings</h1>

		<button
			on:click={toggleTheme}
			class="p-2 rounded-lg bg-[var(--bg-secondary)] hover:bg-[var(--bg-tertiary)] border border-[var(--border-color)]"
			aria-label="Toggle theme"
		>
			{#if $isDark}
				<!-- Sun icon (day mode) -->
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
				</svg>
			{:else}
				<!-- Moon icon (night mode) -->
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
				</svg>
			{/if}
		</button>
	</header>

	<!-- Main Content Card -->
	<div class="max-w-2xl mx-auto bg-[var(--bg-secondary)] rounded-xl shadow-lg p-6 border border-[var(--border-color)]">
		<h2 class="text-2xl font-semibold mb-6">Account Settings</h2>

		{#if loading}
			<div class="text-center py-8">
				<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--border-focus)]"></div>
				<p class="mt-2 text-[var(--text-secondary)]">Loading...</p>
			</div>
		{:else if error}
			<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-4">
				<p class="text-red-800 dark:text-red-200">{error}</p>
			</div>
		{/if}

		<!-- Settings Display -->
		<div class="space-y-4">
			<div class="flex justify-between items-center py-3 border-b border-[var(--border-color)]">
				<span class="text-[var(--text-secondary)]">Equity:</span>
				<span class="font-semibold text-lg">${settings.equity.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
			</div>

			<div class="flex justify-between items-center py-3 border-b border-[var(--border-color)]">
				<span class="text-[var(--text-secondary)]">Risk % per Unit:</span>
				<span class="font-semibold text-lg">{settings.riskPct.toFixed(2)}%</span>
			</div>

			<div class="flex justify-between items-center py-3 border-b border-[var(--border-color)]">
				<span class="text-[var(--text-secondary)]">Portfolio Heat Cap:</span>
				<span class="font-semibold text-lg">{settings.portfolioCap.toFixed(2)}%</span>
			</div>

			<div class="flex justify-between items-center py-3">
				<span class="text-[var(--text-secondary)]">Sector Bucket Cap:</span>
				<span class="font-semibold text-lg">{settings.bucketCap.toFixed(2)}%</span>
			</div>
		</div>

		<!-- Action Buttons -->
		<div class="mt-6 flex gap-4">
			<button
				on:click={fetchSettings}
				disabled={loading}
				class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
			>
				Refresh
			</button>

			<button
				disabled
				class="flex-1 px-4 py-2 bg-gray-300 dark:bg-gray-700 text-gray-500 dark:text-gray-400 font-medium rounded-lg cursor-not-allowed"
			>
				Update (Not Implemented)
			</button>
		</div>

		<!-- Status -->
		<p class="mt-4 text-sm text-[var(--text-tertiary)] text-center">
			Status: Ready • API: {loading ? 'Loading...' : error ? 'Error' : 'Connected'}
		</p>
	</div>

	<!-- Demo Banner (shows gradient transitions) -->
	<div class="max-w-2xl mx-auto mt-8">
		<h3 class="text-xl font-semibold mb-4">Banner Demo (Anti-Impulsivity Core)</h3>
		<div class="banner banner-green rounded-xl p-8 text-white text-center shadow-lg">
			<p class="text-3xl font-bold mb-2">✓ OK TO TRADE ✓</p>
			<p class="text-lg opacity-90">All gates pass • Quality score met</p>
		</div>
	</div>
</div>
