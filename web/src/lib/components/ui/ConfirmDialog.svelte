<script lang="ts">
    import { AlertTriangle } from 'lucide-svelte';
    import { fade, fly } from 'svelte/transition';

    let {
        show = $bindable(false),
        title = 'Konfirmasi',
        message = 'Apakah Anda yakin?',
        confirmLabel = 'Hapus',
        cancelLabel = 'Batal',
        onConfirm = undefined,
        onCancel = undefined
    } = $props();

    function handleConfirm() {
        show = false;
        onConfirm?.();
    }

    function handleCancel() {
        show = false;
        onCancel?.();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Escape') handleCancel();
    }
</script>

{#if show}
    <div
        class="backdrop"
        transition:fade={{ duration: 180 }}
        onclick={handleCancel}
        onkeydown={handleKeydown}
        role="button"
        tabindex="0"
    >
        <div
            class="dialog"
            transition:fly={{ y: -30, duration: 250 }}
            onclick={(e) => e.stopPropagation()}
            onkeydown={(e) => e.stopPropagation()}
            role="alertdialog"
            aria-modal="true"
            tabindex="-1"
        >
            <div class="icon-wrap">
                <AlertTriangle size={28} />
            </div>
            <h2 class="d-title">{title}</h2>
            <p class="d-message">{message}</p>
            <div class="d-actions">
                <button class="btn btn-ghost" onclick={handleCancel}>{cancelLabel}</button>
                <button class="btn btn-danger" onclick={handleConfirm}>{confirmLabel}</button>
            </div>
        </div>
    </div>
{/if}

<style>
    .backdrop {
        position: fixed;
        inset: 0;
        background: rgba(15, 23, 42, 0.65);
        backdrop-filter: blur(4px);
        z-index: 2000;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: var(--space-md);
    }

    .dialog {
        background: white;
        border-radius: var(--radius-xl);
        padding: var(--space-xl);
        max-width: 380px;
        width: 100%;
        box-shadow: 0 25px 50px -12px rgba(0,0,0,0.35);
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--space-sm);
        text-align: center;
    }

    .icon-wrap {
        width: 56px;
        height: 56px;
        border-radius: var(--radius-full);
        background: #fef2f2;
        color: var(--color-error);
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: var(--space-xs);
    }

    .d-title {
        font-size: 1.1rem;
        font-weight: 700;
        color: var(--text-primary);
        margin: 0;
    }

    .d-message {
        font-size: 0.9rem;
        color: var(--text-secondary);
        line-height: 1.5;
        margin: 0;
        white-space: pre-wrap;
    }

    .d-actions {
        display: flex;
        gap: var(--space-sm);
        width: 100%;
        margin-top: var(--space-sm);
    }

    .btn-ghost {
        flex: 1;
        background: var(--bg-main);
        color: var(--text-secondary);
        border: 1px solid var(--border-color);
        padding: var(--space-sm) var(--space-md);
        border-radius: var(--radius-md);
        font-weight: 600;
        cursor: pointer;
        font-size: 0.9rem;
        transition: all var(--transition-fast);
    }
    .btn-ghost:hover { background: #e2e8f0; }

    .btn-danger {
        flex: 1;
        background: var(--color-error);
        color: white;
        border: 1px solid transparent;
        padding: var(--space-sm) var(--space-md);
        border-radius: var(--radius-md);
        font-weight: 600;
        cursor: pointer;
        font-size: 0.9rem;
        transition: all var(--transition-fast);
    }
    .btn-danger:hover { background: #b91c1c; }
</style>
