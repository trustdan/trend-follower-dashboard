<script lang="ts">
	import { fade, fly } from 'svelte/transition';

	export let type: 'success' | 'error' | 'info' = 'info';
	export let message: string;
	export let duration = 3000; // Auto-dismiss after 3s
	export let onDismiss: () => void;

	let visible = true;

	if (duration > 0) {
		setTimeout(() => {
			visible = false;
			setTimeout(onDismiss, 300);
		}, duration);
	}
</script>

{#if visible}
	<div
		class="notification notification-{type}"
		transition:fly={{ y: -20, duration: 200 }}
		on:click={() => {
			visible = false;
			setTimeout(onDismiss, 300);
		}}
		role="alert"
	>
		<div class="icon">
			{#if type === 'success'}✓{/if}
			{#if type === 'error'}✗{/if}
			{#if type === 'info'}ℹ{/if}
		</div>
		<div class="message">{message}</div>
		<button
			class="dismiss"
			on:click|stopPropagation={() => {
				visible = false;
				setTimeout(onDismiss, 300);
			}}
			aria-label="Dismiss notification"
		>
			×
		</button>
	</div>
{/if}

<style>
	.notification {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-4);
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-lg);
		cursor: pointer;
		min-width: 300px;
		max-width: 500px;
	}

	.notification-success {
		background: var(--gradient-green);
		color: white;
	}

	.notification-error {
		background: var(--gradient-red);
		color: white;
		animation: shake var(--transition-slow) ease-in-out;
	}

	.notification-info {
		background: var(--gradient-blue);
		color: white;
	}

	@keyframes shake {
		0%,
		100% {
			transform: translateX(0);
		}
		25% {
			transform: translateX(-10px);
		}
		75% {
			transform: translateX(10px);
		}
	}

	.icon {
		font-size: var(--text-2xl);
		font-weight: bold;
		flex-shrink: 0;
	}

	.message {
		flex: 1;
		font-weight: 500;
		font-size: var(--text-base);
	}

	.dismiss {
		background: none;
		border: none;
		color: white;
		font-size: var(--text-3xl);
		line-height: 1;
		cursor: pointer;
		opacity: 0.7;
		transition: opacity var(--transition-fast) ease;
		padding: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.dismiss:hover {
		opacity: 1;
	}
</style>
