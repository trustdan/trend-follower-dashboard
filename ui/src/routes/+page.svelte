<script lang="ts">
	import { logger } from '$lib/utils/logger';
	import { api, type Settings, type Position, type Candidate } from '$lib/api/client';
	import { onMount } from 'svelte';
	import Card from '$lib/components/Card.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import PositionTable from '$lib/components/PositionTable.svelte';
	import TradingViewLink from '$lib/components/TradingViewLink.svelte';
	import { TrendingUp, Flame, AlertCircle, ScanSearch, ListChecks } from 'lucide-svelte';

	let settings: Settings | null = $state(null);
	let positions: Position[] = $state([]);
	let candidates: Candidate[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	async function loadDashboardData() {
		try {
			loading = true;
			error = null;

			// Fetch all data in parallel
			const [settingsData, positionsData, candidatesData] = await Promise.all([
				api.getSettings(),
				api.getPositions(),
				api.getCandidates()
			]);

			settings = settingsData;
			positions = positionsData;
			candidates = candidatesData;

			logger.info('Dashboard data loaded', {
				positions: positions.length,
				candidates: candidates.length
			});
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load dashboard data';
			logger.error('Dashboard load failed', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		logger.info('Dashboard page loaded');
		loadDashboardData();
	});

	// Calculate portfolio heat
	let portfolioHeat = $state(0);
	let portfolioCap = $state(0);
	let heatPercentage = $state(0);
	let availableCapacity = $state(0);

	$effect(() => {
		portfolioHeat = positions.reduce((sum, p) => sum + p.risk_dollars, 0);
		portfolioCap = settings ? settings.equity * settings.portfolioCap : 0;
		heatPercentage = portfolioCap > 0 ? (portfolioHeat / portfolioCap) * 100 : 0;
		availableCapacity = portfolioCap - portfolioHeat;
	});
</script>

<svelte:head>
	<title>Dashboard - TF-Engine</title>
</svelte:head>

<div class="space-y-6">
	<!-- Page Header -->
	<div>
		<h1 class="text-3xl font-bold text-[var(--text-primary)]">Dashboard</h1>
		<p class="text-[var(--text-secondary)] mt-2">Portfolio overview and quick actions</p>
	</div>

	{#if loading}
		<div class="py-12">
			<LoadingSpinner size="lg" message="Loading dashboard data..." />
		</div>
	{:else if error}
		<Card>
			<div class="text-center py-8">
				<AlertCircle size={48} class="mx-auto mb-4 text-red-500" />
				<p class="text-lg font-semibold text-[var(--text-primary)] mb-2">Failed to Load Data</p>
				<p class="text-sm text-[var(--text-secondary)] mb-4">{error}</p>
				<button
					onclick={loadDashboardData}
					class="px-4 py-2 bg-[var(--border-focus)] text-white rounded-lg hover:opacity-90 transition-opacity"
				>
					Retry
				</button>
			</div>
		</Card>
	{:else if settings}
		<!-- Portfolio Summary Cards -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<!-- Equity Card -->
			<Card>
				<div class="flex items-start justify-between">
					<div>
						<p class="text-sm text-[var(--text-secondary)] mb-1">Account Equity</p>
						<p class="text-2xl font-bold text-[var(--text-primary)]">
							${settings.equity.toLocaleString()}
						</p>
						<p class="text-xs text-[var(--text-tertiary)] mt-1">
							Risk per trade: {(settings.riskPct * 100).toFixed(2)}%
						</p>
					</div>
					<div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
						<TrendingUp size={24} class="text-blue-600 dark:text-blue-400" />
					</div>
				</div>
			</Card>

			<!-- Portfolio Heat Card -->
			<Card>
				<div class="flex items-start justify-between">
					<div class="flex-1">
						<p class="text-sm text-[var(--text-secondary)] mb-1">Portfolio Heat</p>
						<p class="text-2xl font-bold text-[var(--text-primary)]">
							${portfolioHeat.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 })}
						</p>
						<div class="mt-2">
							<div class="flex items-center gap-2">
								<div class="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
									<div
										class="h-full {heatPercentage >= 90
											? 'bg-red-500'
											: heatPercentage >= 70
												? 'bg-yellow-500'
												: 'bg-green-500'} transition-all duration-300"
										style="width: {Math.min(heatPercentage, 100)}%"
									></div>
								</div>
								<span class="text-xs text-[var(--text-tertiary)]"
									>{heatPercentage.toFixed(0)}%</span
								>
							</div>
							<p class="text-xs text-[var(--text-tertiary)] mt-1">
								Cap: ${portfolioCap.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 })}
							</p>
						</div>
					</div>
					<div
						class="p-2 {heatPercentage >= 90
							? 'bg-red-100 dark:bg-red-900/30'
							: 'bg-orange-100 dark:bg-orange-900/30'} rounded-lg"
					>
						<Flame
							size={24}
							class="{heatPercentage >= 90
								? 'text-red-600 dark:text-red-400'
								: 'text-orange-600 dark:text-orange-400'}"
						/>
					</div>
				</div>
			</Card>

			<!-- Available Capacity Card -->
			<Card>
				<div class="flex items-start justify-between">
					<div>
						<p class="text-sm text-[var(--text-secondary)] mb-1">Available Capacity</p>
						<p class="text-2xl font-bold text-[var(--text-primary)]">
							${availableCapacity.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 })}
						</p>
						<p class="text-xs text-[var(--text-tertiary)] mt-1">
							{positions.length} open position{positions.length !== 1 ? 's' : ''}
						</p>
					</div>
					<div class="flex flex-col gap-1">
						<Badge variant={availableCapacity > 0 ? 'green' : 'red'}>
							{availableCapacity > 0 ? 'Room' : 'Full'}
						</Badge>
					</div>
				</div>
			</Card>
		</div>

		<!-- Quick Actions -->
		<Card title="Quick Actions">
			<div class="flex gap-4">
				<a
					href="/scanner"
					class="flex-1 flex items-center justify-center gap-2 px-6 py-4 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg hover:from-blue-600 hover:to-blue-700 transition-all duration-200 shadow-md hover:shadow-lg"
				>
					<ScanSearch size={20} />
					<span class="font-medium">Run FINVIZ Scan</span>
				</a>
				<a
					href="/checklist"
					class="flex-1 flex items-center justify-center gap-2 px-6 py-4 bg-gradient-to-r from-green-500 to-green-600 text-white rounded-lg hover:from-green-600 hover:to-green-700 transition-all duration-200 shadow-md hover:shadow-lg"
				>
					<ListChecks size={20} />
					<span class="font-medium">Start Checklist</span>
				</a>
			</div>
		</Card>

		<!-- Open Positions -->
		<Card title="Open Positions" subtitle="{positions.length} active trade{positions.length !== 1 ? 's' : ''}">
			<PositionTable {positions} />
		</Card>

		<!-- Today's Candidates -->
		<Card title="Today's Candidates" subtitle="{candidates.length} ticker{candidates.length !== 1 ? 's' : ''} imported">
			{#if candidates.length === 0}
				<div class="text-center py-8 text-[var(--text-tertiary)]">
					<p class="mb-2">No candidates for today</p>
					<a
						href="/scanner"
						class="inline-block text-sm text-[var(--border-focus)] hover:underline"
					>
						Run a FINVIZ scan →
					</a>
				</div>
			{:else}
				<div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-2">
					{#each candidates as candidate}
						<div class="candidate-card">
							<div
								class="px-3 py-2 bg-[var(--bg-tertiary)] rounded border border-[var(--border-color)] text-center hover:border-[var(--border-focus)] transition-colors"
							>
								<p class="font-medium text-[var(--text-primary)]">{candidate.ticker}</p>
								{#if candidate.bucket}
									<p class="text-xs text-[var(--text-tertiary)] mt-0.5">{candidate.bucket}</p>
								{/if}
							</div>
							<div class="tv-link-center">
								<TradingViewLink ticker={candidate.ticker} variant="icon" />
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</Card>
	{/if}
</div>

<style>
	.candidate-card {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.tv-link-center {
		display: flex;
		justify-content: center;
	}
</style>
