<script lang="ts">
	import { logger } from '$lib/utils/logger';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api, type HeatCheckResult, type Position } from '$lib/api/client';
	import HeatBar from '$lib/components/HeatBar.svelte';
	import Card from '$lib/components/Card.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import WorkflowProgress from '$lib/components/WorkflowProgress.svelte';
	import { workflow } from '$lib/stores/workflow';
	import { Flame, AlertTriangle, CheckCircle, Calculator, ArrowRight } from 'lucide-svelte';

	// State
	let loading = $state(true);
	let calculating = $state(false);
	let error = $state<string | null>(null);

	// Settings
	let equity = $state(0);
	let portfolioCap = $state(0);
	let bucketCap = $state(0);

	// Current heat (before proposed trade)
	let currentPortfolioHeat = $state(0);
	let currentBucketHeat = $state(0);

	// Form inputs
	let proposedRisk = $state<number | null>(null);
	let selectedBucket = $state('Tech/Comm');

	// Heat check result
	let heatResult = $state<HeatCheckResult | null>(null);

	// Available buckets (from CLAUDE.md: 10 predefined sectors)
	const buckets = [
		'Tech/Comm',
		'Finance',
		'Healthcare',
		'Consumer Cyclical',
		'Consumer Defensive',
		'Energy',
		'Industrials',
		'Materials',
		'Real Estate',
		'Utilities',
	];

	// Load initial data
	async function loadData() {
		try {
			loading = true;
			error = null;

			// Get settings and positions in parallel
			const [settingsData, positions] = await Promise.all([
				api.getSettings(),
				api.getPositions(),
			]);

			logger.info('Heat check loaded', { settings: settingsData, positions });

			// Settings
			equity = settingsData.equity;
			portfolioCap = settingsData.portfolioCap;
			bucketCap = settingsData.bucketCap;

			// Calculate current heat from positions
			currentPortfolioHeat = positions.reduce((sum: number, p: Position) => sum + p.risk_dollars, 0);

			logger.info('Initial heat calculated', {
				currentPortfolioHeat,
				portfolioCap,
				bucketCap,
			});

		} catch (err) {
			logger.error('Failed to load heat check data', err);
			error = err instanceof Error ? err.message : 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	// Calculate heat with proposed trade
	async function checkHeat() {
		if (proposedRisk === null || proposedRisk <= 0) {
			error = 'Please enter a valid risk amount';
			return;
		}

		try {
			calculating = true;
			error = null;

			logger.info('Checking heat', { proposedRisk, selectedBucket });

			const result = await api.checkHeat({
				add_risk_dollars: proposedRisk,
				add_bucket: selectedBucket,
			});

			heatResult = result;

			logger.info('Heat check result', result);

			// Save heat check results to workflow
			const passed = result.result.approved;
			const portfolioHeat = result.result.portfolio_heat_after;
			const bucketHeat = result.result.bucket_heat_after;

			workflow.saveHeatResults(passed, portfolioHeat, bucketHeat);
			logger.info('Heat results saved to workflow', { passed, portfolioHeat, bucketHeat });

		} catch (err) {
			logger.error('Heat check failed', err);
			error = err instanceof Error ? err.message : 'Heat check failed';
		} finally {
			calculating = false;
		}
	}

	// Clear result
	function clearResult() {
		heatResult = null;
		proposedRisk = null;
		error = null;
		logger.info('Heat check result cleared');
	}

	// Load data on mount
	// Navigate to trade entry
	function proceedToEntry() {
		if ($workflow.heatCheckComplete && $workflow.heatCheckPassed) {
			logger.info('Proceeding to trade entry');
			goto('/entry');
		}
	}

	onMount(() => {
		logger.info('Heat Check page loaded');
		loadData();

		// Pre-fill from workflow if available
		if ($workflow.riskDollars) {
			proposedRisk = $workflow.riskDollars;
			logger.info('Pre-filling risk from workflow', proposedRisk);
		}
		if ($workflow.sector) {
			selectedBucket = $workflow.sector;
			logger.info('Pre-filling bucket from workflow', selectedBucket);
		}

		// Update workflow step
		if ($workflow.workflowStarted && $workflow.currentStep !== 'heat') {
			workflow.goToStep('heat');
		}
	});
</script>

<svelte:head><title>Heat Check - TF-Engine</title></svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div>
		<h1 class="text-3xl font-bold text-[var(--text-primary)] flex items-center gap-3">
			<Flame class="w-8 h-8 text-orange-500" />
			Heat Check
		</h1>
		<p class="text-[var(--text-secondary)] mt-2">
			Portfolio and sector heat management - verify new trades don't exceed caps
		</p>
	</div>

	<!-- Loading state -->
	{#if loading}
		<LoadingSpinner message="Loading heat data..." />
	{/if}

	<!-- Error state -->
	{#if error && !loading}
		<Card>
			<div class="text-center py-8">
				<p class="text-red-600 dark:text-red-400 mb-4">⚠️ {error}</p>
				<button
					onclick={loadData}
					class="px-4 py-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg hover:shadow-lg transition-all"
				>
					Retry
				</button>
			</div>
		</Card>
	{/if}

	<!-- Main content -->
	{#if !loading && !error}
		<!-- Current Heat Status -->
		<Card title="Current Portfolio Heat">
			<div class="space-y-6">
				<HeatBar
					label="Portfolio Heat (4% Cap)"
					current={currentPortfolioHeat}
					max={portfolioCap}
					exceeded={currentPortfolioHeat > portfolioCap}
				/>

				<div class="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4 border-t border-gray-200 dark:border-gray-700">
					<div class="text-center">
						<p class="text-sm text-gray-600 dark:text-gray-400">Current Heat</p>
						<p class="text-2xl font-bold text-gray-900 dark:text-gray-100">
							${currentPortfolioHeat.toLocaleString()}
						</p>
					</div>
					<div class="text-center">
						<p class="text-sm text-gray-600 dark:text-gray-400">Portfolio Cap</p>
						<p class="text-2xl font-bold text-gray-900 dark:text-gray-100">
							${portfolioCap.toLocaleString()}
						</p>
					</div>
					<div class="text-center">
						<p class="text-sm text-gray-600 dark:text-gray-400">Available</p>
						<p class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">
							${(portfolioCap - currentPortfolioHeat).toLocaleString()}
						</p>
					</div>
				</div>
			</div>
		</Card>

		<!-- Proposed Trade Form -->
		<Card>
			<div class="space-y-6">
				<div class="flex items-center gap-3">
					<Calculator class="w-6 h-6 text-blue-600 dark:text-blue-400" />
					<h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">Check Proposed Trade</h2>
				</div>

				<p class="text-sm text-gray-600 dark:text-gray-400">
					Enter your proposed trade details to verify it won't exceed heat caps
				</p>

				<!-- Input form -->
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<!-- Risk Amount -->
					<div>
						<label for="risk" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
							Proposed Risk Amount ($)
						</label>
						<input
							id="risk"
							type="number"
							bind:value={proposedRisk}
							placeholder="750"
							min="0"
							step="10"
							class="w-full px-4 py-2 rounded-lg border-2 border-gray-300 dark:border-gray-600
								bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100
								focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent
								transition-all"
						/>
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
							From position sizing calculator
						</p>
					</div>

					<!-- Bucket Selector -->
					<div>
						<label for="bucket" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
							Sector/Bucket
						</label>
						<select
							id="bucket"
							bind:value={selectedBucket}
							class="w-full px-4 py-2 rounded-lg border-2 border-gray-300 dark:border-gray-600
								bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100
								focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent
								transition-all"
						>
							{#each buckets as bucket}
								<option value={bucket}>{bucket}</option>
							{/each}
						</select>
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
							Sector for proposed trade
						</p>
					</div>
				</div>

				<!-- Action buttons -->
				<div class="flex gap-3">
					<button
						onclick={checkHeat}
						disabled={calculating || proposedRisk === null || proposedRisk <= 0}
						class="flex-1 px-6 py-3 bg-gradient-to-r from-blue-500 to-blue-600 text-white
							rounded-lg font-semibold hover:shadow-lg transition-all
							disabled:opacity-50 disabled:cursor-not-allowed
							flex items-center justify-center gap-2"
					>
						{#if calculating}
							<div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
							Checking...
						{:else}
							<Calculator class="w-5 h-5" />
							Check Heat
						{/if}
					</button>

					{#if heatResult}
						<button
							onclick={clearResult}
							class="px-6 py-3 bg-gray-500 text-white rounded-lg font-semibold
								hover:bg-gray-600 transition-all"
						>
							Clear
						</button>
					{/if}
				</div>
			</div>
		</Card>

		<!-- Heat Check Result -->
		{#if heatResult}
			<div class="space-y-4">
				<!-- Approval/Rejection Banner -->
				{#if heatResult.allowed}
					<div class="bg-gradient-to-r from-emerald-500 to-emerald-600 rounded-lg p-6 shadow-lg">
						<div class="flex items-center gap-4">
							<CheckCircle class="w-12 h-12 text-white" />
							<div class="flex-1">
								<h3 class="text-2xl font-bold text-white mb-1">✓ TRADE APPROVED</h3>
								<p class="text-emerald-100">
									Both portfolio and bucket heat are within caps. Trade is allowed.
								</p>
							</div>
						</div>
					</div>
				{:else}
					<div class="bg-gradient-to-r from-red-600 to-red-800 rounded-lg p-6 shadow-lg">
						<div class="flex items-center gap-4">
							<AlertTriangle class="w-12 h-12 text-white" />
							<div class="flex-1">
								<h3 class="text-2xl font-bold text-white mb-1">✗ TRADE REJECTED</h3>
								<p class="text-red-100">{heatResult.rejection_reason}</p>
							</div>
						</div>
					</div>
				{/if}

				<!-- Detailed Heat Breakdown -->
				<Card title="Heat Analysis">
					<div class="space-y-6">
						<!-- Portfolio Heat -->
						<div>
							<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
								Portfolio Heat (4% Cap)
							</h3>
							<HeatBar
								label="With Proposed Trade"
								current={heatResult.new_portfolio_heat}
								max={heatResult.portfolio_cap}
								exceeded={heatResult.portfolio_cap_exceeded}
							/>
							<div class="mt-4 text-sm text-gray-600 dark:text-gray-400">
								Current: ${heatResult.current_portfolio_heat.toLocaleString()}
								→ New: ${heatResult.new_portfolio_heat.toLocaleString()}
								(+${proposedRisk?.toLocaleString()})
							</div>
						</div>

						<!-- Bucket Heat -->
						<div>
							<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
								{selectedBucket} Bucket Heat (1.5% Cap)
							</h3>
							<HeatBar
								label={`${selectedBucket} Bucket`}
								current={heatResult.new_bucket_heat}
								max={heatResult.bucket_cap}
								exceeded={heatResult.bucket_cap_exceeded}
								showIcon={false}
							/>
							<div class="mt-4 text-sm text-gray-600 dark:text-gray-400">
								Current: ${heatResult.current_bucket_heat.toLocaleString()}
								→ New: ${heatResult.new_bucket_heat.toLocaleString()}
								(+${proposedRisk?.toLocaleString()})
							</div>
						</div>

						<!-- Summary Stats -->
						<div class="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t border-gray-200 dark:border-gray-700">
							<div class="text-center">
								<p class="text-xs text-gray-600 dark:text-gray-400">Portfolio %</p>
								<p class="text-lg font-bold" class:text-red-600={heatResult.portfolio_cap_exceeded}>
									{heatResult.portfolio_heat_pct.toFixed(1)}%
								</p>
							</div>
							<div class="text-center">
								<p class="text-xs text-gray-600 dark:text-gray-400">Bucket %</p>
								<p class="text-lg font-bold" class:text-red-600={heatResult.bucket_cap_exceeded}>
									{heatResult.bucket_heat_pct.toFixed(1)}%
								</p>
							</div>
							<div class="text-center">
								<p class="text-xs text-gray-600 dark:text-gray-400">Portfolio Margin</p>
								<p class="text-lg font-bold">
									${(heatResult.portfolio_cap - heatResult.new_portfolio_heat).toLocaleString()}
								</p>
							</div>
							<div class="text-center">
								<p class="text-xs text-gray-600 dark:text-gray-400">Bucket Margin</p>
								<p class="text-lg font-bold">
									${(heatResult.bucket_cap - heatResult.new_bucket_heat).toLocaleString()}
								</p>
							</div>
						</div>
					</div>
				</Card>
			</div>
		{/if}
	{/if}
</div>
