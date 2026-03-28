<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { api } from '$lib/api';
    import CanvasPanZoom from '$lib/components/tree/CanvasPanZoom.svelte';
    import TreeGraph from '$lib/components/tree/TreeGraph.svelte';
    import Modal from '$lib/components/ui/Modal.svelte';
    import ConfirmDialog from '$lib/components/ui/ConfirmDialog.svelte';
    import PersonForm from '$lib/components/tree/PersonForm.svelte';
    import RelationshipForm from '$lib/components/tree/RelationshipForm.svelte';
    import { Settings, Share2, Plus, Home, RefreshCw, Link2, Trash2, Check } from 'lucide-svelte';

    const accessCode = $derived($page.params.accessCode);
    
    let session = $state<any>(null);
    let people = $state<any[]>([]);
    let relationships = $state<any[]>([]);
    let accessType = $state('view');
    let loading = $state(true);
    let syncing = $state(false);
    let error = $state('');
    let lastEtag = $state('');
    let copied = $state(false);

    // Person modal
    let showPersonForm = $state(false);
    let formMode = $state<'add' | 'edit'>('add');
    let selectedPerson = $state<any>(null);

    // Relationship modal
    let showRelForm = $state(false);
    let relPreselectedPerson = $state<any>(null);

    // Relationships panel
    let showRelPanel = $state(false);

    // Confirm dialog state
    let confirmShow   = $state(false);
    let confirmTitle  = $state('');
    let confirmMsg    = $state('');
    let confirmAction = $state<() => void>(() => {});

    function openConfirm(title: string, msg: string, action: () => void) {
        confirmTitle  = title;
        confirmMsg    = msg;
        confirmAction = action;
        confirmShow   = true;
    }

    async function loadTree(isPolling = false, force = false) {
        if (!isPolling) loading = true;
        else syncing = true;

        try {
            if (!accessCode) throw new Error('Kode akses tidak valid');
            
            if (!session) {
                const auth = await api.sessions.verifyCode(accessCode);
                accessType = auth.access_type;
                session = await api.sessions.get(auth.session_id);
            }
            
            const res = await api.sessions.getTree(session.id, force ? '' : lastEtag);

            if (res.status === 200) {
                const data = await res.json();
                people = (data.data ?? data).people || [];
                relationships = (data.data ?? data).relationships || [];
                lastEtag = res.headers.get('ETag') || '';
            } 
        } catch (err: any) {
            if (!isPolling) error = err.message;
            console.error('Load error:', err);
        } finally {
            loading = false;
            syncing = false;
        }
    }

    onMount(() => {
        loadTree();
        const interval = setInterval(() => loadTree(true), 5000);
        return () => clearInterval(interval);
    });

    // --- Share ---
    async function handleShare() {
        try {
            await navigator.clipboard.writeText(window.location.href);
            copied = true;
            setTimeout(() => copied = false, 2000);
        } catch {
            prompt('Salin link ini:', window.location.href);
        }
    }

    // --- Person handlers ---
    function handleEdit(person: any) {
        selectedPerson = person;
        formMode = 'edit';
        showPersonForm = true;
    }

    function handleDelete(person: any) {
        openConfirm(
            'Hapus Anggota',
            `Hapus "${person.name}"?\nSeluruh data hubungan orang ini juga akan ikut terhapus.`,
            async () => {
                try {
                    await api.people.delete(person.id);
                    people = people.filter(p => p.id !== person.id);
                    relationships = relationships.filter(r => r.person_a_id !== person.id && r.person_b_id !== person.id);
                    loadTree(true, true);
                } catch (err: any) {
                    openConfirm('Gagal', 'Gagal menghapus: ' + err.message, () => {});
                }
            }
        );
    }

    function handleAddFamily(person: any) {
        relPreselectedPerson = person;
        showRelForm = true;
    }

    function handleAddRoot() {
        selectedPerson = null;
        formMode = 'add';
        showPersonForm = true;
    }

    function onPersonFormSuccess() {
        showPersonForm = false;
        loadTree(true, true);
    }

    // --- Relationship handlers ---
    function onRelFormSuccess() {
        showRelForm = false;
        loadTree(true, true);
    }

    function handleDeleteRelationship(rel: any) {
        const personA = people.find(p => p.id === rel.person_a_id);
        const personB = people.find(p => p.id === rel.person_b_id);
        const typeLabel = rel.type === 'parent_child' ? 'Orang Tua-Anak' : 'Pasangan';
        const arrow    = rel.type === 'parent_child' ? '→' : '↔';
        const label    = `${personA?.name || '?'} ${arrow} ${personB?.name || '?'}`;

        openConfirm(
            'Hapus Hubungan',
            `Hapus hubungan ${typeLabel}:\n${label}?`,
            async () => {
                try {
                    await api.relationships.delete(rel.id);
                    relationships = relationships.filter(r => r.id !== rel.id);
                    loadTree(true, true);
                } catch (err: any) {
                    openConfirm('Gagal', 'Gagal menghapus: ' + err.message, () => {});
                }
            }
        );
    }

    function getPersonName(id: string) {
        return people.find(p => p.id === id)?.name || '?';
    }
