<script lang="ts">
	import { logger } from '$lib/utils/logger';
	import { api, type SizingRequest, type SizingResult, type Settings } from '$lib/api/client';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Card from '$lib/components/Card.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import WorkflowProgress from '$lib/components/WorkflowProgress.svelte';
	import { workflow } from '$lib/stores/workflow';
	import { Calculator, TrendingUp, AlertTriangle, BarChart3, ArrowRight } from 'lucide-svelte';

	// State
	let settings: Settings | null = $state(null);
	let loading = $state(false);
	let calculating = $state(false);
	let error = $state('');
	let result: SizingResult | null = $state(null);

	// Form fields
	let ticker = $state('');
	let entry = $state<number | null>(null);
	let atrN = $state<number | null>(null);
	let k = $state(2);
	let maxUnits = $state(4);
	let method = $state<'stock' | 'opt-delta-atr' | 'opt-maxloss'>('stock');
	let delta = $state<number | null>(null);
	let maxLoss = $state<number | null>(null);

	// Form validation errors
	let formErrors: Record<string, string> = $state({});

	onMount(async () => {
		logger.info('Position Sizing page loaded');
		await loadSettings();

		// Pre-fill form fields from workflow if available
		if ($workflow.ticker) {
			ticker = $workflow.ticker;
			logger.info('Pre-filling ticker from workflow', ticker);
		}
		if ($workflow.entryPrice) {
			entry = $workflow.entryPrice;
			logger.info('Pre-filling entry from workflow', entry);
		}
		if ($workflow.atrN) {
			atrN = $workflow.atrN;
			logger.info('Pre-filling ATR from workflow', atrN);
		}
		if ($workflow.method) {
			method = $workflow.method;
			logger.info('Pre-filling method from workflow', method);
		}

		// Update workflow to indicate we're on sizing step
		if ($workflow.workflowStarted && $workflow.currentStep !== 'sizing') {
			workflow.goToStep('sizing');
		}
	});

	async function loadSettings() {
		loading = true;
		error = '';
		try {
			settings = await api.getSettings();
			logger.info('Settings loaded', settings);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load settings';
			logger.error('Failed to load settings', err);
		} finally {
			loading = false;
		}
	}

	function validateForm(): boolean {
		const errors: Record<string, string> = {};

		if (!ticker || ticker.length === 0) {
			errors.ticker = 'Ticker is required';
		} else if (!/^[A-Z]{1,5}$/.test(ticker)) {
			errors.ticker = 'Ticker must be 1-5 uppercase letters';
		}

		if (entry === null || entry <= 0) {
			errors.entry = 'Entry price must be positive';
		}

		if (atrN === null || atrN <= 0) {
			errors.atrN = 'ATR must be positive';
		}

		if (k <= 0) {
			errors.k = 'K multiple must be positive';
		}

		if (maxUnits < 1 || maxUnits > 10) {
			errors.maxUnits = 'Max units must be between 1 and 10';
		}

		if (method === 'opt-delta-atr' && (delta === null || delta <= 0 || delta > 1)) {
			errors.delta = 'Delta must be between 0 and 1';
		}

		if (method === 'opt-maxloss' && (maxLoss === null || maxLoss <= 0)) {
			errors.maxLoss = 'Max loss must be positive';
		}

		formErrors = errors;
		return Object.keys(errors).length === 0;
	}

	async function handleCalculate() {
		logger.info('Calculate button clicked');

		if (!validateForm()) {
			logger.warn('Form validation failed', formErrors);
			return;
		}

		if (!settings) {
			error = 'Settings not loaded';
			return;
		}

		calculating = true;
		error = '';
		result = null;

		try {
			const request: SizingRequest = {
				equity: settings.equity,
				risk_pct: settings.riskPct,
				entry: entry!,
				atr_n: atrN!,
				k: k,
				method: method
			};

			if (method === 'opt-delta-atr' && delta !== null) {
				request.delta = delta;
			}

			if (method === 'opt-maxloss' && maxLoss !== null) {
				request.max_loss = maxLoss;
			}

			logger.info('Calculating position size', request);
			result = await api.calculateSize(request);
			logger.info('Position size calculated', result);

			// Save sizing results to workflow
			workflow.saveSizingResults(result);
			logger.info('Sizing results saved to workflow', { workflowState: $workflow });
		} catch (err) {
			error = err instanceof Error ? err.message : 'Calculation failed';
			logger.error('Position size calculation failed', err);
		} finally {
			calculating = false;
		}
	}

	// Navigate to heat check
	function proceedToHeatCheck() {
		if ($workflow.sizingComplete) {
			logger.info('Proceeding to heat check');
			goto('/heat');
		}
	}

	function handleFieldChange(field: string, value: string | number) {
		logger.info('Field changed', { field, value });
		// Clear error for this field when user starts typing
		if (formErrors[field]) {
			const { [field]: _, ...rest } = formErrors;
			formErrors = rest;
		}
	}

	function handleTickerChange(e: Event) {
		const target = e.target as HTMLInputElement;
		const value = target.value.toUpperCase();
		ticker = value;
		handleFieldChange('ticker', value);
	}

	// Calculate add-on levels (0.5N increments)
	function getAddOnLevels(): { unit: number; price: number }[] {
		if (!result || entry === null || atrN === null || !settings) return [];

		const levels: { unit: number; price: number }[] = [];
		const increment = atrN * 0.5;

		for (let i = 0; i < maxUnits; i++) {
			levels.push({
				unit: i + 1,
				price: entry + (i * increment)
			});
		}

		return levels;
	}

	// Calculate concentration percentage
	function getConcentration(): number {
		if (!result || !settings) return 0;
		const positionValue = result.shares * (entry || 0);
		return (positionValue / settings.equity) * 100;
	}

	function isHighConcentration(): boolean {
		return getConcentration() > 25;
	}
