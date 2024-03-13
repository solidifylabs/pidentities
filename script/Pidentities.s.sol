// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

import {Script, console} from "forge-std/Script.sol";
import {Pidentities} from "../src/Pidentities.sol";

contract PidentitiesDeploy is Script {
    address constant ARRAN = 0xFaaDaaB725709f9Ac6d5C03d9C6A6F5E3511FD70;

    function deploy() public returns (Pidentities) {
        vm.broadcast(ARRAN);
        return new Pidentities(ARRAN, false);
    }

    function airdrop(Pidentities nft) public returns (address[] memory) {
        Pidentities.Mint[] memory airdrops = new Pidentities.Mint[](7);

        // BBP must be first for tests. Otherwise the order doesn't matter.

        airdrops[0] = Pidentities.Mint({
            to: 0x34202f199ef058302DcceD326a0105fe2Db53E12,
            c: Pidentities.Contract({
                name: unicode"Bailey–Borwein–Plouffe%20formula",
                salt: 0,
                deployedCode: hex"5f5f5f600160046005600660088460fc1b600281028582025f5b858b018404878c01850401888c01840401898c018304038a1c01848b019a508a891b9b508a891c996103f1116019575f5260fc60205260405ff3"
            })
        });

        airdrops[1] = Pidentities.Mint({
            to: ARRAN, // TODO
            c: Pidentities.Contract({
                name: "Basel%20problem",
                salt: 0,
                deployedCode: hex"60016006607e1b815f5b81820283040183820191620710401160095792505050607e1b600160808182821b038383828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b9150508086048101851c90508086048101851c90508086048101851c90508086048101851c90508086048101851c90508086048101851c90508086048101851c90505f52607e60205260405ff3"
            })
        });

        airdrops[2] = Pidentities.Mint({
            to: 0x46622E91F95F274f4f76460B38d1F5E00905f767,
            c: Pidentities.Contract({
                name: unicode"Madhava–Leibniz%20formula",
                salt: 0,
                deployedCode: hex"620833c0600180811b8160781b805b84841b8401820401838503841b840182049003828503948310600e57821b5f52607860205260405ff3"
            })
        });

        airdrops[3] = Pidentities.Mint({
            to: 0x16cCd2a1346978e27FDCbda43569E251C4227341,
            c: Pidentities.Contract({
                name: "Limiting%20sequence",
                salt: 0,
                deployedCode: hex"5f58607a81811b805b8385851b018204820102821c6103e885850195106008578002821c8490045f52607a60205260405ff3"
            })
        });

        airdrops[4] = Pidentities.Mint({
            to: 0x2dD3A04105b25adEeb6dd356d7835e3d0B069BB2,
            c: Pidentities.Contract({
                name: "Monte%20Carlo",
                salt: 0,
                deployedCode: hex"62029cc0805f6001808160801b03818260401b035b60205f208060401c821680028183168002018311850194508386039515603a575f526014565b868560e21b045f5260e060205260405ff3"
            })
        });

        airdrops[5] = Pidentities.Mint({
            to: 0xD1edDfcc4596CC8bD0bd7495beaB9B979fc50336,
            c: Pidentities.Contract({
                name: unicode"Viète's%20formula",
                salt: 0,
                deployedCode: hex"6168a8607f6001811b6002821b5f5b8101831b600160808182821b038383828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b915050828211828186021c925084861c94508586861b039350818186021b9150508086048101851c90508086048101851c90508086048101851c90508086048101851c90508086048101851c90508086048101851c90508086048101851c905093505050509050808302600185011c925060018503945f10600e578282851b045f52607f60205260405ff3"
            })
        });

        airdrops[6] = Pidentities.Mint({
            to: 0xA85572Cd96f1643458f17340b6f0D6549Af482F5,
            c: Pidentities.Contract({
                name: "Wallis%20product",
                salt: 0,
                deployedCode: hex"607f620498806001821b5b81820260021b6001810390841b0402821c6001820391600110600a5760011b5f52607f60205260405ff3"
            })
        });

        vm.broadcast(ARRAN);
        return nft.mint(airdrops);
    }
}
