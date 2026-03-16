<script lang="ts">
    import { api } from '$lib/api';
    import { Link2, User } from 'lucide-svelte';

    let { sessionId, people = [], preselectedPerson = null, onSuccess } = $props();

    let personAId = $state('');
    let personBId = $state('');
    let relType = $state('parent_child');
    let loading = $state(false);
    let error = $state('');

    // Sync personAId whenever preselectedPerson changes (including initial mount)
    $effect(() => {
        personAId = preselectedPerson?.id || '';
    });

    const relationshipLabels: Record<string, string> = {
        parent_child: 'Orang Tua → Anak',
        spouse: 'Pasangan / Suami-Istri'
    };

    async function handleSubmit(e: SubmitEvent) {
        e.preventDefault();
        if (!personAId || !personBId || personAId === personBId) {
            error = 'Pilih dua orang yang berbeda.';
            return;
        }

        loading = true;
        error = '';
        try {
            const rel = await api.relationships.create(sessionId, {
                person_a_id: personAId,
                person_b_id: personBId,
                type: relType
            });
            onSuccess(rel);
        } catch (err: any) {
            error = err.message;
        } finally {
            loading = false;
        }
    }

    function getPersonName(id: string) {
        return people.find((p: any) => p.id === id)?.name || '—';
    }
</script>

<form onsubmit={handleSubmit} class="rel-form">
    {#if error}
        <div class="alert error">{error}</div>
    {/if}

    <!-- Relationship Type -->
    <div class="form-group">
        <span class="group-label">Jenis Hubungan</span>
        <div class="type-grid">
            {#each Object.entries(relationshipLabels) as [value, label]}
                <label class="type-card {relType === value ? 'selected' : ''}">
                    <input type="radio" bind:group={relType} {value} class="sr-only" />
                    <Link2 size={16} />
                    <span>{label}</span>
                </label>
            {/each}
        </div>
    </div>

    <!-- Person A -->
    <div class="form-group">
        <label for="person-a" class="group-label">
            {relType === 'parent_child' ? '👤 Orang Tua (Person A)' : '👤 Orang Pertama (Person A)'}
        </label>
        <div class="select-wrapper">
            <User size={16} class="select-icon" />
            <select id="person-a" bind:value={personAId} required>
                <option value="">-- Pilih orang --</option>
                {#each people as person}
                    <option value={person.id}>{person.name}</option>
                {/each}
            </select>
        </div>
    </div>

    <!-- Arrow indicator -->
    <div class="arrow-indicator">
        {relType === 'parent_child' ? '⬇ adalah orang tua dari' : '↔ berpasangan dengan'}
    </div>

    <!-- Person B -->
    <div class="form-group">
        <label for="person-b" class="group-label">
            {relType === 'parent_child' ? '👤 Anak (Person B)' : '👤 Orang Kedua (Person B)'}
        </label>
        <div class="select-wrapper">
            <User size={16} class="select-icon" />
            <select id="person-b" bind:value={personBId} required>
                <option value="">-- Pilih orang --</option>
                {#each people.filter((p: any) => p.id !== personAId) as person}
                    <option value={person.id}>{person.name}</option>
                {/each}
            </select>
        </div>
    </div>

    <button type="submit" class="btn btn-primary btn-block" disabled={loading || !personAId || !personBId}>
        {loading ? 'Menyimpan...' : 'Tambah Hubungan'}
    </button>
</form>

<style>
    .rel-form {
        display: flex;
        flex-direction: column;
        gap: var(--space-md);
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: var(--space-xs);
    }

    .group-label {
        font-size: 0.875rem;
        font-weight: 600;
        color: var(--text-secondary);
    }

    .type-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--space-sm);
    }

    .type-card {
        display: flex;
        align-items: center;
        gap: var(--space-xs);
        padding: var(--space-sm);
        border: 2px solid var(--border-color);
        border-radius: var(--radius-md);
        cursor: pointer;
        font-size: 0.8rem;
        font-weight: 500;
        transition: all var(--transition-fast);
        background: white;
    }

    .type-card.selected {
        border-color: var(--color-primary);
        background: #eff6ff;
        color: var(--color-primary);
    }

    .type-card:hover:not(.selected) {
        border-color: #94a3b8;
    }

    .select-wrapper {
        position: relative;
        display: flex;
        align-items: center;
    }

    :global(.select-icon) {
        position: absolute;
        left: var(--space-sm);
        color: var(--text-muted);
        pointer-events: none;
    }

    select {
        width: 100%;
        padding: var(--space-sm) var(--space-sm) var(--space-sm) 2.2rem;
        border: 1px solid var(--border-color);
        border-radius: var(--radius-md);
        font-family: inherit;
        font-size: 0.9rem;
        background: white;
        appearance: none;
        cursor: pointer;
    }

    select:focus {
        outline: none;
        border-color: var(--color-primary);
        box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
    }

    .arrow-indicator {
        text-align: center;
        font-size: 0.8rem;
        color: var(--text-muted);
        padding: var(--space-xs) 0;
        font-style: italic;
    }

    .btn-block {
        width: 100%;
        margin-top: var(--space-sm);
    }

    .alert {
        padding: var(--space-sm);
        border-radius: var(--radius-md);
        font-size: 0.875rem;
    }

    .alert.error {
        background: #fef2f2;
        color: var(--color-error);
        border: 1px solid #fee2e2;
    }

    .sr-only {
        position: absolute;
        width: 1px;
        height: 1px;
        overflow: hidden;
        clip: rect(0,0,0,0);
    }
</style>
