# Design System Audit Checklist

**TF-Engine:** Trend Following Trading Platform
**Created:** 2025-10-29
**Purpose:** Ensure visual consistency across all UI components

---

## Colors

- [ ] All gradients use defined CSS variables (--gradient-blue, --gradient-green, etc.)
- [ ] No hard-coded hex colors in components
- [ ] Theme toggle updates all colors correctly
- [ ] Text contrast meets WCAG AA standards (4.5:1 for normal text)
- [ ] Banner colors (RED/YELLOW/GREEN) are consistent across all screens
- [ ] Error states use consistent red coloring
- [ ] Success states use consistent green coloring
- [ ] Info/neutral states use consistent blue coloring

## Spacing

- [ ] All spacing uses --space-N variables (no arbitrary values like "20px")
- [ ] Consistent padding in cards (24px or var(--space-6))
- [ ] Consistent gaps in flex/grid layouts (--space-4 or --space-5)
- [ ] Form field spacing is uniform (--space-4 between fields)
- [ ] Section spacing is consistent (--space-8 between major sections)
- [ ] Button spacing within groups is uniform (--space-3)

## Typography

- [ ] Font sizes use --text-N variables (no hard-coded px values)
- [ ] Headings follow hierarchy (h1 > h2 > h3)
- [ ] Line height is readable (1.5 for body text, 1.2 for headings)
- [ ] Font weights are consistent:
  - 400 (regular) for body text
  - 500 (medium) for labels
  - 600 (semibold) for buttons and emphasis
  - 700 (bold) for headings
- [ ] Letter spacing is appropriate for all text sizes
- [ ] No font family overrides (all use system font stack)

## Borders

- [ ] Border radius consistent:
  - 6px for inputs
  - 8px for buttons
  - 12px for cards
  - 16px for large containers
- [ ] Border colors use --border-color variable
- [ ] Focus states use --border-focus color (blue accent)
- [ ] Error states show red border (--color-red-500)
- [ ] Disabled states show muted border (--border-color with opacity)

## Shadows

- [ ] Shadows follow elevation system:
  - Small: `0 1px 3px rgba(0, 0, 0, 0.1)`
  - Medium: `0 4px 6px rgba(0, 0, 0, 0.1)`
  - Large: `0 10px 25px rgba(0, 0, 0, 0.15)`
  - Extra large: `0 20px 40px rgba(0, 0, 0, 0.2)`
- [ ] No arbitrary shadow values
- [ ] Gradient buttons have colored shadows matching gradient
- [ ] Hover states increase shadow elevation

## Buttons

- [ ] Primary buttons use gradient backgrounds (--gradient-blue)
- [ ] Secondary buttons have gradient borders with transparent bg
- [ ] Danger buttons use red gradient (--gradient-red)
- [ ] Disabled buttons have:
  - Reduced opacity (0.5)
  - Grayscale filter (0.5)
  - Cursor: not-allowed
  - No hover effects
- [ ] All buttons have hover effects:
  - Transform: translateY(-2px)
  - Increased shadow
- [ ] All buttons have active/pressed states
- [ ] Loading states show spinner with disabled appearance
- [ ] Button sizes are consistent (small, medium, large)

## Forms

- [ ] Input fields have consistent padding (12px 16px)
- [ ] Labels are bold (font-weight: 500) and positioned above inputs
- [ ] Focus states are visually clear (blue border, glow effect)
- [ ] Error states show:
  - Red border
  - Error message below field
  - Error icon (optional)
- [ ] Placeholder text is muted (--text-tertiary)
- [ ] Disabled inputs have:
  - Muted background
  - Muted text
  - Cursor: not-allowed
- [ ] Select dropdowns match input styling
- [ ] Checkboxes have custom styling consistent with design system

## Cards & Containers

- [ ] All cards use consistent padding (24px or --space-6)
- [ ] Card backgrounds use --bg-secondary
- [ ] Card borders use --border-color (1px solid)
- [ ] Card shadows use medium elevation
- [ ] Card hover states (if interactive) show slight lift
- [ ] Card headers have consistent styling
- [ ] Card footers have consistent styling and spacing

## Layout

- [ ] Container max-width is consistent across pages
- [ ] Page padding is consistent (--space-6 or --space-8)
- [ ] Grid layouts use consistent gap values
- [ ] Responsive breakpoints are consistent:
  - Mobile: < 640px
  - Tablet: 640px - 1024px
  - Desktop: > 1024px
