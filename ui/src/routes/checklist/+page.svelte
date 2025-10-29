<script lang="ts">
	import { logger } from '$lib/utils/logger';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Banner from '$lib/components/Banner.svelte';
	import Card from '$lib/components/Card.svelte';
	import CoolOffTimer from '$lib/components/CoolOffTimer.svelte';
	import WorkflowProgress from '$lib/components/WorkflowProgress.svelte';
	import TradingViewLink from '$lib/components/TradingViewLink.svelte';
	import { api } from '$lib/api/client';
	import { startCoolOff, timer } from '$lib/stores/timer';
	import { workflow } from '$lib/stores/workflow';
	import { Save, ArrowRight } from 'lucide-svelte';

	// Banner state management
	type BannerState = 'RED' | 'YELLOW' | 'GREEN';
	let bannerState = $state<BannerState>('RED');
	let bannerMessage = $state('DO NOT TRADE');
	let bannerDetails = $state('Complete the form and check all required gates');

	// Form data
	let formData = $state({
		ticker: '',
		entry: '',
		atr: '',
		sector: '',
		structure: 'Stock',
		journalNote: ''
	});

	// Sector options (from backend - these are the standard buckets)
	const sectors = [
		'Tech/Comm',
		'Finance',
		'Healthcare',
		'Consumer',
		'Industrial',
		'Energy',
		'Materials',
		'Utilities',
		'Real Estate',
		'Other'
	];

	// Structure options
	const structures = [
		'Stock',
		'Call',
		'Put',
		'Call Spread',
		'Put Spread'
	];

	// Required gates (5 gates)
	let gates = $state({
		signal: false,      // Signal: 55-bar Donchian breakout
		risk: false,        // Risk/Size: 2×N stop, proper sizing
		options: false,     // Options: 60-90 DTE, liquidity required
		exits: false,       // Exits: 10-bar Donchian OR 2×N
		behavior: false     // Behavior: Not on cooldown, 2-min timer
	});

	// Optional quality items (4 items)
	let quality = $state({
		regime: false,      // Regime OK (SPY > 200 SMA)
		noChase: false,     // No chase (entry within 2N of 20-EMA)
		earnings: false,    // Earnings OK (no earnings in 2 weeks)
		journal: false      // Journal note written
	});

	// Quality threshold (TODO: fetch from settings API when backend supports it)
	const qualityThreshold = 3;

	// Settings from backend
	let settings = $state<{
		equity: number;
		riskPct: number;
		portfolioCap: number;
		bucketCap: number;
		maxUnits: number;
	} | null>(null);

	// Form validation state
	let formErrors = $state<{[key: string]: string}>({});

	// Validate form data
	function validateForm(): boolean {
		formErrors = {};
		let isValid = true;

		if (!formData.ticker.trim()) {
			formErrors.ticker = 'Ticker is required';
			isValid = false;
		} else if (!/^[A-Z]{1,5}$/.test(formData.ticker.toUpperCase())) {
			formErrors.ticker = 'Invalid ticker format (1-5 uppercase letters)';
			isValid = false;
		}

		if (!formData.entry) {
			formErrors.entry = 'Entry price is required';
			isValid = false;
		} else if (parseFloat(formData.entry) <= 0) {
			formErrors.entry = 'Entry price must be positive';
			isValid = false;
		}

		if (!formData.atr) {
			formErrors.atr = 'ATR/N is required';
			isValid = false;
		} else if (parseFloat(formData.atr) <= 0) {
			formErrors.atr = 'ATR/N must be positive';
			isValid = false;
		}

		if (!formData.sector) {
			formErrors.sector = 'Sector is required';
			isValid = false;
		}

		logger.info('Form validation', { isValid, errors: formErrors, formData });
		return isValid;
	}

	// Handle form field changes
	function handleFieldChange(field: string, value: string) {
		(formData as any)[field] = value;

		// Clear error for this field when user types
		if (formErrors[field]) {
			const newErrors = { ...formErrors };
			delete newErrors[field];
			formErrors = newErrors;
		}

		// Auto-uppercase ticker
		if (field === 'ticker') {
			formData.ticker = value.toUpperCase();
		}

		logger.info('Form field changed', { field, value });
	}

	// Calculate banner state based on gates and quality
	function updateBanner() {
		const requiredCount = Object.values(gates).filter(Boolean).length;
		const qualityCount = Object.values(quality).filter(Boolean).length;
		const formValid = validateForm();

		logger.info('Checklist state updated', {
			requiredGates: requiredCount,
			qualityScore: qualityCount,
			threshold: qualityThreshold,
			formValid
		});

		// If form incomplete or any required gate fails → RED
		if (!formValid || requiredCount < 5) {
			bannerState = 'RED';
			bannerMessage = 'DO NOT TRADE';
			if (!formValid) {
				bannerDetails = 'Complete the form with valid data';
			} else {
				bannerDetails = `${5 - requiredCount} required gate(s) not met`;
			}
		}
		// All required pass, but quality score below threshold → YELLOW
		else if (qualityCount < qualityThreshold) {
			bannerState = 'YELLOW';
			bannerMessage = 'CAUTION';
			bannerDetails = `Quality score: ${qualityCount}/${qualityThreshold} - Consider improving quality items`;
		}
		// All required pass, quality score meets threshold → GREEN
		else {
			bannerState = 'GREEN';
			bannerMessage = 'OK TO TRADE';
			bannerDetails = `All gates pass. Quality score: ${qualityCount}/4. You may proceed.`;
		}
	}

	// Toggle gate
	function toggleGate(gateName: keyof typeof gates) {
		const previousValue = gates[gateName];
		gates[gateName] = !gates[gateName];

		logger.info('Gate toggled', {
			gate: gateName,
			from: previousValue,
			to: gates[gateName],
			timestamp: new Date().toISOString()
		});

		updateBanner();
	}

	// Toggle quality item
	function toggleQuality(itemName: keyof typeof quality) {
		const previousValue = quality[itemName];
		quality[itemName] = !quality[itemName];

		// Auto-check journal if note is written
		if (itemName === 'journal' && formData.journalNote.trim()) {
			quality.journal = true;
		}

		logger.info('Quality item toggled', {
			item: itemName,
			from: previousValue,
			to: quality[itemName],
			timestamp: new Date().toISOString()
		});

		updateBanner();
	}

	// Handle journal note change
	function handleJournalChange(value: string) {
		formData.journalNote = value;
		// Auto-check journal quality item if note has content
		quality.journal = value.trim().length > 0;
		logger.info('Journal note updated', {
			length: value.length,
			autoChecked: quality.journal
		});
		updateBanner();
	}

	// Fetch settings from backend
	async function loadSettings() {
		try {
			const data = await api.getSettings();
			settings = data;
			logger.info('Settings loaded', { settings });
		} catch (error) {
			logger.error('Failed to load settings', { error });
		}
	}

	// State for save evaluation
	let isSaving = $state(false);
	let saveSuccess = $state(false);
	let saveError = $state<string | null>(null);

	// Save evaluation and start timer
	async function saveEvaluation() {
		// Only allow save if banner is GREEN
		if (bannerState !== 'GREEN') {
			logger.warn('Attempted to save evaluation with non-GREEN banner', { bannerState });
			saveError = 'Evaluation can only be saved when banner is GREEN';
			setTimeout(() => saveError = null, 3000);
			return;
		}

		isSaving = true;
		saveError = null;
		saveSuccess = false;

		try {
			// TODO: Call backend API to save evaluation timestamp
			// For now, we'll simulate the API call
			logger.info('Saving checklist evaluation', {
				ticker: formData.ticker,
				gates,
				quality,
				formData,
				timestamp: new Date().toISOString()
			});

			// Simulate API call delay
			await new Promise(resolve => setTimeout(resolve, 500));

			// Update workflow store with checklist results
			const requiredCount = Object.values(gates).filter(Boolean).length;
			const qualityCount = Object.values(quality).filter(Boolean).length;

			// If workflow not started, start it now
			if (!$workflow.workflowStarted) {
				workflow.startTrade(formData.ticker, formData.sector);
			}

			// Update trade information
			workflow.updateTradeInfo(
				formData.ticker,
				parseFloat(formData.entry),
				parseFloat(formData.atr),
				formData.sector
			);

			// Save checklist results
			workflow.saveChecklistResults(bannerState, requiredCount, qualityCount);

			// Start the 2-minute timer
			startCoolOff();

			saveSuccess = true;
			logger.info('Evaluation saved successfully, timer started', {
				workflowState: $workflow
			});

			// Clear success message after 3 seconds
			setTimeout(() => saveSuccess = false, 3000);

		} catch (error) {
			logger.error('Failed to save evaluation', { error });
			saveError = error instanceof Error ? error.message : 'Failed to save evaluation';
			setTimeout(() => saveError = null, 5000);
		} finally {
			isSaving = false;
		}
	}

	// Navigate to next step (sizing)
	function proceedToSizing() {
		if (bannerState === 'GREEN' && $workflow.checklistComplete) {
			logger.info('Proceeding to position sizing');
			goto('/sizing');
		}
	}

	onMount(() => {
		logger.info('Checklist page loaded', {
			timestamp: new Date().toISOString(),
			initialState: {
				gates,
				quality,
				formData
			}
		});
		loadSettings();
		updateBanner();
	});