</script>

<svelte:head>
    <title>{session ? session.title : 'Memuat...'} - Brayat</title>
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
                    <h1 class="session-title">{session.title}</h1>
                    <div class="flex items-center gap-sm">
                        <span class="badge {accessType}">
                            {accessType === 'admin' ? 'Admin' : accessType === 'edit' ? 'Editor' : 'Pembaca'}
                        </span>
                        {#if syncing}
                            <span class="sync-indicator">
                                <RefreshCw size={10} class="spin" /> <span class="nav-label">sinkron...</span>
                            </span>
                        {/if}
                    </div>
                </div>
            </div>

            <div class="nav-actions">
                {#if accessType !== 'view'}
                    <button class="btn btn-primary btn-sm" onclick={handleAddRoot} title="Tambah Orang">
                        <Plus size={18} /> <span class="nav-label">Tambah Orang</span>
                    </button>
                    <button class="btn btn-secondary btn-sm" onclick={() => { relPreselectedPerson = null; showRelForm = true; }} title="Tambah Hubungan">
                        <Link2 size={18} /> <span class="nav-label">Hubungan</span>
                    </button>
                    <button class="btn btn-sm rel-panel-btn" onclick={() => showRelPanel = !showRelPanel} title="Daftar Hubungan">
                        <span class="nav-label">Daftar</span>
                        {#if relationships.length > 0}<span class="rel-count">{relationships.length}</span>{/if}
                    </button>
                {/if}

                <button class="btn btn-secondary btn-sm" onclick={handleShare} title="Bagikan">
                    {#if copied}
                        <Check size={18} /> <span class="nav-label">Tersalin!</span>
                    {:else}
                        <Share2 size={18} /> <span class="nav-label">Bagikan</span>
                    {/if}
                </button>

                {#if accessType === 'admin' && session?.admin_code}
                    <a href="/admin/{session.admin_code}" class="btn btn-sm admin-btn" title="Admin">
                        <Settings size={18} /> <span class="nav-label">Admin</span>
                    </a>
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

            <!-- Relationship Side Panel -->
            {#if showRelPanel}
                <div class="rel-panel">
                    <div class="rel-panel-header">
                        <h3>Daftar Hubungan</h3>
                        <button onclick={() => showRelPanel = false} class="close-btn">✕</button>
                    </div>
                    {#if relationships.length === 0}
                        <p class="empty-msg">Belum ada hubungan.</p>
                    {:else}
                        <ul class="rel-list">
                            {#each relationships as rel}
                                <li class="rel-item">
                                    <div class="rel-info">
                                        <span class="rel-badge {rel.type}">
                                            {rel.type === 'parent_child' ? 'Orang Tua-Anak' : 'Pasangan'}
                                        </span>
                                        <span class="rel-names">
                                            {getPersonName(rel.person_a_id)}
                                            {rel.type === 'parent_child' ? '→' : '↔'}
                                            {getPersonName(rel.person_b_id)}
                                        </span>
                                    </div>
                                    {#if accessType !== 'view'}
                                        <button class="del-btn" onclick={() => handleDeleteRelationship(rel)} title="Hapus hubungan">
                                            <Trash2 size={14} />
                                        </button>
                                    {/if}
                                </li>
                            {/each}
                        </ul>
                    {/if}
                </div>
            {/if}
        </main>

        <!-- Empty State -->
        {#if people.length === 0}
            <div class="empty-hint card">
                <h3>Mulai Silsilah Anda</h3>
                <p>Belum ada data orang. Tambahkan orang pertama untuk memulai.</p>
                {#if accessType !== 'view'}
                    <button class="btn btn-primary" onclick={handleAddRoot}>Tambah Orang Pertama</button>
                {/if}
            </div>
        {/if}

        {#if people.length > 1 && relationships.length === 0 && accessType !== 'view'}
            <div class="rel-hint">
                <Link2 size={14} />
                Klik "Tambah Hubungan" atau tombol <strong>+</strong> di kartu untuk menghubungkan orang.
            </div>
        {/if}

        <!-- Person Form Modal -->
        <Modal bind:show={showPersonForm} title={formMode === 'add' ? 'Tambah Anggota Keluarga' : 'Edit Data Anggota'}>
            <PersonForm
                sessionId={session.id}
                person={selectedPerson}
                mode={formMode}
                onSuccess={onPersonFormSuccess}
            />
        </Modal>

        <!-- Relationship Form Modal -->
        <Modal bind:show={showRelForm} title="Tambah Hubungan">
            <RelationshipForm
                sessionId={session.id}
                {people}
                {relationships}
                preselectedPerson={relPreselectedPerson}
                onSuccess={onRelFormSuccess}
            />
        </Modal>
    {/if}
</div>

<!-- Confirm Dialog (outside tree-page so it always renders on top) -->
<ConfirmDialog
    bind:show={confirmShow}
    title={confirmTitle}
    message={confirmMsg}
    confirmLabel="Hapus"
    cancelLabel="Batal"
    onConfirm={confirmAction}
/>

<style>
    .tree-page {
        height: 100vh;
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    /* ---- Top Nav ---- */
    .top-nav {
        background: white;
        padding: var(--space-sm) var(--space-md);
        display: flex;
        align-items: center;
        justify-content: space-between;
        box-shadow: var(--shadow-sm);
        z-index: 10;
        border-bottom: 1px solid var(--border-color);
        gap: var(--space-sm);
        flex-shrink: 0;
    }

    .nav-actions {
        display: flex;
        align-items: center;
        gap: var(--space-xs);
        flex-wrap: nowrap;
    }

    /* hide text labels on small screens, keep icons */
    @media (max-width: 640px) {
        .nav-label { display: none; }
        .top-nav { padding: var(--space-xs) var(--space-sm); }
        .btn-sm { padding: 6px 8px; min-width: 36px; justify-content: center; }
        .session-title { font-size: 0.9rem; }
        .rel-panel { width: calc(100vw - 2 * var(--space-sm)); right: var(--space-sm); left: var(--space-sm); }
        .rel-hint { display: none; }
        .empty-hint { width: calc(100vw - 2rem); }
    }

    .logo-btn {
        width: 36px;
        height: 36px;
        background: var(--bg-main);
        border-radius: var(--radius-md);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-primary);
        transition: all var(--transition-fast);
        flex-shrink: 0;
    }
    .logo-btn:hover { background: #e0e7ff; }

    .session-title {
        font-size: 1rem;
        font-weight: 700;
        line-height: 1.1;
        margin-bottom: 2px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 140px;
    }

    @media (min-width: 768px) {
        .session-title { font-size: 1.1rem; max-width: none; }
    }

    .badge {
        font-size: 0.6rem;
        font-weight: 800;
        text-transform: uppercase;
        padding: 2px 6px;
        border-radius: var(--radius-full);
    }
    .badge.admin { background: #d1fae5; color: #065f46; }
    .badge.edit  { background: #fee2e2; color: #991b1b; }
    .badge.view  { background: #e0e7ff; color: #3730a3; }

    .sync-indicator {
        font-size: 0.65rem;
        color: var(--text-muted);
        display: flex;
        align-items: center;
        gap: 4px;
    }

    .rel-count {
        background: var(--color-primary);
        color: white;
        font-size: 0.6rem;
        font-weight: 800;
        border-radius: var(--radius-full);
        padding: 1px 5px;
        margin-left: 2px;
    }

    /* ---- Canvas ---- */
    .canvas-area {
        flex: 1;
        position: relative;
        min-height: 0;
    }

    /* ---- Button helpers ---- */
    .btn-sm {
        padding: var(--space-xs) var(--space-sm);
        font-size: 0.8rem;
        gap: 4px;
    }

    .rel-panel-btn {
        background: #f1f5f9;
        border: 1px solid var(--border-color);
    }

    .admin-btn {
        background: #e2e8f0;
        color: var(--text-primary);
        text-decoration: none;
        display: flex;
        align-items: center;
        gap: 4px;
        border-radius: var(--radius-md);
        font-weight: 600;
    }
    .admin-btn:hover { background: #cbd5e1; }

    /* ---- Relationship Side Panel ---- */
    .rel-panel {
        position: absolute;
        top: var(--space-md);
        right: var(--space-md);
        width: 300px;
        background: white;
        border-radius: var(--radius-lg);
        box-shadow: var(--shadow-lg);
        border: 1px solid var(--border-color);
        z-index: 20;
        max-height: 60vh;
        display: flex;
        flex-direction: column;
    }

    .rel-panel-header {
        padding: var(--space-sm) var(--space-md);
        border-bottom: 1px solid var(--border-color);
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    .rel-panel-header h3 { font-size: 0.95rem; font-weight: 700; }

    .close-btn {
        background: none;
        border: none;
        cursor: pointer;
        font-size: 1rem;
        color: var(--text-muted);
        padding: 2px 6px;
        border-radius: var(--radius-sm);
        min-width: 28px;
        min-height: 28px;
    }
    .close-btn:hover { background: #f1f5f9; }

    .rel-list {
        list-style: none;
        overflow-y: auto;
        flex: 1;
        padding: var(--space-sm);
        display: flex;
        flex-direction: column;
        gap: var(--space-xs);
        -webkit-overflow-scrolling: touch;
    }

    .rel-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: var(--space-sm);
        background: var(--bg-main);
        border: 1px solid var(--border-color);
        border-radius: var(--radius-md);
        gap: var(--space-sm);
        transition: border-color var(--transition-fast);
    }
    
    .rel-item:hover {
        border-color: #cbd5e1;
    }

    .rel-info {
        display: flex;
        flex-direction: column;
        gap: 4px;
        min-width: 0;
    }

    .rel-badge {
        font-size: 0.6rem;
        font-weight: 700;
        text-transform: uppercase;
        padding: 2px 6px;
        border-radius: var(--radius-full);
        width: fit-content;
    }
    .rel-badge.parent_child { background: #dbeafe; color: #1e40af; }
    .rel-badge.spouse       { background: #fce7f3; color: #9d174d; }

    .rel-names {
        font-size: 0.85rem;
        font-weight: 500;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .del-btn {
        background: white;
        border: 1px solid var(--border-color);
        cursor: pointer;
        color: var(--color-error);
        padding: 6px;
        border-radius: var(--radius-sm);
        display: flex;
        flex-shrink: 0;
        transition: all var(--transition-fast);
        min-width: 32px;
        min-height: 32px;
        align-items: center;
        justify-content: center;
    }
    .del-btn:hover { background: #fee2e2; border-color: #fca5a5; }

    .empty-msg {
        padding: var(--space-md);
        text-align: center;
        color: var(--text-muted);
        font-size: 0.875rem;
    }

    /* ---- Overlays ---- */
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
        padding: var(--space-lg);
        text-align: center;
        gap: var(--space-md);
    }
    .overlay.error { color: var(--color-error); }

    :global(.spin) { animation: spin 1s linear infinite; }
    @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

    /* ---- Empty state ---- */
    .empty-hint {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        text-align: center;
        max-width: 360px;
        width: calc(100% - 2rem);
        border-top: 5px solid var(--color-primary);
        z-index: 5;
    }

    .rel-hint {
        position: absolute;
        bottom: var(--space-lg);
        left: 50%;
        transform: translateX(-50%);
        background: white;
        border: 1px solid var(--border-color);
        border-radius: var(--radius-full);
        padding: var(--space-xs) var(--space-md);
        font-size: 0.75rem;
        color: var(--text-secondary);
        box-shadow: var(--shadow-sm);
        display: flex;
        align-items: center;
        gap: var(--space-xs);
        z-index: 5;
        white-space: nowrap;
        max-width: calc(100vw - 2rem);
        overflow: hidden;
        text-overflow: ellipsis;
    }
</style>

