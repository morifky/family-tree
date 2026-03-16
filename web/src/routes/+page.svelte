<script lang="ts">
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { LogIn, PlusCircle, TreePine } from 'lucide-svelte';

	let title = $state('');
	let accessCode = $state('');
	let loading = $state(false);
	let error = $state('');

	async function handleCreateSession(e: SubmitEvent) {
		e.preventDefault();
		if (!title) return;
		loading = true;
		error = '';
		try {
			const session = await api.sessions.create(title);
			goto(`/admin/${session.admin_code}`);
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	async function handleEnterCode(e: SubmitEvent) {
		e.preventDefault();
		if (!accessCode) return;
		loading = true;
		error = '';
		try {
			const res = await api.sessions.verifyCode(accessCode);
			if (res.access_type === 'admin') {
				goto(`/admin/${accessCode}`);
			} else {
				goto(`/tree/${accessCode}`);
			}
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Brayat - Silsilah Keluarga Modern</title>
	<meta name="description" content="Abadikan sejarah keluarga Anda dengan cara yang modern." />
</svelte:head>

<div class="home">
	<header class="hero">
		<div class="brand-icon">
			<TreePine size={36} />
		</div>
		<h1 class="brand-name">Brayat</h1>
		<p class="brand-tagline">Abadikan sejarah keluarga Anda dengan cara yang modern.</p>
	</header>

	{#if error}
		<div class="error-banner">{error}</div>
	{/if}

	<div class="cards">
		<!-- Create Session -->
		<section class="card">
			<div class="card-icon create-icon">
				<PlusCircle size={28} />
			</div>
			<h2 class="card-title">Buat Sesi Baru</h2>
			<p class="card-desc">
				Mulai pohon keluarga baru. Anda akan mendapatkan kode admin untuk mengelola sesi ini.
			</p>
			<form onsubmit={handleCreateSession} class="form">
				<div class="field">
					<label for="title">Nama Silsilah</label>
					<input
						type="text"
						id="title"
						bind:value={title}
						placeholder="Contoh: Keluarga Besar Trah..."
						required
					/>
				</div>
				<button type="submit" class="btn btn-primary" disabled={loading}>
					{loading ? 'Memproses...' : 'Buat Sekarang'}
				</button>
			</form>
		</section>

		<!-- Enter Code -->
		<section class="card">
			<div class="card-icon enter-icon">
				<LogIn size={28} />
			</div>
			<h2 class="card-title">Masukkan Kode</h2>
			<p class="card-desc">
				Punya kode akses? Masukkan di bawah ini untuk melihat atau mengedit pohon keluarga.
			</p>
			<form onsubmit={handleEnterCode} class="form">
				<div class="field">
					<label for="code">Kode Akses / Admin</label>
					<input
						type="text"
						id="code"
						bind:value={accessCode}
						placeholder="Masukkan 10 digit kode"
						autocomplete="off"
						autocapitalize="none"
						required
					/>
				</div>
				<button type="submit" class="btn btn-secondary" disabled={loading}>
					{loading ? 'Memproses...' : 'Masuk'}
				</button>
			</form>
		</section>
	</div>

	<footer class="footer">
		&copy; 2026 Brayat - Silsilah Keluarga Indonesia
	</footer>
</div>

<style>
	.home {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: var(--space-2xl) var(--space-md);
		gap: var(--space-xl);
	}

	/* ---- Hero ---- */
	.hero {
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-sm);
	}

	.brand-icon {
		width: 72px;
		height: 72px;
		background: linear-gradient(135deg, #eef2ff, #fdf2f8);
		border-radius: var(--radius-xl);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-primary);
		margin-bottom: var(--space-xs);
		box-shadow: var(--shadow-md);
	}

	.brand-name {
		font-size: clamp(2.5rem, 8vw, 4rem);
		font-weight: 800;
		background: linear-gradient(135deg, var(--color-primary), var(--color-secondary));
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
		line-height: 1;
	}

	.brand-tagline {
		font-size: clamp(0.95rem, 3vw, 1.2rem);
		color: var(--text-secondary);
		max-width: 400px;
	}

	/* ---- Error ---- */
	.error-banner {
		width: 100%;
		max-width: 680px;
		background: #fef2f2;
		border: 1px solid #fecaca;
		color: var(--color-error);
		border-radius: var(--radius-md);
		padding: var(--space-sm) var(--space-md);
		font-size: 0.875rem;
	}

	/* ---- Cards grid ---- */
	.cards {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--space-lg);
		width: 100%;
		max-width: 680px;
	}

	@media (max-width: 600px) {
		.cards {
			grid-template-columns: 1fr;
		}
	}

	/* ---- Card ---- */
	.card {
		background: white;
		border-radius: var(--radius-xl);
		padding: var(--space-lg);
		box-shadow: var(--shadow-md);
		border: 1px solid var(--border-color);
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}

	.card-icon {
		width: 52px;
		height: 52px;
		border-radius: var(--radius-md);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.create-icon { background: #eef2ff; color: var(--color-primary); }
	.enter-icon  { background: #fdf2f8; color: var(--color-secondary); }

	.card-title {
		font-size: 1.25rem;
		font-weight: 700;
	}

	.card-desc {
		font-size: 0.875rem;
		color: var(--text-secondary);
		line-height: 1.6;
		flex: 1;
	}

	/* ---- Form ---- */
	.form {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
	}

	.field label {
		font-size: 0.875rem;
		font-weight: 600;
	}

	.field input {
		padding: var(--space-sm) var(--space-md);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
		font-family: inherit;
		font-size: 1rem; /* ≥16px prevents iOS zoom */
		transition: border-color var(--transition-fast);
		width: 100%;
	}

	.field input:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
	}

	.btn {
		width: 100%;
		padding: var(--space-sm) var(--space-md);
		font-size: 1rem;
		font-weight: 700;
		border-radius: var(--radius-md);
		border: none;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn:disabled { opacity: 0.6; cursor: not-allowed; }

	.btn-primary {
		background: var(--color-primary);
		color: white;
	}
	.btn-primary:hover:not(:disabled) {
		background: var(--color-primary-hover);
		transform: translateY(-1px);
		box-shadow: var(--shadow-md);
	}

	.btn-secondary {
		background: var(--color-secondary);
		color: white;
	}
	.btn-secondary:hover:not(:disabled) {
		background: var(--color-secondary-hover);
		transform: translateY(-1px);
		box-shadow: var(--shadow-md);
	}

	/* ---- Footer ---- */
	.footer {
		color: var(--text-muted);
		font-size: 0.8rem;
		padding-bottom: var(--space-md);
	}
</style>
