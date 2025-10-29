<script lang="ts">
	interface Props {
		variant?: 'green' | 'red' | 'yellow' | 'blue' | 'gray';
		children?: any;
	}

	let { variant = 'gray', children }: Props = $props();

	const variantClasses = {
		green: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
		red: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
		yellow: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
		blue: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300',
		gray: 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-300'
	};
</script>

<span
	class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {variantClasses[
		variant
	]}"
>
	{@render children?.()}
</span>
