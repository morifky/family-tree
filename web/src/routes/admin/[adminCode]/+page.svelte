<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { api } from '$lib/api';
    import {
        Home, Share2, Eye, Edit3, Copy, Check, Plus,
        Trash2, RefreshCw, Lock, Unlock, TreePine, ExternalLink, Users
    } from 'lucide-svelte';

    const adminCode = $derived($page.params.adminCode);

    let session = $state<any>(null);
    let links = $state<any[]>([]);
    let loading = $state(true);
    let error = $state('');
    let saving = $state(false);
    let copiedId = $state<string | null>(null);

    // Tree summary counters (fetched separately)
    let personCount = $state(0);
    let relationshipCount = $state(0);

    async function loadDashboard() {
        loading = true;
        error = '';
        try {
            if (!adminCode) throw new Error('Kode admin tidak valid');
            // 1. Get session via admin code
            session = await api.sessions.getByAdminCode(adminCode);

            // 2. Load access links
            const linksData = await api.sessions.getLinks(session.id);
            links = Array.isArray(linksData) ? linksData : [];

            // 3. Load tree counts
            const treeRes = await api.sessions.getTree(session.id, '');
            if (treeRes.status === 200) {
                const data = await treeRes.json();
                const tree = data.data ?? data;
                personCount = (tree.people || []).length;
                relationshipCount = (tree.relationships || []).length;
            }
        } catch (err: any) {
            error = err.message || 'Gagal memuat data sesi.';
        } finally {
            loading = false;
        }
    }

    onMount(loadDashboard);

    async function createLink(type: 'view' | 'edit') {
        saving = true;
        try {
            const link = await api.sessions.createLink(session.id, type);
            links = [...links, link];
        } catch (err: any) {
            alert('Gagal membuat link: ' + err.message);
        } finally {
            saving = false;
        }
    }

    async function copyLink(link: any) {
        const url = `${window.location.origin}/tree/${link.code}`;
        try {
            await navigator.clipboard.writeText(url);
            copiedId = link.id;
            setTimeout(() => copiedId = null, 2000);
        } catch {
            prompt('Salin link ini:', url);
        }
    }

    function getLinkUrl(link: any) {
        return `${window.location.origin}/tree/${link.code}`;
    }

    function formatDate(iso: string) {
        return new Date(iso).toLocaleDateString('id-ID', {
            day: 'numeric', month: 'short', year: 'numeric'
        });
    }
</script>

<svelte:head>
    <title>{session ? `Admin – ${session.title}` : 'Admin'} - Brayat</title>
</svelte:head>

