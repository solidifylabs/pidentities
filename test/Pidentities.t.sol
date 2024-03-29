// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {Test, console} from "forge-std/Test.sol";
import {PidentitiesDeploy} from "../script/Pidentities.s.sol";
import {Pidentities, Ownable} from "../src/Pidentities.sol";

contract PidentitiesTest is Test {
    PidentitiesDeploy public deploy;
    Pidentities public nft;

    address constant OWNER = 0xFaaDaaB725709f9Ac6d5C03d9C6A6F5E3511FD70;

    function setUp() public {
        deploy = new PidentitiesDeploy();
        nft = deploy.deploy();
        vm.label(OWNER, "Arran");
    }

    function testVanityAddress() public view {
        // forgefmt: disable-next-item
        console.logBytes32(keccak256(abi.encodePacked(
            type(Pidentities).creationCode,
            abi.encode(OWNER, true)
        )));

        console.log(address(nft));

        bytes4 prefix = bytes4(bytes20(uint160(address(nft))));
        bytes4 vanity = hex"31415926";
        assertEq(prefix, vanity);
    }

    function testAirdrops() public {
        deploy.setBaseImageURI(nft);

        address[] memory deployed = deploy.airdrop(nft);

        string memory actual = "3.14159265358979323846264338327950288419716939937510582097494459230781640628";

        bytes32[7] memory hashes = [
            bytes32(0xeda528834408a2fb8fb98505502130c59a1f3b032e07da0dd8536fbfd8ffbfc3),
            0x216c212a285ec0e44fbda4a66018196caed3df48773c30fef5391df93919dc8a,
            0xd2d75e2fc5b9ea8a5ec924be1fe5233570be1e6ff11516e1051289ca8a37fd32,
            0xbd28dc882908901c41c14789265fb43b96edcfe58c9fb61737e2451e8708ce5b,
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

            // Parsing a data URI in Solidity is a pain, so this is here for inspection only and I'll "test" with a
            // testnet deployment.
            console.log(nft.tokenURI{gas: 26e6}(tokenId));

            vm.prank(OWNER);
            nft.setTokenName(tokenId, "Gary"); // no, not V
            console.log(nft.tokenURI{gas: 26e6}(tokenId));
        }
    }

    function testMintAuth(address vandal, Pidentities.Mint memory mint, string memory graffiti) public {
        vm.assume(vandal != OWNER);

        address[] memory deployed = deploy.airdrop(nft);

        Pidentities.Mint[] memory mints = new Pidentities.Mint[](1);
        mints[0] = mint;

        bytes memory err = abi.encodeWithSelector(Ownable.Unauthorized.selector);

        vm.startPrank(vandal);

        vm.expectRevert(err);
        nft.mint(mints);

        vm.expectRevert(err);
        nft.setBaseImageURI(graffiti);

        for (uint256 i = 0; i < deployed.length; ++i) {
            vm.expectRevert(err);
            nft.setTokenName(uint256(uint160(deployed[i])), graffiti);
        }

        vm.stopPrank();
    }
}
