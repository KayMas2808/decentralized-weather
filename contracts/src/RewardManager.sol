pragma solidity ^0.8.19;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./DeviceRegistry.sol";

contract RewardManager is Ownable {
    IERC20 public weatherToken;
    DeviceRegistry public deviceRegistry;
    
    uint256 public baseReward = 10 * 10**18;
    uint256 public bonusMultiplier = 150;
    uint256 public bonusThreshold = 24 hours;
    uint256 public dailyRewardLimit = 1000 * 10**18;
    
    mapping(bytes32 => uint256) public deviceEarnings;
    mapping(bytes32 => uint256) public lastRewardTime;
    mapping(uint256 => uint256) public dailyRewardsDistributed;
    
    event RewardDistributed(bytes32 indexed deviceId, address indexed recipient, uint256 amount);
    event RewardConfigUpdated(uint256 baseReward, uint256 bonusMultiplier, uint256 bonusThreshold);
    
    constructor(address _weatherToken, address _deviceRegistry) Ownable(msg.sender) {
        weatherToken = IERC20(_weatherToken);
        deviceRegistry = DeviceRegistry(_deviceRegistry);
    }
    
    function distributeReward(bytes32 deviceId) external onlyOwner {
        require(deviceRegistry.isDeviceActive(deviceId), "Device not active");
        
        DeviceRegistry.Device memory device = deviceRegistry.getDevice(deviceId);
        require(device.registrationTime > 0, "Device not registered");
        
        uint256 currentDay = block.timestamp / 1 days;
        require(dailyRewardsDistributed[currentDay] < dailyRewardLimit, "Daily reward limit reached");
        
        uint256 rewardAmount = calculateReward(deviceId);
        require(rewardAmount > 0, "No reward to distribute");
        
        deviceEarnings[deviceId] += rewardAmount;
        lastRewardTime[deviceId] = block.timestamp;
        dailyRewardsDistributed[currentDay] += rewardAmount;
        
        require(weatherToken.transfer(device.owner, rewardAmount), "Token transfer failed");
        
        emit RewardDistributed(deviceId, device.owner, rewardAmount);
    }
    
    function calculateReward(bytes32 deviceId) public view returns (uint256) {
        DeviceRegistry.Device memory device = deviceRegistry.getDevice(deviceId);
        
        if (!deviceRegistry.isDeviceActive(deviceId) || device.registrationTime == 0) {
            return 0;
        }
        
        uint256 reward = baseReward;
        
        if (device.totalSubmissions > 0 && 
            block.timestamp - device.lastSubmission <= bonusThreshold) {
            reward = (reward * bonusMultiplier) / 100;
        }
        
        if (device.totalSubmissions >= 100) {
            reward = (reward * 120) / 100;
        } else if (device.totalSubmissions >= 50) {
            reward = (reward * 110) / 100;
        }
        
        uint256 currentDay = block.timestamp / 1 days;
        uint256 remainingDailyReward = dailyRewardLimit - dailyRewardsDistributed[currentDay];
        
        if (reward > remainingDailyReward) {
            reward = remainingDailyReward;
        }
        
        return reward;
    }
    
    function updateRewardConfig(
        uint256 _baseReward,
        uint256 _bonusMultiplier,
        uint256 _bonusThreshold,
        uint256 _dailyRewardLimit
    ) external onlyOwner {
        baseReward = _baseReward;
        bonusMultiplier = _bonusMultiplier;
        bonusThreshold = _bonusThreshold;
        dailyRewardLimit = _dailyRewardLimit;
        
        emit RewardConfigUpdated(_baseReward, _bonusMultiplier, _bonusThreshold);
    }
    
    function withdrawTokens(uint256 amount) external onlyOwner {
        require(weatherToken.transfer(owner(), amount), "Token transfer failed");
    }
    
    function getDeviceEarnings(bytes32 deviceId) external view returns (uint256) {
        return deviceEarnings[deviceId];
    }
    
    function getDailyRewardsDistributed(uint256 day) external view returns (uint256) {
        return dailyRewardsDistributed[day];
    }
    
    function getCurrentDayRewards() external view returns (uint256) {
        uint256 currentDay = block.timestamp / 1 days;
        return dailyRewardsDistributed[currentDay];
    }
} 