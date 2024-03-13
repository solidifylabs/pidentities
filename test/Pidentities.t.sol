// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

import {Test, console} from "forge-std/Test.sol";
import {PidentitiesDeploy} from "../script/Pidentities.s.sol";
import {Pidentities} from "../src/Pidentities.sol";

contract PidentitiesTest is Test {
    PidentitiesDeploy public deploy;
    Pidentities public nft;

    address owner = makeAddr("owner");

    function setUp() public {
        deploy = new PidentitiesDeploy();
        nft = deploy.deploy(owner);
    }

    function testAirdrops() public {
        vm.startPrank(owner);
        address[] memory deployed = deploy.airdrop(nft);
        vm.stopPrank();

        string memory actual = "3.14159265358979323846264338327950288419716939937510582097494459230781640628";

        for (uint256 i = 0; i < deployed.length; ++i) {
            uint256 tokenId = uint256(uint160(deployed[i]));
            string memory pi = nft.piString(tokenId);
            console.log(pi);
            console.log(actual);
            console.log(nft.tokenURI(tokenId));
        }
    }
}
