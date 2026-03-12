<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { page } from '$app/stores';
    import { api } from '$lib/api';
    import CanvasPanZoom from '$lib/components/tree/CanvasPanZoom.svelte';
    import TreeGraph from '$lib/components/tree/TreeGraph.svelte';
    import Modal from '$lib/components/ui/Modal.svelte';
    import PersonForm from '$lib/components/tree/PersonForm.svelte';
    import { Settings, Share2, Plus, Home, RefreshCw } from 'lucide-svelte';

    const accessCode = $derived($page.params.accessCode);
    
    let session = $state<any>(null);
    let people = $state<any[]>([]);
    let relationships = $state<any[]>([]);
    let accessType = $state('view');
    let loading = $state(true);
    let syncing = $state(false);
    let error = $state('');
    let lastEtag = $state('');

    // Modal state
    let showPersonForm = $state(false);
    let formMode = $state<'add' | 'edit'>('add');
    let selectedPerson = $state<any>(null);

    async function loadTree(isPolling = false) {
        if (!isPolling) loading = true;
        else syncing = true;

        try {
            if (!accessCode) throw new Error('Kode akses tidak valid');
            
            // 1. Verify access if not already known
            if (!session) {
                const auth = await api.sessions.verifyCode(accessCode);
                accessType = auth.access_type;
                session = await api.sessions.get(auth.session_id);
            }
            
            // 2. Fetch full tree data with ETag
            const res = await api.sessions.getTree(session.ID, lastEtag);

            if (res.status === 200) {
                const data = await res.json();
                people = data.people || [];
                relationships = data.relationships || [];
                lastEtag = res.headers.get('ETag') || '';
            } 
            // 304 Not Modified -> do nothing, data is up to date
        } catch (err: any) {
            if (!isPolling) error = err.message;
            console.error('Polling error:', err);
        } finally {
            loading = false;
            syncing = false;
        }
    }

    onMount(() => {
        loadTree();
        
        // Setup polling
        const interval = setInterval(() => {
            loadTree(true);
        }, 5000);

        return () => clearInterval(interval);
    });

    // Handlers
    function handleEdit(person: any) {
        selectedPerson = person;
        formMode = 'edit';
        showPersonForm = true;
    }

    async function handleDelete(person: any) {
        if (!confirm(`Hapus ${person.Name}? Seluruh data hubungan juga akan hilang.`)) return;
        try {
            await api.people.delete(person.ID);
            // Optimization: update local state immediately
            people = people.filter(p => p.ID !== person.ID);
        } catch (err: any) {
            alert(err.message);
        }
    }

    function handleAddFamily(person: any) {
        // Feature: auto-linking family will be in next task, 
        // for now just open the form
        selectedPerson = null;
        formMode = 'add';
        showPersonForm = true;
    }

    function handleAddRoot() {
        selectedPerson = null;
        formMode = 'add';
        showPersonForm = true;
    }

    function onFormSuccess(newOrUpdatedPerson: any) {
        showPersonForm = false;
        // The next poll will fetch the updated data, 
        // but let's trigger it immediately for better UX
        loadTree(true);
    }
</script>

<svelte:head>
    <title>{session ? session.Title : 'Memuat...'} - Brayat</title>
</svelte:head>

