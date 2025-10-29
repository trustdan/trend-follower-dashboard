<script lang="ts">
/**
 * GateStatus Component
 *
 * Displays the status of one of the 5 hard gates.
 * Shows a checkmark (✓) if passed, X (✗) if failed.
 * Green gradient for pass, red for fail.
 */

interface Props {
  gateNumber: number;
  description: string;
  passed: boolean;
}

let { gateNumber, description, passed }: Props = $props();
</script>

<div class="gate-item flex items-center gap-3 p-3 rounded-lg border-2 transition-all duration-300"
     class:border-emerald-500={passed}
     class:bg-emerald-50={passed}
     class:border-red-500={!passed}
     class:bg-red-50={!passed}>

  <!-- Gate Number -->
  <div class="flex-shrink-0 flex items-center justify-center w-8 h-8 rounded-full font-bold text-sm"
       class:bg-emerald-500={passed}
       class:text-white={passed}
       class:bg-red-500={!passed}>
    {gateNumber}
  </div>

  <!-- Gate Description -->
  <div class="flex-1">
    <p class="font-medium"
       class:text-emerald-900={passed}
       class:text-red-900={!passed}>
      {description}
    </p>
  </div>

  <!-- Status Icon -->
  <div class="flex-shrink-0 text-2xl font-bold"
       class:text-emerald-600={passed}
       class:text-red-600={!passed}>
    {passed ? '✓' : '✗'}
  </div>
</div>

<style>
  .gate-item {
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }
</style>
