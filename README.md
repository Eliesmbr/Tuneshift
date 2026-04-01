# <img src="web/public/tuneshift.svg" width="32" height="32" alt=""> Tuneshift

**Free, open-source tool to migrate your Spotify playlists to Tidal.**

No account needed. No data stored. Just upload, connect, and migrate.

<br>

## How it works

```
Export playlists from Spotify  →  Upload CSV  →  Connect Tidal  →  Done
         (via Exportify)
```

1. Export your Spotify playlists at [exportify.app](https://exportify.app)
2. Upload the CSV files to Tuneshift
3. Select which playlists to migrate
4. Connect your Tidal account
5. Hit migrate — tracks are matched by ISRC, playlists created on Tidal

<br>

## Track matching

Tuneshift uses a two-step matching strategy:

- **ISRC lookup** — Most tracks have an International Standard Recording Code. This gives exact matches.
- **Fuzzy search** — Falls back to searching by track name + artist with smart normalization (strips remaster tags, handles spelling variations, duration matching).

In testing, **91/91 tracks** matched successfully with ISRC data from Exportify.

<br>

## Self-hosting

Tuneshift is designed to run on your own server via Docker.

### Prerequisites

- A server with Docker installed
- A [Tidal Developer](https://developer.tidal.com) app (free, takes 2 minutes)

### Setup

```bash
git clone https://github.com/Eliesmbr/Tuneshift.git
cd Tuneshift
cp .env.example .env
```

Edit `.env`:

```env
TIDAL_CLIENT_ID=your_tidal_client_id
BASE_URL=https://yourdomain.com
```

Set the Tidal redirect URI to `https://yourdomain.com/api/auth/tidal/callback`

### Run

```bash
docker compose up --build -d
```

The app runs on port `8080`. Put a reverse proxy (Caddy, Nginx) in front for HTTPS.

<br>

## Architecture

```
┌──────────────────────────────────────┐
│           Docker Container           │
│                                      │
│   React SPA  ←──  Go Backend :8080   │
│   (static)       │                   │
│                  API + OAuth + SSE   │
└──────────────────┼───────────────────┘
                   │
          Tidal API (openapi.tidal.com/v2)
```

- **Backend:** Go (standard library, zero dependencies)
- **Frontend:** React + Tailwind CSS
- **Auth:** Tidal OAuth 2.0 with PKCE
- **Sessions:** AES-256-GCM encrypted HTTP-only cookies
- **Progress:** Server-Sent Events (SSE) for real-time updates
- **Image size:** ~15 MB

<br>

## Security

- No database, no user accounts — fully stateless
- OAuth tokens encrypted in HTTP-only cookies, never exposed to JavaScript
- PKCE flow — no client secret needed
- CSRF protection via OAuth state parameter + SameSite cookies
- Rate limiting on all API endpoints
- Container runs as non-root user
- CSV files parsed in memory, auto-deleted after 30 minutes

<br>

## Why not use the Spotify API directly?

Since February 2026, Spotify requires Premium for API access and limits development apps to 5 users. Extended quota mode is only available to companies with 250k+ monthly active users. Using Exportify's CSV export is the only viable path for a free, public tool.

<br>

## Tech stack

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.22+ (net/http) |
| Frontend | React 19, Tailwind CSS 3 |
| Build | Multi-stage Docker (Node + Go + Alpine) |
| Auth | OAuth 2.0 + PKCE |
| Encryption | AES-256-GCM |
| Progress | Server-Sent Events |
| Matching | ISRC + fuzzy search |

<br>

## License

MIT
