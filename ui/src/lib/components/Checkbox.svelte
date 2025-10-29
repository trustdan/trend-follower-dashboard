<script lang="ts">
	export let checked = false;
	export let label = '';
	export let gradient: 'blue' | 'green' | 'red' = 'blue';
	export let disabled = false;
	export let name = '';
</script>

<label class="checkbox-container" class:disabled>
	<input type="checkbox" {name} {disabled} bind:checked on:change />
	<span class="checkmark {gradient}" class:checked></span>
	{#if label}
		<span class="label-text">{label}</span>
	{:else}
		<slot />
	{/if}
</label>

<style>
	.checkbox-container {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		cursor: pointer;
		user-select: none;
	}

	.checkbox-container.disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	input[type='checkbox'] {
		position: absolute;
		opacity: 0;
		width: 0;
		height: 0;
	}

	.checkmark {
		position: relative;
		width: 24px;
		height: 24px;
		border: 2px solid var(--border-color);
		border-radius: var(--radius-md);
		transition: all var(--transition-base) ease;
		flex-shrink: 0;
	}

	.checkmark.checked {
		border-color: transparent;
	}

	.checkmark.blue.checked {
		background: var(--gradient-blue);
	}

	.checkmark.green.checked {
		background: var(--gradient-green);
	}

	.checkmark.red.checked {
		background: var(--gradient-red);
	}

	/* Checkmark icon */
	.checkmark::after {
		content: '';
		position: absolute;
		display: none;
		left: 7px;
		top: 3px;
		width: 6px;
		height: 12px;
		border: solid white;
		border-width: 0 2px 2px 0;
		transform: rotate(45deg);
		animation: checkmark-draw var(--transition-slow) ease-in-out;
	}

	.checkmark.checked::after {
		display: block;
	}

	@keyframes checkmark-draw {
		0% {
			height: 0;
		}
		100% {
			height: 12px;
		}
	}

	.label-text {
		font-size: var(--text-base);
		color: var(--text-primary);
	}

	.checkbox-container:hover .checkmark:not(.checked) {
		border-color: var(--border-focus);
	}

	.checkbox-container.disabled .checkmark {
		background: var(--bg-tertiary);
	}
</style>
