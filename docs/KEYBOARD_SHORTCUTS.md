# Keyboard Shortcuts

**TF-Engine Trading Platform**
**Last Updated:** 2025-10-29

---

## Global Shortcuts

These shortcuts work from any screen in the application:

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Escape` | Close Modal / Clear Focus | Closes any open modal dialogs and removes focus from active input fields |
| `Ctrl/Cmd + K` | Focus Ticker Input | Instantly focuses the ticker symbol input field and selects any existing text |
| `Ctrl/Cmd + S` | Save/Submit | Triggers the primary save or submit button on the current screen (if enabled) |
| `Ctrl/Cmd + Shift + D` | Toggle Debug Panel | Opens/closes the debug panel for viewing application logs (dev mode only) |

---

## Form-Specific Shortcuts

### Checklist Screen
- `Tab` - Navigate between form fields
- `Space` - Toggle checkboxes
- `Enter` - Submit form (when focused on submit button)

### Position Sizing Screen
- `Tab` - Navigate between input fields
- `Enter` - Calculate position size (when focused on calculate button)

### Trade Entry Screen
- `Tab` - Navigate between form fields
- `Enter` - Submit decision (when save button is enabled and focused)

---

## Navigation Shortcuts (Planned for Future Release)

| Shortcut | Action |
|----------|--------|
| `Ctrl/Cmd + 1` | Go to Dashboard |
| `Ctrl/Cmd + 2` | Go to Scanner |
| `Ctrl/Cmd + 3` | Go to Checklist |
| `Ctrl/Cmd + 4` | Go to Position Sizing |
| `Ctrl/Cmd + 5` | Go to Heat Check |
| `Ctrl/Cmd + 6` | Go to Trade Entry |
| `Ctrl/Cmd + 7` | Go to Calendar |

---

## Debug Panel Shortcuts

When debug panel is open:
- `Escape` - Close debug panel
- Filter logs by level using dropdown
- Click "Clear" to remove all logs
- Click "Export" to download logs as JSON

---

## Browser Standard Shortcuts

These work across all web applications:
- `Ctrl/Cmd + R` - Refresh page
- `Ctrl/Cmd + +` - Zoom in
- `Ctrl/Cmd + -` - Zoom out
- `Ctrl/Cmd + 0` - Reset zoom
- `F11` - Toggle fullscreen
- `F12` - Open browser DevTools

---

## Accessibility Features

- **Keyboard Navigation:** All interactive elements can be reached via Tab key
- **Focus Indicators:** Clear visual focus states on all inputs and buttons
- **ARIA Labels:** Screen reader support for all form elements
- **Skip Links:** Jump to main content (planned)

---

## Tips

1. **Ticker Input Focus:** Use `Ctrl/Cmd + K` to quickly jump to the ticker input without reaching for your mouse. This works from any screen.

2. **Quick Save:** After completing a form, press `Ctrl/Cmd + S` instead of clicking the save button. The shortcut only works if the button is enabled.

3. **Debug Logs:** If something isn't working as expected, open the debug panel (`Ctrl/Cmd + Shift + D`) to see real-time logs. Export the logs if you need to report an issue.

4. **Escape to Reset:** Hit `Escape` to quickly clear focus from any input field. This is useful when you want to ensure no field is selected before using other shortcuts.

---

## Planned Enhancements (Step 23+)

- Command palette (Ctrl/Cmd + P) for fuzzy search navigation
- Customizable keyboard shortcuts
- Shortcut cheat sheet overlay (Shift + ?)
- Quick actions menu (Ctrl/Cmd + J)

---

**Note:** Mac users should use `Cmd` where Windows/Linux users use `Ctrl`.
