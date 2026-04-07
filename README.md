# <img src="web/public/tuneshift.svg" width="32" height="32" alt=""> Tuneshift

**Free, open-source tool to migrate your playlists to Tidal.**

No account needed. No data stored. Just upload, connect, and migrate.

<p align="center">
  <img src="assets/screenshot.png?v=2" alt="Tuneshift screenshot" width="700">
</p>

## Supported sources

| Source | Method |
|--------|--------|
| Spotify | CSV export via [Exportify](https://exportify.app) |
| YouTube Music | ZIP export via [Google Takeout](https://takeout.google.com) |
| Apple Music | Coming soon |
| Amazon Music | No API available |

## How it works

```
Export your playlists  ->  Upload file  ->  Connect Tidal  ->  Done
```

### Spotify

1. Export your playlists at [exportify.app](https://exportify.app)
2. Upload the CSV files to Tuneshift
3. Select which playlists to migrate
4. Connect your Tidal account
5. Hit migrate - tracks are matched by ISRC, playlists created on Tidal

### YouTube Music

1. Go to [takeout.google.com](https://takeout.google.com)
2. Deselect all, then select only "YouTube and YouTube Music"
3. Under "All YouTube data included", select **Music library** and **Playlists**
4. Create the export, wait for the email, download the ZIP
5. Upload the ZIP to Tuneshift and continue as above

## Track matching

Tuneshift uses a two-step matching strategy:

- **ISRC lookup** - Most tracks have an International Standard Recording Code. This gives exact matches (Spotify CSVs include ISRCs).
- **Fuzzy search** - Falls back to searching by track name + artist with smart normalization (strips remaster tags, handles spelling variations, duration matching).

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

## Architecture

```
┌──────────────────────────────────────┐
│           Docker Container           │
│                                      │
│   React SPA  <-──  Go Backend :8080   │
│   (static)       │                   │
│                  API + OAuth + SSE   │
└──────────────────┼───────────────────┘
                   │
          Tidal API (openapi.tidal.com/v2)
```

- **Backend:** Go (standard library, zero dependencies)
- **Frontend:** React + Tailwind CSS v4
- **Auth:** Tidal OAuth 2.0 with PKCE
- **Sessions:** AES-256-GCM encrypted HTTP-only cookies
- **Progress:** Server-Sent Events (SSE) for real-time updates
- **Image size:** ~15 MB

## Security

- No database, no user accounts - fully stateless
- OAuth tokens encrypted in HTTP-only cookies, never exposed to JavaScript
- PKCE flow - no client secret needed
- CSRF protection via OAuth state parameter + SameSite cookies
- Rate limiting on all API endpoints
- Container runs as non-root user
- Uploaded files parsed in memory, auto-deleted after 30 minutes

## Why not use Spotify/YouTube APIs directly?

> [!NOTE]
> **Spotify:** Since February 2026, Spotify requires a Premium subscription for Web API access and limits new developer apps to just 5 manually allowlisted users. [Exportify](https://exportify.app) was registered before these restrictions and is grandfathered in with full API access.
>
> **YouTube Music:** Google requires an extended CASA (Cloud Application Security Assessment) security audit for production OAuth apps - disproportionate for a small open-source project. [Google Takeout](https://takeout.google.com) lets users export their own data directly, no API needed.

## Tech stack

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.24 (net/http) |
| Frontend | React 19, Tailwind CSS v4 |
| Build | Multi-stage Docker (Node + Go + Alpine) |
| Auth | OAuth 2.0 + PKCE |
| Encryption | AES-256-GCM |
| Progress | Server-Sent Events |
| Matching | ISRC + fuzzy search |

## License

MIT
