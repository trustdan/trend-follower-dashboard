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
    console.log('[Timer] Starting 2-minute cool-off period');
    timer.set({
        startTime: Date.now(),
        duration: 120,
        running: true
    });
}

export function stopTimer() {
    console.log('[Timer] Stopping timer');
    timer.update(t => ({ ...t, running: false }));
}

export function resetTimer() {
    console.log('[Timer] Resetting timer');
    timer.set({
        startTime: null,
        duration: 120,
        running: false
    });
}

// Derived store that calculates remaining time
export const timeRemaining = derived(timer, ($timer, set) => {
    if (!$timer.running || !$timer.startTime) {
        set(0);
        return;
    }

    // Initial calculation
    const elapsed = Math.floor((Date.now() - $timer.startTime!) / 1000);
    const remaining = Math.max(0, $timer.duration - elapsed);
    set(remaining);

    // Set up interval to update every second
    const interval = setInterval(() => {
        const elapsed = Math.floor((Date.now() - $timer.startTime!) / 1000);
        const remaining = Math.max(0, $timer.duration - elapsed);
        set(remaining);

        if (remaining === 0) {
            console.log('[Timer] Cool-off period complete');
            timer.update(t => ({ ...t, running: false }));
            clearInterval(interval);
        }
    }, 1000);

    // Cleanup function
    return () => clearInterval(interval);
});

// Format as MM:SS
export const timeRemainingFormatted = derived(timeRemaining, $time => {
    const minutes = Math.floor($time / 60);
    const seconds = $time % 60;
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
});
