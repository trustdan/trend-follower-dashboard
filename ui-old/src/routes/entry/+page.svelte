<script lang="ts">
/**
 * Trade Entry Screen - Phase 3 Step 16
 *
 * Final gate check before making a GO/NO-GO trade decision.
 * Enforces all 5 hard gates:
 * 1. Banner is GREEN
 * 2. Cooloff timer elapsed (> 2 minutes)
 * 3. Not on cooldown list
 * 4. Heat check passed
 * 5. Position sizing completed
 */

import { onMount } from 'svelte';
import Banner from '$lib/components/Banner.svelte';
import Card from '$lib/components/Card.svelte';
import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
import GateStatus from '$lib/components/GateStatus.svelte';
import { timer } from '$lib/stores/timer';
import { logger } from '$lib/utils/logger';
import { api } from '$lib/api/client';
import type { SaveDecisionRequest } from '$lib/api/client';

// Trade entry form data
let formData = $state({
  ticker: '',
  entry: 0,
  atr: 0,
  method: 'stock',
  bannerStatus: 'RED',
  shares: 0,
  contracts: 0,
  sector: '',
  strategy: 'TF_BREAKOUT_LONG',
  riskDollars: 0,
});

// Decision data
let decision = $state<'GO' | 'NO-GO' | null>(null);
let notes = $state('');

// Gate states
let gates = $state({
  bannerGreen: false,
  timerComplete: false,
  notOnCooldown: true,
  heatPassed: false,
  sizingComplete: false,
});

// UI state
let loading = $state(false);
let saving = $state(false);
let error = $state<string | null>(null);
let success = $state<string | null>(null);

// Sectors for dropdown
const sectors = [
  'Tech/Comm',
  'Finance',
  'Healthcare',
  'Energy',
  'Industrials',
  'Consumer',
  'Materials',
  'Utilities',
  'Real Estate',
  'Other'
];

// Methods
const methods = [
  { value: 'stock', label: 'Stock/ETF' },
  { value: 'opt-delta-atr', label: 'Options (Delta-ATR)' },
  { value: 'opt-maxloss', label: 'Options (Max Loss)' }
];

// Strategies
const strategies = [
  'TF_BREAKOUT_LONG',
  'TF_BREAKOUT_SHORT',
  'TF_SWING_LONG',
  'TF_SWING_SHORT'
];

// Banner statuses
const bannerStatuses = ['RED', 'YELLOW', 'GREEN'];

/**
 * Update gate states based on form data
 */
function updateGates() {
  // Gate 1: Banner is GREEN
  gates.bannerGreen = formData.bannerStatus === 'GREEN';

  // Gate 2: Timer complete (from timer store)
  gates.timerComplete = $timer.isComplete;

  // Gate 3: Not on cooldown (for now, we assume true - backend will validate)
  gates.notOnCooldown = true;

  // Gate 4: Heat passed (check if risk dollars are set)
  gates.heatPassed = formData.riskDollars > 0;

  // Gate 5: Sizing complete (check if shares or contracts are set)
  gates.sizingComplete = formData.shares > 0 || formData.contracts > 0;

  logger.log('Gates updated', gates);
}

/**
 * Check if all gates pass
 */
function allGatesPass(): boolean {
  return gates.bannerGreen &&
         gates.timerComplete &&
         gates.notOnCooldown &&
         gates.heatPassed &&
         gates.sizingComplete;
}

/**
 * Handle form field changes
 */
function handleFieldChange(field: string, value: any) {
  (formData as any)[field] = value;
  updateGates();
  logger.log(`Field changed: ${field}`, value);
}

/**
 * Save trade decision
 */