- [ ] Navigation spacing is consistent
- [ ] Footer spacing is consistent

## Animations & Transitions

- [ ] All transitions use consistent duration:
  - Fast: 150ms (hover, focus)
  - Medium: 200ms (slide, fade)
  - Slow: 300ms (large movements)
- [ ] All transitions use consistent easing (ease, ease-in-out)
- [ ] Loading states have smooth animations
- [ ] Page transitions are consistent
- [ ] Modal animations (fade in backdrop, slide in modal)
- [ ] Notification animations (slide in from top)
- [ ] No janky or laggy animations

## Icons & Images

- [ ] Icon sizes are consistent (16px, 20px, 24px, 32px)
- [ ] Icon colors use CSS variables
- [ ] Icons have proper aria-labels for accessibility
- [ ] SVG icons are optimized
- [ ] Loading spinners are consistent in size and color
- [ ] Image placeholders show skeleton loaders

## Accessibility

- [ ] All interactive elements have focus states
- [ ] Focus states are visible and clear
- [ ] Tab order is logical
- [ ] All form inputs have associated labels
- [ ] ARIA labels present where needed
- [ ] Color is not the only indicator (icons/text accompany colors)
- [ ] Text contrast meets WCAG AA (4.5:1 minimum)
- [ ] Keyboard navigation works throughout app
- [ ] Screen reader support verified

## Theme Support

- [ ] All components work in day mode
- [ ] All components work in night mode
- [ ] Theme toggle updates all components immediately
- [ ] No theme-specific hard-coded values
- [ ] Gradients adapt to theme (if applicable)
- [ ] Images/icons adapt to theme (if applicable)

## Component-Specific Checks

### Banner Component
- [ ] RED banner: Large, impossible to miss, gradient background
- [ ] YELLOW banner: Amber gradient, clear warning
- [ ] GREEN banner: Green gradient, positive confirmation
- [ ] Banner text is bold and large
- [ ] Banner icons are visible and appropriate
- [ ] Banner transitions smoothly between states

### Position Table
- [ ] Table headers are bold and distinct
- [ ] Table rows have hover states
- [ ] Table cells are properly aligned
- [ ] Table spacing is consistent
- [ ] Table borders use --border-color
- [ ] Empty state shows helpful message

### Heat Bars
- [ ] Progress bars animate smoothly
- [ ] Color transitions (green → yellow → red) are clear
- [ ] Percentage labels are readable
- [ ] Bar height is consistent
- [ ] Bar borders and shadows are appropriate

### Calendar Grid
- [ ] Grid cells are evenly sized
- [ ] Cell borders are consistent
- [ ] Ticker badges are styled consistently
- [ ] Hover states are clear
- [ ] Color coding is intuitive

---

## Audit Process

1. **Initial Review:** Go through each component file visually
2. **Browser Testing:** Test in Chrome, Firefox, Safari
3. **Contrast Check:** Use WebAIM Contrast Checker
4. **Theme Toggle:** Verify all components in both themes
5. **Responsive Test:** Check mobile, tablet, desktop layouts
6. **Accessibility Test:** Use screen reader and keyboard-only navigation
7. **Performance Test:** Check for animation performance (60fps target)

---

## Identified Issues

### Component: [Name]
- **Issue:** [Description]
- **Location:** [File:line]
- **Fix:** [Action to take]
- **Priority:** [Low/Medium/High]
- **Status:** [To Do/In Progress/Fixed]

---

## Sign-Off

- [ ] All color checks passed
- [ ] All spacing checks passed
- [ ] All typography checks passed
- [ ] All border checks passed
- [ ] All shadow checks passed
- [ ] All button checks passed
- [ ] All form checks passed
- [ ] All card checks passed
- [ ] All layout checks passed
- [ ] All animation checks passed
- [ ] All icon checks passed
- [ ] All accessibility checks passed
- [ ] All theme checks passed
- [ ] All component-specific checks passed

**Audited By:** Claude Code
**Date:** 2025-10-29
**Result:** [PASS / FAIL / PARTIAL]

---

**Next Steps:** Apply fixes to all identified issues, then re-audit before marking Step 22 complete.
