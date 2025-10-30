/**
 * Simple memoization utility for expensive function calls
 * Caches results based on JSON-stringified arguments
 */

type AnyFunction = (...args: any[]) => any;

interface MemoizeOptions {
	maxSize?: number; // Maximum cache size (default: 100)
	ttl?: number; // Time to live in milliseconds (default: no expiration)
}

interface CacheEntry<T> {
	value: T;
	timestamp: number;
}

export function memoize<T extends AnyFunction>(
	fn: T,
	options: MemoizeOptions = {}
): T {
	const { maxSize = 100, ttl } = options;
	const cache = new Map<string, CacheEntry<ReturnType<T>>>();

	return ((...args: Parameters<T>): ReturnType<T> => {
		const key = JSON.stringify(args);
		const now = Date.now();

		// Check if cached value exists and is still valid
		if (cache.has(key)) {
			const entry = cache.get(key)!;

			// Check TTL if specified
			if (!ttl || now - entry.timestamp < ttl) {
				return entry.value;
			}

			// Expired, remove from cache
			cache.delete(key);
		}

		// Compute new value
		const value = fn(...args);

		// Add to cache
		cache.set(key, { value, timestamp: now });

		// Enforce max size (LRU - remove oldest)
		if (cache.size > maxSize) {
			const firstKey = cache.keys().next().value;
			cache.delete(firstKey);
		}

		return value;
	}) as T;
}

/**
 * Memoize a promise-returning function
 * Useful for API calls or async operations
 */
export function memoizeAsync<T extends (...args: any[]) => Promise<any>>(
	fn: T,
	options: MemoizeOptions = {}
): T {
	const { maxSize = 100, ttl } = options;
	const cache = new Map<string, CacheEntry<Promise<ReturnType<T>>>>();
	const pending = new Map<string, Promise<ReturnType<T>>>();

	return (async (...args: Parameters<T>): Promise<ReturnType<T>> => {
		const key = JSON.stringify(args);
		const now = Date.now();

		// Check if cached value exists and is still valid
		if (cache.has(key)) {
			const entry = cache.get(key)!;

			// Check TTL if specified
			if (!ttl || now - entry.timestamp < ttl) {
				return entry.value;
			}

			// Expired, remove from cache
			cache.delete(key);
		}

		// Check if request is already pending
		if (pending.has(key)) {
			return pending.get(key)!;
		}

		// Create promise
		const promise = fn(...args);
		pending.set(key, promise);

		try {
			const value = await promise;

			// Add to cache
			cache.set(key, { value: Promise.resolve(value), timestamp: now });

			// Enforce max size (LRU - remove oldest)
			if (cache.size > maxSize) {
				const firstKey = cache.keys().next().value;
				cache.delete(firstKey);
			}

			return value;
		} finally {
			// Remove from pending
			pending.delete(key);
		}
	}) as T;
}

/**
 * Clear memoization cache for a specific function
 * (This requires storing the cache externally)
 */
export function clearMemoCache<T extends AnyFunction>(fn: T): void {
	// Note: This is a placeholder. In practice, you'd need to export the cache
	// or use a more sophisticated memoization library like lodash.memoize
	console.warn('clearMemoCache not implemented for this memoize utility');
}
