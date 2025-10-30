/**
 * Keyboard Shortcuts for TF-Engine
 *
 * Global shortcuts:
 * - Escape: Close modals, clear focus
 * - Ctrl/Cmd + K: Focus ticker input
 * - Ctrl/Cmd + S: Save/Submit form
 */

export function setupKeyboardShortcuts() {
	window.addEventListener('keydown', (e) => {
		// Escape: Close modals, clear focus
		if (e.key === 'Escape') {
			// Close any open modals
			const modals = document.querySelectorAll('.modal-backdrop');
			modals.forEach((modal) => {
				const closeButton = modal.querySelector('.close-btn') as HTMLButtonElement;
				if (closeButton) closeButton.click();
			});

			// Clear focus from any focused input
			const activeElement = document.activeElement as HTMLElement;
			if (activeElement && activeElement.blur) {
				activeElement.blur();
			}
		}

		// Ctrl/Cmd + K: Focus ticker/search input
		if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
			e.preventDefault();
			const tickerInput = document.querySelector(
				'input[name="ticker"], input[type="search"]'
			) as HTMLInputElement;
			if (tickerInput) {
				tickerInput.focus();
				tickerInput.select();
			}
		}

		// Ctrl/Cmd + S: Save (prevent browser save)
		if ((e.ctrlKey || e.metaKey) && e.key === 's') {
			e.preventDefault();
			const saveButton = document.querySelector(
				'button[type="submit"], button[data-action="save"]'
			) as HTMLButtonElement;
			if (saveButton && !saveButton.disabled) {
				saveButton.click();
			}
		}

		// Ctrl/Cmd + Shift + D: Toggle debug panel (in development mode only)
		if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'D') {
			const event = new CustomEvent('toggle-debug-panel');
			window.dispatchEvent(event);
		}
	});
}

/**
 * Add keyboard shortcut hints to elements
 */
export function addShortcutHint(element: HTMLElement, shortcut: string) {
	const hint = document.createElement('span');
	hint.className = 'keyboard-hint';
	hint.textContent = shortcut;
	element.appendChild(hint);
}
