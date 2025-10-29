<script lang="ts">
    import { timer, timeRemaining, timeRemainingFormatted } from '$lib/stores/timer';
    import { Clock, CheckCircle2 } from 'lucide-svelte';

    // Track if timer was just completed
    let justCompleted = $state(false);

    // Watch for timer completion
    $effect(() => {
        const remaining = $timeRemaining;
        const running = $timer.running;

        // If timer just stopped and reached 0
        if (!running && remaining === 0 && !justCompleted) {
            justCompleted = true;
            console.log('[CoolOffTimer] Timer completed at', new Date().toLocaleTimeString());

            // Reset the "just completed" flag after 5 seconds
            setTimeout(() => {
                justCompleted = false;
            }, 5000);
        }
    });
</script>

{#if $timer.running}
    <!-- Timer is running - show countdown -->
    <div
        class="cool-off-timer bg-gradient-to-r from-amber-500 to-yellow-400 text-white rounded-lg p-6 shadow-xl border-2 border-amber-600"
        role="status"
        aria-live="polite"
        aria-label="Cool-off timer: {$timeRemainingFormatted} remaining"
    >
        <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
                <Clock class="w-8 h-8 animate-pulse" />
                <div>
                    <div class="text-sm font-medium opacity-90">Impulse Brake Active</div>
                    <div class="text-2xl font-bold tracking-wide">
                        {$timeRemainingFormatted}
                    </div>
                </div>
            </div>
            <div class="text-right">
                <div class="text-sm opacity-90">remaining</div>
                <div class="text-xs opacity-75 mt-1">
                    {Math.floor($timeRemaining / 60)}m {$timeRemaining % 60}s
                </div>
            </div>
        </div>
        <div class="mt-4 text-sm opacity-90">
            <p>This 2-minute pause prevents impulsive decisions.</p>
            <p class="mt-1">Use this time to review your analysis.</p>
        </div>
    </div>
{:else if justCompleted}
    <!-- Timer just completed - show success message -->
    <div
        class="cool-off-timer bg-gradient-to-r from-emerald-500 to-emerald-600 text-white rounded-lg p-6 shadow-xl border-2 border-emerald-700 animate-pulse"
        role="status"
        aria-live="assertive"
        aria-label="Cool-off period complete"
    >
        <div class="flex items-center gap-3">
            <CheckCircle2 class="w-8 h-8" />
            <div>
                <div class="text-sm font-medium opacity-90">Cool-Off Complete</div>
                <div class="text-2xl font-bold">
                    Ready to Proceed
                </div>
            </div>
        </div>
        <div class="mt-4 text-sm opacity-90">
            <p>You may now proceed with trade entry if all gates pass.</p>
        </div>
    </div>
{:else if $timer.startTime}
    <!-- Timer completed (not just recently) - show completed state -->
    <div
        class="cool-off-timer bg-gradient-to-r from-gray-600 to-gray-700 text-white rounded-lg p-4 shadow-lg border border-gray-500"
        role="status"
        aria-label="Cool-off period completed"
    >
        <div class="flex items-center gap-3">
            <CheckCircle2 class="w-6 h-6 text-emerald-400" />
            <div class="text-sm">
                <div class="font-medium">Cool-off completed</div>
                <div class="opacity-75 text-xs">You may proceed with trade entry</div>
            </div>
        </div>
    </div>
{/if}

<style>
    .cool-off-timer {
        transition: all 0.3s ease-in-out;
    }

    @keyframes pulse-glow {
        0%, 100% {
            box-shadow: 0 0 20px rgba(251, 191, 36, 0.4);
        }
        50% {
            box-shadow: 0 0 30px rgba(251, 191, 36, 0.6);
        }
    }

    .cool-off-timer:first-child {
        animation: pulse-glow 2s ease-in-out infinite;
    }
</style>
