import { useEffect, useState } from 'react';
import { Smartphone, MapPin, Clock, Activity, Thermometer, Droplets, Wind } from 'lucide-react';
import { backendUrl } from '../config/blockchain';

interface DeviceData {
  device_id: string;
  location: string;
  temperature: number;
  humidity: number;
  pressure: number;
  wind_speed: number;
  wind_direction: string;
  timestamp: string;
  ipfs_hash: string;
}

interface DeviceListProps {
  limit?: number;
  showAll?: boolean;
}

const DeviceList = ({ limit, showAll = false }: DeviceListProps) => {
  const [data, setData] = useState<DeviceData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const endpoint = showAll ? '/data' : '/data/latest';
        const url = `${backendUrl}${endpoint}${limit ? `?limit=${limit}` : ''}`;
        
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Failed to fetch data');
        }
        
        const result = await response.json();
        setData(result.data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
        setData([
          {
            device_id: '0x1a2b3c4d5e6f',
            location: 'New York, NY',
            temperature: 22.5,
            humidity: 65,
            pressure: 1013.2,
            wind_speed: 12.5,
            wind_direction: 'SW',
            timestamp: new Date().toISOString(),
            ipfs_hash: 'QmYjtig7VJQ6XsnUjqqJvj7QaMcCAwtrgNdahSiFofrE7o',
          },
          {
            device_id: '0x9a8b7c6d5e4f',
            location: 'London, UK',
            temperature: 18.3,
            humidity: 72,
            pressure: 1009.8,
            wind_speed: 8.2,
            wind_direction: 'W',
            timestamp: new Date(Date.now() - 300000).toISOString(),
            ipfs_hash: 'QmPK1s3pNYLi9ERiq3BDxKa4XosgWwFRQUydHUtz4YgpqB',
          },
        ].slice(0, limit));
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [limit, showAll]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span className="ml-2 text-gray-600">Loading weather data...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <Smartphone className="mx-auto h-12 w-12 text-red-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">Error loading data</h3>
        <p className="mt-1 text-sm text-gray-500">{error}</p>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="text-center py-8">
        <Activity className="mx-auto h-12 w-12 text-gray-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">No weather data</h3>
        <p className="mt-1 text-sm text-gray-500">
          No weather data submissions found.
        </p>
      </div>
    );
  }

  const formatTimestamp = (timestamp: string) => {
    try {
      return new Date(timestamp).toLocaleString();
    } catch {
      return timestamp;
    }
  };

  const getStatusColor = (timestamp: string) => {
    const now = Date.now();
    const submissionTime = new Date(timestamp).getTime();
    const diffMinutes = (now - submissionTime) / (1000 * 60);
    
    if (diffMinutes < 10) return 'bg-green-100 text-green-800';
    if (diffMinutes < 60) return 'bg-yellow-100 text-yellow-800';
    return 'bg-red-100 text-red-800';
  };

  const getStatusText = (timestamp: string) => {
    const now = Date.now();
    const submissionTime = new Date(timestamp).getTime();
    const diffMinutes = (now - submissionTime) / (1000 * 60);
    
    if (diffMinutes < 10) return 'Active';
    if (diffMinutes < 60) return 'Recent';
    return 'Inactive';
  };

  return (
    <div className="space-y-4">
      {data.map((item, index) => (
        <div key={`${item.device_id}-${index}`} className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <div className="flex items-center space-x-2 mb-2">
                <Smartphone className="h-4 w-4 text-blue-600" />
                <span className="font-medium text-gray-900">
                  {item.device_id}
                </span>
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(item.timestamp)}`}>
                  {getStatusText(item.timestamp)}
                </span>
              </div>
              
              <div className="flex items-center text-sm text-gray-600 mb-2">
                <MapPin className="h-4 w-4 mr-1" />
                <span>{item.location}</span>
                <Clock className="h-4 w-4 ml-4 mr-1" />
                <span>{formatTimestamp(item.timestamp)}</span>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                <div className="flex items-center">
                  <Thermometer className="h-4 w-4 text-red-500 mr-1" />
                  <span className="text-gray-600">Temp:</span>
                  <span className="font-medium ml-1">{item.temperature.toFixed(1)}°C</span>
                </div>
                
                <div className="flex items-center">
                  <Droplets className="h-4 w-4 text-blue-500 mr-1" />
                  <span className="text-gray-600">Humidity:</span>
                  <span className="font-medium ml-1">{item.humidity.toFixed(0)}%</span>
                </div>
                
                <div className="flex items-center">
                  <Wind className="h-4 w-4 text-green-500 mr-1" />
                  <span className="text-gray-600">Wind:</span>
                  <span className="font-medium ml-1">{item.wind_speed.toFixed(1)} km/h {item.wind_direction}</span>
                </div>
                
                <div className="flex items-center">
                  <Activity className="h-4 w-4 text-purple-500 mr-1" />
                  <span className="text-gray-600">Pressure:</span>
                  <span className="font-medium ml-1">{item.pressure.toFixed(0)} hPa</span>
                </div>
              </div>

              {showAll && (
                <div className="mt-2 text-xs text-gray-500">
                  <span className="font-medium">IPFS:</span> {item.ipfs_hash}
                </div>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default DeviceList; 