</script>

<svelte:head><title>Position Sizing - TF-Engine</title></svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div>
		<h1 class="text-3xl font-bold text-[var(--text-primary)]">Position Sizing Calculator</h1>
		<p class="text-[var(--text-secondary)] mt-2">
			Calculate shares/contracts using ATR-based risk (Van Tharp method)
		</p>
	</div>

	<!-- Workflow Progress -->
	<WorkflowProgress />

	{#if loading}
		<LoadingSpinner message="Loading settings..." />
	{:else if error && !settings}
		<Card>
			<div class="text-center py-8">
				<AlertTriangle class="w-16 h-16 mx-auto text-red-500 mb-4" />
				<p class="text-red-600 dark:text-red-400 mb-4">{error}</p>
				<button
					onclick={loadSettings}
					class="px-4 py-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg hover:shadow-lg transition-all"
				>
					Retry
				</button>
			</div>
		</Card>
	{:else if settings}
		<!-- Account Info Banner -->
		<div
			class="bg-gradient-to-r from-blue-500/10 to-blue-600/10 border border-blue-500/30 rounded-lg p-4"
		>
			<div class="flex items-center gap-4 flex-wrap">
				<div>
					<span class="text-[var(--text-secondary)] text-sm">Account Equity:</span>
					<span class="text-[var(--text-primary)] font-semibold ml-2"
						>${settings.equity.toLocaleString()}</span
					>
				</div>
				<div>
					<span class="text-[var(--text-secondary)] text-sm">Risk per Unit:</span>
					<span class="text-[var(--text-primary)] font-semibold ml-2"
						>{(settings.riskPct * 100).toFixed(2)}%</span
					>
				</div>
				<div>
					<span class="text-[var(--text-secondary)] text-sm">Risk Dollars:</span>
					<span class="text-[var(--text-primary)] font-semibold ml-2"
						>${(settings.equity * settings.riskPct).toLocaleString()}</span
					>
				</div>
			</div>
		</div>

		<!-- Input Form -->
		<Card>
			<div class="space-y-6">
				<h2 class="text-xl font-semibold text-[var(--text-primary)] flex items-center gap-2">
					<Calculator class="w-6 h-6" />
					Trade Information
				</h2>

				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					<!-- Ticker -->
					<div>
						<label for="ticker" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
							Ticker Symbol <span class="text-red-500">*</span>
						</label>
						<input
							id="ticker"
							type="text"
							value={ticker}
							oninput={handleTickerChange}
							placeholder="AAPL"
							maxlength="5"
							class="w-full px-3 py-2 bg-[var(--bg-primary)] border {formErrors.ticker
								? 'border-red-500'
								: 'border-[var(--border-color)]'} rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)] uppercase"
						/>
						{#if formErrors.ticker}
							<p class="text-red-500 text-sm mt-1">{formErrors.ticker}</p>
						{/if}
					</div>

					<!-- Entry Price -->
					<div>
						<label for="entry" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
							Entry Price <span class="text-red-500">*</span>
						</label>
						<input
							id="entry"
							type="number"
							bind:value={entry}
							oninput={() => handleFieldChange('entry', entry || 0)}
							placeholder="180.00"
							step="0.01"
							class="w-full px-3 py-2 bg-[var(--bg-primary)] border {formErrors.entry
								? 'border-red-500'
								: 'border-[var(--border-color)]'} rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)]"
						/>
						{#if formErrors.entry}
							<p class="text-red-500 text-sm mt-1">{formErrors.entry}</p>
						{/if}
					</div>

					<!-- ATR (N) -->
					<div>
						<label for="atrN" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
							ATR (N) <span class="text-red-500">*</span>
						</label>
						<input
							id="atrN"
							type="number"
							bind:value={atrN}
							oninput={() => handleFieldChange('atrN', atrN || 0)}
							placeholder="1.50"
							step="0.01"
							class="w-full px-3 py-2 bg-[var(--bg-primary)] border {formErrors.atrN
								? 'border-red-500'
								: 'border-[var(--border-color)]'} rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)]"
						/>
						{#if formErrors.atrN}
							<p class="text-red-500 text-sm mt-1">{formErrors.atrN}</p>
						{/if}
					</div>

					<!-- K Multiple -->
					<div>
						<label for="k" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
							K Multiple (Stop)
						</label>
						<input
							id="k"
							type="number"
							bind:value={k}
							oninput={() => handleFieldChange('k', k)}
							placeholder="2"
							step="0.1"
							min="0.1"
							class="w-full px-3 py-2 bg-[var(--bg-primary)] border {formErrors.k
								? 'border-red-500'
								: 'border-[var(--border-color)]'} rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)]"
						/>
						{#if atrN}
							<p class="text-xs text-[var(--text-tertiary)] mt-1">
								Stop Distance: {(k * atrN).toFixed(2)}
							</p>
						{/if}
						{#if formErrors.k}
							<p class="text-red-500 text-sm mt-1">{formErrors.k}</p>
						{/if}
					</div>

					<!-- Max Units -->
					<div>
						<label for="maxUnits" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
							Max Units (Pyramid)
						</label>
						<input
							id="maxUnits"
							type="number"
							bind:value={maxUnits}
							oninput={() => handleFieldChange('maxUnits', maxUnits)}
							placeholder="4"
							min="1"
							max="10"
							class="w-full px-3 py-2 bg-[var(--bg-primary)] border {formErrors.maxUnits
								? 'border-red-500'
								: 'border-[var(--border-color)]'} rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)]"
						/>
						{#if formErrors.maxUnits}
							<p class="text-red-500 text-sm mt-1">{formErrors.maxUnits}</p>
						{/if}
					</div>

					<!-- Method -->
					<div>
						<label for="method" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
							Sizing Method
						</label>
						<select
							id="method"
							bind:value={method}
							onchange={() => {
								result = null;
								formErrors = {};
							}}
							class="w-full px-3 py-2 bg-[var(--bg-primary)] border border-[var(--border-color)] rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)]"
						>
							<option value="stock">Stock/ETF</option>
							<option value="opt-delta-atr">Options (Delta-ATR)</option>
							<option value="opt-maxloss">Options (Max Loss)</option>
						</select>
					</div>
				</div>

				<!-- Options-specific fields -->
				{#if method === 'opt-delta-atr'}
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t border-[var(--border-color)]">
						<div>
							<label for="delta" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
								Option Delta <span class="text-red-500">*</span>
							</label>
							<input
								id="delta"
								type="number"
								bind:value={delta}
								oninput={() => handleFieldChange('delta', delta || 0)}
								placeholder="0.50"
								step="0.01"
								min="0"
								max="1"
								class="w-full px-3 py-2 bg-[var(--bg-primary)] border {formErrors.delta
									? 'border-red-500'
									: 'border-[var(--border-color)]'} rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)]"
							/>
							{#if formErrors.delta}
								<p class="text-red-500 text-sm mt-1">{formErrors.delta}</p>
							{/if}
						</div>
					</div>
				{/if}

				{#if method === 'opt-maxloss'}
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t border-[var(--border-color)]">
						<div>
							<label for="maxLoss" class="block text-sm font-medium text-[var(--text-secondary)] mb-1">
								Max Loss per Contract <span class="text-red-500">*</span>
							</label>
							<input
								id="maxLoss"
								type="number"
								bind:value={maxLoss}
								oninput={() => handleFieldChange('maxLoss', maxLoss || 0)}
								placeholder="500"
								step="1"
								class="w-full px-3 py-2 bg-[var(--bg-primary)] border {formErrors.maxLoss
									? 'border-red-500'
									: 'border-[var(--border-color)]'} rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500 text-[var(--text-primary)]"
							/>
							{#if formErrors.maxLoss}
								<p class="text-red-500 text-sm mt-1">{formErrors.maxLoss}</p>
							{/if}
						</div>
					</div>
				{/if}

				<!-- Calculate Button -->
				<div class="flex gap-4 pt-4 border-t border-[var(--border-color)]">
					<button
						onclick={handleCalculate}
						disabled={calculating}
						class="px-6 py-3 bg-gradient-to-r from-emerald-500 to-emerald-600 text-white rounded-lg font-semibold hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
					>
						{#if calculating}
							<div class="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
							Calculating...
						{:else}
							<Calculator class="w-5 h-5" />
							Calculate Position
						{/if}
					</button>

					{#if result}
						<button
							onclick={() => {
								result = null;
								error = '';
							}}
							class="px-6 py-3 bg-[var(--bg-secondary)] border border-[var(--border-color)] text-[var(--text-primary)] rounded-lg font-semibold hover:shadow-lg transition-all"
						>
							Clear
						</button>
					{/if}
				</div>
			</div>
		</Card>

		<!-- Error Display -->
		{#if error}
			<div
				class="bg-gradient-to-r from-red-600 to-red-800 text-white rounded-lg p-4 flex items-center gap-3"
			>
				<AlertTriangle class="w-6 h-6 flex-shrink-0" />
				<div>
					<p class="font-semibold">Calculation Error</p>
					<p class="text-sm">{error}</p>
				</div>
			</div>
		{/if}

		<!-- Results Section -->
		{#if result}
			<div class="space-y-6">
				<!-- Concentration Warning -->
				{#if isHighConcentration()}
					<div
						class="bg-gradient-to-r from-amber-500 to-yellow-400 text-gray-900 rounded-lg p-4 flex items-center gap-3"
					>
						<AlertTriangle class="w-6 h-6 flex-shrink-0" />
						<div>
							<p class="font-semibold">High Concentration Warning</p>
							<p class="text-sm">
								This position represents {getConcentration().toFixed(1)}% of your equity (>${(
									result.shares *
									(entry || 0)
								).toLocaleString()}). Consider reducing size if > 25% of account.
							</p>
						</div>
					</div>
				{/if}

				<!-- Results Cards -->
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					<!-- Risk Dollars -->
					<Card>
						<div class="space-y-2">
							<p class="text-sm text-[var(--text-secondary)]">Risk Dollars</p>
							<p class="text-2xl font-bold text-[var(--text-primary)]">
								${result.risk_dollars.toLocaleString()}
							</p>
							<p class="text-xs text-[var(--text-tertiary)]">
								{((result.risk_dollars / settings.equity) * 100).toFixed(2)}% of equity
							</p>
						</div>
					</Card>

					<!-- Shares/Contracts -->
					<Card>
						<div class="space-y-2">
							<p class="text-sm text-[var(--text-secondary)]">
								{method === 'stock' ? 'Shares' : 'Contracts'}
							</p>
							<p class="text-2xl font-bold text-[var(--text-primary)]">
								{method === 'stock' ? result.shares.toLocaleString() : result.contracts}
							</p>
							{#if method === 'stock' && entry}
								<p class="text-xs text-[var(--text-tertiary)]">
									Position Value: ${(result.shares * entry).toLocaleString()}
								</p>
							{/if}
						</div>
					</Card>

					<!-- Initial Stop -->
					<Card>
						<div class="space-y-2">
							<p class="text-sm text-[var(--text-secondary)]">Initial Stop</p>
							<p class="text-2xl font-bold text-[var(--text-primary)]">
								${result.initial_stop.toFixed(2)}
							</p>
							<p class="text-xs text-[var(--text-tertiary)]">
								Stop Distance: ${result.stop_distance.toFixed(2)} ({k}×N)
							</p>
						</div>
					</Card>
				</div>

				<!-- Add-On Schedule Table -->
				<Card>
					<div class="space-y-4">
						<h2 class="text-xl font-semibold text-[var(--text-primary)] flex items-center gap-2">
							<TrendingUp class="w-6 h-6" />
							Add-On Schedule (0.5N Pyramid)
						</h2>
						<p class="text-sm text-[var(--text-secondary)]">
							Add units every 0.5×N = ${atrN ? (atrN * 0.5).toFixed(2) : '0.00'} above entry
						</p>

						<div class="overflow-x-auto">
							<table class="w-full text-sm">
								<thead>
									<tr class="border-b border-[var(--border-color)]">
										<th class="text-left py-2 px-4 text-[var(--text-secondary)] font-medium"
											>Unit</th
										>
										<th class="text-left py-2 px-4 text-[var(--text-secondary)] font-medium"
											>Entry Price</th
										>
										<th class="text-left py-2 px-4 text-[var(--text-secondary)] font-medium"
											>{method === 'stock' ? 'Shares' : 'Contracts'}</th
										>
										<th class="text-left py-2 px-4 text-[var(--text-secondary)] font-medium"
											>Cumulative {method === 'stock' ? 'Shares' : 'Contracts'}</th
										>
										<th class="text-left py-2 px-4 text-[var(--text-secondary)] font-medium"
											>Cumulative Risk</th
										>
									</tr>
								</thead>
								<tbody>
									{#each getAddOnLevels() as level, i}
										<tr
											class="border-b border-[var(--border-color)] hover:bg-[var(--bg-secondary)] transition-colors"
										>
											<td class="py-2 px-4 text-[var(--text-primary)] font-semibold"
												>Unit {level.unit}</td
											>
											<td class="py-2 px-4 text-[var(--text-primary)]"
												>${level.price.toFixed(2)}</td
											>
											<td class="py-2 px-4 text-[var(--text-primary)]">
												{method === 'stock' ? result.shares.toLocaleString() : result.contracts}
											</td>
											<td class="py-2 px-4 text-[var(--text-primary)] font-semibold">
												{method === 'stock'
													? (result.shares * (i + 1)).toLocaleString()
													: result.contracts * (i + 1)}
											</td>
											<td class="py-2 px-4 text-[var(--text-primary)] font-semibold">
												${(result.risk_dollars * (i + 1)).toLocaleString()}
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>
				</Card>

				<!-- Pyramid Visualization -->
				<Card>
					<div class="space-y-4">
						<h2 class="text-xl font-semibold text-[var(--text-primary)] flex items-center gap-2">
							<BarChart3 class="w-6 h-6" />
							Pyramid Visualization
						</h2>

						<div class="space-y-2">
							{#each getAddOnLevels() as level, i}
								{@const barWidth = ((maxUnits - i) / maxUnits) * 100}
								<div class="flex items-center gap-3">
									<span class="text-sm font-semibold text-[var(--text-primary)] w-16"
										>Unit {level.unit}</span
									>
									<div class="flex-1 bg-[var(--bg-secondary)] rounded-full h-8 relative overflow-hidden">
										<div
											class="bg-gradient-to-r from-emerald-500 to-emerald-600 h-full rounded-full flex items-center justify-end pr-3 text-white text-xs font-semibold transition-all duration-500"
											style="width: {barWidth}%"
										>
											${level.price.toFixed(2)}
										</div>
									</div>
									<span class="text-xs text-[var(--text-secondary)] w-20">
										{method === 'stock' ? result.shares : result.contracts} {method === 'stock'
											? 'shares'
											: 'contracts'}
									</span>
								</div>
							{/each}
						</div>

						<div class="pt-4 border-t border-[var(--border-color)] text-sm text-[var(--text-secondary)]">
							<p>
								Total Max Exposure: ${maxUnits > 0
									? (result.risk_dollars * maxUnits).toLocaleString()
									: '0'} ({maxUnits > 0
									? (((result.risk_dollars * maxUnits) / settings.equity) * 100).toFixed(2)
									: '0.00'}% of equity)
							</p>
						</div>
					</div>
				</Card>

				<!-- Proceed to Heat Check Button -->
				{#if $workflow.sizingComplete}
					<Card>
						<div class="space-y-4">
							<button
								onclick={proceedToHeatCheck}
								class="w-full px-6 py-4 rounded-lg font-semibold text-white text-lg
									   transition-all duration-300 ease-in-out transform
									   bg-gradient-to-r from-blue-500 to-blue-600
									   hover:from-blue-600 hover:to-blue-700
									   shadow-lg hover:shadow-xl active:scale-95
									   flex items-center justify-center gap-2"
							>
								<span>Proceed to Heat Check</span>
								<ArrowRight class="w-5 h-5" />
							</button>
							<p class="text-sm text-blue-600 dark:text-blue-400 text-center font-medium">
								Next: Verify portfolio and sector heat caps are not exceeded
							</p>
						</div>
					</Card>
				{/if}
			</div>
		{/if}
	{/if}
</div>
