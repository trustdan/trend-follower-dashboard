<script lang="ts">
	import { page } from '$app/stores';

	$: pathSegments = $page.url.pathname.split('/').filter(Boolean);
	$: breadcrumbs = pathSegments.map((segment, i) => ({
		name: capitalize(segment.replace(/-/g, ' ')),
		path: '/' + pathSegments.slice(0, i + 1).join('/')
	}));

	function capitalize(str: string) {
		return str.charAt(0).toUpperCase() + str.slice(1);
	}
</script>

<nav class="breadcrumbs" aria-label="Breadcrumb">
	<a href="/" class="breadcrumb-item">Dashboard</a>
	{#each breadcrumbs as crumb, i}
		<span class="separator" aria-hidden="true">›</span>
		<a
			href={crumb.path}
			class="breadcrumb-item"
			class:current={i === breadcrumbs.length - 1}
			aria-current={i === breadcrumbs.length - 1 ? 'page' : undefined}
		>
			{crumb.name}
		</a>
	{/each}
</nav>

<style>
	.breadcrumbs {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-3) var(--space-6);
		background: var(--bg-secondary);
		border-bottom: 1px solid var(--border-color);
		font-size: var(--text-sm);
	}

	.breadcrumb-item {
		color: var(--text-secondary);
		text-decoration: none;
		transition: color var(--transition-fast) ease;
	}

	.breadcrumb-item:hover {
		color: var(--border-focus);
	}

	.breadcrumb-item.current {
		color: var(--text-primary);
		font-weight: 500;
		pointer-events: none;
	}

	.separator {
		color: var(--text-tertiary);
		user-select: none;
	}
</style>