<div class="admin-page">
    <!-- Sidebar -->
    <aside class="sidebar">
        <div class="brand">
            <TreePine size={22} />
            <span>Brayat</span>
        </div>
        <nav class="side-nav">
            <a href="/" class="side-link">
                <Home size={16} /> Beranda
            </a>
            {#if session}
                <a href="/tree/{session.admin_code}" class="side-link">
                    <ExternalLink size={16} /> Lihat Pohon
                </a>
            {/if}
        </nav>
    </aside>

    <!-- Main content -->
    <main class="main">
        {#if loading}
            <div class="center-state">
                <RefreshCw size={40} class="spin" />
                <p>Memuat data admin...</p>
            </div>
        {:else if error}
            <div class="center-state error">
                <h2>Gagal Memuat</h2>
                <p>{error}</p>
                <button class="btn btn-primary" onclick={loadDashboard}>Coba Lagi</button>
            </div>
        {:else}
            <!-- Header -->
            <header class="page-header">
                <div>
                    <p class="breadcrumb">Panel Admin</p>
                    <h1 class="page-title">{session.title}</h1>
                </div>
                <a href="/tree/{session.admin_code}" class="btn btn-primary btn-sm">
                    <ExternalLink size={16} /> Buka Pohon
                </a>
            </header>

            <!-- Stats row -->
            <div class="stats-row">
                <div class="stat-card">
                    <div class="stat-icon people">
                        <Users size={20} />
                    </div>
                    <div>
                        <p class="stat-value">{personCount}</p>
                        <p class="stat-label">Anggota</p>
                    </div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon rels">
                        <Share2 size={20} />
                    </div>
                    <div>
                        <p class="stat-value">{relationshipCount}</p>
                        <p class="stat-label">Hubungan</p>
                    </div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon links">
                        <ExternalLink size={20} />
                    </div>
                    <div>
                        <p class="stat-value">{links.length}</p>
                        <p class="stat-label">Link Akses</p>
                    </div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon expiry">
                        <Lock size={20} />
                    </div>
                    <div>
                        <p class="stat-value">{session.status}</p>
                        <p class="stat-label">Kadaluarsa {formatDate(session.expires_at)}</p>
                    </div>
                </div>
            </div>

            <!-- Link Management -->
            <section class="section">
                <div class="section-header">
                    <div>
                        <h2 class="section-title">Link Berbagi</h2>
                        <p class="section-desc">Buat dan kelola link akses untuk dibagikan kepada anggota keluarga.</p>
                    </div>
                    <div class="flex gap-sm">
                        <button
                            class="btn btn-sm view-btn"
                            onclick={() => createLink('view')}
                            disabled={saving}
                        >
                            <Eye size={15} /> + Link Lihat
                        </button>
                        <button
                            class="btn btn-primary btn-sm"
                            onclick={() => createLink('edit')}
                            disabled={saving}
                        >
                            <Edit3 size={15} /> + Link Edit
                        </button>
                    </div>
                </div>

                {#if links.length === 0}
                    <div class="empty-links">
                        <Share2 size={36} />
                        <p>Belum ada link akses. Buat link di atas untuk mulai berbagi.</p>
                    </div>
                {:else}
                    <div class="links-list">
                        {#each links as link}
                            <div class="link-card">
                                <div class="link-left">
                                    <span class="access-badge {link.type}">
                                        {#if link.type === 'view'}
                                            <Eye size={11} /> Hanya Lihat
                                        {:else}
                                            <Edit3 size={11} /> Bisa Edit
                                        {/if}
                                    </span>
                                    <span class="link-url">{getLinkUrl(link)}</span>
                                    <span class="link-meta">Dibuat {formatDate(link.created_at)}</span>
                                </div>
                                <div class="link-actions">
                                    <a
                                        href={getLinkUrl(link)}
                                        target="_blank"
                                        class="icon-btn"
                                        title="Buka di tab baru"
                                    >
                                        <ExternalLink size={15} />
                                    </a>
                                    <button
                                        class="icon-btn copy-btn {copiedId === link.id ? 'copied' : ''}"
                                        onclick={() => copyLink(link)}
                                        title="Salin link"
                                    >
                                        {#if copiedId === link.id}
                                            <Check size={15} />
                                        {:else}
                                            <Copy size={15} />
                                        {/if}
                                    </button>
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
            </section>

            <!-- Admin code section -->
            <section class="section">
                <h2 class="section-title">Kode Admin</h2>
                <p class="section-desc">Simpan kode ini. Digunakan untuk membuka halaman admin ini kembali.</p>
                <div class="admin-code-box">
                    <code class="admin-code">{session.admin_code}</code>
                    <button class="icon-btn" onclick={() => { navigator.clipboard.writeText(session.admin_code); }} title="Salin kode admin">
                        <Copy size={15} />
                    </button>
                </div>
            </section>
        {/if}
    </main>
</div>

<style>
    .admin-page {
        display: flex;
        min-height: 100vh;
        background: var(--bg-main, #f8fafc);
        font-family: inherit;
    }

    /* ---- Sidebar (desktop) ---- */
    .sidebar {
        width: 220px;
        flex-shrink: 0;
        background: white;
        border-right: 1px solid var(--border-color, #e2e8f0);
        display: flex;
        flex-direction: column;
        padding: var(--space-lg, 1.5rem) var(--space-md, 1rem);
        gap: var(--space-lg, 1.5rem);
    }

    /* On mobile, sidebar becomes a compact top header */
    @media (max-width: 640px) {
        .admin-page   { flex-direction: column; }
        .sidebar {
            width: 100%;
            flex-direction: row;
            align-items: center;
            justify-content: space-between;
            border-right: none;
            border-bottom: 1px solid var(--border-color, #e2e8f0);
            padding: var(--space-sm) var(--space-md);
            gap: var(--space-md);
        }
        .side-nav { flex-direction: row; gap: var(--space-xs); }
    }

    .brand {
        display: flex;
        align-items: center;
        gap: var(--space-sm, 0.5rem);
        font-size: 1.1rem;
        font-weight: 800;
        color: var(--color-primary, #6366f1);
        white-space: nowrap;
    }

    .side-nav {
        display: flex;
        flex-direction: column;
        gap: var(--space-xs, 0.25rem);
    }

    .side-link {
        display: flex;
        align-items: center;
        gap: var(--space-sm, 0.5rem);
        padding: var(--space-xs, 0.25rem) var(--space-sm, 0.5rem);
        border-radius: var(--radius-md, 8px);
        color: var(--text-secondary, #64748b);
        font-size: 0.875rem;
        font-weight: 500;
        text-decoration: none;
        transition: all 0.15s;
        white-space: nowrap;
    }
    .side-link:hover {
        background: var(--bg-main, #f8fafc);
        color: var(--color-primary, #6366f1);
    }

    /* ---- Main content ---- */
    .main {
        flex: 1;
        padding: var(--space-xl, 2rem) var(--space-2xl, 3rem);
        overflow-y: auto;
        max-width: 900px;
    }

    @media (max-width: 768px) {
        .main { padding: var(--space-md); }
    }

    .center-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 60vh;
        color: var(--color-primary, #6366f1);
        gap: var(--space-md, 1rem);
        text-align: center;
        padding: var(--space-lg);
    }
    .center-state.error { color: var(--color-error, #ef4444); }

    /* ---- Page header ---- */
    .page-header {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        margin-bottom: var(--space-xl, 2rem);
        gap: var(--space-md);
        flex-wrap: wrap;
    }

    .breadcrumb {
        font-size: 0.75rem;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--text-muted, #94a3b8);
        margin-bottom: 4px;
    }

    .page-title {
        font-size: clamp(1.25rem, 4vw, 1.75rem);
        font-weight: 800;
        line-height: 1.2;
    }

    /* ---- Stats row ---- */
    .stats-row {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: var(--space-md, 1rem);
        margin-bottom: var(--space-xl, 2rem);
    }

    @media (max-width: 900px) {
        .stats-row { grid-template-columns: repeat(2, 1fr); }
    }

    @media (max-width: 480px) {
        .stats-row { grid-template-columns: 1fr 1fr; gap: var(--space-sm); }
    }

    .stat-card {
        background: white;
        border: 1px solid var(--border-color, #e2e8f0);
        border-radius: var(--radius-lg, 12px);
        padding: var(--space-md, 1rem);
        display: flex;
        align-items: center;
        gap: var(--space-sm, 0.5rem);
    }

    .stat-icon {
        width: 40px;
        height: 40px;
        border-radius: var(--radius-md, 8px);
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }
    .stat-icon.people { background: #eff6ff; color: #3b82f6; }
    .stat-icon.rels   { background: #fdf4ff; color: #a855f7; }
    .stat-icon.links  { background: #f0fdf4; color: #22c55e; }
    .stat-icon.expiry { background: #fff7ed; color: #f97316; }

    .stat-value {
        font-size: 1.25rem;
        font-weight: 800;
        line-height: 1;
        margin-bottom: 2px;
        text-transform: capitalize;
    }

    .stat-label {
        font-size: 0.7rem;
        color: var(--text-muted, #94a3b8);
        line-height: 1.3;
    }

    /* ---- Sections ---- */
    .section {
        background: white;
        border: 1px solid var(--border-color, #e2e8f0);
        border-radius: var(--radius-lg, 12px);
        padding: var(--space-lg, 1.5rem);
        margin-bottom: var(--space-lg, 1.5rem);
    }

    @media (max-width: 640px) {
        .section { padding: var(--space-md); }
    }

    .section-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: var(--space-lg, 1.5rem);
        gap: var(--space-md);
        flex-wrap: wrap;
    }

    .section-title {
        font-size: 1rem;
        font-weight: 700;
        margin-bottom: 4px;
    }

    .section-desc {
        font-size: 0.8rem;
        color: var(--text-muted, #94a3b8);
    }

    /* ---- Link cards ---- */
    .empty-links {
        text-align: center;
        padding: var(--space-xl, 2rem) var(--space-md);
        color: var(--text-muted, #94a3b8);
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--space-sm, 0.5rem);
        font-size: 0.875rem;
    }

    .links-list {
        display: flex;
        flex-direction: column;
        gap: var(--space-sm, 0.5rem);
    }

    .link-card {
        display: flex;
        align-items: center;
        justify-content: space-between;
        border: 1px solid var(--border-color, #e2e8f0);
        border-radius: var(--radius-md, 8px);
        padding: var(--space-sm, 0.5rem) var(--space-md, 1rem);
        gap: var(--space-sm, 0.5rem);
        background: var(--bg-main, #f8fafc);
    }

    .link-left {
        display: flex;
        align-items: center;
        gap: var(--space-sm, 0.5rem);
        min-width: 0;
        flex: 1;
        flex-wrap: wrap;
    }

    /* on very small screens, hide the date meta */
    @media (max-width: 480px) {
        .link-meta { display: none; }
    }

    .access-badge {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 0.7rem;
        font-weight: 700;
        text-transform: uppercase;
        padding: 3px 8px;
        border-radius: var(--radius-full, 999px);
        white-space: nowrap;
        flex-shrink: 0;
    }
    .access-badge.view { background: #dbeafe; color: #1d4ed8; }
    .access-badge.edit { background: #dcfce7; color: #15803d; }

    .link-url {
        font-size: 0.75rem;
        color: var(--text-secondary, #64748b);
        font-family: monospace;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        min-width: 0;
    }

    .link-meta {
        font-size: 0.7rem;
        color: var(--text-muted, #94a3b8);
        white-space: nowrap;
        flex-shrink: 0;
    }

    .link-actions {
        display: flex;
        gap: var(--space-xs, 0.25rem);
        flex-shrink: 0;
    }

    .icon-btn {
        width: 36px;
        height: 36px;
        border-radius: var(--radius-sm, 6px);
        border: 1px solid var(--border-color, #e2e8f0);
        background: white;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        color: var(--text-muted, #94a3b8);
        transition: all 0.15s;
        text-decoration: none;
    }
    .icon-btn:hover { border-color: var(--color-primary, #6366f1); color: var(--color-primary, #6366f1); }
    .icon-btn.copy-btn.copied { background: #f0fdf4; border-color: #22c55e; color: #22c55e; }

    /* ---- CTA Buttons ---- */
    .view-btn {
        background: #eff6ff;
        color: #1d4ed8;
        border: 1px solid #bfdbfe;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 4px;
        white-space: nowrap;
    }
    .view-btn:hover { background: #dbeafe; }

    .btn-sm {
        padding: 0.45rem 0.75rem;
        font-size: 0.8rem;
        display: flex;
        align-items: center;
        gap: 4px;
        white-space: nowrap;
    }

    /* ---- flex helper (not global, used internally) ---- */
    :global(.flex.gap-sm) { gap: var(--space-sm); }

    /* ---- Admin code box ---- */
    .admin-code-box {
        display: flex;
        align-items: center;
        gap: var(--space-md, 1rem);
        background: var(--bg-main, #f8fafc);
        border: 1px solid var(--border-color, #e2e8f0);
        border-radius: var(--radius-md, 8px);
        padding: var(--space-sm, 0.5rem) var(--space-md, 1rem);
        margin-top: var(--space-md, 1rem);
        width: fit-content;
        max-width: 100%;
    }

    .admin-code {
        font-family: monospace;
        font-size: 1.1rem;
        font-weight: 700;
        letter-spacing: 0.1em;
        color: var(--color-primary, #6366f1);
        word-break: break-all;
    }

    /* ---- Spin ---- */
    :global(.spin) { animation: spin 1s linear infinite; }
    @keyframes spin {
        from { transform: rotate(0deg); }
        to   { transform: rotate(360deg); }
    }
</style>
