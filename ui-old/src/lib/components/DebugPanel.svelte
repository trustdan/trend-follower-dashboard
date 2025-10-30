<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { logger, type LogLevel } from '$lib/utils/logger';
	import { onMount } from 'svelte';

	let visible = false;
	let logs = logger.getLogs();
	let filter: 'all' | LogLevel = 'all';
	let refreshInterval: number;

	function togglePanel() {
		visible = !visible;
	}

	function clearLogs() {
		logger.clear();
		logs = [];
	}

	function exportLogs() {
		const jsonData = logger.export();
		const blob = new Blob([jsonData], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `tf-engine-logs-${new Date().toISOString()}.json`;
		a.click();
		URL.revokeObjectURL(url);
	}

	function refreshLogs() {
		logs = logger.getLogs();
	}

	$: filteredLogs = filter === 'all' ? logs : logs.filter((log) => log.level === filter);

	// Auto-refresh logs when panel is visible
	$: if (visible && !refreshInterval) {
		refreshInterval = setInterval(refreshLogs, 500) as unknown as number;
	} else if (!visible && refreshInterval) {
		clearInterval(refreshInterval);
		refreshInterval = 0 as unknown as number;
	}

	// Listen for keyboard shortcut and custom event
	onMount(() => {
		const handleToggle = () => togglePanel();
		window.addEventListener('toggle-debug-panel', handleToggle);

		return () => {
			window.removeEventListener('toggle-debug-panel', handleToggle);
			if (refreshInterval) clearInterval(refreshInterval);
		};
	});
</script>

<!-- Toggle button (only visible in dev mode) -->
{#if import.meta.env.DEV}
	<button
		class="debug-toggle"
		on:click={togglePanel}
		title="Debug Panel (Ctrl+Shift+D)"
		aria-label="Toggle debug panel"
	>
		🛠️
	</button>
{/if}

{#if visible}
	<div class="debug-panel" transition:fly={{ x: 300, duration: 200 }}>
		<div class="panel-header">
			<h3>Debug Panel</h3>
			<button class="close-btn" on:click={togglePanel} aria-label="Close debug panel">×</button>
		</div>

		<div class="panel-toolbar">
			<select bind:value={filter} aria-label="Filter logs by level">
				<option value="all">All Logs</option>
				<option value="debug">Debug</option>
				<option value="info">Info</option>
				<option value="warn">Warnings</option>
				<option value="error">Errors</option>
			</select>
			<button on:click={clearLogs}>Clear</button>
			<button on:click={exportLogs}>Export</button>
			<span class="log-count">{filteredLogs.length} logs</span>
		</div>

		<div class="logs-container">
			{#each filteredLogs as log (log.timestamp)}
				<div class="log-entry log-{log.level}">
					<span class="timestamp">{new Date(log.timestamp).toLocaleTimeString()}</span>
					<span class="level">{log.level.toUpperCase()}</span>
					<span class="message">{log.message}</span>
					{#if log.data}
						<pre class="data">{JSON.stringify(log.data, null, 2)}</pre>
					{/if}
				</div>
			{:else}
				<div class="empty-state">No logs yet. Interact with the app to see logs appear.</div>
			{/each}
		</div>
	</div>
{/if}

<style>
	.debug-toggle {
		position: fixed;
		bottom: 20px;
		right: 20px;
		width: 50px;
		height: 50px;
		border-radius: var(--radius-full);
		background: var(--gradient-purple);
		color: white;
		border: none;
		font-size: 24px;
		cursor: pointer;
		box-shadow: var(--shadow-lg);
		z-index: 999;
		transition: transform var(--transition-base) ease;
	}

	.debug-toggle:hover {
		transform: scale(1.1);
	}

	.debug-panel {
		position: fixed;
		top: 0;
		right: 0;
		width: 450px;
		height: 100vh;
		background: var(--bg-primary);
		border-left: 1px solid var(--border-color);
		box-shadow: var(--shadow-xl);
		z-index: 1000;
		display: flex;
		flex-direction: column;
	}

	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--space-4);
		border-bottom: 1px solid var(--border-color);
		background: var(--bg-secondary);
	}

	.panel-header h3 {
		font-size: var(--text-xl);
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

	.panel-toolbar {
		display: flex;
		gap: var(--space-2);
		padding: var(--space-3);
		border-bottom: 1px solid var(--border-color);
		background: var(--bg-secondary);
		align-items: center;
	}

	.panel-toolbar select,
	.panel-toolbar button {
		padding: 6px 12px;
		border-radius: var(--radius-md);
		border: 1px solid var(--border-color);
		background: var(--bg-primary);
		color: var(--text-primary);
		cursor: pointer;
		font-size: var(--text-sm);
		transition: all var(--transition-fast) ease;
	}

	.panel-toolbar button:hover {
		background: var(--bg-tertiary);
	}

	.log-count {
		margin-left: auto;
		font-size: var(--text-sm);
		color: var(--text-secondary);
	}

	.logs-container {
		flex: 1;
		overflow-y: auto;
		padding: var(--space-3);
	}

	.log-entry {
		margin-bottom: var(--space-2);
		padding: var(--space-2);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		border-left: 3px solid;
		font-family: 'Courier New', monospace;
	}

	.log-debug {
		background: rgba(148, 163, 184, 0.1);
		border-color: #94a3b8;
	}

	.log-info {
		background: rgba(59, 130, 246, 0.1);
		border-color: #3b82f6;
	}

	.log-warn {
		background: rgba(245, 158, 11, 0.1);
		border-color: #f59e0b;
	}

	.log-error {
		background: rgba(220, 38, 38, 0.1);
		border-color: #dc2626;
	}

	.timestamp {
		color: var(--text-tertiary);
		margin-right: var(--space-2);
		font-size: var(--text-xs);
	}

	.level {
		font-weight: 600;
		margin-right: var(--space-2);
		font-size: var(--text-xs);
	}

	.message {
		color: var(--text-primary);
	}

	.data {
		margin-top: var(--space-2);
		padding: var(--space-2);
		background: var(--bg-tertiary);
		border-radius: var(--radius-sm);
		font-size: var(--text-xs);
		overflow-x: auto;
		white-space: pre-wrap;
		word-wrap: break-word;
	}

	.empty-state {
		text-align: center;
		padding: var(--space-8);
		color: var(--text-tertiary);
		font-style: italic;
	}
</style>
