<script lang="ts">
	import { logger } from '$lib/utils/logger';
	import { api } from '$lib/api/client';
	import { onMount } from 'svelte';
	import Card from '$lib/components/Card.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import TradingViewLink from '$lib/components/TradingViewLink.svelte';
	import { Search, CheckCircle, XCircle, AlertCircle } from 'lucide-svelte';

	// State
	let selectedPreset = $state('TF_BREAKOUT_LONG');
	let scanning = $state(false);
	let scanResults: string[] = $state([]);
	let scanDate = $state('');
	let scanCount = $state(0);
	let error = $state('');
	let importing = $state(false);
	let selectedTickers = $state<Set<string>>(new Set());
	let importSuccess = $state(false);

	// Available presets
	const presets = [
		{ value: 'TF_BREAKOUT_LONG', label: 'TF Breakout Long', description: 'Channel up + 1w performance > 10%' }
	];

	onMount(() => logger.info('Scanner page loaded'));

	async function runScan() {
		scanning = true;
		error = '';
		scanResults = [];
		importSuccess = false;

		logger.info('Starting FINVIZ scan', { preset: selectedPreset });

		try {
			const result = await api.scanCandidates(selectedPreset);
			scanResults = result.tickers;
			scanDate = result.date;
			scanCount = result.count;

			logger.info('Scan completed', { count: result.count, date: result.date });

			// Auto-select all tickers by default
			selectedTickers = new Set(result.tickers);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to scan FINVIZ';
			logger.error('Scan failed', err);
		} finally {
			scanning = false;
		}
	}

	async function importSelected() {
		if (selectedTickers.size === 0) {
			error = 'Please select at least one ticker to import';
			return;
		}

		importing = true;
		error = '';

		logger.info('Importing candidates', { count: selectedTickers.size });

		try {
			const tickerArray = Array.from(selectedTickers);
			const result = await api.importCandidates(tickerArray, scanDate);

			logger.info('Import completed', { imported: result.imported, date: result.date });

			importSuccess = true;
			// Clear results after successful import
			setTimeout(() => {
				scanResults = [];
				selectedTickers = new Set();
				importSuccess = false;
			}, 3000);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to import candidates';
			logger.error('Import failed', err);
		} finally {
			importing = false;
		}
	}

	function toggleTicker(ticker: string) {
		if (selectedTickers.has(ticker)) {
			selectedTickers.delete(ticker);
		} else {
			selectedTickers.add(ticker);
		}
		selectedTickers = selectedTickers; // Trigger reactivity
	}

	function selectAll() {
		selectedTickers = new Set(scanResults);
	}

	function deselectAll() {
		selectedTickers = new Set();
	}
</script>

