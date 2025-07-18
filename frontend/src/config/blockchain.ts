import { getDefaultConfig } from '@rainbow-me/rainbowkit';
import { arbitrumSepolia } from 'wagmi/chains';

export const config = getDefaultConfig({
  appName: 'Weather Network',
  projectId: import.meta.env.VITE_WALLET_CONNECT_PROJECT_ID || 'demo-project-id',
  chains: [arbitrumSepolia],
  ssr: false,
});

export const contractAddresses = {
  weatherToken: import.meta.env.VITE_WEATHER_TOKEN_ADDRESS || '0x0000000000000000000000000000000000000000',
  deviceRegistry: import.meta.env.VITE_DEVICE_REGISTRY_ADDRESS || '0x0000000000000000000000000000000000000000',
  weatherData: import.meta.env.VITE_WEATHER_DATA_ADDRESS || '0x0000000000000000000000000000000000000000',
  rewardManager: import.meta.env.VITE_REWARD_MANAGER_ADDRESS || '0x0000000000000000000000000000000000000000',
};

export const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080/api'; 