import { useState } from 'react';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { useAccount } from 'wagmi';
import WeatherMap from './WeatherMap';
import DeviceRegistration from './DeviceRegistration';
import WeatherChart from './WeatherChart';
import DeviceList from './DeviceList';
import { Cloud, Thermometer, Droplets, Wind } from 'lucide-react';

const Dashboard = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const { isConnected } = useAccount();

  const tabs = [
    { id: 'overview', label: 'Overview', icon: Cloud },
    { id: 'devices', label: 'Devices', icon: Thermometer },
    { id: 'data', label: 'Data', icon: Droplets },
    { id: 'register', label: 'Register Device', icon: Wind },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <header className="bg-white shadow-lg border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div className="flex items-center">
              <Cloud className="h-8 w-8 text-blue-600 mr-3" />
              <h1 className="text-3xl font-bold text-gray-900">
                Decentralized Weather Network
              </h1>
            </div>
            <ConnectButton />
          </div>
        </div>
      </header>

      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex space-x-8">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`flex items-center py-4 px-2 text-sm font-medium border-b-2 transition-colors ${
                    activeTab === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <Icon className="h-5 w-5 mr-2" />
                  {tab.label}
                </button>
              );
            })}
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {!isConnected && activeTab !== 'overview' && (
          <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4 mb-6">
            <div className="flex">
              <div className="ml-3">
                <h3 className="text-sm font-medium text-yellow-800">
                  Wallet Connection Required
                </h3>
                <div className="mt-2 text-sm text-yellow-700">
                  <p>Please connect your wallet to access this feature.</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'overview' && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="bg-white rounded-lg shadow p-6">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">
                  Live Weather Map
                </h2>
                <WeatherMap />
              </div>
              <div className="bg-white rounded-lg shadow p-6">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">
                  Temperature Trends
                </h2>
                <WeatherChart />
              </div>
            </div>
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">
                Recent Submissions
              </h2>
              <DeviceList limit={5} />
            </div>
          </div>
        )}

        {activeTab === 'devices' && (
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              Registered Devices
            </h2>
            <DeviceList />
          </div>
        )}

        {activeTab === 'data' && (
          <div className="space-y-6">
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">
                Weather Data Analytics
              </h2>
              <WeatherChart />
            </div>
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">
                All Weather Data
              </h2>
              <DeviceList showAll />
            </div>
          </div>
        )}

        {activeTab === 'register' && (
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              Register New Device
            </h2>
            <DeviceRegistration />
          </div>
        )}
      </main>
    </div>
  );
};

export default Dashboard; 