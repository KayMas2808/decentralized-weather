import { useState } from 'react';
import { MapPin, Thermometer } from 'lucide-react';

interface WeatherLocation {
  id: string;
  location: string;
  lat: number;
  lng: number;
  temperature: number;
  humidity: number;
  deviceId: string;
}

const WeatherMap = () => {
  const [locations] = useState<WeatherLocation[]>([
    {
      id: '1',
      location: 'New York, NY',
      lat: 40.7128,
      lng: -74.0060,
      temperature: 22.5,
      humidity: 65,
      deviceId: '0x1a2b3c4d5e6f',
    },
    {
      id: '2',
      location: 'London, UK',
      lat: 51.5074,
      lng: -0.1278,
      temperature: 18.3,
      humidity: 72,
      deviceId: '0x9a8b7c6d5e4f',
    },
    {
      id: '3',
      location: 'Tokyo, Japan',
      lat: 35.6762,
      lng: 139.6503,
      temperature: 25.1,
      humidity: 58,
      deviceId: '0x3f2e1d0c9b8a',
    },
  ]);

  return (
    <div className="h-64 bg-gradient-to-br from-blue-100 to-green-100 rounded-lg relative overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-r from-blue-500/10 to-green-500/10"></div>
      
      <div className="relative z-10 p-4">
        <div className="grid grid-cols-1 gap-4">
          {locations.map((location) => (
            <div
              key={location.id}
              className="bg-white/90 backdrop-blur-sm rounded-lg p-3 shadow-sm border border-white/50"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <MapPin className="h-4 w-4 text-blue-600 mr-2" />
                  <span className="text-sm font-medium text-gray-900">
                    {location.location}
                  </span>
                </div>
                <div className="flex items-center">
                  <Thermometer className="h-4 w-4 text-red-500 mr-1" />
                  <span className="text-sm font-bold text-gray-900">
                    {location.temperature}°C
                  </span>
                </div>
              </div>
              <div className="mt-1 text-xs text-gray-600">
                Device: {location.deviceId} • Humidity: {location.humidity}%
              </div>
            </div>
          ))}
        </div>
      </div>
      
      <div className="absolute bottom-2 right-2 text-xs text-gray-500 bg-white/80 px-2 py-1 rounded">
        Live Weather Stations
      </div>
    </div>
  );
};

export default WeatherMap; 