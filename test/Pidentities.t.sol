// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {Test, console} from "forge-std/Test.sol";
import {PidentitiesDeploy} from "../script/Pidentities.s.sol";
import {Pidentities} from "../src/Pidentities.sol";

contract PidentitiesTest is Test {
    PidentitiesDeploy public deploy;
    Pidentities public nft;

    address constant ARRAN = 0xFaaDaaB725709f9Ac6d5C03d9C6A6F5E3511FD70;

    function setUp() public {
        deploy = new PidentitiesDeploy();
        nft = deploy.deploy();
    }

    function testAirdrops() public {
        vm.startPrank(ARRAN);
        address[] memory deployed = deploy.airdrop(nft);
        vm.stopPrank();

        string memory actual = "3.14159265358979323846264338327950288419716939937510582097494459230781640628";
        console.log(actual);

        for (uint256 i = 0; i < deployed.length; ++i) {
            uint256 tokenId = uint256(uint160(deployed[i]));
            string memory pi = nft.piString(tokenId);

            if (i == 0) {
                assertEq(pi, actual);
            }
            console.log(nft.tokenURI(tokenId));
        }
    }
}
