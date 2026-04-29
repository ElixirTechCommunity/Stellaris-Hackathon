# Equifund

## Team Name
OrbitX

## Team Members
| Name | Role | GitHub |
|------|------|--------|
| K. Naga Jethin | Backend | [@Jethin10](https://github.com/Jethin10) |
| Kavaya Jaiswal | Frontend | [@kavyajaiswal007](https://github.com/kavyajaiswal007) |
| Ishani Varshney | Product Thinker | [@ivarshney29](https://github.com/ivarshney29) |
| Divyanshi Yadav | Research | [@divyanshi131107-ai](https://github.com/divyanshi131107-ai) |

## Problem Statement
Crowdfunding platforms are good at helping founders collect money, but weak at enforcing accountability after fundraising. Backers usually lose visibility and control once they contribute, while founders who are genuinely building struggle to prove progress in a structured way.

Equifund solves this with a hybrid crowdfunding protocol where funds are tied to milestones and released only when progress is shown and approved. It supports:

- an India-compliant fiat rail using providers like Razorpay, UPI, and PhonePe
- a global crypto-native rail using USDC and smart-escrow-style logic
- a shared governance layer with milestone proof submission, backer voting, and validator arbitration

The core idea is simple: backers should not just fund a promise, they should fund progress.

In short, Equifund rethinks crowdfunding as a trust system, not just a payment flow.

## Tech Stack
- Frontend: React, TypeScript, Vite, Motion
- Backend: Fastify, TypeScript, Zod
- Database: SQLite for local demo, Postgres-ready for deployment
- Auth: Email auth with provider-ready Google, Facebook, and Apple flows
- Payments: Razorpay-style India rail integration architecture
- Web3 / Escrow modeling: viem, smart-contract escrow architecture, USDC rail modeling
- Deployment: Render Blueprint (`render.yaml`)

## Links
- **Live Demo:** [https://jethin10-abes-hackathon-web.onrender.com](https://jethin10-abes-hackathon-web.onrender.com)
- **Source Repo:** [https://github.com/Jethin10/ABES-hackathon](https://github.com/Jethin10/ABES-hackathon)
- **Video Demo:** Add link
- **Presentation (PPT/PDF):** Add link

## Screenshots
The current live product includes the following app views:

1. **Landing Page**  
   Premium protocol-style homepage with the Equifund brand, zero-fee positioning, and founder / investor navigation.

2. **Core Feature Section**  
   Highlights smart contract escrow, zero platform fees, quadratic voting, and instant liquidity.

3. **Founder Hub**  
   Dashboard showing raised capital, backer metrics, milestones, and escrow unlock tracking.

4. **Investor Portfolio**  
   Dashboard showing invested amount, active projects, and contribution activity synced with the backend.

5. **Judge Demo Flows**  
   Dedicated founder and investor login paths designed for smooth live presentation during judging.

## How to Run Locally
### 1. Backend setup
```bash
npm install
cp .env.example .env
```

On Windows PowerShell:
```powershell
Copy-Item .env.example .env
```

### 2. Start the backend
```bash
npm run dev
```

The backend runs on:
```text
http://localhost:4000
```

### 3. Start the frontend
Open a new terminal and run:

```bash
cd "Frontend/original from gemini"
npm install
cp .env.example .env.local
npm run dev
```

On Windows PowerShell:
```powershell
cd "Frontend/original from gemini"
npm install
Copy-Item .env.example .env.local
npm run dev
```

The frontend runs on:
```text
http://localhost:3000
```

### 4. Demo accounts
Use the built-in judge demo buttons on the login screen, or use seeded credentials such as:

- Founder: `founder@stellaris.dev`
- Investor: `backer1@stellaris.dev`

## Architecture Summary
Equifund uses two parallel rails with one common trust engine.

### India rail
- Backers pay in INR
- Funds are modeled in a compliant fiat flow
- Internal non-transferable ledger tokens represent governance power
- Capital is split into yield deployment and a liquidity buffer

### Global rail
- Backers fund with USDC
- Escrow is modeled through smart-contract-style architecture
- Yield routing is designed for protocols like Aave and Morpho

### Common layer
- milestone engine
- proof submission
- quadratic voting
- quorum thresholds
- validator arbitration
- payout release logic

## Key Innovation
The most important design decision in Equifund is the separation of the funding rail from the governance rail.

- In India, the product stays compliance-aware through a fiat-first flow using familiar providers.
- Globally, the system can support crypto-native capital and smart escrow logic.
- In both cases, the same milestone release engine governs trust.

This gives Equifund a clear real-world deployment path while still preserving the programmability and transparency benefits of modern protocol design.

## Why This Project Matters
Equifund turns crowdfunding from a leap of faith into a system of earned trust. Instead of handing over funds upfront, the platform enforces progress-linked release and creates a more transparent, accountable fundraising model for both founders and backers.
