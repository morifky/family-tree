import { browser } from '$app/environment';

const API_BASE = '/api/v1';

async function request(path: string, options: RequestInit = {}) {
    const res = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers || {})
        }
    });

    if (!res.ok) {
        const error = await res.json().catch(() => ({ message: 'Terjadi kesalahan sistem' }));
        throw new Error(error.message || 'Gagal terhubung ke server');
    }

    if (res.status === 204) return null;
    return res.json();
}

export const api = {
    sessions: {
        create: (title: string) => request('/sessions', { 
            method: 'POST', 
            body: JSON.stringify({ title }) 
        }),
        get: (id: string) => request(`/sessions/${id}`),
        updateStatus: (id: string, status: string) => request(`/sessions/${id}/status`, {
            method: 'PUT',
            body: JSON.stringify({ status })
        }),
        extend: (id: string) => request(`/sessions/${id}/extend`, {
            method: 'POST'
        }),
        getLinks: (id: string) => request(`/sessions/${id}/links`, {
            method: 'GET'
        }),
        createLink: (id: string, type: string) => request(`/sessions/${id}/links`, {
            method: 'POST',
            body: JSON.stringify({ type })
        }),
        // Verify code (admin, edit, or view)
        verifyCode: (code: string) => request(`/sessions/verify/${code}`),
        getTree: (id: string) => request(`/sessions/${id}/tree`)
    },
    people: {
        create: (sessionId: string, data: any) => request(`/sessions/${sessionId}/people`, {
            method: 'POST',
            body: JSON.stringify(data)
        }),
        update: (id: string, data: any) => request(`/people/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        }),
        delete: (id: string) => request(`/people/${id}`, {
            method: 'DELETE'
        })
    },
    relationships: {
        create: (sessionId: string, data: any) => request(`/sessions/${sessionId}/relationships`, {
            method: 'POST',
            body: JSON.stringify(data)
        }),
        delete: (id: string) => request(`/relationships/${id}`, {
            method: 'DELETE'
        })
    }
};
