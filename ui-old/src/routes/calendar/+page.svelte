<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type CalendarData, type WeekData, type PositionInfo } from '$lib/api/client';
	import { logger } from '$lib/utils/logger';
	import Card from '$lib/components/Card.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	// TradingView integration
	function openTradingView(ticker: string) {
		const url = `https://www.tradingview.com/chart/?symbol=${ticker}`;
		window.open(url, '_blank', 'noopener,noreferrer');
	}

	let calendarData: CalendarData | null = null;
	let loading = false;
	let error: string | null = null;

	// Format date for display (e.g., "Oct 28")
	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
	}

	// Format date range (e.g., "Oct 28 - Nov 3")
	function formatWeekRange(weekStart: string, weekEnd: string): string {
		return `${formatDate(weekStart)} - ${formatDate(weekEnd)}`;
	}

	// Get week indicator (e.g., "Past", "This Week", "Future")
	function getWeekIndicator(weekStart: string): string {
		const today = new Date();
		const startDate = new Date(weekStart);
		const endDate = new Date(weekStart);
		endDate.setDate(endDate.getDate() + 6);

		if (today >= startDate && today <= endDate) {
			return 'current';
		} else if (today > endDate) {
			return 'past';
		} else {
			return 'future';
		}
	}

	// Load calendar data
	async function loadCalendar() {
		loading = true;
		error = null;
		try {
			logger.info('Loading calendar data...');
			calendarData = await api.getCalendar();
			logger.success('Calendar data loaded successfully');
			logger.debug('Calendar data:', calendarData);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load calendar';
			logger.error('Failed to load calendar:', err);
		} finally {
			loading = false;
		}
	}

	// Get color for position badge based on days held
	function getPositionColor(daysHeld: number): string {
		if (daysHeld < 7) return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30';
		if (daysHeld < 30) return 'bg-blue-500/10 text-blue-400 border-blue-500/30';
		return 'bg-purple-500/10 text-purple-400 border-purple-500/30';
	}

	onMount(() => {
		logger.info('Calendar page mounted');
		loadCalendar();
	});
</script>

