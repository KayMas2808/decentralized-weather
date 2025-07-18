import { useState } from 'react';
import { useAccount } from 'wagmi';
import { Smartphone, MapPin, Key, CheckCircle } from 'lucide-react';

const DeviceRegistration = () => {
  const [deviceId, setDeviceId] = useState('');
  const [publicKey, setPublicKey] = useState('');
  const [location, setLocation] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [registrationStatus, setRegistrationStatus] = useState<'idle' | 'success' | 'error'>('idle');
  
  const { isConnected } = useAccount();

  const generateDeviceId = () => {
    const randomBytes = new Uint8Array(16);
    crypto.getRandomValues(randomBytes);
    const deviceId = '0x' + Array.from(randomBytes, byte => byte.toString(16).padStart(2, '0')).join('');
    setDeviceId(deviceId);
  };

  const generateKeyPair = () => {
    const randomBytes = new Uint8Array(32);
    crypto.getRandomValues(randomBytes);
    const publicKey = '0x' + Array.from(randomBytes, byte => byte.toString(16).padStart(2, '0')).join('');
    setPublicKey(publicKey);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isConnected) return;

    setIsSubmitting(true);
    try {
      setRegistrationStatus('success');
      setTimeout(() => setRegistrationStatus('idle'), 3000);
    } catch (error) {
      console.error('Registration failed:', error);
      setRegistrationStatus('error');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isConnected) {
    return (
      <div className="text-center py-8">
        <Smartphone className="mx-auto h-12 w-12 text-gray-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">Wallet not connected</h3>
        <p className="mt-1 text-sm text-gray-500">
          Please connect your wallet to register a device.
        </p>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto">
      {registrationStatus === 'success' && (
        <div className="mb-4 bg-green-50 border border-green-200 rounded-md p-4">
          <div className="flex">
            <CheckCircle className="h-5 w-5 text-green-400" />
            <div className="ml-3">
              <h3 className="text-sm font-medium text-green-800">Registration Successful</h3>
              <div className="mt-2 text-sm text-green-700">
                <p>Your device has been registered successfully!</p>
              </div>
            </div>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="device-id" className="block text-sm font-medium text-gray-700">
            Device ID
          </label>
          <div className="mt-1 flex rounded-md shadow-sm">
            <input
              type="text"
              id="device-id"
              value={deviceId}
              onChange={(e) => setDeviceId(e.target.value)}
              className="flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-l-md border border-gray-300 focus:ring-blue-500 focus:border-blue-500"
              placeholder="0x..."
              required
            />
            <button
              type="button"
              onClick={generateDeviceId}
              className="inline-flex items-center px-3 py-2 border border-l-0 border-gray-300 rounded-r-md bg-gray-50 text-gray-500 text-sm hover:bg-gray-100"
            >
              <Smartphone className="h-4 w-4" />
            </button>
          </div>
        </div>

        <div>
          <label htmlFor="public-key" className="block text-sm font-medium text-gray-700">
            Public Key
          </label>
          <div className="mt-1 flex rounded-md shadow-sm">
            <input
              type="text"
              id="public-key"
              value={publicKey}
              onChange={(e) => setPublicKey(e.target.value)}
              className="flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-l-md border border-gray-300 focus:ring-blue-500 focus:border-blue-500"
              placeholder="0x..."
              required
            />
            <button
              type="button"
              onClick={generateKeyPair}
              className="inline-flex items-center px-3 py-2 border border-l-0 border-gray-300 rounded-r-md bg-gray-50 text-gray-500 text-sm hover:bg-gray-100"
            >
              <Key className="h-4 w-4" />
            </button>
          </div>
        </div>

        <div>
          <label htmlFor="location" className="block text-sm font-medium text-gray-700">
            Location
          </label>
          <div className="mt-1 relative rounded-md shadow-sm">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <MapPin className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              id="location"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              placeholder="New York, NY"
              required
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            disabled={isSubmitting || !deviceId || !publicKey || !location}
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            {isSubmitting ? 'Registering...' : 'Register Device'}
          </button>
        </div>
      </form>

      <div className="mt-6 text-sm text-gray-600">
        <h4 className="font-medium text-gray-900">Instructions:</h4>
        <ul className="mt-2 list-disc list-inside space-y-1">
          <li>Generate a unique device ID for your weather station</li>
          <li>Create a cryptographic key pair for secure data submission</li>
          <li>Specify the physical location of your device</li>
          <li>After registration, use the client software to submit weather data</li>
        </ul>
      </div>
    </div>
  );
};

export default DeviceRegistration; 