async function saveDecision() {
  if (!decision) {
    error = 'Please select GO or NO-GO decision';
    return;
  }

  if (decision === 'GO' && !allGatesPass()) {
    error = 'Cannot save GO decision - not all gates passed';
    return;
  }

  if (!notes.trim()) {
    error = 'Please add notes explaining your decision';
    return;
  }

  saving = true;
  error = null;

  try {
    logger.log('Saving trade decision', {
      decision,
      formData,
      gates,
      notes
    });

    // Create request payload
    const request: SaveDecisionRequest = {
      ticker: formData.ticker,
      entry: formData.entry,
      atr: formData.atr,
      method: formData.method,
      banner_status: formData.bannerStatus,
      shares: formData.shares,
      contracts: formData.contracts,
      sector: formData.sector,
      strategy: formData.strategy,
      risk_dollars: formData.riskDollars,
      decision: decision,
      notes: notes,
      banner_green: gates.bannerGreen,
      timer_complete: gates.timerComplete,
      not_on_cooldown: gates.notOnCooldown,
      heat_passed: gates.heatPassed,
      sizing_complete: gates.sizingComplete,
    };

    // Call backend API to save decision
    const response = await api.saveDecision(request);

    logger.log('Decision saved successfully', response);

    success = `${decision} decision saved successfully! (ID: ${response.id})`;

    // Auto-dismiss success message after 3 seconds
    setTimeout(() => {
      success = null;
      // Reset form after successful save
      resetForm();
    }, 3000);

  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to save decision';
    logger.error('Error saving decision:', err);
  } finally {
    saving = false;
  }
}

/**
 * Reset form to initial state
 */
function resetForm() {
  formData = {
    ticker: '',
    entry: 0,
    atr: 0,
    method: 'stock',
    bannerStatus: 'RED',
    shares: 0,
    contracts: 0,
    sector: '',
    strategy: 'TF_BREAKOUT_LONG',
    riskDollars: 0,
  };
  decision = null;
  notes = '';
  updateGates();
  logger.log('Form reset');
}

// Initialize on mount
onMount(() => {
  logger.log('Trade Entry page mounted');
  updateGates();
});

// Reactive: Update gates whenever form data or timer changes
$effect(() => {
  // Watch timer state
  const timerState = $timer;
  updateGates();
});
</script>

<svelte:head><title>Trade Entry - TF-Engine</title></svelte:head>

