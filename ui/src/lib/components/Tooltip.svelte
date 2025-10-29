<script lang="ts">
	export let text: string;
	export let position: 'top' | 'bottom' | 'left' | 'right' = 'top';

	let visible = false;
</script>

<div
	class="tooltip-container"
	on:mouseenter={() => (visible = true)}
	on:mouseleave={() => (visible = false)}
	role="tooltip"
>
	<slot />
	{#if visible}
		<div class="tooltip tooltip-{position}">
			{text}
		</div>
	{/if}
</div>

<style>
	.tooltip-container {
		position: relative;
		display: inline-block;
	}

	.tooltip {
		position: absolute;
		background: var(--bg-primary);
		color: var(--text-primary);
		padding: 8px 12px;
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		white-space: nowrap;
		box-shadow: var(--shadow-lg);
		border: 1px solid var(--border-color);
		z-index: 1000;
		pointer-events: none;
		animation: fade-in var(--transition-base) ease;
	}

	.tooltip-top {
		bottom: calc(100% + 8px);
		left: 50%;
		transform: translateX(-50%);
	}

	.tooltip-bottom {
		top: calc(100% + 8px);
		left: 50%;
		transform: translateX(-50%);
	}

	.tooltip-left {
		right: calc(100% + 8px);
		top: 50%;
		transform: translateY(-50%);
	}

	.tooltip-right {
		left: calc(100% + 8px);
		top: 50%;
		transform: translateY(-50%);
	}

	/* Arrow indicators (optional, can be added later) */
	.tooltip::before {
		content: '';
		position: absolute;
		width: 0;
		height: 0;
		border: 6px solid transparent;
	}

	.tooltip-top::before {
		bottom: -12px;
		left: 50%;
		transform: translateX(-50%);
		border-top-color: var(--bg-primary);
	}

	.tooltip-bottom::before {
		top: -12px;
		left: 50%;
		transform: translateX(-50%);
		border-bottom-color: var(--bg-primary);
	}

	.tooltip-left::before {
		right: -12px;
		top: 50%;
		transform: translateY(-50%);
		border-left-color: var(--bg-primary);
	}

	.tooltip-right::before {
		left: -12px;
		top: 50%;
		transform: translateY(-50%);
		border-right-color: var(--bg-primary);
	}

	@keyframes fade-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
</style>
