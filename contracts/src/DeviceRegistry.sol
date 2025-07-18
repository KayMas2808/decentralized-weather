pragma solidity ^0.8.19;

import "@openzeppelin/contracts/access/Ownable.sol";

contract DeviceRegistry is Ownable {
    struct Device {
        address owner;
        string publicKey;
        uint256 registrationTime;
        bool isActive;
        uint256 lastSubmission;
        uint256 totalSubmissions;
    }
    
    mapping(bytes32 => Device) public devices;
    mapping(address => bytes32[]) public ownerDevices;
    bytes32[] public allDevices;
    
    event DeviceRegistered(bytes32 indexed deviceId, address indexed owner, string publicKey);
    event DeviceDeactivated(bytes32 indexed deviceId);
    event DeviceActivated(bytes32 indexed deviceId);
    event SubmissionRecorded(bytes32 indexed deviceId, uint256 timestamp);
    
    constructor() Ownable(msg.sender) {}
    
    function registerDevice(bytes32 deviceId, string calldata publicKey) external {
        require(devices[deviceId].registrationTime == 0, "Device already registered");
        require(bytes(publicKey).length > 0, "Public key required");
        
        devices[deviceId] = Device({
            owner: msg.sender,
            publicKey: publicKey,
            registrationTime: block.timestamp,
            isActive: true,
            lastSubmission: 0,
            totalSubmissions: 0
        });
        
        ownerDevices[msg.sender].push(deviceId);
        allDevices.push(deviceId);
        
        emit DeviceRegistered(deviceId, msg.sender, publicKey);
    }
    
    function deactivateDevice(bytes32 deviceId) external {
        require(devices[deviceId].owner == msg.sender || msg.sender == owner(), "Not authorized");
        require(devices[deviceId].isActive, "Device already inactive");
        
        devices[deviceId].isActive = false;
        emit DeviceDeactivated(deviceId);
    }
    
    function activateDevice(bytes32 deviceId) external {
        require(devices[deviceId].owner == msg.sender || msg.sender == owner(), "Not authorized");
        require(!devices[deviceId].isActive, "Device already active");
        
        devices[deviceId].isActive = true;
        emit DeviceActivated(deviceId);
    }
    
    function recordSubmission(bytes32 deviceId) external onlyOwner {
        require(devices[deviceId].isActive, "Device not active");
        
        devices[deviceId].lastSubmission = block.timestamp;
        devices[deviceId].totalSubmissions++;
        
        emit SubmissionRecorded(deviceId, block.timestamp);
    }
    
    function isDeviceActive(bytes32 deviceId) external view returns (bool) {
        return devices[deviceId].isActive;
    }
    
    function getDevice(bytes32 deviceId) external view returns (Device memory) {
        return devices[deviceId];
    }
    
    function getOwnerDevices(address owner) external view returns (bytes32[] memory) {
        return ownerDevices[owner];
    }
    
    function getTotalDevices() external view returns (uint256) {
        return allDevices.length;
    }
} 