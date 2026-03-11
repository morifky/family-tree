<script lang="ts">
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { LogIn, PlusCircle, LayoutDashboard } from 'lucide-svelte';

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
			goto(`/admin/${session.AdminCode}`);
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
				// For edit/view links, we might have a different route later
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
</svelte:head>

<div class="container" style="padding: var(--space-2xl) var(--space-md); max-width: 800px;">
	<header style="text-align: center; margin-bottom: var(--space-2xl);">
		<h1 style="font-size: 3.5rem; color: var(--color-primary); margin-bottom: var(--space-sm);">Brayat</h1>
		<p style="font-size: 1.25rem; color: var(--text-muted);">Abadikan sejarah keluarga Anda dengan cara yang modern.</p>
	</header>

	{#if error}
		<div class="card" style="border-color: var(--color-error); background: #fef2f2; color: var(--color-error); margin-bottom: var(--space-lg); padding: var(--space-md);">
			{error}
		</div>
	{/if}

	<div style="display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-xl);">
		<!-- Create Session -->
		<section class="card flex flex-col gap-lg">
			<div class="flex items-center gap-md">
				<div style="background: #eef2ff; padding: var(--space-sm); border-radius: var(--radius-md); color: var(--color-primary);">
					<PlusCircle size={32} />
				</div>
				<h2 style="font-size: 1.5rem;">Buat Sesi Baru</h2>
			</div>
			
			<p style="color: var(--text-muted); font-size: 0.95rem;">
				Mulai pohon keluarga baru. Anda akan mendapatkan kode admin untuk mengelola sesi ini.
			</p>

			<form onsubmit={handleCreateSession} class="flex flex-col gap-md">
				<div class="flex flex-col gap-sm">
					<label for="title" style="font-weight: 600; font-size: 0.875rem;">Nama Silsilah</label>
					<input 
						type="text" 
						id="title" 
						bind:value={title} 
						placeholder="Contoh: Keluarga Besar Trah..."
						style="padding: var(--space-sm); border: 1px solid var(--border-color); border-radius: var(--radius-md); font-family: inherit;"
						required
					/>
				</div>
				<button type="submit" class="btn btn-primary" disabled={loading}>
					{loading ? 'Memproses...' : 'Buat Sekarang'}
				</button>
			</form>
		</section>

		<!-- Enter Code -->
		<section class="card flex flex-col gap-lg">
			<div class="flex items-center gap-md">
				<div style="background: #fdf2f8; padding: var(--space-sm); border-radius: var(--radius-md); color: var(--color-secondary);">
					<LogIn size={32} />
				</div>
				<h2 style="font-size: 1.5rem;">Masukkan Kode</h2>
			</div>

			<p style="color: var(--text-muted); font-size: 0.95rem;">
				Punya kode akses? Masukkan di bawah ini untuk melihat atau mengedit pohon keluarga.
			</p>

			<form onsubmit={handleEnterCode} class="flex flex-col gap-md">
				<div class="flex flex-col gap-sm">
					<label for="code" style="font-weight: 600; font-size: 0.875rem;">Kode Akses / Admin</label>
					<input 
						type="text" 
						id="code" 
						bind:value={accessCode} 
						placeholder="Masukkan 10 digit kode"
						style="padding: var(--space-sm); border: 1px solid var(--border-color); border-radius: var(--radius-md); font-family: inherit;"
						required
					/>
				</div>
				<button type="submit" class="btn btn-secondary" disabled={loading}>
					{loading ? 'Memproses...' : 'Masuk'}
				</button>
			</form>
		</section>
	</div>

	<footer style="margin-top: var(--space-2xl); text-align: center; color: var(--text-muted); font-size: 0.875rem;">
		&copy; 2026 Brayat - Silsilah Keluarga Indonesia
	</footer>
</div>
