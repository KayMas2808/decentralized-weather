pragma solidity ^0.8.19;

import "@openzeppelin/contracts/access/Ownable.sol";

contract WeatherData is Ownable {
    struct WeatherEntry {
        bytes32 deviceId;
        string ipfsHash;
        uint256 timestamp;
        bytes32 dataHash;
        bool verified;
    }
    
    mapping(uint256 => WeatherEntry) public weatherEntries;
    mapping(bytes32 => uint256[]) public deviceSubmissions;
    mapping(string => bool) public usedIPFSHashes;
    uint256 public totalEntries;
    
    event WeatherDataSubmitted(
        uint256 indexed entryId,
        bytes32 indexed deviceId,
        string ipfsHash,
        bytes32 dataHash,
        uint256 timestamp
    );
    
    event DataVerified(uint256 indexed entryId);
    
    constructor() Ownable(msg.sender) {}
    
    function submitWeatherData(
        bytes32 deviceId,
        string calldata ipfsHash,
        bytes32 dataHash
    ) external onlyOwner returns (uint256) {
        require(bytes(ipfsHash).length > 0, "IPFS hash required");
        require(!usedIPFSHashes[ipfsHash], "IPFS hash already used");
        
        uint256 entryId = totalEntries++;
        
        weatherEntries[entryId] = WeatherEntry({
            deviceId: deviceId,
            ipfsHash: ipfsHash,
            timestamp: block.timestamp,
            dataHash: dataHash,
            verified: true
        });
        
        deviceSubmissions[deviceId].push(entryId);
        usedIPFSHashes[ipfsHash] = true;
        
        emit WeatherDataSubmitted(entryId, deviceId, ipfsHash, dataHash, block.timestamp);
        emit DataVerified(entryId);
        
        return entryId;
    }
    
    function getWeatherEntry(uint256 entryId) external view returns (WeatherEntry memory) {
        require(entryId < totalEntries, "Entry does not exist");
        return weatherEntries[entryId];
    }
    
    function getDeviceSubmissions(bytes32 deviceId) external view returns (uint256[] memory) {
        return deviceSubmissions[deviceId];
    }
    
    function getLatestEntries(uint256 count) external view returns (WeatherEntry[] memory) {
        if (count > totalEntries) {
            count = totalEntries;
        }
        
        WeatherEntry[] memory entries = new WeatherEntry[](count);
        
        for (uint256 i = 0; i < count; i++) {
            entries[i] = weatherEntries[totalEntries - 1 - i];
        }
        
        return entries;
    }
    
    function getEntriesByTimeRange(uint256 startTime, uint256 endTime) 
        external 
        view 
        returns (WeatherEntry[] memory) 
    {
        require(startTime <= endTime, "Invalid time range");
        
        uint256 count = 0;
        for (uint256 i = 0; i < totalEntries; i++) {
            if (weatherEntries[i].timestamp >= startTime && weatherEntries[i].timestamp <= endTime) {
                count++;
            }
        }
        
        WeatherEntry[] memory entries = new WeatherEntry[](count);
        uint256 index = 0;
        
        for (uint256 i = 0; i < totalEntries; i++) {
            if (weatherEntries[i].timestamp >= startTime && weatherEntries[i].timestamp <= endTime) {
                entries[index] = weatherEntries[i];
                index++;
            }
        }
        
        return entries;
    }
} 