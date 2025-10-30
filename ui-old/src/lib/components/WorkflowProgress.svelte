<script lang="ts">
    import { workflow, workflowProgress, nextStep, readyForEntry } from '$lib/stores/workflow';
    import { Check, Circle, ChevronRight } from 'lucide-svelte';

    const steps = [
        { id: 'scanner', name: 'Scanner', description: 'Find candidates' },
        { id: 'checklist', name: 'Checklist', description: '5 gates + quality' },
        { id: 'sizing', name: 'Position Sizing', description: 'Calculate shares' },
        { id: 'heat', name: 'Heat Check', description: 'Verify caps' },
        { id: 'entry', name: 'Trade Entry', description: 'Final decision' }
    ];

    function getStepStatus(stepId: string): 'completed' | 'current' | 'pending' {
        const state = $workflow;

        if (state.currentStep === stepId) return 'current';

        // Determine if step is completed
        switch (stepId) {
            case 'scanner':
                return state.ticker ? 'completed' : 'pending';
            case 'checklist':
                return state.checklistComplete ? 'completed' : 'pending';
            case 'sizing':
                return state.sizingComplete ? 'completed' : 'pending';
            case 'heat':
                return state.heatCheckComplete ? 'completed' : 'pending';
            case 'entry':
                return 'pending'; // Entry is the final step
            default:
                return 'pending';
        }
    }

    function getStepColor(status: string): string {
        switch (status) {
            case 'completed':
                return 'bg-emerald-500 border-emerald-600 text-white';
            case 'current':
                return 'bg-blue-500 border-blue-600 text-white animate-pulse';
            case 'pending':
            default:
                return 'bg-gray-300 dark:bg-gray-700 border-gray-400 dark:border-gray-600 text-gray-600 dark:text-gray-400';
        }
    }

    function getConnectorColor(completedBefore: boolean): string {
        return completedBefore
            ? 'bg-emerald-500'
            : 'bg-gray-300 dark:bg-gray-700';
    }
</script>

{#if $workflow.workflowStarted}
<div class="mb-6">
    <!-- Header with progress percentage -->
    <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-3">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                Trade Workflow Progress
            </h3>
            {#if $workflow.ticker}
                <span class="px-2 py-1 text-sm font-mono font-semibold bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded">
                    {$workflow.ticker}
                </span>
            {/if}
        </div>
        <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600 dark:text-gray-400">
                {$workflowProgress}% Complete
            </span>
            {#if $readyForEntry}
                <span class="px-2 py-1 text-xs font-semibold bg-emerald-100 dark:bg-emerald-900 text-emerald-800 dark:text-emerald-200 rounded-full">
                    Ready for Entry
                </span>
            {/if}
        </div>
    </div>

    <!-- Progress bar -->
    <div class="relative mb-6">
        <div class="overflow-hidden h-2 text-xs flex rounded-full bg-gray-200 dark:bg-gray-700">
            <div
                style="width: {$workflowProgress}%"
                class="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-gradient-to-r from-emerald-500 to-emerald-600 transition-all duration-500"
            ></div>
        </div>
    </div>

    <!-- Step indicators -->
    <div class="flex items-center justify-between">
        {#each steps as step, index}
            {@const status = getStepStatus(step.id)}
            {@const stepColor = getStepColor(status)}
            {@const isCompleted = status === 'completed'}
            {@const isCurrent = status === 'current'}

            <!-- Step circle -->
            <div class="flex flex-col items-center relative" style="flex: 1">
                <div
                    class="w-12 h-12 rounded-full border-2 flex items-center justify-center transition-all duration-300 {stepColor}"
                    class:shadow-lg={isCurrent}
                >
                    {#if isCompleted}
                        <Check class="w-6 h-6" />
                    {:else if isCurrent}
                        <Circle class="w-6 h-6 fill-current" />
                    {:else}
                        <Circle class="w-6 h-6" />
                    {/if}
                </div>

                <!-- Step name and description -->
                <div class="text-center mt-2">
                    <div class="text-sm font-semibold" class:text-emerald-600={isCompleted} class:text-blue-600={isCurrent} class:text-gray-600={status === 'pending'} class:dark:text-emerald-400={isCompleted} class:dark:text-blue-400={isCurrent} class:dark:text-gray-400={status === 'pending'}>
                        {step.name}
                    </div>
                    <div class="text-xs text-gray-500 dark:text-gray-400">
                        {step.description}
                    </div>
                </div>

                <!-- Connector line (not for last step) -->
                {#if index < steps.length - 1}
                    {@const prevStepCompleted = getStepStatus(steps[index].id) === 'completed'}
                    <div
                        class="absolute top-6 left-1/2 w-full h-0.5 -z-10 transition-all duration-500 {getConnectorColor(prevStepCompleted)}"
                        style="transform: translateX(0%);"
                    ></div>
                {/if}
            </div>
        {/each}
    </div>

    <!-- Next step hint -->
    {#if $nextStep && $nextStep !== 'entry'}
        <div class="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
            <div class="flex items-center gap-2 text-sm text-blue-800 dark:text-blue-200">
                <ChevronRight class="w-4 h-4" />
                <span class="font-semibold">Next:</span>
                <span>Complete the {steps.find(s => s.id === $nextStep)?.name} step</span>
            </div>
        </div>
    {:else if $readyForEntry}
        <div class="mt-4 p-3 bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 rounded-lg">
            <div class="flex items-center gap-2 text-sm text-emerald-800 dark:text-emerald-200">
                <Check class="w-4 h-4" />
                <span class="font-semibold">Ready:</span>
                <span>All gates passed - proceed to Trade Entry for final decision</span>
            </div>
        </div>
    {/if}
</div>
{/if}

<style>
    @keyframes pulse {
        0%, 100% {
            opacity: 1;
        }
        50% {
            opacity: .7;
        }
    }

    .animate-pulse {
        animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
    }
</style>