<div class="tree-page">
    {#if loading}
        <div class="overlay">
            <RefreshCw size={48} class="spin" />
            <p style="margin-top: var(--space-md);">Menyiapkan silsilah keluarga...</p>
        </div>
    {:else if error}
        <div class="overlay error">
            <h2>Gagal Memuat</h2>
            <p>{error}</p>
            <a href="/" class="btn btn-primary">Kembali ke Beranda</a>
        </div>
    {:else}
        <!-- Top Navigation -->
        <nav class="top-nav">
            <div class="flex items-center gap-md">
                <a href="/" class="logo-btn" title="Beranda">
                    <Home size={20} />
                </a>
                <div>
                    <h1 class="session-title">{session.Title}</h1>
                    <div class="flex items-center gap-sm">
                        <span class="badge {accessType}">
                            {accessType === 'admin' ? 'Admin' : accessType === 'edit' ? 'Editor' : 'Pembaca'}
                        </span>
                        {#if syncing}
                            <span class="sync-indicator">
                                <RefreshCw size={10} class="spin" /> sinkron...
                            </span>
                        {/if}
                    </div>
                </div>
            </div>

            <div class="flex gap-sm">
                {#if accessType !== 'view'}
                    <button class="btn btn-primary btn-sm" onclick={handleAddRoot}>
                        <Plus size={18} /> Tambah Orang
                    </button>
                {/if}
                <button class="btn btn-secondary btn-sm" onclick={() => alert('Link share: ' + window.location.href)}>
                    <Share2 size={18} /> Bagikan
                </button>
                {#if accessType === 'admin'}
                    <button class="btn btn-sm" style="background: #e2e8f0;" onclick={() => window.location.href = `/admin/${accessCode}`}>
                        <Settings size={18} /> Admin
                    </button>
                {/if}
            </div>
        </nav>

        <!-- Visualization Area -->
        <main class="canvas-area">
            <CanvasPanZoom>
                {#snippet children()}
                    <TreeGraph 
                        {people} 
                        {relationships} 
                        {accessType} 
                        onEdit={handleEdit}
                        onDelete={handleDelete}
                        onAddFamily={handleAddFamily}
                    />
                {/snippet}
            </CanvasPanZoom>
        </main>

        <!-- Empty State Hint -->
        {#if people.length === 0}
            <div class="empty-hint card">
                <h3>Mulai Silsilah Anda</h3>
                <p>Belum ada data orang. Tambahkan orang pertama untuk memulai.</p>
                <button class="btn btn-primary" onclick={handleAddRoot}>Tambah Orang Pertama</button>
            </div>
        {/if}

        <!-- Person Form Modal -->
        <Modal bind:show={showPersonForm} title={formMode === 'add' ? 'Tambah Anggota Keluarga' : 'Edit Data Anggota'}>
            <PersonForm
                sessionId={session.ID}
                person={selectedPerson}
                mode={formMode}
                onSuccess={onFormSuccess}
            />
        </Modal>
    {/if}
</div>

<style>
    .tree-page {
        height: 100vh;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    .top-nav {
        background: white;
        padding: var(--space-sm) var(--space-lg);
        display: flex;
        align-items: center;
        justify-content: space-between;
        box-shadow: var(--shadow-sm);
        z-index: 10;
        border-bottom: 1px solid var(--border-color);
    }

    .logo-btn {
        width: 40px;
        height: 40px;
        background: var(--bg-main);
        border-radius: var(--radius-md);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-primary);
        transition: all var(--transition-fast);
    }

    .logo-btn:hover {
        background: #e0e7ff;
    }

    .session-title {
        font-size: 1.1rem;
        font-weight: 700;
        line-height: 1.1;
        margin-bottom: 2px;
    }

    .badge {
        font-size: 0.6rem;
        font-weight: 800;
        text-transform: uppercase;
        padding: 2px 6px;
        border-radius: var(--radius-full);
    }

    .badge.admin { background: #d1fae5; color: #065f46; }
    .badge.edit { background: #fee2e2; color: #991b1b; }
    .badge.view { background: #e0e7ff; color: #3730a3; }

    .sync-indicator {
        font-size: 0.65rem;
        color: var(--text-muted);
        display: flex;
        align-items: center;
        gap: 4px;
    }

    .canvas-area {
        flex: 1;
        position: relative;
    }

    .btn-sm {
        padding: var(--space-xs) var(--space-md);
        font-size: 0.875rem;
    }

    .overlay {
        position: fixed;
        inset: 0;
        background: white;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        z-index: 100;
        color: var(--color-primary);
    }

    .overlay.error {
        color: var(--color-error);
    }

    .spin {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }

    .empty-hint {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        text-align: center;
        max-width: 400px;
        border-top: 5px solid var(--color-primary);
        z-index: 5;
    }
</style>
