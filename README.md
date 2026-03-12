# 🌳 Brayat — Family Tree Visualizer

A fun, colorful web app that lets families build and explore their family tree together — no tech skills required.

Targeted toward Indonesian families, the minimal app provides a privacy-first, simplified model focusing heavily on ease of use.

## Features

- **Instant Sessions**: Create isolated sessions requiring no user sign-ups or complex email flows.
- **Role-Based Links**: Generate `View-Only` or `Edit` share links for easy collaboration.
- **Smart Photo Compression**: High-quality photo management handled client-side using WebP converting, reducing server load.
- **Real-Time Visualization**: Auto-refreshing visual tree optimized for panning and zooming on mobile and desktop.

## Technology Stack

- **Backend & API**: Go 1.26+, `gin-gonic`, `gorm`, and `sqlite`
- **Frontend App**: SvelteKit (`adapter-static`), Svelte 5 Runes, and Vanilla CSS
- **Deployment**: Single process Go server acting as both the API and Static Asset host, deployed smoothly onto Fly.io

## Quickstart

```bash
# 1. Install Dependencies
go mod download

# 2. Run tests to ensure everything compiles nicely
make test

# 3. Start development server
```

## Deployment

The application is designed to be deployed as a single Docker container.

### Deploy to Fly.io

1. **Install Fly CTL**: Follow the instructions at [fly.io/docs/hands-on/install-flyctl/](https://fly.io/docs/hands-on/install-flyctl/).
2. **Login**: `fly auth login`.
3. **Create Volume**: The app requires a volume for persistent storage (SQLite & Photos):
   ```bash
   fly volumes create brayat_data --size 1
   ```
4. **Deploy**:
   ```bash
   fly deploy
   ```

The `Dockerfile` handles the multi-stage build (SvelteKit → Go), and the Go server serves both the API and the static frontend assets.
