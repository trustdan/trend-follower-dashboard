<script lang="ts">
  import { onMount } from 'svelte';

  // Props
  type BannerState = 'RED' | 'YELLOW' | 'GREEN';

  interface Props {
    state: BannerState;
    message: string;
    details?: string;
  }

  let { state = $bindable('GREEN'), message, details }: Props = $props();

  // Track previous state for pulse animation
  let previousState = $state<BannerState | null>(null);
  let shouldPulse = $state(false);

  // Watch for state changes to trigger pulse effect
  $effect(() => {
    if (previousState !== null && previousState !== state) {
      shouldPulse = true;
      setTimeout(() => {
        shouldPulse = false;
      }, 600); // Pulse duration
    }
    previousState = state;
  });

  // Define gradient classes for each state
  const stateClasses = {
    RED: 'bg-gradient-to-br from-red-600 to-red-800 shadow-red-500/50',
    YELLOW: 'bg-gradient-to-br from-amber-500 to-yellow-400 shadow-yellow-500/50',
    GREEN: 'bg-gradient-to-br from-emerald-500 to-emerald-600 shadow-green-500/50'
  };

  // Icons for each state
  const stateIcons = {
    RED: '🛑',
    YELLOW: '⚠️',
    GREEN: '✓'
  };

  // Accessibility labels
  const stateLabels = {
    RED: 'Do Not Trade',
    YELLOW: 'Caution',
    GREEN: 'OK to Trade'
  };
</script>

<!--
  Banner Component
  Large 3-state gradient banner for immediate feedback
  States: RED (stop), YELLOW (caution), GREEN (go)
-->
<div
  class="banner-container relative w-full min-h-[150px] h-[20vh] rounded-xl overflow-hidden transition-all duration-300 ease-in-out
         {stateClasses[state]} shadow-2xl
         {shouldPulse ? 'animate-pulse' : ''}"
  role="status"
  aria-live="polite"
  aria-label="{stateLabels[state]}: {message}"
>
  <!-- Glow effect overlay -->
  <div class="absolute inset-0 bg-gradient-to-t from-black/10 to-transparent pointer-events-none"></div>

  <!-- Content -->
  <div class="relative h-full flex flex-col items-center justify-center px-8 py-6 text-center">
    <!-- Icon -->
    <div class="text-6xl mb-4 drop-shadow-lg transition-transform duration-300 ease-in-out
                {shouldPulse ? 'scale-125' : 'scale-100'}">
      {stateIcons[state]}
    </div>

    <!-- Main message -->
    <h1 class="text-4xl md:text-5xl lg:text-6xl font-bold text-white drop-shadow-lg mb-2 transition-opacity duration-300">
      {message}
    </h1>

    <!-- Details (optional subtext) -->
    {#if details}
      <p class="text-lg md:text-xl text-white/90 drop-shadow-md max-w-3xl transition-opacity duration-300">
        {details}
      </p>
    {/if}

    <!-- State label for accessibility -->
    <div class="sr-only">{stateLabels[state]}</div>
  </div>

  <!-- Subtle animated gradient overlay on pulse -->
  {#if shouldPulse}
    <div class="absolute inset-0 bg-gradient-radial from-white/20 via-transparent to-transparent animate-ping pointer-events-none"></div>
  {/if}
</div>

<style>
  /* Custom pulse animation that's more subtle than default */
  @keyframes subtle-pulse {
    0%, 100% {
      opacity: 1;
      transform: scale(1);
    }
    50% {
      opacity: 0.95;
      transform: scale(1.02);
    }
  }

  .animate-pulse {
    animation: subtle-pulse 0.6s ease-in-out;
  }

  /* Radial gradient for pulse overlay */
  .bg-gradient-radial {
    background: radial-gradient(circle at center, currentColor 0%, transparent 70%);
  }

  /* Screen reader only class */
  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border-width: 0;
  }

  /* Smooth transitions for all interactive elements */
  .banner-container * {
    transition-property: transform, opacity;
    transition-timing-function: ease-in-out;
    transition-duration: 0.3s;
  }
</style>
