<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
	export let size: 'small' | 'medium' | 'large' = 'medium';
	export let disabled = false;
	export let type: 'button' | 'submit' = 'button';
	export let loading = false;
</script>

<button {type} class="btn btn-{variant} btn-{size}" {disabled} on:click>
	{#if loading}
		<span class="spinner"></span>
	{/if}
	<slot />
</button>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		font-weight: 600;
		border: none;
		cursor: pointer;
		transition: all var(--transition-base) ease;
		white-space: nowrap;
		font-family: inherit;
	}

	/* Sizes */
	.btn-small {
		padding: 8px 16px;
		font-size: var(--text-sm);
		border-radius: var(--radius-md);
	}

	.btn-medium {
		padding: 12px 24px;
		font-size: var(--text-base);
		border-radius: var(--radius-lg);
	}

	.btn-large {
		padding: 16px 32px;
		font-size: var(--text-lg);
		border-radius: var(--radius-lg);
	}

	/* Variants */
	.btn-primary {
		background: var(--gradient-blue);
		color: white;
		box-shadow: 0 2px 4px rgba(59, 130, 246, 0.2);
	}

	.btn-primary:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(59, 130, 246, 0.3);
	}

	.btn-secondary {
		background: transparent;
		color: var(--text-primary);
		border: 2px solid var(--border-color);
	}

	.btn-secondary:hover:not(:disabled) {
		border-color: var(--border-focus);
		background: var(--bg-tertiary);
	}

	.btn-danger {
		background: var(--gradient-red);
		color: white;
		box-shadow: 0 2px 4px rgba(220, 38, 38, 0.2);
	}

	.btn-danger:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(220, 38, 38, 0.3);
	}

	/* Active state */
	.btn:active:not(:disabled) {
		transform: translateY(0);
	}

	/* Disabled state */
	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		filter: grayscale(0.5);
	}

	.btn:disabled:hover {
		transform: none;
		box-shadow: none;
	}

	/* Loading spinner */
	.spinner {
		display: inline-block;
		width: 14px;
		height: 14px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
