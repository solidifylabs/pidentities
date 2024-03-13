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
        vm.label(ARRAN, "Arran");
    }

    function testAirdrops() public {
        vm.startPrank(ARRAN);
        address[] memory deployed = deploy.airdrop(nft);
        vm.stopPrank();

        string memory actual = "3.14159265358979323846264338327950288419716939937510582097494459230781640628";

        bytes32[7] memory hashes = [
            bytes32(0xeda528834408a2fb8fb98505502130c59a1f3b032e07da0dd8536fbfd8ffbfc3),
            0x216c212a285ec0e44fbda4a66018196caed3df48773c30fef5391df93919dc8a,
            0xd2d75e2fc5b9ea8a5ec924be1fe5233570be1e6ff11516e1051289ca8a37fd32,
            0x98e823b40930b3d81c153973537bca9c9c63de3aa694000d0e945a918627fcd6,
            0xdbd87bae04a2a5421a7f55ff8377ca84f5f9969359d68a40f195dd245350fe5a,
            0x88fdf694ef8524bde7a8a996680317164eccc5e2ad46063f6cb841c487214e08,
            0x49c74e3bf1922e804be0030e6c7f20179a5556e10b82ae2126a877cbcf265c06
        ];

        assertEq(deployed.length, hashes.length, "number deployed");

        for (uint256 i = 0; i < deployed.length; ++i) {
            assertEq(deployed[i].codehash, hashes[i], "code hash");

            uint256 tokenId = uint256(uint160(deployed[i]));
            string memory pi = nft.piString(tokenId);

            if (i == 0) {
                assertEq(pi, actual, "exact value of pi for BBP");
            }
            console.log(nft.tokenURI(tokenId));
        }
    }
}
