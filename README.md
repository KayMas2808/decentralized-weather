# Decentralized Weather Network

A full-stack decentralized application that incentivizes contributors to run weather data collection clients and rewards them with tokens for verified data submissions.

## Project Architecture

This project consists of four main components:

### 1. Smart Contracts (Solidity)
- **WeatherToken.sol**: ERC-20 token ($WTHR) for rewards
- **DeviceRegistry.sol**: Manages registered weather data contributors
- **WeatherData.sol**: On-chain logbook storing IPFS hashes of weather data
- **RewardManager.sol**: Handles token distribution logic

### 2. Go Client Application
A lightweight client that contributors run to submit weather data:
- Generates and manages cryptographic keypairs
- Simulates weather sensor readings
- Signs and submits data to the backend verifier
- Configurable submission intervals

### 3. Go Backend API
A Gin Gonic-based verifier service that acts as a gatekeeper:
- Verifies cryptographic signatures from devices
- Implements rate limiting (Proof-of-Physical-Work)
- Uploads verified data to IPFS via Pinata
- Records data hashes on the blockchain
- Provides REST API for frontend

### 4. React Frontend Dashboard
A modern web interface built with Vite, Wagmi, and RainbowKit:
- Wallet connection via RainbowKit
- Device registration interface
- Real-time weather data visualization
- Interactive dashboard with charts and maps

## Technology Stack

- **Blockchain**: Arbitrum Sepolia (low fees, fast transactions)
- **Smart Contracts**: Solidity with Foundry framework
- **Backend**: Go with Gin Gonic
- **Frontend**: React with TypeScript, Vite, Wagmi/Viem, RainbowKit
- **Storage**: IPFS via Pinata API
- **Styling**: Tailwind CSS
- **Charts**: Recharts
- **Icons**: Lucide React

## Quick Start

### Prerequisites
- Node.js 18+ and npm
- Go 1.19+
- Git

### 1. Clone and Setup
```bash
git clone <repository-url>
cd decentralized-weather
```

### 2. Deploy Smart Contracts
```bash
cd contracts
forge build
# Configure .env with your private key and RPC URL
forge script script/Deploy.s.sol --rpc-url arbitrum_sepolia --broadcast --verify
```

### 3. Start Backend Service
```bash
cd backend
go mod tidy
# Configure .env with contract addresses and API keys
go run .
```

### 4. Run Weather Client
```bash
cd client
go mod tidy
# Configure .env with backend URL and device location
go run . register  # Register device first
go run .           # Start submitting data
```

### 5. Launch Frontend Dashboard
```bash
cd frontend
npm install
# Configure .env with contract addresses
npm run dev
```

## Configuration

### Backend Environment Variables
```
ETHEREUM_RPC=https://sepolia-rollup.arbitrum.io/rpc
PRIVATE_KEY=your-private-key-here
DEVICE_REGISTRY_ADDRESS=deployed-contract-address
WEATHER_DATA_ADDRESS=deployed-contract-address
REWARD_MANAGER_ADDRESS=deployed-contract-address
PINATA_API_KEY=your-pinata-api-key
PINATA_SECRET_KEY=your-pinata-secret-key
RATE_LIMIT_WINDOW=3600
MAX_SUBMISSIONS_PER_WINDOW=12
PORT=8080
```

### Frontend Environment Variables
```
VITE_BACKEND_URL=http://localhost:8080/api
VITE_WEATHER_TOKEN_ADDRESS=deployed-contract-address
VITE_DEVICE_REGISTRY_ADDRESS=deployed-contract-address
VITE_WEATHER_DATA_ADDRESS=deployed-contract-address
VITE_REWARD_MANAGER_ADDRESS=deployed-contract-address
VITE_WALLET_CONNECT_PROJECT_ID=your-wallet-connect-project-id
```

### Client Environment Variables
```
BACKEND_URL=http://localhost:8080/api
SUBMISSION_INTERVAL=300
KEYS_PATH=./device_keys.json
DEVICE_LOCATION=New York, NY
```

## API Endpoints

### Backend REST API
- `POST /api/register` - Register a new weather device
- `POST /api/submit` - Submit weather data from client
- `GET /api/data` - Get all weather data
- `GET /api/data/latest` - Get latest weather submissions
- `GET /api/devices` - Get registered devices
- `GET /api/health` - Health check


## Reward System

### Token Economics
- Base reward: 10 WTHR per submission
- Bonus multiplier: 150% for consistent devices
- Loyalty bonus: Up to 120% for devices with 100+ submissions
- Daily limit: 1000 WTHR per day network-wide

### Distribution Logic
- Automatic rewards for verified submissions
- Pro-rated based on submission frequency
- Device reliability scoring
- Anti-spam protection

## IPFS Integration

### Data Storage
- Weather data JSON stored on IPFS
- Pinata service for reliable pinning
- Content addressing for integrity
- Only IPFS hashes stored on-chain

### Data Structure
```json
{
  "device_id": "0x...",
  "location": "New York, NY",
  "temperature": 22.5,
  "humidity": 65,
  "pressure": 1013.2,
  "wind_speed": 12.5,
  "wind_direction": "SW",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Development

### Smart Contract Testing
```bash
cd contracts
forge test
```

### Backend Testing
```bash
cd backend
go test ./...
```

### Frontend Development
```bash
cd frontend
npm run dev
```
