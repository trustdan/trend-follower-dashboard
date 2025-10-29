<script lang="ts">
	export let label = '';
	export let value = '';
	export let type: 'text' | 'number' | 'email' | 'password' = 'text';
	export let placeholder = '';
	export let error = '';
	export let disabled = false;
	export let required = false;
	export let name = '';
	export let id = name || label.toLowerCase().replace(/\s+/g, '-');
</script>

<div class="input-container">
	{#if label}
		<label for={id} class="label">
			{label}
			{#if required}
				<span class="required">*</span>
			{/if}
		</label>
	{/if}

	<input
		{id}
		{name}
		{type}
		{placeholder}
		{disabled}
		{required}
		class="input"
		class:error
		bind:value
		on:input
		on:change
		on:focus
		on:blur
	/>

	{#if error}
		<span class="error-message">{error}</span>
	{/if}
</div>

<style>
	.input-container {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.label {
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--text-primary);
		display: flex;
		align-items: center;
		gap: var(--space-1);
	}

	.required {
		color: var(--color-red-500);
	}

	.input {
		padding: 12px 16px;
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
		font-size: var(--text-base);
		color: var(--text-primary);
		background: var(--bg-primary);
		transition: all var(--transition-base) ease;
		font-family: inherit;
	}

	.input:focus {
		outline: none;
		border-color: var(--border-focus);
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.input:disabled {
		background: var(--bg-tertiary);
		color: var(--text-tertiary);
		cursor: not-allowed;
	}

	.input.error {
		border-color: var(--color-red-500);
	}

	.input.error:focus {
		box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
	}

	.input::placeholder {
		color: var(--text-tertiary);
	}

	.error-message {
		font-size: var(--text-sm);
		color: var(--color-red-500);
		display: flex;
		align-items: center;
		gap: var(--space-1);
	}
</style>
