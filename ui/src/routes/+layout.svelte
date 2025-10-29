<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';
	import { theme } from '$lib/stores/theme';
	import { logger } from '$lib/utils/logger';
	import { setupKeyboardShortcuts } from '$lib/utils/keyboard';
	import Header from '$lib/components/Header.svelte';
	import Navigation from '$lib/components/Navigation.svelte';
	import Breadcrumbs from '$lib/components/Breadcrumbs.svelte';
	import DebugPanel from '$lib/components/DebugPanel.svelte';

	let { children } = $props();

	// Initialize theme and keyboard shortcuts on mount
	onMount(() => {
		theme.initialize();
		setupKeyboardShortcuts();
		logger.info('TF-Engine application started');
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Main application layout -->
<div class="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)]">
	<!-- Header -->
	<Header />

	<!-- Main content area with sidebar -->
	<div class="flex h-[calc(100vh-64px)]">
		<!-- Sidebar Navigation -->
		<Navigation />

		<!-- Page Content with Breadcrumbs -->
		<div class="flex-1 flex flex-col">
			<!-- Breadcrumbs -->
			<Breadcrumbs />

			<!-- Main Content -->
			<main class="flex-1 overflow-y-auto p-6 page-transition">
				{@render children?.()}
			</main>
		</div>
	</div>

	<!-- Debug Panel (dev mode only) -->
	<DebugPanel />
</div>