<svelte:head><title>Calendar - TF-Engine</title></svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-[var(--text-primary)]">Calendar</h1>
			<p class="text-[var(--text-secondary)] mt-2">
				Rolling 10-week view (2 weeks back + 8 weeks forward)
			</p>
		</div>
		<button
			on:click={loadCalendar}
			class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
			disabled={loading}
		>
			{loading ? 'Loading...' : 'Refresh'}
		</button>
	</div>

	<!-- Error State -->
	{#if error}
		<Card>
			<div class="text-center py-8">
				<div class="text-red-400 text-lg mb-2">❌ Error Loading Calendar</div>
				<p class="text-[var(--text-tertiary)]">{error}</p>
				<button
					on:click={loadCalendar}
					class="mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
				>
					Try Again
				</button>
			</div>
		</Card>
	{/if}

	<!-- Loading State -->
	{#if loading}
		<Card>
			<div class="flex items-center justify-center py-12">
				<LoadingSpinner />
				<span class="ml-3 text-[var(--text-secondary)]">Loading calendar...</span>
			</div>
		</Card>
	{/if}

	<!-- Calendar Grid -->
	{#if !loading && !error && calendarData}
		<Card>
			<div class="overflow-x-auto">
				<div class="min-w-full">
					<!-- Legend -->
					<div class="mb-4 flex items-center gap-4 text-sm text-[var(--text-tertiary)]">
						<div class="flex items-center gap-2">
							<div class="w-3 h-3 rounded bg-emerald-500/20 border border-emerald-500/30"></div>
							<span>New (&lt; 7 days)</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-3 h-3 rounded bg-blue-500/20 border border-blue-500/30"></div>
							<span>Active (&lt; 30 days)</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-3 h-3 rounded bg-purple-500/20 border border-purple-500/30"></div>
							<span>Mature (30+ days)</span>
						</div>
					</div>

					<!-- Table -->
					<table class="w-full border-collapse">
						<thead>
							<tr class="border-b border-[var(--border-color)]">
								<th
									class="text-left py-3 px-4 text-[var(--text-secondary)] font-semibold sticky left-0 bg-[var(--bg-secondary)] z-10"
								>
									Sector
								</th>
								{#each calendarData.weeks as week}
									{@const indicator = getWeekIndicator(week.week_start)}
									<th
										class="text-center py-3 px-4 text-[var(--text-secondary)] font-semibold min-w-[140px] {indicator ===
										'current'
											? 'bg-blue-500/10'
											: ''}"
									>
										<div class="text-sm">{formatDate(week.week_start)}</div>
										<div class="text-xs text-[var(--text-tertiary)] mt-1">
											{indicator === 'current' ? '← This Week →' : ''}
										</div>
									</th>
								{/each}
							</tr>
						</thead>
						<tbody>
							{#each calendarData.sectors as sector}
								<tr class="border-b border-[var(--border-color)] hover:bg-[var(--bg-tertiary)]">
									<td
										class="py-3 px-4 font-medium text-[var(--text-primary)] sticky left-0 bg-[var(--bg-secondary)] z-10"
									>
										{sector}
									</td>
									{#each calendarData.weeks as week}
										{@const indicator = getWeekIndicator(week.week_start)}
										<td
											class="py-3 px-4 align-top {indicator === 'current'
												? 'bg-blue-500/5'
												: ''}"
										>
											{#if week.sectors[sector] && week.sectors[sector].length > 0}
												<div class="flex flex-col gap-1">
													{#each week.sectors[sector] as position}
														<button
															class="ticker-badge text-xs px-2 py-1 rounded border {getPositionColor(
																position.days_held
															)} font-mono transition-all hover:scale-105 hover:shadow-md cursor-pointer"
															title="Click to open {position.ticker} in TradingView | Entry: ${position.entry_price.toFixed(
																2
															)} | Risk: ${position.risk_dollars.toFixed(2)} | Days: {position.days_held}"
															on:click={() => openTradingView(position.ticker)}
														>
															{position.ticker} ↗
														</button>
													{/each}
												</div>
											{:else}
												<div class="text-center text-[var(--text-tertiary)] text-xs">—</div>
											{/if}
										</td>
									{/each}
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>

			<!-- Summary -->
			{#if calendarData}
				<div class="mt-6 pt-6 border-t border-[var(--border-color)]">
					<div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
						<div>
							<div class="text-[var(--text-tertiary)]">Total Sectors</div>
							<div class="text-2xl font-bold text-[var(--text-primary)] mt-1">
								{calendarData.sectors.length}
							</div>
						</div>
						<div>
							<div class="text-[var(--text-tertiary)]">Time Window</div>
							<div class="text-2xl font-bold text-[var(--text-primary)] mt-1">10 Weeks</div>
						</div>
						<div>
							<div class="text-[var(--text-tertiary)]">View Range</div>
							<div class="text-sm font-medium text-[var(--text-primary)] mt-1">
								{calendarData.weeks[0]
									? formatWeekRange(calendarData.weeks[0].week_start, calendarData.weeks[0].week_end)
									: 'N/A'}
								<br />
								to
								<br />
								{calendarData.weeks[9]
									? formatWeekRange(calendarData.weeks[9].week_start, calendarData.weeks[9].week_end)
									: 'N/A'}
							</div>
						</div>
					</div>
				</div>
			{/if}
		</Card>

		<!-- Help Text -->
		<Card>
			<div class="text-sm text-[var(--text-secondary)]">
				<h3 class="font-semibold text-[var(--text-primary)] mb-2">📅 Calendar Purpose</h3>
				<p class="mb-3">
					The calendar helps you visualize position distribution across sectors and time to:
				</p>
				<ul class="list-disc list-inside space-y-1 ml-2">
					<li><strong>Avoid sector crowding:</strong> Limit exposure to single sector (1.5% bucket cap)</li>
					<li>
						<strong>Monitor diversification:</strong> Spread risk across different sectors and entry dates
					</li>
					<li>
						<strong>Track position aging:</strong> Color-coded badges show how long you've held each
						position
					</li>
					<li>
						<strong>Plan entries:</strong> See which weeks/sectors have capacity for new positions
					</li>
				</ul>
			</div>
		</Card>
	{/if}

	<!-- Empty State -->
	{#if !loading && !error && calendarData && calendarData.weeks.every((w) => Object.keys(w.sectors).length === 0)}
		<Card>
			<div class="text-center py-12">
				<div class="text-4xl mb-4">📭</div>
				<h3 class="text-xl font-semibold text-[var(--text-primary)] mb-2">No Open Positions</h3>
				<p class="text-[var(--text-tertiary)]">
					When you open positions, they'll appear in this calendar view.
				</p>
			</div>
		</Card>
	{/if}
</div>