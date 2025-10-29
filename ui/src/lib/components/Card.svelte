<script lang="ts">
	interface Props {
		title?: string;
		subtitle?: string;
		padding?: string;
		gradient?: boolean;
		children?: any;
	}

	let { title, subtitle, padding = 'p-6', gradient = false, children }: Props = $props();
</script>

<div
	class="bg-[var(--bg-secondary)] rounded-lg border border-[var(--border-color)] {padding} {gradient
		? 'bg-gradient-to-br from-[var(--bg-secondary)] to-[var(--bg-tertiary)]'
		: ''}"
>
	{#if title}
		<div class="mb-4">
			<h3 class="text-lg font-semibold text-[var(--text-primary)]">{title}</h3>
			{#if subtitle}
				<p class="text-sm text-[var(--text-secondary)] mt-1">{subtitle}</p>
			{/if}
		</div>
	{/if}
	{@render children?.()}
</div>
