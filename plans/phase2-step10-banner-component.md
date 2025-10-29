# Phase 2 - Step 10: Banner Component

**Phase:** 2 - Checklist & Position Sizing
**Step:** 10 of 28 (overall), 1 of 5 (Phase 2)
**Duration:** 2 days
**Dependencies:** Phase 1 complete (Layout & Navigation established)

---

## Objectives

Build the centerpiece of the anti-impulsivity design: the large gradient banner that displays RED/YELLOW/GREEN states.

This component is **critical** - it's the visual enforcement of the 5 gates system.

---

## Prerequisites

- Phase 1 completed (UI structure, TailwindCSS configured)
- CSS gradient classes defined in `app.css`
- Understanding of banner state logic from `overview-plan.md`

---

## Banner Requirements (from overview-plan.md)

**Size:** Minimum 20% of screen height (at least 150px)
**States:** RED (DO NOT TRADE), YELLOW (CAUTION), GREEN (OK TO TRADE)
**Gradients:**
- RED: `linear-gradient(135deg, #DC2626 0%, #991B1B 100%)`
- YELLOW: `linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%)`
- GREEN: `linear-gradient(135deg, #10B981 0%, #059669 100%)`

**Transitions:** 0.3s ease-in-out with pulse animation on state change
**Cannot be missed:** Large, obvious, prominent placement

---

## Step-by-Step Implementation

### 1. Create Banner Component

Create `ui/src/lib/components/common/Banner.svelte`:

```svelte
<script lang="ts">
    export let state: 'red' | 'yellow' | 'green' = 'red';
    export let message: string = '';
    export let details: string = '';
    export let animate: boolean = true;

    // Map state to gradient class
    $: gradientClass = {
        red: 'banner-red',
        yellow: 'banner-yellow',
        green: 'banner-green'
    }[state];

    // Map state to icon
    $: icon = {
        red: '🛑',
        yellow: '⚠️',
        green: '✓'
    }[state];

    // Map state to title
    $: title = {
        red: 'DO NOT TRADE',
        yellow: 'CAUTION',
        green: 'OK TO TRADE'
    }[state];
</script>

<div
    class="banner {gradientClass} rounded-2xl shadow-2xl mx-auto w-[90%] min-h-[150px] flex items-center justify-center text-white {animate ? 'banner-animate' : ''}"
    style="min-height: max(150px, 20vh);"
>
    <div class="text-center px-8 py-6">
        <!-- Icon + Title -->
        <div class="flex items-center justify-center space-x-4 mb-3">
            <span class="text-5xl" role="img" aria-label={state}>{icon}</span>
            <h2 class="text-4xl font-bold tracking-wide drop-shadow-lg">
                {title}
            </h2>
            <span class="text-5xl" role="img" aria-label={state}>{icon}</span>
        </div>

        <!-- Message -->
        {#if message}
            <p class="text-xl font-medium opacity-90 drop-shadow">
                {message}
            </p>
        {/if}

        <!-- Details -->
        {#if details}
            <p class="text-base opacity-80 mt-2 drop-shadow">
                {details}
            </p>
        {/if}
    </div>
</div>

<style>
    /* Banner animation keyframe is in app.css */
    .banner-animate {
        animation: bannerPulse 0.5s ease-in-out;
    }

    /* Add glow effect */
    .banner {
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    }

    .banner-red {
        box-shadow: 0 8px 32px rgba(220, 38, 38, 0.4);
    }

    .banner-yellow {
        box-shadow: 0 8px 32px rgba(245, 158, 11, 0.4);
    }

    .banner-green {
        box-shadow: 0 8px 32px rgba(16, 185, 129, 0.4);
    }
</style>
```

---

### 2. Test Banner Component

Create a demo page: `ui/src/routes/banner-demo/+page.svelte`:

