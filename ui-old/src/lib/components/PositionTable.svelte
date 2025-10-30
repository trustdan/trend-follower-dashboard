<script lang="ts">
	import Badge from './Badge.svelte';
	import TradingViewLink from './TradingViewLink.svelte';
	import type { Position } from '$lib/api/client';

	interface Props {
		positions: Position[];
	}

	let { positions }: Props = $props();

	function formatCurrency(value: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 2
		}).format(value);
	}

	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		const now = new Date();
		const diffDays = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24));
		return `${diffDays} day${diffDays !== 1 ? 's' : ''}`;
	}
</script>

{#if positions.length === 0}
	<div class="text-center py-8 text-[var(--text-tertiary)]">
		<p>No open positions</p>
	</div>
{:else}
	<div class="overflow-x-auto">
		<table class="w-full">
			<thead>
				<tr class="border-b border-[var(--border-color)]">
					<th class="text-left py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Ticker</th
					>
					<th class="text-left py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Bucket</th
					>
					<th class="text-right py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Entry</th
					>
					<th class="text-right py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Stop</th
					>
					<th class="text-right py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Shares</th
					>
					<th class="text-right py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Risk</th
					>
					<th class="text-left py-3 px-4 text-sm font-semibold text-[var(--text-primary)]">Days</th
					>
					<th class="text-left py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Status</th
					>
					<th class="text-center py-3 px-4 text-sm font-semibold text-[var(--text-primary)]"
						>Chart</th
					>
				</tr>
			</thead>
			<tbody>
				{#each positions as position}
					<tr
						class="border-b border-[var(--border-color)] hover:bg-[var(--bg-tertiary)] transition-colors"
					>
						<td class="py-3 px-4 text-sm font-medium text-[var(--text-primary)]"
							>{position.ticker}</td
						>
						<td class="py-3 px-4 text-sm text-[var(--text-secondary)]"
							>{position.bucket || '—'}</td
						>
						<td class="py-3 px-4 text-sm text-right text-[var(--text-primary)]"
							>{formatCurrency(position.entry_price)}</td
						>
						<td class="py-3 px-4 text-sm text-right text-[var(--text-primary)]"
							>{formatCurrency(position.current_stop)}</td
						>
						<td class="py-3 px-4 text-sm text-right text-[var(--text-primary)]"
							>{position.shares.toLocaleString()}</td
						>
						<td class="py-3 px-4 text-sm text-right text-[var(--text-primary)]"
							>{formatCurrency(position.risk_dollars)}</td
						>
						<td class="py-3 px-4 text-sm text-[var(--text-secondary)]"
							>{formatDate(position.opened_at)}</td
						>
						<td class="py-3 px-4">
							<Badge variant={position.status === 'OPEN' ? 'green' : 'gray'}>
								{position.status}
							</Badge>
						</td>
						<td class="py-3 px-4 text-center">
							<TradingViewLink ticker={position.ticker} variant="icon" />
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