<svelte:head><title>Scanner - TF-Engine</title></svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div>
		<h1 class="text-3xl font-bold text-[var(--text-primary)]">FINVIZ Scanner</h1>
		<p class="text-[var(--text-secondary)] mt-2">Daily screening for trend-following candidates</p>
	</div>

	<!-- Preset Selector -->
	<Card title="Scan Configuration">
		<div class="space-y-4">
			<div>
				<label for="preset" class="block text-sm font-medium text-[var(--text-primary)] mb-2">
					Scan Preset
				</label>
				<select
					id="preset"
					bind:value={selectedPreset}
					class="w-full px-4 py-2 bg-[var(--bg-primary)] border border-[var(--border-color)] rounded-lg text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					{#each presets as preset}
						<option value={preset.value}>{preset.label}</option>
					{/each}
				</select>
				<p class="text-sm text-[var(--text-tertiary)] mt-1">
					{presets.find((p) => p.value === selectedPreset)?.description}
				</p>
			</div>

			<!-- Scan Button -->
			<button
				onclick={runScan}
				disabled={scanning}
				class="w-full bg-gradient-to-r from-blue-500 to-blue-600 text-white px-8 py-4 rounded-lg font-semibold text-lg hover:from-blue-600 hover:to-blue-700 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 shadow-lg hover:shadow-xl"
			>
				{#if scanning}
					<LoadingSpinner message="Scanning FINVIZ..." />
				{:else}
					<Search size={24} />
					<span>Run Daily Scan</span>
				{/if}
			</button>
		</div>
	</Card>

	<!-- Error Message -->
	{#if error}
		<div
			class="bg-gradient-to-r from-red-500/20 to-red-600/20 border border-red-500 rounded-lg p-4 flex items-start gap-3"
		>
			<XCircle class="text-red-500 flex-shrink-0" size={24} />
			<div>
				<h3 class="text-red-500 font-semibold mb-1">Scan Error</h3>
				<p class="text-[var(--text-secondary)]">{error}</p>
			</div>
		</div>
	{/if}

	<!-- Import Success Message -->
	{#if importSuccess}
		<div
			class="bg-gradient-to-r from-green-500/20 to-green-600/20 border border-green-500 rounded-lg p-4 flex items-start gap-3"
		>
			<CheckCircle class="text-green-500 flex-shrink-0" size={24} />
			<div>
				<h3 class="text-green-500 font-semibold mb-1">Import Successful</h3>
				<p class="text-[var(--text-secondary)]">
					{selectedTickers.size} candidates imported for {scanDate}
				</p>
			</div>
		</div>
	{/if}

	<!-- Scan Results -->
	{#if scanResults.length > 0}
		<Card title="Scan Results">
			<div class="space-y-4">
				<!-- Results Summary -->
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-3">
						<AlertCircle class="text-blue-500" size={20} />
						<p class="text-[var(--text-primary)]">
							Found <strong>{scanCount}</strong> candidates on {scanDate}
						</p>
					</div>
					<div class="flex gap-2">
						<button
							onclick={selectAll}
							class="px-3 py-1 text-sm bg-[var(--bg-secondary)] border border-[var(--border-color)] rounded hover:bg-[var(--bg-primary)] transition-colors"
						>
							Select All
						</button>
						<button
							onclick={deselectAll}
							class="px-3 py-1 text-sm bg-[var(--bg-secondary)] border border-[var(--border-color)] rounded hover:bg-[var(--bg-primary)] transition-colors"
						>
							Deselect All
						</button>
					</div>
				</div>

				<!-- Results Grid -->
				<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
					{#each scanResults as ticker}
						<div class="ticker-card-wrapper">
							<button
								onclick={() => toggleTicker(ticker)}
								class="ticker-card px-4 py-3 rounded-lg border-2 transition-all duration-200 {selectedTickers.has(
									ticker
								)
									? 'border-blue-500 bg-blue-500/10 text-blue-600 font-semibold'
									: 'border-[var(--border-color)] bg-[var(--bg-secondary)] text-[var(--text-primary)] hover:border-blue-400'}"
							>
								{ticker}
							</button>
							<div class="tv-link-overlay">
								<TradingViewLink {ticker} variant="icon" />
							</div>
						</div>
					{/each}
				</div>

				<!-- Selection Summary -->
				<div class="text-sm text-[var(--text-tertiary)] text-center">
					{selectedTickers.size} of {scanCount} selected
				</div>

				<!-- Import Button -->
				<button
					onclick={importSelected}
					disabled={importing || selectedTickers.size === 0}
					class="w-full bg-gradient-to-r from-green-500 to-green-600 text-white px-8 py-4 rounded-lg font-semibold text-lg hover:from-green-600 hover:to-green-700 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 shadow-lg hover:shadow-xl"
				>
					{#if importing}
						<LoadingSpinner message="Importing..." />
					{:else}
						<CheckCircle size={24} />
						<span>Import Selected ({selectedTickers.size})</span>
					{/if}
				</button>
			</div>
		</Card>
	{/if}
</div>

<style>
	.ticker-card-wrapper {
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.ticker-card {
		width: 100%;
	}

	.tv-link-overlay {
		display: flex;
		justify-content: center;
	}
</style>
