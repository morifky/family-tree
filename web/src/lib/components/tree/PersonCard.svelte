<script lang="ts">
    import { User, Edit, Trash2, Plus } from 'lucide-svelte';

    let { person, accessType, onEdit, onDelete, onAddFamily } = $props();

    const isMale = $derived(person.gender === 'male');
    const accentColor = $derived(isMale ? 'var(--color-primary)' : 'var(--color-secondary)');
    const bgAccent = $derived(isMale ? '#eff6ff' : '#fdf2f8');
</script>

<div class="person-card" style="--accent: {accentColor}; --bg-accent: {bgAccent}">
    <div class="avatar-container">
        {#if person.photo_path}
            <img src="/photos/{person.photo_path}" alt={person.name} class="avatar" />
        {:else}
            <div class="avatar-placeholder">
                <User size={32} />
            </div>
        {/if}
    </div>

    <div class="info">
        <h3 class="name">{person.name}</h3>
        {#if person.nickname}
            <p class="nickname">"{person.nickname}"</p>
        {/if}
    </div>

    {#if accessType !== 'view'}
        <div class="actions">
            <button class="action-btn edit" onclick={() => onEdit(person)} title="Edit">
                <Edit size={14} />
            </button>
            <button class="action-btn delete" onclick={() => onDelete(person)} title="Hapus">
                <Trash2 size={14} />
            </button>
            <button class="action-btn add" onclick={() => onAddFamily(person)} title="Tambah Keluarga">
                <Plus size={14} />
            </button>
        </div>
    {/if}
</div>

<style>
    .person-card {
        width: 180px;
        background: var(--bg-card);
        border-radius: var(--radius-lg);
        padding: var(--space-md);
        box-shadow: var(--shadow-md);
        border: 2px solid var(--border-color);
        display: flex;
        flex-direction: column;
        align-items: center;
        text-align: center;
        transition: transform var(--transition-normal), border-color var(--transition-normal);
        position: relative;
        overflow: visible;
        border-top: 4px solid var(--accent);
    }

    .person-card:hover {
        transform: translateY(-4px);
        border-color: var(--accent);
        box-shadow: var(--shadow-lg);
    }

    .avatar-container {
        width: 64px;
        height: 64px;
        border-radius: var(--radius-full);
        margin-bottom: var(--space-sm);
        overflow: hidden;
        border: 2px solid var(--accent);
        background: var(--bg-accent);
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .avatar {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .avatar-placeholder {
        color: var(--accent);
    }

    .info {
        width: 100%;
    }

    .name {
        font-size: 0.95rem;
        margin-bottom: 2px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .nickname {
        font-size: 0.75rem;
        color: var(--text-muted);
        font-style: italic;
    }

    .actions {
        display: flex;
        gap: var(--space-xs);
        margin-top: var(--space-md);
    }

    .action-btn {
        width: 28px;
        height: 28px;
        border-radius: var(--radius-sm);
        border: 1px solid var(--border-color);
        background: white;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        color: var(--text-muted);
        transition: all var(--transition-fast);
    }

    .action-btn:hover {
        border-color: var(--accent);
        color: var(--accent);
        background: var(--bg-accent);
    }

    .action-btn.delete:hover {
        border-color: var(--color-error);
        color: var(--color-error);
        background: #fef2f2;
    }
</style>
