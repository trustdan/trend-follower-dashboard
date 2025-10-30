// Frontend logging utility
// Logs to browser console with color-coding and timestamps

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

interface LogEntry {
	timestamp: string;
	level: LogLevel;
	message: string;
	data?: unknown;
}

class Logger {
	private logs: LogEntry[] = [];
	private maxLogs = 1000;

	private formatTimestamp(): string {
		const now = new Date();
		return now.toISOString();
	}

	private log(level: LogLevel, message: string, data?: unknown) {
		const entry: LogEntry = {
			timestamp: this.formatTimestamp(),
			level,
			message,
			data
		};

		this.logs.push(entry);

		// Keep only last maxLogs entries
		if (this.logs.length > this.maxLogs) {
			this.logs = this.logs.slice(-this.maxLogs);
		}

		// Console output with colors
		const styles = {
			debug: 'color: #94A3B8; font-weight: normal',
			info: 'color: #3B82F6; font-weight: bold',
			warn: 'color: #F59E0B; font-weight: bold',
			error: 'color: #DC2626; font-weight: bold'
		};

		const timestamp = entry.timestamp.split('T')[1].slice(0, 12);
		const prefix = `[${timestamp}] [${level.toUpperCase()}]`;

		if (data !== undefined) {
			console.log(`%c${prefix} ${message}`, styles[level], data);
		} else {
			console.log(`%c${prefix} ${message}`, styles[level]);
		}
	}

	debug(message: string, data?: unknown) {
		this.log('debug', message, data);
	}

	info(message: string, data?: unknown) {
		this.log('info', message, data);
	}

	warn(message: string, data?: unknown) {
		this.log('warn', message, data);
	}

	error(message: string, data?: unknown) {
		this.log('error', message, data);
	}

	// Navigation logging
	navigate(from: string, to: string) {
		this.info(`Navigation: ${from} → ${to}`);
	}

	// Theme logging
	themeChange(theme: string) {
		this.info(`Theme changed to: ${theme}`);
	}

	// API logging
	apiRequest(method: string, url: string) {
		this.debug(`API ${method} ${url}`);
	}

	apiResponse(method: string, url: string, status: number, duration: number) {
		if (status >= 400) {
			this.error(`API ${method} ${url} → ${status} (${duration}ms)`);
		} else {
			this.debug(`API ${method} ${url} → ${status} (${duration}ms)`);
		}
	}

	// Get all logs
	getLogs(): LogEntry[] {
		return [...this.logs];
	}

	// Clear logs
	clear() {
		this.logs = [];
		console.clear();
		this.info('Logs cleared');
	}

	// Export logs as JSON
	export(): string {
		return JSON.stringify(this.logs, null, 2);
	}
}

// Global logger instance
export const logger = new Logger();
