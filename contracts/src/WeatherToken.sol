pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract WeatherToken is ERC20, Ownable {
    uint256 public constant MAX_SUPPLY = 1000000000 * 10**18;
    
    constructor() ERC20("Weather Token", "WTHR") Ownable(msg.sender) {
        _mint(msg.sender, MAX_SUPPLY);
    }
    
    function mint(address to, uint256 amount) external onlyOwner {
        require(totalSupply() + amount <= MAX_SUPPLY, "Max supply exceeded");
        _mint(to, amount);
    }
    
    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }
} 