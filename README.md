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
# Decentralized Weather Network: Local Setup & Running Instructions

This guide will walk you through setting up and running the full-stack Decentralized Weather Network locally. We'll cover everything from prerequisites to launching the frontend dashboard, incorporating all necessary configurations and detailed troubleshooting steps encountered during your setup.

### **Phase 0: Initial Preparations (System Setup)**

Before you begin, ensure you have the following tools installed on your system:

1.  **Node.js (version 18+):** Download and install from the official [Node.js website](https://nodejs.org/). `npm` (Node Package Manager) is usually bundled with Node.js.
2.  **Go (version 1.19+):** Download and install from the official [Go website](https://go.dev/dl/).
3.  **Git:** Download and install from the official [Git website](https://git-scm.com/downloads).
4.  **Foundry (Forge & Anvil):** This is the toolkit for Solidity smart contract development.
    * Open your terminal or command prompt.
    * Run the installation command:
        ```bash
        curl -L [https://foundry.paradigm.xyz](https://foundry.paradigm.xyz) | bash
        ```
    * Then, run `foundryup` to install the latest version of Foundry and its components (including Forge and Anvil).
        ```bash
        foundryup
        ```

### **Phase 1: Repository & Environment Setup**

Assuming you have already cloned the `decentralized-weather` repository.

1.  **Navigate into the Cloned Repository:**
    * Open your terminal or command prompt.
    * Change your current directory to the root of the cloned repository:
        ```bash
        cd decentralized-weather
        ```

2.  **Prepare a Local Blockchain for Contracts (Anvil):**
    * Open a **new, dedicated terminal window** (keep your main terminal open for other steps).
    * In this new terminal, start Anvil, your local Ethereum development blockchain:
        ```bash
        anvil
        ```
    * **Crucial Step:** Anvil will display details including its `Listening on` address (your **RPC URL**, typically `http://127.0.0.1:8545`) and several test accounts with their **Private Keys**.
    * **Copy these down:** You'll need the **RPC URL** and at least one **Private Key** for later steps.

3.  **Obtain External API Keys:**
    * **Pinata API Keys (for IPFS storage):**
        1.  Go to the [Pinata Cloud website](https://www.pinata.cloud/).
        2.  Sign up or log in to your account.
        3.  In your dashboard, navigate to the "API Keys" section.
        4.  Generate a new API key. **Copy both the API Key and the API Secret Key immediately**, as the secret key is usually shown only once.
    * **WalletConnect Project ID (for Frontend Wallet Connection):**
        1.  Go to the [WalletConnect Cloud dashboard](https://cloud.walletconnect.com/).
        2.  Sign up or log in.
        3.  Create a new project in your dashboard.
        4.  You will be provided with a unique **Project ID**. Copy this ID.

### **Phase 2: Smart Contract Deployment**

This phase involves configuring, compiling, and deploying your smart contracts to your local Anvil blockchain.

1.  **Navigate to the Contracts Directory:**
    * In your **main terminal window** (where you are currently in `decentralized-weather`):
        ```bash
        cd contracts
        ```

2.  **Update `foundry.toml` Configuration:**
    * Open the file `foundry.toml` located in `decentralized-weather/contracts` using a text editor.
    * **Modify the `[profile.default]` section:**
        * **Comment out or remove** the `solc_version` line.
        * **Add** `auto_detect_solc = true`.
        * **Add** the `excluded_paths` array to prevent compilation of OpenZeppelin's internal test/certora files.
    * **Modify the `[etherscan]` section:**
        * **Comment out** the entire `[etherscan]` section. This prevents Forge from trying to verify on a public explorer, which is not needed for local deployment and avoids `ARBISCAN_API_KEY` errors.

    Your `foundry.toml` should look *exactly* like this (ensure all previous `remappings` are also present):

    ```toml
    [profile.default]
    src = "src"
    out = "out"
    libs = ["lib"]
    remappings = [
        "@openzeppelin/contracts/=lib/openzeppelin-contracts/contracts/",
        "forge-std/=lib/forge-std/src/",
        "erc4626-tests/=lib/openzeppelin-contracts/lib/erc4626-tests/",
        "halmos-cheatcodes/=lib/openzeppelin-contracts/lib/halmos-cheatcodes/src/"
    ]
    # solc_version = "0.8.19" # This line should be commented out or removed
    auto_detect_solc = true
    optimizer = true
    optimizer_runs = 200
    via_ir = false
    excluded_paths = [
        "lib/openzeppelin-contracts/certora",
        "lib/openzeppelin-contracts/test",
        "lib/forge-std/test"
    ]

    [rpc_endpoints]
    arbitrum_sepolia = "[https://sepolia-rollup.arbitrum.io/rpc](https://sepolia-rollup.arbitrum.io/rpc)"

    #[etherscan] # Comment out this line
    #arbitrum_sepolia = { key = "${ARBISCAN_API_KEY}", url = "[https://api-sepolia.arbiscan.io/api](https://api-sepolia.arbiscan.io/api)" } # Comment out this line
    ```
    * **Save the `foundry.toml` file.**

3.  **Clean and Re-initialize Submodules:**
    * This ensures your `lib` dependencies are in a perfect state, resolving stubborn "Source not found" issues.
        ```bash
        rm -rf lib
        git submodule init
        git submodule update
        ```

4.  **Clean Build Artifacts:**
    * Remove any old compilation output:
        ```bash
        forge clean
        ```

5.  **Build the Smart Contracts:**
    * Compile the Solidity contracts using Forge:
        ```bash
        forge build
        ```
    * This command should now compile successfully.

6.  **Deploy the Smart Contracts:**
    * Provide your Anvil RPC URL and one of the private keys (from Anvil's output) as environment variables to the deployment script.

    ```bash
    RPC_URL="[http://127.0.0.1:8545](http://127.0.0.1:8545)" PRIVATE_KEY="YOUR_ANVIL_PRIVATE_KEY" forge script script/Deploy.s.sol --rpc-url $RPC_URL --broadcast
    ```
    * **Replace `http://127.0.0.1:8545`** with the exact RPC URL from your Anvil terminal.
    * **Replace `YOUR_ANVIL_PRIVATE_KEY`** with one of the test private keys provided by Anvil (e.g., `0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`).
    * **Crucial Step:** After successful deployment, the terminal output will display the **deployed addresses** for each smart contract (`WeatherToken`, `DeviceRegistry`, `WeatherData`, `RewardManager`). **Copy these four addresses down carefully.** You will need them for the backend and frontend configurations. Example output:
        ```
        WeatherToken deployed at: 0x5FbDB2315678afecb367f032d93F642f64180aa3
        DeviceRegistry deployed at: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
        WeatherData deployed at: 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0
        RewardManager deployed at: 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
        ```

### **Phase 3: Backend Setup and Execution**

The Go backend acts as the bridge between your clients and the blockchain/IPFS.

1.  **Navigate to the Backend Directory:**
    * In your **main terminal window** (from `contracts` directory):
        ```bash
        cd ../backend
        ```

2.  **Install Go Dependencies:**
    * Clean up and install required Go modules:
        ```bash
        go mod tidy
        ```

3.  **Create and Configure the Backend Environment File:**
    * Create a new file named `.env` inside the `backend` directory: `decentralized-weather/backend/.env`.
    * Open `backend/.env` with a text editor and add the following lines, replacing placeholders with your actual values:

    ```env
    ETHEREUM_RPC=[http://127.0.0.1:8545](http://127.0.0.1:8545)  # Your Anvil RPC URL
    PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 # Your Anvil private key (WITHOUT "0x" PREFIX)
    DEVICE_REGISTRY_ADDRESS=0x0165878A594ca255338adfa4d48449f69242Eb8F      # Deployed DeviceRegistry.sol address (from Phase 2, Step 6)
    WEATHER_DATA_ADDRESS=0xa513E6E4b8f2a923D98304ec87F64353C4D5C853         # Deployed WeatherData.sol address (from Phase 2, Step 6)
    REWARD_MANAGER_ADDRESS=0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6       # Deployed RewardManager.sol address (from Phase 2, Step 6)
    PINATA_API_KEY=YOUR_PINATA_API_KEY # Your Pinata API Key (from Phase 1, Step 3)
    PINATA_SECRET_KEY=YOUR_PINATA_SECRET_KEY # Your Pinata Secret Key (from Phase 1, Step 3)
    RATE_LIMIT_WINDOW=3600             # Default rate limit window in seconds (e.g., 1 hour)
    MAX_SUBMISSIONS_PER_WINDOW=12      # Default max submissions per window
    PORT=8080                          # Port for the backend API
    ```
    * **Important:** Ensure no spaces around the `=` signs. **Remove the `0x` prefix** from the `PRIVATE_KEY` for the Go backend.
    * **Save the `backend/.env` file.**

4.  **Run the Backend Service:**
    * Execute the backend application:
        ```bash
        go run .
        ```
    * The backend should start and display `[GIN-debug] Listening and serving HTTP on :8080`. Keep this terminal window open.

### **Phase 4: Client Setup and Execution**

The Go client simulates weather data submission.

1.  **Navigate to the Client Directory:**
    * Open a **new terminal window** (keep the Anvil and backend terminals open).
    * Change directory to the `client` folder:
        ```bash
        cd decentralized-weather/client
        ```

2.  **Install Go Dependencies:**
    * Clean up and install required Go modules:
        ```bash
        go mod tidy
        ```

3.  **Create and Configure the Client Environment File:**
    * Create a new file named `.env` inside the `client` directory: `decentralized-weather/client/.env`.
    * Open `client/.env` with a text editor and add the following lines:

    ```env
    BACKEND_URL=http://localhost:8080/api # URL of your running backend API
    SUBMISSION_INTERVAL=300              # Data submission interval in seconds (e.g., 5 minutes)
    KEYS_PATH=./device_keys.json         # Path for storing client's cryptographic keys (will be created)
    DEVICE_LOCATION="New York, NY"       # Location for simulated weather data
    ```
    * **Important:** Ensure no spaces around the `=` signs.
    * **Save the `client/.env` file.**

4.  **Export Client Environment Variables (Crucial for Go Apps):**
    * The Go client does not automatically load variables from `.env`. You must export them to your shell session.
        ```bash
        export BACKEND_URL="http://localhost:8080/api"
        export SUBMISSION_INTERVAL="300"
        export KEYS_PATH="./device_keys.json"
        export DEVICE_LOCATION="New York, NY"
        ```

5.  **Register the Device:**
    * Ensure your backend is running in its dedicated terminal.
    * Run the client to register itself with the backend:
        ```bash
        go run . register
        ```
    * This will create `device_keys.json` and register the client. You should see "Creating new device keys..." or "Loaded existing device keys..." and no errors.

6.  **Start Submitting Data:**
    * Now, run the client to start simulating and submitting weather data to the backend:
        ```bash
        go run .
        ```
    * You should see output indicating data submissions to the backend. Keep this terminal window open.

### **Phase 5: Frontend Setup and Execution**

The React frontend provides the dashboard to visualize data.

1.  **Navigate to the Frontend Directory:**
    * Open a **new terminal window** (keep all previous terminals open for Anvil, Backend, and Client).
    * Change directory to the `frontend` folder:
        ```bash
        cd decentralized-weather/frontend
        ```

2.  **Install Node.js (npm) Dependencies:**
    * Install all the required frontend packages:
        ```bash
        npm install
        ```

3.  **Create and Configure the Frontend Environment File:**
    * Create a new file named `.env` inside the `frontend` directory: `decentralized-weather/frontend/.env`.
    * Open `frontend/.env` with a text editor and add the following lines, replacing placeholders:

    ```env
    VITE_BACKEND_URL=http://localhost:8080/api # URL of your running backend API
    VITE_WEATHER_TOKEN_ADDRESS=0x5FC8d32690cc91D4c39d9d3abcBD16989F875707           # Deployed WeatherToken.sol address (from Phase 2, Step 6)
    VITE_DEVICE_REGISTRY_ADDRESS=0x0165878A594ca255338adfa4d48449f69242Eb8F        # Deployed DeviceRegistry.sol address (from Phase 2, Step 6)
    VITE_WEATHER_DATA_ADDRESS=0xa513E6E4b8f2a923D98304ec87F64353C4D5C853           # Deployed WeatherData.sol address (from Phase 2, Step 6)
    VITE_REWARD_MANAGER_ADDRESS=0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6         # Deployed RewardManager.sol address (from Phase 2, Step 6)
    VITE_WALLET_CONNECT_PROJECT_ID=YOUR_WALLET_CONNECT_PROJECT_ID # Your WalletConnect Project ID (from Phase 1, Step 3)
    ```
    * **Important:** Ensure no spaces around the `=` signs. Frontend environment variables in Vite projects need the `VITE_` prefix.
    * **Save the `frontend/.env` file.**

4.  **Launch the Frontend Dashboard:**
    * Start the development server for the React application:
        ```bash
        npm run dev
        ```
    * This command will typically provide a local URL (e.g., `http://localhost:5173`). Open this URL in your web browser.

You should now have all components of the Decentralized Weather Network running locally!
