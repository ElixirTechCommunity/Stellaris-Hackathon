# airpay

**GSM-assisted offline payment infrastructure for Android.**  
Built for the dead zones — no internet, no problem.

---

## Overview

airpay is a native Android application that rethinks the payment flow for environments where internet connectivity is unreliable or entirely absent. Instead of depending on data networks, airpay routes payment interactions through the GSM layer — the same infrastructure your carrier uses for voice and SMS — enabling users to complete transactions from virtually anywhere a cellular signal exists.

The project targets a real and underserved gap: hundreds of millions of users in India and emerging markets face dropped connections, high-latency networks, and intermittent data coverage daily. airpay treats offline-first as a hard constraint, not an afterthought.

---

## Team

**Team 7**

| Name | Role | GitHub |
|------|------|--------|
| Chirag Gupta | Developer | [@cgchiraggupta](https://github.com/cgchiraggupta) |
| Sparshika | Team Member | — |
| Meet Khurana | Team Member | — |

---

## Problem

Standard UPI and digital payment flows assume a stable data connection. When that connection degrades — patchy 2G, network congestion, rural dead zones — the entire payment stack fails. Users are left with no fallback, no feedback, and no transaction.

The infrastructure to solve this already exists in every device: GSM. airpay leverages it.

---

## Solution

airpay wraps a guided, minimal mobile interface around a GSM-assisted payment flow. The user experience is deliberately stripped down to reduce failure points. Payment information is reviewed on-device, the interaction is routed through the GSM channel, and the user receives confirmation — all without touching a data packet.

The frontend is engineered for low-cognitive-load operation: large tap targets, clear state transitions, and no dependency on real-time API calls for core payment steps.

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Platform | Android (native) |
| Language | Kotlin |
| Payment Transport | GSM-assisted flow |
| Frontend Reference | Web (hosted on Vercel) |

---

## Architecture

```
User Device
    │
    ├── airpay Android App (Kotlin)
    │       ├── Payment UI Layer         ← guided screens, minimal state
    │       ├── Input Validation         ← on-device, no network required
    │       └── GSM Interaction Module   ← routes payment through carrier layer
    │
    └── GSM Network (carrier infrastructure)
            └── Payment Processing
```

The GSM interaction module handles the transport layer entirely. No REST calls. No WebSocket. No dependency on the user's data plan.

---

## Payment Flow

```
1. Launch airpay
        │
        ▼
2. Enter payment details
   (validated on-device)
        │
        ▼
3. Review transaction screen
        │
        ▼
4. Initiate GSM-assisted flow
        │
        ▼
5. Carrier routes interaction
        │
        ▼
6. Confirmation returned to device
```

---

## Links

| Resource | URL |
|----------|-----|
| Live Demo | [airpaywebsite.vercel.app/about](https://airpaywebsite.vercel.app/about) |
| Video Demo | [YouTube Shorts](https://youtube.com/shorts/viAPPQx1r9U?si=aYHEXjK7b6zZsqi7) |
| Presentation | [Google Drive](https://drive.google.com/file/d/1iFc6vUQjXlpI-kgZgBZYivRLXWyG0q9p/view?usp=sharing) |

---

## Screenshots

<p float="left">
  <img src="./screenshots/airpay-1.jpeg" width="30%" />
  <img src="./screenshots/airpay-2.jpeg" width="30%" />
  <img src="./screenshots/airpay-3.jpeg" width="30%" />
</p>

---

## Running Locally

Demo assets, APK builds, and additional technical documentation are available through the links listed above.

Frontend reference files are located in the `app/` directory. These correspond to the public-facing website and can be served locally with any static file server:

```bash
cd app/
npx serve .
# or
python3 -m http.server 8080
```

For the Android application, build requirements and setup instructions will be published in a follow-up release.

---

## Why GSM

| Factor | Data-dependent apps | airpay |
|--------|---------------------|--------|
| Requires internet | Yes | No |
| Works on 2G | Partial | Yes |
| Works with zero data | No | Yes |
| Rural coverage | Poor | GSM coverage |
| Failure mode | Silent / timeout | Defined fallback |

GSM infrastructure in India covers regions where 4G and even 3G are inconsistent. By treating the carrier layer as the transport rather than IP networking, airpay inherits decades of telecom reliability.

---

## Status

Hackathon build. Core payment flow implemented. Active development ongoing.

---

*Built at [Hackathon Name] · Team 7*
