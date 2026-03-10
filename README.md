# 🌳 Brayat — Family Tree Visualizer

A fun, colorful web app that lets families build and explore their family tree together — no tech skills required.

Targeted toward Indonesian families, the minimal app provides a privacy-first, simplified model focusing heavily on ease of use.

## Features

- **Instant Sessions**: Create isolated sessions requiring no user sign-ups or complex email flows.
- **Role-Based Links**: Generate `View-Only` or `Edit` share links for easy collaboration.
- **Smart Photo Compression**: High-quality photo management handled client-side using WebP converting, reducing server load.
- **Real-Time Visualization**: Auto-refreshing visual tree optimized for panning and zooming on mobile and desktop.

## Technology Stack

- **Backend & API**: Go 1.23+, `gin-gonic`, `gorm`, and `sqlite`
- **Frontend App**: SvelteKit (`adapter-static`), Svelte 5 Runes, and Vanilla CSS
- **Deployment**: Single process Go server acting as both the API and Static Asset host, deployed smoothly onto Fly.io

## Quickstart

```bash
# 1. Install Dependencies
go mod download

# 2. Run tests to ensure everything compiles nicely
make test

# 3. Start development server
make dev
```
