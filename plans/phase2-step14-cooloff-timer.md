# Phase 2 - Step 14: 2-Minute Cool-Off Timer

**Phase:** 2 - Checklist & Position Sizing
**Step:** 14 of 28, 5 of 5 (Phase 2)
**Duration:** 1-2 days
**Dependencies:** Step 13 (Position Sizing)

---

## Objectives

Implement the 2-minute impulse brake timer. When user saves checklist evaluation (banner GREEN), start countdown. Timer displayed prominently. Certain actions disabled during countdown.

---

## Key Components

1. **Timer Store** - Countdown from 120 seconds to 0
2. **Timer Display** - Large, prominent: "Cool-off period: 2:00 remaining"
3. **Countdown** - Updates every second: 2:00 → 1:59 → ... → 0:01 → 0:00
4. **Persistence** - Timer state survives navigation between screens (Svelte store)
5. **Backend Validation** - Backend validates 2-min elapsed when gates checked
6. **Disable Actions** - "SAVE GO DECISION" button disabled until timer complete

---

## Implementation

**Create** `ui/src/lib/stores/timer.ts`:
```typescript
import { writable, derived } from 'svelte/store';

interface TimerState {
    startTime: number | null;
    duration: number; // 120 seconds
    running: boolean;
}

export const timer = writable<TimerState>({
    startTime: null,
    duration: 120,
    running: false
});

export function startCoolOff() {
    timer.set({
        startTime: Date.now(),
        duration: 120,
        running: true
    });
}

export const timeRemaining = derived(timer, ($timer, set) => {
    if (!$timer.running || !$timer.startTime) {
        set(0);
        return;
    }

    const interval = setInterval(() => {
        const elapsed = Math.floor((Date.now() - $timer.startTime!) / 1000);
        const remaining = Math.max(0, $timer.duration - elapsed);
        set(remaining);

        if (remaining === 0) {
            timer.update(t => ({ ...t, running: false }));
            clearInterval(interval);
        }
    }, 1000);

    return () => clearInterval(interval);
});

// Format as MM:SS
export const timeRemainingFormatted = derived(timeRemaining, $time => {
    const minutes = Math.floor($time / 60);
    const seconds = $time % 60;
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
});
```

**Update Checklist Screen:**
- "Save Evaluation" button calls `startCoolOff()` after API success
- Display timer: "Cool-off: {$timeRemainingFormatted} remaining"
- Show completion message when timer reaches 0:00

**Create Timer Component** `ui/src/lib/components/common/CoolOffTimer.svelte`:
- Displays countdown with gradient styling
- Shows completion message
- Logs when timer starts/completes

**Backend stores timestamp:**
- `POST /api/checklist/evaluate` saves evaluation timestamp
- `POST /api/gates/check` validates 2+ minutes elapsed

---

## Expected Outcome

User saves GREEN checklist evaluation. Timer starts: "2:00 remaining". Counts down every second. After 2 minutes: "0:00 - Cool-off complete". Backend validates elapsed time when user tries to save GO decision.

---

## Phase 2 Complete!

After this step, users can:
- See large gradient banner (RED/YELLOW/GREEN)
- Complete checklist with 5 required gates
- Add optional quality items (4 items)
- Calculate position sizing (Van Tharp method)
- Honor 2-minute cool-off period

**Next Phase:** [Phase 3: Heat Check & Trade Entry](../plans/phase3-step15-heat-check.md)

---

## Time Estimate

~4-6 hours (1 day)

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