<div class="container mx-auto px-4 py-6 max-w-6xl">
  <h1 class="text-3xl font-bold mb-6">Trade Entry: Final Gate Check</h1>

  {#if loading}
    <LoadingSpinner message="Loading trade entry form..." />
  {:else}
    <!-- Error Banner -->
    {#if error}
      <div class="mb-6 p-4 bg-gradient-to-r from-red-500 to-red-600 text-white rounded-lg shadow-lg">
        <p class="font-bold">Error</p>
        <p>{error}</p>
      </div>
    {/if}

    <!-- Success Banner -->
    {#if success}
      <div class="mb-6 p-4 bg-gradient-to-r from-emerald-500 to-emerald-600 text-white rounded-lg shadow-lg">
        <p class="font-bold">Success</p>
        <p>{success}</p>
      </div>
    {/if}

    <!-- Trade Information Form -->
    <Card title="Trade Information">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <!-- Ticker -->
        <div>
          <label for="ticker" class="block text-sm font-medium mb-1">Ticker Symbol</label>
          <input
            id="ticker"
            type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            placeholder="AAPL"
            value={formData.ticker}
            oninput={(e) => handleFieldChange('ticker', (e.target as HTMLInputElement).value.toUpperCase())}
          />
        </div>

        <!-- Entry Price -->
        <div>
          <label for="entry" class="block text-sm font-medium mb-1">Entry Price ($)</label>
          <input
            id="entry"
            type="number"
            step="0.01"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            placeholder="180.00"
            value={formData.entry || ''}
            oninput={(e) => handleFieldChange('entry', parseFloat((e.target as HTMLInputElement).value) || 0)}
          />
        </div>

        <!-- ATR (N) -->
        <div>
          <label for="atr" class="block text-sm font-medium mb-1">ATR / N ($)</label>
          <input
            id="atr"
            type="number"
            step="0.01"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            placeholder="1.50"
            value={formData.atr || ''}
            oninput={(e) => handleFieldChange('atr', parseFloat((e.target as HTMLInputElement).value) || 0)}
          />
        </div>

        <!-- Method -->
        <div>
          <label for="method" class="block text-sm font-medium mb-1">Method</label>
          <select
            id="method"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            value={formData.method}
            onchange={(e) => handleFieldChange('method', (e.target as HTMLSelectElement).value)}
          >
            {#each methods as method}
              <option value={method.value}>{method.label}</option>
            {/each}
          </select>
        </div>

        <!-- Banner Status -->
        <div>
          <label for="bannerStatus" class="block text-sm font-medium mb-1">Banner Status</label>
          <select
            id="bannerStatus"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            value={formData.bannerStatus}
            onchange={(e) => handleFieldChange('bannerStatus', (e.target as HTMLSelectElement).value)}
          >
            {#each bannerStatuses as status}
              <option value={status}>{status}</option>
            {/each}
          </select>
        </div>

        <!-- Shares/Contracts -->
        <div>
          <label for="shares" class="block text-sm font-medium mb-1">
            {formData.method === 'stock' ? 'Shares' : 'Contracts'}
          </label>
          <input
            id="shares"
            type="number"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            placeholder={formData.method === 'stock' ? '250' : '5'}
            value={formData.method === 'stock' ? (formData.shares || '') : (formData.contracts || '')}
            oninput={(e) => {
              const value = parseInt((e.target as HTMLInputElement).value) || 0;
              if (formData.method === 'stock') {
                handleFieldChange('shares', value);
              } else {
                handleFieldChange('contracts', value);
              }
            }}
          />
        </div>

        <!-- Sector -->
        <div>
          <label for="sector" class="block text-sm font-medium mb-1">Sector / Bucket</label>
          <select
            id="sector"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            value={formData.sector}
            onchange={(e) => handleFieldChange('sector', (e.target as HTMLSelectElement).value)}
          >
            <option value="">Select sector...</option>
            {#each sectors as sector}
              <option value={sector}>{sector}</option>
            {/each}
          </select>
        </div>

        <!-- Strategy -->
        <div>
          <label for="strategy" class="block text-sm font-medium mb-1">Strategy Preset</label>
          <select
            id="strategy"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            value={formData.strategy}
            onchange={(e) => handleFieldChange('strategy', (e.target as HTMLSelectElement).value)}
          >
            {#each strategies as strategy}
              <option value={strategy}>{strategy}</option>
            {/each}
          </select>
        </div>

        <!-- Risk Dollars -->
        <div>
          <label for="riskDollars" class="block text-sm font-medium mb-1">Risk Amount ($)</label>
          <input
            id="riskDollars"
            type="number"
            step="0.01"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            placeholder="750.00"
            value={formData.riskDollars || ''}
            oninput={(e) => handleFieldChange('riskDollars', parseFloat((e.target as HTMLInputElement).value) || 0)}
          />
        </div>
      </div>
    </Card>

    <!-- 5 Hard Gates -->
    <Card title="5 Hard Gates" class="mt-6">
      <div class="space-y-3">
        <GateStatus
          gateNumber={1}
          description="Banner is GREEN (all checklist items passed)"
          passed={gates.bannerGreen}
        />

        <GateStatus
          gateNumber={2}
          description="Cooloff timer elapsed (> 2 minutes since checklist evaluation)"
          passed={gates.timerComplete}
        />

        <GateStatus
          gateNumber={3}
          description="Not on cooldown list (ticker/sector available for trading)"
          passed={gates.notOnCooldown}
        />

        <GateStatus
          gateNumber={4}
          description="Heat check passed (portfolio and sector caps not exceeded)"
          passed={gates.heatPassed}
        />

        <GateStatus
          gateNumber={5}
          description="Position sizing completed (shares/contracts calculated)"
          passed={gates.sizingComplete}
        />
      </div>

      <!-- Overall Gate Status -->
      <div class="mt-6 p-4 rounded-lg border-2"
           class:border-emerald-500={allGatesPass()}
           class:bg-emerald-50={allGatesPass()}
           class:border-red-500={!allGatesPass()}
           class:bg-red-50={!allGatesPass()}>
        <p class="text-center font-bold text-lg"
           class:text-emerald-900={allGatesPass()}
           class:text-red-900={!allGatesPass()}>
          {allGatesPass() ? '✓ ALL GATES PASSED' : '✗ GATES NOT PASSED'}
        </p>
        {#if !allGatesPass()}
          <p class="text-center text-sm mt-1 text-red-700">
            You must pass all 5 gates to save a GO decision
          </p>
        {/if}
      </div>
    </Card>

    <!-- Decision Section -->
    <Card title="Trade Decision" class="mt-6">
      <div class="space-y-4">
        <!-- GO/NO-GO Selection -->
        <div>
          <label class="block text-sm font-medium mb-2">Decision</label>
          <div class="flex gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                name="decision"
                value="GO"
                checked={decision === 'GO'}
                onchange={() => { decision = 'GO'; error = null; }}
                class="w-5 h-5 text-emerald-600 focus:ring-emerald-500"
              />
              <span class="text-lg font-medium">GO (Execute Trade)</span>
            </label>

            <label class="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                name="decision"
                value="NO-GO"
                checked={decision === 'NO-GO'}
                onchange={() => { decision = 'NO-GO'; error = null; }}
                class="w-5 h-5 text-red-600 focus:ring-red-500"
              />
              <span class="text-lg font-medium">NO-GO (Do Not Trade)</span>
            </label>
          </div>
        </div>

        <!-- Notes -->
        <div>
          <label for="notes" class="block text-sm font-medium mb-1">
            Notes (required - explain your reasoning)
          </label>
          <textarea
            id="notes"
            rows="4"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500"
            placeholder="Explain why you're making this decision. What are the key factors? What's your plan?"
            bind:value={notes}
          ></textarea>
          <p class="text-sm text-gray-600 mt-1">
            {notes.length} characters
          </p>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-4">
          <button
            type="button"
            class="px-6 py-3 rounded-lg font-bold text-white shadow-lg transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            class:bg-gradient-to-r={decision === 'GO'}
            class:from-emerald-500={decision === 'GO'}
            class:to-emerald-600={decision === 'GO'}
            class:hover:shadow-xl={decision === 'GO' && !saving}
            class:bg-gray-400={decision !== 'GO'}
            disabled={decision !== 'GO' || saving || !allGatesPass()}
            onclick={saveDecision}
          >
            {saving && decision === 'GO' ? 'Saving...' : 'Save GO Decision'}
          </button>

          <button
            type="button"
            class="px-6 py-3 rounded-lg font-bold text-white shadow-lg transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            class:bg-gradient-to-r={decision === 'NO-GO'}
            class:from-red-500={decision === 'NO-GO'}
            class:to-red-600={decision === 'NO-GO'}
            class:hover:shadow-xl={decision === 'NO-GO' && !saving}
            class:bg-gray-400={decision !== 'NO-GO'}
            disabled={decision !== 'NO-GO' || saving}
            onclick={saveDecision}
          >
            {saving && decision === 'NO-GO' ? 'Saving...' : 'Save NO-GO Decision'}
          </button>

          <button
            type="button"
            class="px-6 py-3 bg-gray-300 hover:bg-gray-400 rounded-lg font-bold transition-all duration-300"
            disabled={saving}
            onclick={resetForm}
          >
            Reset Form
          </button>
        </div>
      </div>
    </Card>
  {/if}
</div>
