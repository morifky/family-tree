<script lang="ts">
    import { api } from '$lib/api';
    import imageCompression from 'browser-image-compression';
    import { Camera, Loader2 } from 'lucide-svelte';

    /**
     * @typedef {Object} Props
     * @property {string} sessionId
     * @property {any} person
     * @property {string} mode - 'add' or 'edit'
     * @property {Function} onSuccess
     */
    let { sessionId, person = null, mode = 'add', onSuccess } = $props();

    let name = $state('');
    let nickname = $state('');
    let gender = $state('male');
    let photo: File | null = $state(null);
    let photoPreview = $state('');
    
    let loading = $state(false);
    let error = $state('');

    // Sync state when person changes (e.g. modal re-opened with different person)
    $effect(() => {
        name = person?.Name || '';
        nickname = person?.Nickname || '';
        gender = person?.Gender || 'male';
        photoPreview = person?.PhotoPath || '';
        photo = null;
    });

    async function handlePhotoChange(e: Event) {
        const input = e.target as HTMLInputElement;
        if (!input.files?.length) return;

        const originalFile = input.files[0];
        
        // Setup compression options
        const options = {
            maxSizeMB: 0.1, // < 100KB
            maxWidthOrHeight: 400,
            useWebWorker: true,
            fileType: 'image/webp'
        };

        try {
            loading = true;
            const compressedFile = await imageCompression(originalFile, options);
            photo = compressedFile;
            photoPreview = URL.createObjectURL(compressedFile);
        } catch (err: any) {
            error = 'Gagal mengompres gambar: ' + err.message;
        } finally {
            loading = false;
        }
    }

    async function handleSubmit(e: SubmitEvent) {
        e.preventDefault();
        if (!name) return;

        loading = true;
        error = '';
        
        try {
            const formData = new FormData();
            formData.append('name', name);
            formData.append('nickname', nickname);
            formData.append('gender', gender);
            if (photo) {
                formData.append('photo', photo, 'avatar.webp');
            }

            if (mode === 'add') {
                const newPerson = await api.people.create(sessionId, formData);
                onSuccess(newPerson);
            } else {
                await api.people.update(person.ID, formData);
                onSuccess({ ...person, Name: name, Nickname: nickname, Gender: gender, PhotoPath: photoPreview });
            }
        } catch (err: any) {
            error = err.message;
        } finally {
            loading = false;
        }
    }
</script>

<form onsubmit={handleSubmit} class="person-form">
    {#if error}
        <div class="alert error">{error}</div>
    {/if}

    <div class="photo-upload">
        <div class="preview-container">
            {#if photoPreview}
                <img src={photoPreview} alt="Preview" class="avatar-preview" />
            {:else}
                <div class="avatar-placeholder">
                    <Camera size={40} />
                </div>
            {/if}
            {#if loading}
                <div class="compression-overlay">
                    <Loader2 size={32} class="spin" />
                </div>
            {/if}
        </div>
        <label class="upload-btn" for="photo-input">
            <span>{photoPreview ? 'Ubah Foto' : 'Unggah Foto'}</span>
        </label>
        <input type="file" id="photo-input" accept="image/*" onchange={handlePhotoChange} class="sr-only" />
        <p class="hint">Maks. 100KB (WebP)</p>
    </div>

    <div class="form-group">
        <label for="name">Nama Lengkap *</label>
        <input type="text" id="name" bind:value={name} placeholder="Nama lengkap" required />
    </div>

    <div class="form-group">
        <label for="nickname">Nama Panggilan</label>
        <input type="text" id="nickname" bind:value={nickname} placeholder="Panggilan akrab" />
    </div>

    <div class="form-group">
        <span class="label">Jenis Kelamin</span>
        <div class="radio-group">
            <label class="radio-label">
                <input type="radio" value="male" bind:group={gender} />
                Laki-laki
            </label>
            <label class="radio-label">
                <input type="radio" value="female" bind:group={gender} />
                Perempuan
            </label>
        </div>
    </div>

    <button type="submit" class="btn btn-primary btn-block" disabled={loading}>
        {#if loading}
            <Loader2 size={18} class="spin" /> Menyimpan...
        {:else}
            {mode === 'add' ? 'Tambah Anggota' : 'Simpan Perubahan'}
        {/if}
    </button>
</form>

<style>
    .person-form {
        display: flex;
        flex-direction: column;
        gap: var(--space-lg);
    }

    .photo-upload {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--space-sm);
        margin-bottom: var(--space-md);
    }

    .preview-container {
        width: 120px;
        height: 120px;
        border-radius: var(--radius-full);
        border: 3px solid var(--color-primary);
        overflow: hidden;
        position: relative;
        background: var(--bg-main);
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .avatar-preview {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .avatar-placeholder {
        color: var(--text-muted);
    }

    .compression-overlay {
        position: absolute;
        inset: 0;
        background: rgba(255, 255, 255, 0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-primary);
    }

    :global(.spin) {
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }

    .upload-btn {
        background: var(--color-primary);
        color: white;
        padding: var(--space-xs) var(--space-md);
        border-radius: var(--radius-full);
        font-size: 0.875rem;
        font-weight: 600;
        cursor: pointer;
        transition: all var(--transition-fast);
    }

    .upload-btn:hover {
        background: var(--color-primary-hover);
        transform: scale(1.05);
    }

    .hint {
        font-size: 0.75rem;
        color: var(--text-muted);
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: var(--space-xs);
    }

    .form-group label {
        font-weight: 600;
        font-size: 0.875rem;
    }

    .form-group input {
        padding: var(--space-sm);
        border: 1px solid var(--border-color);
        border-radius: var(--radius-md);
        font-family: inherit;
    }

    .radio-group {
        display: flex;
        gap: var(--space-lg);
        padding-top: var(--space-xs);
    }

    .radio-label {
        display: flex;
        align-items: center;
        gap: var(--space-xs);
        font-size: 0.95rem;
        cursor: pointer;
    }

    .btn-block {
        width: 100%;
        margin-top: var(--space-md);
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
</style>
