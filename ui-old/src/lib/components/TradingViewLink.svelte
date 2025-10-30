<script lang="ts">
	/**
	 * TradingViewLink Component
	 *
	 * Provides quick access to TradingView charts for signal verification.
	 * Clicking opens a new browser tab with the ticker's chart.
	 *
	 * Anti-Impulsivity Design:
	 * - Manual verification required (no automated signal checking)
	 * - Forces trader to visually confirm 55-bar breakout
	 * - Trader must manually note N (ATR) value from chart
	 */

	export let ticker: string;
	export let variant: 'button' | 'icon' | 'text' = 'icon';

	// Default TradingView URL template
	// TODO: Allow customization via settings store
	const urlTemplate = 'https://www.tradingview.com/chart/?symbol={ticker}';

	function openChart() {
		if (!ticker) return;

		const url = urlTemplate.replace('{ticker}', ticker.toUpperCase());
		window.open(url, '_blank', 'noopener,noreferrer');
	}
</script>

{#if variant === 'button'}
	<button class="tv-button" on:click={openChart} title="Open {ticker} in TradingView">
		<svg class="tv-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
			<rect x="3" y="3" width="18" height="18" stroke-width="2" rx="2"/>
			<polyline points="9 10 12 7 15 10" stroke-width="2"/>
			<line x1="12" y1="7" x2="12" y2="17" stroke-width="2"/>
		</svg>
		<span>Open in TradingView</span>
	</button>
{:else if variant === 'icon'}
	<button class="tv-icon-only" on:click={openChart} title="Open {ticker} in TradingView">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
			<rect x="3" y="3" width="18" height="18" stroke-width="2" rx="2"/>
			<polyline points="9 10 12 7 15 10" stroke-width="2"/>
			<line x1="12" y1="7" x2="12" y2="17" stroke-width="2"/>
		</svg>
	</button>
{:else}
	<button class="tv-text-link" on:click={openChart}>
		{ticker} ↗
	</button>
{/if}

<style>
	.tv-button {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
		color: white;
		border: none;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		box-shadow: 0 2px 4px rgba(59, 130, 246, 0.2);
	}

	.tv-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(59, 130, 246, 0.3);
	}

	.tv-button:active {
		transform: translateY(0);
		box-shadow: 0 2px 4px rgba(59, 130, 246, 0.2);
	}

	.tv-icon {
		width: 1rem;
		height: 1rem;
	}

	.tv-icon-only {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 2rem;
		height: 2rem;
		padding: 0.375rem;
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: 0.375rem;
		cursor: pointer;
		transition: all 0.15s ease;
		color: var(--color-text-secondary);
	}

	.tv-icon-only:hover {
		background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
		color: white;
		border-color: transparent;
		box-shadow: 0 2px 4px rgba(59, 130, 246, 0.2);
	}

	.tv-icon-only svg {
		width: 1.25rem;
		height: 1.25rem;
	}

	.tv-text-link {
		background: none;
		border: none;
		color: #3b82f6;
		font-weight: 500;
		cursor: pointer;
		transition: color 0.15s ease;
		padding: 0;
		text-decoration: none;
	}

	.tv-text-link:hover {
		color: #2563eb;
		text-decoration: underline;
	}
</style>
