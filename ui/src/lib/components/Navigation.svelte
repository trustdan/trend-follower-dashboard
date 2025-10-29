<script lang="ts">
	import { page } from '$app/stores';
	import { logger } from '$lib/utils/logger';
	import {
		LayoutDashboard,
		ScanSearch,
		ListChecks,
		Calculator,
		Flame,
		CheckCircle,
		Calendar
	} from 'lucide-svelte';

	interface NavItem {
		name: string;
		path: string;
		icon: any;
		description: string;
	}

	const navItems: NavItem[] = [
		{
			name: 'Dashboard',
			path: '/',
			icon: LayoutDashboard,
			description: 'Portfolio overview'
		},
		{
			name: 'Scanner',
			path: '/scanner',
			icon: ScanSearch,
			description: 'FINVIZ screening'
		},
		{
			name: 'Checklist',
			path: '/checklist',
			icon: ListChecks,
			description: '5 gates validation'
		},
		{
			name: 'Position Sizing',
			path: '/sizing',
			icon: Calculator,
			description: 'Calculate shares/contracts'
		},
		{
			name: 'Heat Check',
			path: '/heat',
			icon: Flame,
			description: 'Risk management'
		},
		{
			name: 'Trade Entry',
			path: '/entry',
			icon: CheckCircle,
			description: 'Final decision'
		},
		{
			name: 'Calendar',
			path: '/calendar',
			icon: Calendar,
			description: '10-week sector view'
		}
	];

	let currentPath: string;
	page.subscribe(p => {
		if (currentPath && currentPath !== p.url.pathname) {
			logger.navigate(currentPath, p.url.pathname);
		}
		currentPath = p.url.pathname;
	});

	$: isActive = (path: string) => {
		if (path === '/') {
			return $page.url.pathname === '/';
		}
		return $page.url.pathname.startsWith(path);
	};
</script>

<nav class="w-64 bg-[var(--bg-secondary)] border-r border-[var(--border-color)] flex flex-col">
	<!-- Navigation Items -->
	<div class="flex-1 py-6">
		{#each navItems as item}
			<a
				href={item.path}
				class="flex items-center gap-3 px-6 py-3 text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)] transition-all duration-200 {isActive(
					item.path
				)
					? 'bg-[var(--bg-tertiary)] text-[var(--text-primary)] border-l-4 border-[var(--border-focus)] font-medium'
					: ''}"
			>
				<svelte:component this={item.icon} size={20} />
				<div class="flex flex-col">
					<span class="text-sm">{item.name}</span>
					{#if isActive(item.path)}
						<span class="text-xs text-[var(--text-tertiary)]">{item.description}</span>
					{/if}
				</div>
			</a>
		{/each}
	</div>

	<!-- Footer -->
	<div class="p-4 border-t border-[var(--border-color)]">
		<div class="text-xs text-[var(--text-tertiary)] text-center">
			<p>TF-Engine v1.0</p>
			<p class="mt-1">© 2025</p>
		</div>
	</div>
</nav>
