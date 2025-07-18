pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/WeatherToken.sol";
import "../src/DeviceRegistry.sol";
import "../src/WeatherData.sol";
import "../src/RewardManager.sol";

contract DeployScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        WeatherToken weatherToken = new WeatherToken();
        console.log("WeatherToken deployed at:", address(weatherToken));

        DeviceRegistry deviceRegistry = new DeviceRegistry();
        console.log("DeviceRegistry deployed at:", address(deviceRegistry));

        WeatherData weatherData = new WeatherData();
        console.log("WeatherData deployed at:", address(weatherData));

        RewardManager rewardManager = new RewardManager(
            address(weatherToken),
            address(deviceRegistry)
        );
        console.log("RewardManager deployed at:", address(rewardManager));

        weatherToken.transfer(address(rewardManager), weatherToken.totalSupply() / 2);
        console.log("Transferred 50% of tokens to RewardManager");

        vm.stopBroadcast();

        console.log("\nDeployment Summary:");
        console.log("==================");
        console.log("WeatherToken:", address(weatherToken));
        console.log("DeviceRegistry:", address(deviceRegistry));
        console.log("WeatherData:", address(weatherData));
        console.log("RewardManager:", address(rewardManager));
        console.log("\nNext steps:");
        console.log("1. Update backend with contract addresses");
        console.log("2. Set WeatherData owner to backend service address");
        console.log("3. Set RewardManager owner to backend service address");
    }
} 