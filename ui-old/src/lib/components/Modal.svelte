<script lang="ts">
	import { fade, fly } from 'svelte/transition';

	export let isOpen = false;
	export let title = '';
	export let onClose: (() => void) | undefined = undefined;

	function handleClose() {
		isOpen = false;
		if (onClose) onClose();
	}

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			handleClose();
		}
	}

	// Handle Escape key
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && isOpen) {
			handleClose();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen}
	<div class="modal-backdrop" on:click={handleBackdropClick} transition:fade={{ duration: 200 }}>
		<div class="modal" transition:fly={{ y: -20, duration: 200 }} role="dialog" aria-modal="true">
			<div class="modal-header">
				{#if title}
					<h2 class="modal-title">{title}</h2>
				{/if}
				<button class="close-btn" on:click={handleClose} aria-label="Close modal"> × </button>
			</div>

			<div class="modal-content">
				<slot />
			</div>

			<div class="modal-footer">
				<slot name="footer" />
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: var(--space-4);
	}

	.modal {
		background: var(--bg-primary);
		border-radius: var(--radius-xl);
		box-shadow: var(--shadow-xl);
		max-width: 600px;
		width: 100%;
		max-height: 90vh;
		display: flex;
		flex-direction: column;
		border: 1px solid var(--border-color);
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--space-6);
		border-bottom: 1px solid var(--border-color);
	}

	.modal-title {
		font-size: var(--text-2xl);
		font-weight: 600;
		color: var(--text-primary);
		margin: 0;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: var(--text-4xl);
		line-height: 1;
		cursor: pointer;
		color: var(--text-secondary);
		transition: color var(--transition-fast) ease;
		padding: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.close-btn:hover {
		color: var(--text-primary);
	}

	.modal-content {
		padding: var(--space-6);
		overflow-y: auto;
		flex: 1;
	}

	.modal-footer {
		padding: var(--space-6);
		border-top: 1px solid var(--border-color);
		display: flex;
		justify-content: flex-end;
		gap: var(--space-3);
	}

	/* Hide footer if empty */
	.modal-footer:empty {
		display: none;
	}
</style>
