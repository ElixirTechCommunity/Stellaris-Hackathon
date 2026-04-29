# Aurva Dashboard - Quick Start

## Prerequisites
- Node.js 18+ 
- Control plane running on port 9090

## Installation

```bash
cd dashboard
npm install
```

## Configuration

Create `.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:9090
```

## Development

```bash
npm run dev
```

Visit http://localhost:3000

## Production Build

```bash
npm run build
npm start
```

## Routes

- `/` → Redirects to `/connect`
- `/connect` → IAM role onboarding form
- `/scan/:scanId` → Live scan progress terminal
- `/dashboard` → Compliance findings and reports

## API Endpoints Expected

The dashboard expects these endpoints on the control plane:

```
POST /api/scan          → Start new scan
GET  /api/scan/:id      → Poll scan status
GET  /api/findings      → Get all findings
GET  /api/report/pdf    → Download PDF report
GET  /health            → Health check
```

## Mock Data

If control plane is not running, the dashboard will show connection errors.
For demo purposes, ensure the control plane is running on port 9090.

## Design Notes

- Dark theme only (#080C14 background)
- Amber (#F97316) = danger/urgency
- Sky blue (#0EA5E9) = safe/healthy
- Geist font for UI
- JetBrains Mono for all data/numbers

## Troubleshooting

**Port 3000 already in use:**
```bash
PORT=3001 npm run dev
```

**Control plane not responding:**
- Check control plane is running: `curl http://localhost:9090/health`
- Check .env.local has correct API_URL

**Build errors:**
- Delete .next folder: `rm -rf .next`
- Clear node_modules: `rm -rf node_modules && npm install`
