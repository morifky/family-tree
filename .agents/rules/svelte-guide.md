---
description: SvelteKit Development Guidelines for Brayat
---

# Svelte & SvelteKit Guide for Brayat

These rules ensure the Svelte frontend remains fast, maintainable, and statically compiled.

## 1. Core Architecture (Static Site Generation)

- Use **SvelteKit** configured with `@sveltejs/adapter-static`.
- Fully compile to static HTML/JS/CSS assets at build time.
- All dynamic data fetching must happen client-side (`export const ssr = false;`).

## 2. Aesthetics & Styling

- Premium, fun, colorful design targeted at non-tech users in Indonesia.
- Use **Vanilla CSS** with variables (`var(--color-primary)`).
- Use smooth micro-animations.

## 3. State Management

- Use Svelte 5 runes (`$state`, `$derived`, `$effect`) for local component reactivity.
- For global polling state, use Svelte custom stores or runes to manage the loaded Family Tree graph.

## 4. Component Structure

Separate presentation from logic:

- `components/ui/` - Pure, dumb UI components (Buttons, Modals).
- `components/tree/` - Complex components handling the graph logic.

## 5. Client-Side Processing

- **Photo Compression**: Compress (WebP <100KB) client-side before upload.
- **Tree Rendering**: Backend sends flat lists. The client assembles the graph hierarchy.
