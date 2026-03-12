<script lang="ts">
    import { X } from 'lucide-svelte';
    import { fade, fly } from 'svelte/transition';

    let { show = $bindable(false), title, children } = $props();

    function close() {
        show = false;
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Escape') close();
    }
</script>

{#if show}
    <div 
        class="modal-backdrop" 
        transition:fade={{ duration: 200 }} 
        onclick={close}
        onkeydown={handleKeydown}
        role="button"
        tabindex="0"
    >
        <div 
            class="modal-content" 
            transition:fly={{ y: 100, duration: 300 }}
            onclick={(e) => e.stopPropagation()}
            onkeydown={(e) => e.stopPropagation()}
            role="dialog"
            aria-modal="true"
            aria-labelledby="modal-title"
            tabindex="-1"
        >
            <header class="modal-header">
                <h2 id="modal-title">{title}</h2>
                <button class="close-btn" onclick={close}>
                    <X size={24} />
                </button>
            </header>
            
            <div class="modal-body">
                {@render children()}
            </div>
        </div>
    </div>
{/if}

<style>
    .modal-backdrop {
        position: fixed;
        inset: 0;
        background: rgba(15, 23, 42, 0.7);
        backdrop-filter: blur(4px);
        z-index: 1000;
        display: flex;
        align-items: flex-end;
        justify-content: center;
    }

    @media (min-width: 768px) {
        .modal-backdrop {
            align-items: center;
        }
    }

    .modal-content {
        background: white;
        width: 100%;
        max-width: 600px;
        border-radius: var(--radius-xl) var(--radius-xl) 0 0;
        max-height: 90vh;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        box-shadow: var(--shadow-lg);
    }

    @media (min-width: 768px) {
        .modal-content {
            border-radius: var(--radius-xl);
            width: 90%;
        }
    }

    .modal-header {
        padding: var(--space-lg);
        display: flex;
        align-items: center;
        justify-content: space-between;
        border-bottom: 1px solid var(--border-color);
        position: sticky;
        top: 0;
        background: white;
        z-index: 10;
    }

    .modal-header h2 {
        font-size: 1.25rem;
        color: var(--color-primary);
    }

    .close-btn {
        background: var(--bg-main);
        border: none;
        width: 36px;
        height: 36px;
        border-radius: var(--radius-full);
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        color: var(--text-muted);
        transition: all var(--transition-fast);
    }

    .close-btn:hover {
        background: #fee2e2;
        color: var(--color-error);
    }

    .modal-body {
        padding: var(--space-lg);
    }
</style>
