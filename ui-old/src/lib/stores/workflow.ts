// Workflow store for managing trade evaluation state across screens
import { writable, derived } from 'svelte/store';
import type { SizingResult } from '$lib/api/client';

export interface TradeWorkflowState {
    // Basic trade information
    ticker: string;
    entryPrice: number | null;
    atrN: number | null;
    sector: string;
    method: 'stock' | 'opt-delta-atr' | 'opt-maxloss';

    // Checklist results
    checklistComplete: boolean;
    bannerStatus: 'RED' | 'YELLOW' | 'GREEN' | null;
    requiredGatesPassed: number;
    qualityScore: number;
    checklistTimestamp: number | null;

    // Position sizing results
    sizingComplete: boolean;
    sizingResult: SizingResult | null;
    shares: number | null;
    contracts: number | null;
    riskDollars: number | null;

    // Heat check results
    heatCheckComplete: boolean;
    heatCheckPassed: boolean | null;
    portfolioHeat: number | null;
    bucketHeat: number | null;

    // Workflow tracking
    currentStep: 'scanner' | 'checklist' | 'sizing' | 'heat' | 'entry' | null;
    workflowStarted: boolean;
}

const initialState: TradeWorkflowState = {
    ticker: '',
    entryPrice: null,
    atrN: null,
    sector: '',
    method: 'stock',

    checklistComplete: false,
    bannerStatus: null,
    requiredGatesPassed: 0,
    qualityScore: 0,
    checklistTimestamp: null,

    sizingComplete: false,
    sizingResult: null,
    shares: null,
    contracts: null,
    riskDollars: null,

    heatCheckComplete: false,
    heatCheckPassed: null,
    portfolioHeat: null,
    bucketHeat: null,

    currentStep: null,
    workflowStarted: false
};

function createWorkflowStore() {
    const { subscribe, set, update } = writable<TradeWorkflowState>(initialState);

    return {
        subscribe,

        // Start a new trade evaluation
        startTrade: (ticker: string, sector: string = '') => {
            console.log('[Workflow] Starting new trade evaluation:', ticker);
            set({
                ...initialState,
                ticker: ticker.toUpperCase(),
                sector,
                currentStep: 'checklist',
                workflowStarted: true
            });
        },

        // Update basic trade info
        updateTradeInfo: (ticker: string, entryPrice: number, atrN: number, sector: string, method: 'stock' | 'opt-delta-atr' | 'opt-maxloss' = 'stock') => {
            console.log('[Workflow] Updating trade info:', { ticker, entryPrice, atrN, sector, method });
            update(state => ({
                ...state,
                ticker: ticker.toUpperCase(),
                entryPrice,
                atrN,
                sector,
                method
            }));
        },

        // Save checklist results
        saveChecklistResults: (bannerStatus: 'RED' | 'YELLOW' | 'GREEN', requiredGatesPassed: number, qualityScore: number) => {
            console.log('[Workflow] Saving checklist results:', { bannerStatus, requiredGatesPassed, qualityScore });
            update(state => ({
                ...state,
                checklistComplete: true,
                bannerStatus,
                requiredGatesPassed,
                qualityScore,
                checklistTimestamp: Date.now(),
                currentStep: 'sizing'
            }));
        },

        // Save position sizing results
        saveSizingResults: (result: SizingResult) => {
            console.log('[Workflow] Saving sizing results:', result);
            update(state => ({
                ...state,
                sizingComplete: true,
                sizingResult: result,
                shares: result.shares || null,
                contracts: result.contracts || null,
                riskDollars: result.risk_dollars,
                currentStep: 'heat'
            }));
        },

        // Save heat check results
        saveHeatResults: (passed: boolean, portfolioHeat: number, bucketHeat: number) => {
            console.log('[Workflow] Saving heat check results:', { passed, portfolioHeat, bucketHeat });
            update(state => ({
                ...state,
                heatCheckComplete: true,
                heatCheckPassed: passed,
                portfolioHeat,
                bucketHeat,
                currentStep: 'entry'
            }));
        },

        // Move to specific step
        goToStep: (step: TradeWorkflowState['currentStep']) => {
            console.log('[Workflow] Moving to step:', step);
            update(state => ({
                ...state,
                currentStep: step
            }));
        },

        // Reset workflow
        reset: () => {
            console.log('[Workflow] Resetting workflow');
            set(initialState);
        },

        // Get current state snapshot
        getState: () => {
            let currentState: TradeWorkflowState = initialState;
            subscribe(value => { currentState = value; })();
            return currentState;
        }
    };
}

export const workflow = createWorkflowStore();

// Derived store: Is workflow ready for trade entry?
export const readyForEntry = derived(workflow, $workflow => {
    return $workflow.checklistComplete &&
           $workflow.bannerStatus === 'GREEN' &&
           $workflow.sizingComplete &&
           $workflow.heatCheckComplete &&
           $workflow.heatCheckPassed === true;
});

// Derived store: Workflow completion percentage
export const workflowProgress = derived(workflow, $workflow => {
    if (!$workflow.workflowStarted) return 0;

    let completed = 0;
    const total = 4; // checklist, sizing, heat, entry

    if ($workflow.checklistComplete) completed++;
    if ($workflow.sizingComplete) completed++;
    if ($workflow.heatCheckComplete) completed++;

    return Math.round((completed / total) * 100);
});

// Derived store: Next recommended step
export const nextStep = derived(workflow, $workflow => {
    if (!$workflow.workflowStarted) return 'scanner';
    if (!$workflow.checklistComplete) return 'checklist';
    if (!$workflow.sizingComplete) return 'sizing';
    if (!$workflow.heatCheckComplete) return 'heat';
    return 'entry';
});

// Derived store: Can proceed to next step?
export const canProceed = derived(workflow, $workflow => {
    switch ($workflow.currentStep) {
        case 'checklist':
            return $workflow.checklistComplete && $workflow.bannerStatus === 'GREEN';
        case 'sizing':
            return $workflow.sizingComplete;
        case 'heat':
            return $workflow.heatCheckComplete && $workflow.heatCheckPassed === true;
        case 'entry':
            return true;
        default:
            return false;
    }
});
