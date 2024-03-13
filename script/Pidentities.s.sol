// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {Pidentities} from "../src/Pidentities.sol";

contract PidentitiesDeploy is Script {
    function deploy(address owner) public returns (Pidentities) {
        return new Pidentities(owner, false);
    }

    function airdrop(Pidentities nft) public returns (address[] memory) {
        Pidentities.Mint[] memory airdrops = new Pidentities.Mint[](1);

        airdrops[0] = Pidentities.Mint({
            to: address(99),
            c: Pidentities.Contract({
                name: "BBP",
                salt: 0,
                deployedCode: hex"5f5f5f600160046005600660088460fc1b600281028582025f5b858b018404878c01850401888c01840401898c018304038a1c01848b019a508a891b9b508a891c996103f1116019575f5260fc60205260405ff3"
            })
        });

        return nft.mint(airdrops);
    }
}
