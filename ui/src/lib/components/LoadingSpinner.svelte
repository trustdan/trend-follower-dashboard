<script lang="ts">
	interface Props {
		size?: 'sm' | 'md' | 'lg';
		message?: string;
	}

	let { size = 'md', message }: Props = $props();

	const sizeClasses = {
		sm: 'w-4 h-4 border-2',
		md: 'w-8 h-8 border-3',
		lg: 'w-12 h-12 border-4'
	};
</script>

<div class="flex flex-col items-center justify-center gap-3">
	<div
		class="animate-spin rounded-full border-t-[var(--border-focus)] border-r-transparent border-b-transparent border-l-transparent {sizeClasses[
			size
		]}"
	></div>
	{#if message}
		<p class="text-sm text-[var(--text-secondary)]">{message}</p>
	{/if}
</div>