```svelte
<script lang="ts">
    import Banner from '$lib/components/common/Banner.svelte';
    import { logger } from '$lib/utils/logger';

    let state: 'red' | 'yellow' | 'green' = 'red';
    let animate = true;

    function changeState(newState: 'red' | 'yellow' | 'green') {
        logger.info('Banner state changed', { from: state, to: newState });
        state = newState;
        // Re-enable animation
        animate = true;
        setTimeout(() => animate = false, 500);
    }
</script>

<div class="max-w-7xl mx-auto space-y-8">
    <h2 class="text-3xl font-bold mb-6">Banner Component Demo</h2>

    <!-- Banner Display -->
    {#if state === 'red'}
        <Banner
            state="red"
            message="One or more REQUIRED gates failed"
            details="Complete all required checklist items to proceed"
            {animate}
        />
    {:else if state === 'yellow'}
        <Banner
            state="yellow"
            message="Quality score below threshold"
            details="All required gates pass, but consider improving quality score"
            {animate}
        />
    {:else}
        <Banner
            state="green"
            message="All gates pass • Quality met"
            details="You may proceed to trade entry"
            {animate}
        />
    {/if}

    <!-- Control Buttons -->
    <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
        <h3 class="text-xl font-semibold mb-4">Change Banner State</h3>

        <div class="flex gap-4">
            <button
                on:click={() => changeState('red')}
                class="px-6 py-3 bg-gradient-to-r from-red-600 to-red-800 text-white rounded-lg font-semibold hover:from-red-700 hover:to-red-900 transition-all"
            >
                RED (Do Not Trade)
            </button>

            <button
                on:click={() => changeState('yellow')}
                class="px-6 py-3 bg-gradient-to-r from-yellow-500 to-yellow-700 text-white rounded-lg font-semibold hover:from-yellow-600 hover:to-yellow-800 transition-all"
            >
                YELLOW (Caution)
            </button>

            <button
                on:click={() => changeState('green')}
                class="px-6 py-3 bg-gradient-to-r from-green-600 to-green-800 text-white rounded-lg font-semibold hover:from-green-700 hover:to-green-900 transition-all"
            >
                GREEN (OK to Trade)
            </button>
        </div>

        <p class="mt-4 text-sm text-[var(--text-secondary)]">
            Current state: <span class="font-semibold uppercase">{state}</span>
        </p>
    </div>

    <!-- State Logic Explanation -->
    <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
        <h3 class="text-xl font-semibold mb-4">Banner State Logic</h3>

        <div class="space-y-3">
            <div class="flex items-start space-x-3">
                <span class="text-2xl">🔴</span>
                <div>
                    <p class="font-semibold text-red-600 dark:text-red-400">RED: Do Not Trade</p>
                    <p class="text-sm text-[var(--text-secondary)]">
                        Any required gate fails → <strong>Absolutely no trade</strong>
                    </p>
                </div>
            </div>

            <div class="flex items-start space-x-3">
                <span class="text-2xl">🟡</span>
                <div>
                    <p class="font-semibold text-yellow-600 dark:text-yellow-400">YELLOW: Caution</p>
                    <p class="text-sm text-[var(--text-secondary)]">
                        All required pass, quality score &lt; threshold (default 3.0) → <strong>Proceed with caution</strong>
                    </p>
                </div>
            </div>

            <div class="flex items-start space-x-3">
                <span class="text-2xl">🟢</span>
                <div>
                    <p class="font-semibold text-green-600 dark:text-green-400">GREEN: OK to Trade</p>
                    <p class="text-sm text-[var(--text-secondary)]">
                        All required pass, quality score ≥ threshold → <strong>Proceed to trade entry</strong>
                    </p>
                </div>
            </div>
        </div>
    </div>
</div>
```

---

### 3. Test the Banner

```bash
cd /home/kali/fresh-start-trading-platform/ui

# Run dev server
npm run dev

# Open browser to: http://localhost:5173/banner-demo

# Test:
# 1. Click "RED (Do Not Trade)" - banner turns red with pulse animation
# 2. Click "YELLOW (Caution)" - banner turns yellow with pulse animation
# 3. Click "GREEN (OK to Trade)" - banner turns green with pulse animation
# 4. Verify smooth 0.3s transition between states
# 5. Verify glow effect in banner color
# 6. Verify text is readable (white with shadow)
# 7. Check in both day/night modes (toggle theme)
# 8. Verify banner is at least 20% of viewport height
# 9. Check console logs for state changes
```

---

### 4. Create Checklist Store

The banner will be driven by checklist state. Create the store now:

Create `ui/src/lib/stores/checklist.ts`:

```typescript
import { writable, derived } from 'svelte/store';
import type { Writable, Readable } from 'svelte/store';

export interface RequiredGates {
    signal: boolean;          // Gate 1: Signal confirmed
    riskSize: boolean;        // Gate 2: Risk/size calculated
    liquidity: boolean;       // Gate 3: Options requirements met
    exits: boolean;           // Gate 4: Exits defined
    behavior: boolean;        // Gate 5: Behavior constraints honored
}

export interface QualityItems {
    regime: boolean;          // Optional: Regime OK (SPY > 200 SMA)
    noChase: boolean;         // Optional: No chase (entry within 2N of 20-EMA)
    earnings: boolean;        // Optional: Earnings OK (no earnings within 2 weeks)
    journal: string;          // Optional: Journal note
}

export interface ChecklistState {
    ticker: string;
    sector: string;
    entry: number;
    atr: number;
    structure: string;        // "Stock", "Call", "Put", etc.
    requiredGates: RequiredGates;
    qualityItems: QualityItems;
    evaluatedAt: string | null;
}

// Default state
const defaultState: ChecklistState = {
    ticker: '',
    sector: '',
    entry: 0,
    atr: 0,
    structure: 'Stock',
    requiredGates: {
        signal: false,
        riskSize: false,
        liquidity: false,
        exits: false,
        behavior: false
    },
    qualityItems: {
        regime: false,
        noChase: false,
        earnings: false,
        journal: ''
    },
    evaluatedAt: null
};

// Create checklist store
export const checklist: Writable<ChecklistState> = writable(defaultState);

// Derived: count missing required gates
export const missingCount: Readable<number> = derived(checklist, $checklist => {
    const gates = $checklist.requiredGates;
    return Object.values(gates).filter(v => !v).length;
});

// Derived: quality score (0-4)
export const qualityScore: Readable<number> = derived(checklist, $checklist => {
    const items = $checklist.qualityItems;
    let score = 0;
    if (items.regime) score++;
    if (items.noChase) score++;
    if (items.earnings) score++;
    if (items.journal.trim() !== '') score++;
    return score;
});

// Derived: banner state
export const bannerState: Readable<'red' | 'yellow' | 'green'> = derived(
    [missingCount, qualityScore],
    ([$missingCount, $qualityScore]) => {
        // RED: Any required gate fails
        if ($missingCount > 0) return 'red';

        // GREEN: All required pass, quality score >= threshold (3)
        if ($qualityScore >= 3) return 'green';

        // YELLOW: All required pass, but quality score < 3
        return 'yellow';
    }
);

// Derived: banner message
export const bannerMessage: Readable<string> = derived(
    [missingCount, qualityScore, bannerState],
    ([$missingCount, $qualityScore, $bannerState]) => {
        if ($bannerState === 'red') {
            if ($missingCount === 1) return '1 required gate failed';
            return `${$missingCount} required gates failed`;
        }

        if ($bannerState === 'yellow') {
            return `Quality score below threshold (${$qualityScore}/3)`;
        }

        return `All gates pass • Quality score: ${$qualityScore}/4`;
    }
);

// Reset checklist
export function resetChecklist() {
    checklist.set(defaultState);
}
```

---

### 5. Test Banner with Store

Update banner demo to use store:

```svelte
<script lang="ts">
    import Banner from '$lib/components/common/Banner.svelte';
    import { checklist, bannerState, bannerMessage, missingCount, qualityScore } from '$lib/stores/checklist';
    import { logger } from '$lib/utils/logger';

    let animate = false;

    // Watch for state changes to trigger animation
    $: if ($bannerState) {
        animate = true;
        setTimeout(() => animate = false, 500);
    }

    function toggleGate(gate: keyof typeof $checklist.requiredGates) {
        checklist.update(c => ({
            ...c,
            requiredGates: {
                ...c.requiredGates,
                [gate]: !c.requiredGates[gate]
            }
        }));
        logger.info('Gate toggled', { gate, value: $checklist.requiredGates[gate] });
    }

    function toggleQuality(item: keyof Omit<typeof $checklist.qualityItems, 'journal'>) {
        checklist.update(c => ({
            ...c,
            qualityItems: {
                ...c.qualityItems,
                [item]: !c.qualityItems[item]
            }
        }));
        logger.info('Quality item toggled', { item, value: $checklist.qualityItems[item] });
    }
</script>

<div class="max-w-7xl mx-auto space-y-8">
    <h2 class="text-3xl font-bold mb-6">Banner Component with Store</h2>

    <!-- Banner Display (driven by store) -->
    <Banner
        state={$bannerState}
        message={$bannerMessage}
        details={$missingCount > 0 ? 'Complete all required gates to proceed' : 'Ready for position sizing'}
        {animate}
    />

    <!-- Interactive Controls -->
    <div class="grid md:grid-cols-2 gap-6">
        <!-- Required Gates -->
        <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
            <h3 class="text-xl font-semibold mb-4">Required Gates (Must Check All)</h3>

            <div class="space-y-2">
                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.requiredGates.signal}
                        on:change={() => toggleGate('signal')}
                        class="w-5 h-5"
                    />
                    <span>✓ Signal: 55-bar breakout confirmed</span>
                </label>

                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.requiredGates.riskSize}
                        on:change={() => toggleGate('riskSize')}
                        class="w-5 h-5"
                    />
                    <span>✓ Risk/Size: 2×N stop, 0.5×N adds</span>
                </label>

                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.requiredGates.liquidity}
                        on:change={() => toggleGate('liquidity')}
                        class="w-5 h-5"
                    />
                    <span>✓ Liquidity: Volume/OI sufficient</span>
                </label>

                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.requiredGates.exits}
                        on:change={() => toggleGate('exits')}
                        class="w-5 h-5"
                    />
                    <span>✓ Exits: 10-bar Donchian OR 2×N</span>
                </label>

                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.requiredGates.behavior}
                        on:change={() => toggleGate('behavior')}
                        class="w-5 h-5"
                    />
                    <span>✓ Behavior: No cooldown, heat OK</span>
                </label>
            </div>

            <p class="mt-4 text-sm">
                Missing: <span class="font-semibold">{$missingCount}</span>
            </p>
        </div>

        <!-- Quality Items -->
        <div class="bg-[var(--bg-secondary)] rounded-xl p-6 border border-[var(--border-color)]">
            <h3 class="text-xl font-semibold mb-4">Quality Items (Optional)</h3>

            <div class="space-y-2">
                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.qualityItems.regime}
                        on:change={() => toggleQuality('regime')}
                        class="w-5 h-5"
                    />
                    <span>Regime OK (SPY > 200 SMA)</span>
                </label>

                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.qualityItems.noChase}
                        on:change={() => toggleQuality('noChase')}
                        class="w-5 h-5"
                    />
                    <span>No Chase (within 2N of 20-EMA)</span>
                </label>

                <label class="flex items-center space-x-3 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={$checklist.qualityItems.earnings}
                        on:change={() => toggleQuality('earnings')}
                        class="w-5 h-5"
                    />
                    <span>Earnings OK (no earnings soon)</span>
                </label>
            </div>

            <p class="mt-4 text-sm">
                Quality Score: <span class="font-semibold">{$qualityScore} / 4</span>
            </p>
        </div>
    </div>
</div>
```

---

## Verification Checklist

Before proceeding to Step 11, verify:

- [ ] Banner component created (`Banner.svelte`)
- [ ] Banner displays with correct gradient for each state
- [ ] Banner is at least 20% viewport height (150px minimum)
- [ ] Pulse animation works on state change
- [ ] Glow effect in banner color visible
- [ ] Text is readable (white with shadow)
- [ ] Smooth 0.3s transition between states
- [ ] Works in both day and night modes
- [ ] Checklist store created with derived stores
- [ ] Banner state logic correct (RED → YELLOW → GREEN)
- [ ] Banner message updates based on state
- [ ] Demo page works and all interactions functional
- [ ] Console logs state changes
- [ ] No TypeScript errors

---

## Expected Outputs

1. **Banner Component:** `ui/src/lib/components/common/Banner.svelte`
2. **Checklist Store:** `ui/src/lib/stores/checklist.ts`
3. **Demo Page:** `ui/src/routes/banner-demo/+page.svelte`
4. **Working Features:**
   - RED/YELLOW/GREEN gradient transitions
   - Pulse animation on state change
   - Banner driven by checklist store
   - Quality score calculation
   - Missing gates count

---

## Time Estimate

- Banner Component: 2-3 hours
- Checklist Store: 2-3 hours
- Demo Page: 1-2 hours
- Testing & Refinement: 2-3 hours
- **Total:** ~7-11 hours (1-2 days)

---

## References

- [overview-plan.md - Banner System](../plans/overview-plan.md#banner-system-visual-discipline-enforcement)
- [overview-plan.md - The 5 Hard Gates](../plans/overview-plan.md#the-5-hard-gates-enforced-by-backend)
- [roadmap.md - Step 10](../plans/roadmap.md#step-10-banner-component)

---

## Next Step

Proceed to: **[Phase 2 - Step 11: Checklist Form & Required Gates](phase2-step11-checklist-form.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
