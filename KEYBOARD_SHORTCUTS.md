# TF-Engine Keyboard Shortcuts Reference

Quick reference for all keyboard shortcuts in TF-Engine GUI v8.

---

## VIM-Style Navigation

### Hint Mode (Vimium-style)

| Key | Action |
|-----|--------|
| **F** | Toggle link hint mode |
| **a-z** | Activate button with that hint (in hint mode) |
| **aa-zz** | Activate button with multi-letter hint (if needed) |
| **ESC** | Exit hint mode |

**How to use:**
1. Press **F** anywhere in the app
2. See available hints in `tf-gui.log` (or console)
3. Type the letter(s) for the button you want to click
4. Button activates automatically
5. Press **ESC** to cancel

---

### Navigation Keys (Always Active)

| Key | Action |
|-----|--------|
| **J** | Navigate down / scroll down |
| **K** | Navigate up / scroll up |
| **D** | Page down |
| **U** | Page up |

**Note:** Navigation keys have placeholder implementations in v8. Future versions will add proper scroll control.

---

## Mouse Shortcuts

| Action | Shortcut |
|--------|----------|
| **Dark Mode Toggle** | Click "🌙 Dark Mode" in top-left |
| **Help/FAQ** | Click "❓ Help" in top-left |
| **Navigate Screens** | Click menu items on left sidebar |

---

## Screen-Specific Shortcuts

### Dashboard
- **Edit Settings** - Click to modify equity, risk %, caps

### Scanner
- **Quick Presets** - Click preset buttons to auto-fill FINVIZ URL
- **Scan & Import** - Click to scrape FINVIZ and import candidates

### Checklist
- **Info Icons (ℹ️)** - Click to see detailed explanations for each gate
- **Evaluate Checklist** - Click to calculate banner (RED/YELLOW/GREEN)
- **Reset** - Click to clear all checkboxes

### Position Sizing
- **Calculate Position Size** - Click to compute shares/contracts

### Heat Check
- **Check If Trade Allowed** - Click to verify heat caps

### Trade Entry
- **Check All 5 Gates** - Click to run final validation
- **Save NO-GO Decision** - Click to log rejected trade

### Calendar
- **Refresh Calendar** - Click to reload position data

---

## Tips & Tricks

### Fastest Workflow (VIM Power User)

1. **Press F** to enter hint mode
2. **Type letter** to click any button
3. **Press ESC** if you change your mind
4. **Use J/K** to scroll through results

### Mouse User Workflow

1. Click through menu items on left
2. Fill in forms
3. Click action buttons
4. Use Help button when stuck

### First-Time User

1. Read welcome dialog
2. Click "❓ Help" for full guide
3. Follow 6-step workflow:
   - Dashboard → Scanner → Checklist → Sizing → Heat → Entry

---

## Accessibility

- **High Contrast** - All buttons have white text on British Racing Green
- **Dark Mode** - Toggle with "🌙 Dark Mode" button
- **Keyboard-First** - Full VIM navigation support
- **Mouse-First** - All features accessible by mouse
- **Help Always Available** - "❓ Help" button on every screen

---

## Customization (Future)

Planned for future versions:

- **Custom keybindings** - Remap J/K/D/U to your preference
- **VIM mode toggle** - Turn VIM hints on/off
- **Focus indicators** - Highlight currently selected button
- **Visual hints** - Show hint letters overlaid on buttons

---

## Troubleshooting

**Q: I pressed F but don't see hints**
A: Check `tf-gui.log` in the same directory as the executable. Hints are logged there.

**Q: The letter keys don't work**
A: Make sure you pressed **F** first to enter hint mode. You should see "Hint mode activated" in the log.

**Q: How do I exit hint mode?**
A: Press **ESC** key.

**Q: J/K/D/U keys don't scroll**
A: These are placeholder implementations in v8. Full scroll control coming in future versions.

**Q: Can I disable VIM mode?**
A: Not yet. Future versions will add a settings toggle.

---

## Log File Location

Keyboard activity is logged to:
```
tf-gui.log
```

In the same directory as `tf-gui-v8.exe`.

Check this file to see:
- Available hint keys when F is pressed
- Which buttons are assigned which hints
- Keyboard events and handlers

---

## Examples

### Example 1: Quick Scanner Import

```
1. Press F
2. Log shows: "a: Scan FINVIZ & Import"
3. Press a
4. Scanner runs automatically
```

### Example 2: Navigate and Calculate Position Size

```
1. Click "📐 Position Sizing" in menu (or use F + hint)
2. Fill in ticker, entry, ATR
3. Press F
4. Press letter for "Calculate Position Size" button
5. See results
```

### Example 3: Check Checklist

```
1. Click "✅ Checklist" in menu
2. Check off required gates
3. Press F
4. Press letter for "Evaluate Checklist"
5. See banner color
6. If wrong, press F + letter for "Reset"
```

---

## Quick Reference Card (Print This!)

```
╔═══════════════════════════════════════════╗
║   TF-ENGINE KEYBOARD SHORTCUTS v8         ║
╠═══════════════════════════════════════════╣
║ F          Toggle hint mode               ║
║ ESC        Exit hint mode                 ║
║ a-z        Click hinted button            ║
║                                           ║
║ J          Navigate down                  ║
║ K          Navigate up                    ║
║ D          Page down                      ║
║ U          Page up                        ║
║                                           ║
║ Mouse: Click ❓ Help for full guide       ║
╚═══════════════════════════════════════════╝
```

---

**Remember:** The value is in what this system prevents, not what it allows.

Keyboard shortcuts make discipline enforcement faster, not easier to bypass.
