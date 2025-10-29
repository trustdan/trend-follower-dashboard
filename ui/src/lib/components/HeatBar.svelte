<script lang="ts">
	import { Flame } from 'lucide-svelte';

	// Props
	let {
		label = 'Heat',
		current = 0,
		max = 100,
		exceeded = false,
		showIcon = true,
	}: {
		label?: string;
		current?: number;
		max?: number;
		exceeded?: boolean;
		showIcon?: boolean;
	} = $props();

	// Calculate percentage (capped at 100% for display)
	const percentage = $derived(max > 0 ? Math.min((current / max) * 100, 100) : 0);

	// Determine color based on usage
	const getColorClasses = (pct: number, isExceeded: boolean) => {
		if (isExceeded) {
			return {
				bg: 'from-red-500 to-red-600',
				text: 'text-red-700',
				border: 'border-red-500',
			};
		}
		if (pct >= 90) {
			return {
				bg: 'from-amber-500 to-yellow-500',
				text: 'text-amber-700',
				border: 'border-amber-500',
			};
		}
		if (pct >= 70) {
			return {
				bg: 'from-yellow-400 to-yellow-500',
				text: 'text-yellow-700',
				border: 'border-yellow-500',
			};
		}
		return {
			bg: 'from-emerald-500 to-emerald-600',
			text: 'text-emerald-700',
			border: 'border-emerald-500',
		};
	};

	const colors = $derived(getColorClasses(percentage, exceeded));

	// Format currency
	const formatCurrency = (value: number): string => {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 0,
			maximumFractionDigits: 0,
		}).format(value);
	};
</script>

<div class="space-y-2">
	<!-- Label and values -->
	<div class="flex items-center justify-between text-sm">
		<div class="flex items-center gap-2">
			{#if showIcon}
				<Flame class="w-4 h-4 {colors.text}" />
			{/if}
			<span class="font-medium text-gray-700 dark:text-gray-300">{label}</span>
		</div>
		<div class="flex items-center gap-2">
			<span class="font-semibold {colors.text}">
				{formatCurrency(current)} / {formatCurrency(max)}
			</span>
			<span class="text-xs text-gray-500 dark:text-gray-400">
				({percentage.toFixed(1)}%)
			</span>
		</div>
	</div>

	<!-- Progress bar container -->
	<div class="relative w-full h-8 bg-gray-200 dark:bg-gray-700 rounded-lg overflow-hidden border-2 {colors.border}">
		<!-- Filled portion -->
		<div
			class="h-full bg-gradient-to-r {colors.bg} transition-all duration-500 ease-out flex items-center justify-center"
			style="width: {percentage}%"
		>
			{#if percentage > 15}
				<span class="text-xs font-bold text-white drop-shadow">
					{percentage.toFixed(1)}%
				</span>
			{/if}
		</div>

		<!-- Exceeded overlay -->
		{#if exceeded}
			<div class="absolute inset-0 flex items-center justify-center bg-red-500/20 backdrop-blur-[2px]">
				<span class="text-xs font-bold text-red-900 dark:text-red-100 drop-shadow px-2 py-1 bg-red-200/80 dark:bg-red-800/80 rounded">
					EXCEEDED
				</span>
			</div>
		{/if}
	</div>

	<!-- Status message -->
	{#if exceeded}
		<div class="flex items-center gap-2 text-xs {colors.text}">
			<span class="font-semibold">⚠️ Over limit by {formatCurrency(current - max)}</span>
		</div>
	{:else if percentage >= 90}
		<div class="flex items-center gap-2 text-xs {colors.text}">
			<span class="font-semibold">⚠️ Approaching limit - {formatCurrency(max - current)} remaining</span>
		</div>
	{:else if percentage >= 70}
		<div class="flex items-center gap-2 text-xs {colors.text}">
			<span class="font-medium">Moderate usage - {formatCurrency(max - current)} available</span>
		</div>
	{:else}
		<div class="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-400">
			<span>{formatCurrency(max - current)} available capacity</span>
		</div>
	{/if}
</div>