</script>

<svelte:head><title>Checklist - TF-Engine</title></svelte:head>

<div class="space-y-6">
	<!-- Page Header -->
	<div>
		<h1 class="text-3xl font-bold text-[var(--text-primary)]">Checklist</h1>
		<p class="text-[var(--text-secondary)] mt-2">Complete the form and verify all gates before proceeding</p>
	</div>

	<!-- Workflow Progress -->
	<WorkflowProgress />

	<!-- Trade Information Form -->
	<Card>
		<h2 class="text-2xl font-bold text-[var(--text-primary)] mb-6">Trade Information</h2>

		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			<!-- Ticker -->
			<div>
				<label for="ticker" class="block text-sm font-semibold text-[var(--text-primary)] mb-2">
					Ticker Symbol *
				</label>
				<div class="ticker-input-wrapper">
					<input
						id="ticker"
						type="text"
						value={formData.ticker}
						oninput={(e) => handleFieldChange('ticker', e.currentTarget.value)}
						placeholder="AAPL"
						class="w-full px-4 py-2 rounded-lg border-2 transition-colors
							   {formErrors.ticker ? 'border-red-500' : 'border-[var(--border-color)]'}
							   bg-[var(--bg-secondary)] text-[var(--text-primary)]
							   focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 focus:outline-none
							   placeholder:text-[var(--text-tertiary)]"
					/>
					{#if formData.ticker && !formErrors.ticker}
						<div class="tv-button-wrapper">
							<TradingViewLink ticker={formData.ticker} variant="button" />
						</div>
					{/if}
				</div>
				{#if formErrors.ticker}
					<p class="text-red-500 text-sm mt-1">{formErrors.ticker}</p>
				{/if}
			</div>

			<!-- Entry Price -->
			<div>
				<label for="entry" class="block text-sm font-semibold text-[var(--text-primary)] mb-2">
					Entry Price *
				</label>
				<input
					id="entry"
					type="number"
					step="0.01"
					value={formData.entry}
					oninput={(e) => handleFieldChange('entry', e.currentTarget.value)}
					placeholder="180.00"
					class="w-full px-4 py-2 rounded-lg border-2 transition-colors
						   {formErrors.entry ? 'border-red-500' : 'border-[var(--border-color)]'}
						   bg-[var(--bg-secondary)] text-[var(--text-primary)]
						   focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 focus:outline-none
						   placeholder:text-[var(--text-tertiary)]"
				/>
				{#if formErrors.entry}
					<p class="text-red-500 text-sm mt-1">{formErrors.entry}</p>
				{/if}
			</div>

			<!-- ATR/N -->
			<div>
				<label for="atr" class="block text-sm font-semibold text-[var(--text-primary)] mb-2">
					ATR / N (20-period) *
				</label>
				<input
					id="atr"
					type="number"
					step="0.01"
					value={formData.atr}
					oninput={(e) => handleFieldChange('atr', e.currentTarget.value)}
					placeholder="1.50"
					class="w-full px-4 py-2 rounded-lg border-2 transition-colors
						   {formErrors.atr ? 'border-red-500' : 'border-[var(--border-color)]'}
						   bg-[var(--bg-secondary)] text-[var(--text-primary)]
						   focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 focus:outline-none
						   placeholder:text-[var(--text-tertiary)]"
				/>
				{#if formErrors.atr}
					<p class="text-red-500 text-sm mt-1">{formErrors.atr}</p>
				{/if}
			</div>

			<!-- Sector -->
			<div>
				<label for="sector" class="block text-sm font-semibold text-[var(--text-primary)] mb-2">
					Sector / Bucket *
				</label>
				<select
					id="sector"
					value={formData.sector}
					onchange={(e) => handleFieldChange('sector', e.currentTarget.value)}
					class="w-full px-4 py-2 rounded-lg border-2 transition-colors
						   {formErrors.sector ? 'border-red-500' : 'border-[var(--border-color)]'}
						   bg-[var(--bg-secondary)] text-[var(--text-primary)]
						   focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 focus:outline-none"
				>
					<option value="">-- Select Sector --</option>
					{#each sectors as sector}
						<option value={sector}>{sector}</option>
					{/each}
				</select>
				{#if formErrors.sector}
					<p class="text-red-500 text-sm mt-1">{formErrors.sector}</p>
				{/if}
			</div>

			<!-- Structure -->
			<div>
				<label for="structure" class="block text-sm font-semibold text-[var(--text-primary)] mb-2">
					Structure
				</label>
				<select
					id="structure"
					value={formData.structure}
					onchange={(e) => handleFieldChange('structure', e.currentTarget.value)}
					class="w-full px-4 py-2 rounded-lg border-2 border-[var(--border-color)]
						   bg-[var(--bg-secondary)] text-[var(--text-primary)]
						   focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 focus:outline-none"
				>
					{#each structures as structure}
						<option value={structure}>{structure}</option>
					{/each}
				</select>
			</div>
		</div>
	</Card>

	<!-- Large 3-State Banner (20% viewport height) -->
	<Banner state={bannerState} message={bannerMessage} details={bannerDetails} />

	<!-- Required Gates Section -->
	<Card>
		<h2 class="text-2xl font-bold text-[var(--text-primary)] mb-4">Required Gates (5)</h2>
		<p class="text-sm text-[var(--text-secondary)] mb-6">All 5 must pass to proceed. Any failure → RED banner.</p>

		<div class="space-y-4">
			<!-- Gate 1: Signal -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={gates.signal}
					onchange={() => toggleGate('signal')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-emerald-500 checked:to-emerald-600 focus:ring-2 focus:ring-emerald-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-emerald-600 transition-colors">
						Signal: 55-bar Donchian breakout confirmed
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						Long: Price &gt; 55-bar high | Short: Price &lt; 55-bar low
					</div>
				</div>
			</label>

			<!-- Gate 2: Risk/Size -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={gates.risk}
					onchange={() => toggleGate('risk')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-emerald-500 checked:to-emerald-600 focus:ring-2 focus:ring-emerald-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-emerald-600 transition-colors">
						Risk/Size: 2×N stop, 0.5×N adds, max 4 units
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						Per-unit risk = % of equity using ATR-based stop
					</div>
				</div>
			</label>

			<!-- Gate 3: Options/Liquidity -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={gates.options}
					onchange={() => toggleGate('options')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-emerald-500 checked:to-emerald-600 focus:ring-2 focus:ring-emerald-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-emerald-600 transition-colors">
						Liquidity: Avg volume &gt;1M shares OR options OI &gt;100
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						Options: 60-90 DTE, roll/close ~21 DTE
					</div>
				</div>
			</label>

			<!-- Gate 4: Exits -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={gates.exits}
					onchange={() => toggleGate('exits')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-emerald-500 checked:to-emerald-600 focus:ring-2 focus:ring-emerald-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-emerald-600 transition-colors">
						Exits: 10-bar Donchian OR 2×N stop
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						Mechanical exit, no discretion
					</div>
				</div>
			</label>

			<!-- Gate 5: Behavior -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={gates.behavior}
					onchange={() => toggleGate('behavior')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-emerald-500 checked:to-emerald-600 focus:ring-2 focus:ring-emerald-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-emerald-600 transition-colors">
						Behavior: Not on cooldown, heat OK, 2-min timer honored
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						No intraday overrides, cool-off period enforced
					</div>
				</div>
			</label>
		</div>
	</Card>

	<!-- Optional Quality Items Section -->
	<Card>
		<h2 class="text-2xl font-bold text-[var(--text-primary)] mb-4">Optional Quality Items (4)</h2>
		<p class="text-sm text-[var(--text-secondary)] mb-6">
			Each adds 1 point. Score ≥ {qualityThreshold} → GREEN (if all required pass). Score &lt; {qualityThreshold} → YELLOW.
		</p>

		<div class="space-y-4">
			<!-- Quality 1: Regime -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={quality.regime}
					onchange={() => toggleQuality('regime')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-blue-500 checked:to-blue-600 focus:ring-2 focus:ring-blue-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-blue-600 transition-colors">
						Regime OK (SPY &gt; 200 SMA for longs)
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						Market environment favorable for trend following
					</div>
				</div>
			</label>

			<!-- Quality 2: No Chase -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={quality.noChase}
					onchange={() => toggleQuality('noChase')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-blue-500 checked:to-blue-600 focus:ring-2 focus:ring-blue-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-blue-600 transition-colors">
						No Chase (entry within 2N of 20-EMA)
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						Not entering too far from mean
					</div>
				</div>
			</label>

			<!-- Quality 3: Earnings -->
			<label class="flex items-start space-x-3 cursor-pointer group">
				<input
					type="checkbox"
					checked={quality.earnings}
					onchange={() => toggleQuality('earnings')}
					class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-blue-500 checked:to-blue-600 focus:ring-2 focus:ring-blue-500"
				/>
				<div class="flex-1">
					<div class="font-semibold text-[var(--text-primary)] group-hover:text-blue-600 transition-colors">
						Earnings OK (no earnings within next 2 weeks)
					</div>
					<div class="text-sm text-[var(--text-secondary)]">
						Avoid event risk when holding long premium
					</div>
				</div>
			</label>

			<!-- Quality 4: Journal -->
			<div class="space-y-2">
				<label class="flex items-start space-x-3 cursor-pointer group">
					<input
						type="checkbox"
						checked={quality.journal}
						onchange={() => toggleQuality('journal')}
						class="mt-1 w-5 h-5 rounded border-2 border-[var(--border-color)] checked:bg-gradient-to-br checked:from-blue-500 checked:to-blue-600 focus:ring-2 focus:ring-blue-500"
					/>
					<div class="flex-1">
						<div class="font-semibold text-[var(--text-primary)] group-hover:text-blue-600 transition-colors">
							Journal Note (why this trade now?)
						</div>
						<div class="text-sm text-[var(--text-secondary)]">
							Document reasoning and plan
						</div>
					</div>
				</label>

				<!-- Journal Textarea -->
				<div class="ml-8">
					<textarea
						value={formData.journalNote}
						oninput={(e) => handleJournalChange(e.currentTarget.value)}
						placeholder="Why this trade now? What's your thesis? Exit plan for spreads?"
						rows="3"
						class="w-full px-4 py-2 rounded-lg border-2 border-[var(--border-color)]
							   bg-[var(--bg-secondary)] text-[var(--text-primary)]
							   focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none
							   placeholder:text-[var(--text-tertiary)] resize-y"
					></textarea>
					<p class="text-xs text-[var(--text-tertiary)] mt-1">
						Writing a journal note auto-checks the quality item
					</p>
				</div>
			</div>
		</div>

		<!-- Quality Score Display -->
		<div class="mt-6 p-4 bg-[var(--bg-primary)] rounded-lg border border-[var(--border-color)]">
			<div class="flex justify-between items-center">
				<span class="text-[var(--text-secondary)]">Quality Score:</span>
				<span class="text-2xl font-bold text-[var(--text-primary)]">
					{Object.values(quality).filter(Boolean).length} / 4
				</span>
			</div>
			<div class="text-sm text-[var(--text-tertiary)] mt-2">
				Threshold: {qualityThreshold} (score ≥ threshold → GREEN if all required pass)
			</div>
		</div>
	</Card>

	<!-- Save Evaluation & Timer Section -->
	<Card title="Save Evaluation">
		<div class="space-y-6">
			<!-- Save Evaluation Button -->
			<div class="flex flex-col gap-4">
				<button
					onclick={saveEvaluation}
					disabled={bannerState !== 'GREEN' || isSaving || $timer.running}
					class="w-full px-6 py-4 rounded-lg font-semibold text-white text-lg
						   transition-all duration-300 ease-in-out transform
						   bg-gradient-to-r from-emerald-500 to-emerald-600
						   hover:from-emerald-600 hover:to-emerald-700
						   disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed
						   disabled:opacity-50
						   shadow-lg hover:shadow-xl active:scale-95
						   flex items-center justify-center gap-2"
				>
					{#if isSaving}
						<span class="inline-block w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
						Saving...
					{:else if $timer.running}
						<span>Timer Active</span>
					{:else}
						<Save class="w-5 h-5" />
						Save Evaluation
					{/if}
				</button>

				<!-- Help Text -->
				<p class="text-sm text-[var(--text-tertiary)] text-center">
					{#if bannerState !== 'GREEN'}
						Button enabled when banner is GREEN
					{:else if $timer.running}
						Evaluation saved - Timer is running
					{:else}
						Click to save evaluation and start 2-minute impulse timer
					{/if}
				</p>

				<!-- Success/Error Messages -->
				{#if saveSuccess}
					<div class="p-4 bg-gradient-to-r from-emerald-500 to-emerald-600 text-white rounded-lg shadow-lg">
						<p class="font-semibold">✓ Evaluation Saved Successfully</p>
						<p class="text-sm opacity-90 mt-1">2-minute impulse timer has started</p>
					</div>
				{/if}

				{#if saveError}
					<div class="p-4 bg-gradient-to-r from-red-500 to-red-600 text-white rounded-lg shadow-lg">
						<p class="font-semibold">✗ Error</p>
						<p class="text-sm opacity-90 mt-1">{saveError}</p>
					</div>
				{/if}

				<!-- Proceed to Sizing Button -->
				{#if $workflow.checklistComplete && bannerState === 'GREEN'}
					<button
						onclick={proceedToSizing}
						class="w-full px-6 py-4 rounded-lg font-semibold text-white text-lg
							   transition-all duration-300 ease-in-out transform
							   bg-gradient-to-r from-blue-500 to-blue-600
							   hover:from-blue-600 hover:to-blue-700
							   shadow-lg hover:shadow-xl active:scale-95
							   flex items-center justify-center gap-2"
					>
						<span>Proceed to Position Sizing</span>
						<ArrowRight class="w-5 h-5" />
					</button>
					<p class="text-sm text-blue-600 dark:text-blue-400 text-center font-medium">
						Next: Calculate your position size using ATR-based sizing
					</p>
				{/if}
			</div>

			<!-- Cool-Off Timer Display -->
			{#if $timer.running || $timer.startTime}
				<div class="mt-6">
					<CoolOffTimer />
				</div>
			{/if}

			<!-- Timer Explanation -->
			{#if !$timer.running && !$timer.startTime}
				<div class="p-4 bg-[var(--bg-secondary)] rounded-lg border border-[var(--border-color)]">
					<h3 class="font-semibold text-[var(--text-primary)] mb-2">What happens after saving?</h3>
					<ul class="text-sm text-[var(--text-secondary)] space-y-2">
						<li>• A 2-minute impulse timer starts</li>
						<li>• This prevents impulsive trade decisions</li>
						<li>• Use this time to review your analysis</li>
						<li>• Calculate position sizing (see Sizing tab)</li>
						<li>• Check portfolio heat (see Heat tab)</li>
						<li>• After 2 minutes, you may proceed to Trade Entry</li>
					</ul>
				</div>
			{/if}
		</div>
	</Card>
</div>

<style>
	.ticker-input-wrapper {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.tv-button-wrapper {
		display: flex;
		justify-content: flex-start;
	}
</style>
