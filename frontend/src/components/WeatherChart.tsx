import { Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts';

const WeatherChart = () => {
  const temperatureData = [
    { time: '00:00', temperature: 18.5, humidity: 65 },
    { time: '04:00', temperature: 16.2, humidity: 72 },
    { time: '08:00', temperature: 22.1, humidity: 58 },
    { time: '12:00', temperature: 28.3, humidity: 45 },
    { time: '16:00', temperature: 31.2, humidity: 42 },
    { time: '20:00', temperature: 26.8, humidity: 55 },
    { time: '24:00', temperature: 21.4, humidity: 63 },
  ];

  return (
    <div className="h-64">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={temperatureData} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
          <defs>
            <linearGradient id="temperatureGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#3B82F6" stopOpacity={0.8}/>
              <stop offset="95%" stopColor="#3B82F6" stopOpacity={0.1}/>
            </linearGradient>
            <linearGradient id="humidityGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#10B981" stopOpacity={0.8}/>
              <stop offset="95%" stopColor="#10B981" stopOpacity={0.1}/>
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
          <XAxis 
            dataKey="time" 
            stroke="#6B7280"
            fontSize={12}
          />
          <YAxis 
            stroke="#6B7280"
            fontSize={12}
          />
          <Tooltip 
            contentStyle={{
              backgroundColor: 'white',
              border: '1px solid #E5E7EB',
              borderRadius: '8px',
              boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
            }}
            formatter={(value: number, name: string) => [
              `${value}${name === 'temperature' ? '°C' : '%'}`,
              name === 'temperature' ? 'Temperature' : 'Humidity'
            ]}
          />
          <Area
            type="monotone"
            dataKey="temperature"
            stroke="#3B82F6"
            fillOpacity={1}
            fill="url(#temperatureGradient)"
            strokeWidth={2}
          />
          <Line
            type="monotone"
            dataKey="humidity"
            stroke="#10B981"
            strokeWidth={2}
            dot={{ fill: '#10B981', strokeWidth: 2, r: 4 }}
            activeDot={{ r: 6, stroke: '#10B981', strokeWidth: 2 }}
          />
        </AreaChart>
      </ResponsiveContainer>
      
      <div className="flex justify-center mt-4 space-x-6">
        <div className="flex items-center">
          <div className="w-3 h-3 bg-blue-500 rounded-full mr-2"></div>
          <span className="text-sm text-gray-600">Temperature (°C)</span>
        </div>
        <div className="flex items-center">
          <div className="w-3 h-3 bg-green-500 rounded-full mr-2"></div>
          <span className="text-sm text-gray-600">Humidity (%)</span>
        </div>
      </div>
    </div>
  );
};

export default WeatherChart; 