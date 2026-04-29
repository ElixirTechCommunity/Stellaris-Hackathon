# Aurva Dashboard

Production-ready DPDP Act 2023 compliance dashboard for Indian startups.

## Design Philosophy

**Bloomberg Terminal meets Linear App**
- Dark theme (#080C14 background)
- Amber (#F97316) for danger/urgency
- Sky blue (#0EA5E9) for safety/health
- Dense, data-rich interface
- Every animation communicates state

## Tech Stack

- **Framework**: Next.js 14 App Router + TypeScript
- **Styling**: Tailwind CSS 4 with custom tokens
- **State**: 
  - React Query (TanStack) for server state
  - Zustand for client state
- **Visualization**: Recharts (donut + bar charts)
- **Animation**: Framer Motion (purposeful only)
- **Icons**: Lucide React
- **Fonts**: Geist (UI), JetBrains Mono (data)

## Routes

### `/connect` - IAM Role Onboarding
Entry point with split layout:
- Left: Brand statement + trust signals
- Right: AWS credentials form with live validation
- Triggers POST /api/scan → redirects to /scan/:id

### `/scan/[scanId]` - Live Scan Progress
Terminal-style live feed:
- Elapsed time counter
- Real-time resource discovery feed
- Color-coded scan results (amber=PII, blue=clean)
- Live updating stats sidebar
- Auto-redirects to dashboard on completion

### `/dashboard` - Main Compliance View
The money screen:
- 4 animated stat cards (compliance score, PII count, critical, last scan)
- Sortable/filterable findings table
- Risk distribution donut chart
- PII type breakdown bar chart
- Slide-in finding detail drawer
- PDF report download

## API Integration

Base URL: `http://localhost:9090`

```typescript
POST /api/scan          // Trigger scan
GET  /api/scan/:id      // Poll status (2s interval)
GET  /api/findings      // Dashboard data
GET  /api/report/pdf    // Download report
GET  /health            // Connection check
```

React Query handles:
- Automatic polling during scans
- Query invalidation on scan completion
- Stale time management (30s for findings)
- Error/loading states

## Key Components

### Dashboard Components
- `FindingsTable` - Sortable table with virtualization support
- `FindingDrawer` - Slide-in detail panel with remediation steps
- `RiskDonut` - Interactive risk distribution chart
- `PIIBreakdown` - Horizontal bar chart of PII types
- `StatCards` - 4 animated metric cards

### Scan Components
- `ResourceFeed` - Terminal-style live event stream
- `ScanStats` - Live counters for resources/PII

### UI Components
- `StatCard` - Animated number counter
- `RiskBadge` - Color-coded risk level pill
- `PIIBadge` - Monospace PII type tag
- `ConfidenceBar` - Thin gradient progress bar
- `ScanDot` - Pulsing scan indicator

## State Management

### Zustand Store (`lib/store.ts`)
- Selected finding for drawer
- Risk level filters
- PII type filters
- Table sorting state

### React Query
- Server state caching
- Automatic background refetching
- Optimistic updates
- Error retry logic

## Development

```bash
# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

Environment variables:
```env
NEXT_PUBLIC_API_URL=http://localhost:9090
```

## Design Tokens

```css
colors: {
  bg: '#080C14',           // Background
  surface: '#0D1117',      // Cards
  surface2: '#161B22',     // Elevated surfaces
  border: 'rgba(255,255,255,0.06)',
  amber: '#F97316',        // Danger/urgency
  sky: '#0EA5E9',          // Safe/healthy
  critical: '#EF4444',     // Critical risk
  high: '#F97316',         // High risk
  medium: '#EAB308',       // Medium risk
  low: '#22C55E',          // Low risk
  muted: '#6B7280',        // Secondary text
}
```

## Performance

- Route-based code splitting
- React Query caching reduces API calls
- Framer Motion animations use GPU
- Recharts lazy loads chart library
- Table virtualization for 100+ rows

## Quality Bar

Every screen looks YC Demo Day ready:
- ✅ Non-technical founders understand risk in <10 seconds
- ✅ Terminal feed is the wow moment
- ✅ Compliance score is the hook
- ✅ PDF download is the close

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
