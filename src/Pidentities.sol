// SPDX-License-Identifier: MIT
// Copyright 2024 Arran Schlosberg (@divergencearran)
pragma solidity 0.8.24;

import {ERC721} from "solady/tokens/ERC721.sol";
import {Ownable} from "solady/auth/Ownable.sol";
import {FixedPointMathLib} from "solady/utils/FixedPointMathLib.sol";
import {LibString} from "solady/utils/LibString.sol";
import {Bytecode} from "sstore2/utils/Bytecode.sol";

interface Pidentity {
    /**
     * @dev Any function name and calldata suffice as they're ignored.
     */
    function pi() external view returns (uint256 numerator, uint256 bitsOfPrecision);
}

/**
 *
 */
contract Pidentities is ERC721, Ownable {
    using LibString for uint256;

    /**
     *
     */
    error DeploymentFailed();
    error NonPiPrefix(uint160);
    error NonExistentToken(uint256);

    /**
     *
     */
    bool private immutable requirePiPrefix;

    constructor(address owner, bool _requirePiPrefix) {
        _initializeOwner(owner);
        requirePiPrefix = _requirePiPrefix;
    }

    /**
     *
     */
    struct Contract {
        string name;
        bytes deployedCode;
        uint256 salt;
    }

    /**
     *
     */
    struct Mint {
        address to;
        Contract c;
    }

    /**
     *
     */
    function mint(Mint[] memory mints) public returns (address[] memory) {
        address[] memory deployed = new address[](mints.length);
        for (uint256 i = 0; i < mints.length; ++i) {
            deployed[i] = _mint(mints[i].to, mints[i].c);
        }
        return deployed;
    }

    /**
     *
     */
    function _mint(address to, Contract memory c) internal returns (address) {
        address addr = _deploy(c);
        _mint(to, uint256(uint160(addr)));
        return addr;
    }

    /**
     *
     */
    function _deploy(Contract memory c) internal returns (address) {
        bytes memory createCode = Bytecode.creationCodeFor(c.deployedCode);
        uint256 salt = c.salt;

        address addr;
        assembly ("memory-safe") {
            addr := create2(0, add(createCode, 0x20), mload(createCode), salt)
        }

        if (addr == address(0)) {
            revert DeploymentFailed();
        }
        if (requirePiPrefix && uint160(addr) / 1e44 != 3141) {
            revert NonPiPrefix(uint160(addr));
        }

        return addr;
    }

    /**
     *
     */
    function piFrac(uint256 tokenId)
        public
        view
        returns (uint256 numerator, uint256 denominator, uint256 bitsOfPrecision)
    {
        if (!_exists(tokenId)) {
            revert NonExistentToken(tokenId);
        }
        uint256 bits;
        (numerator, bits) = Pidentity(address(uint160(tokenId))).pi();
        return (numerator, 1 << bits, bits);
    }

    /**
     *
     */
    function piString(uint256 tokenId) public view returns (string memory) {
        (uint256 num, uint256 denom, uint256 bits) = piFrac(tokenId);
        uint256 hartleys = (1000 * bits) / 3322 - 1;

        uint256 integer = num / denom;
        uint256 decimal = FixedPointMathLib.fullMulDiv(num % denom, 10 ** hartleys, denom);

        return string.concat(integer.toString(), ".", decimal.toString());
    }

    /**
     * @inheritdoc ERC721
     */
    function name() public pure override returns (string memory) {
        return string.concat(unicode"π", "dentities");
    }

    /**
     * @inheritdoc ERC721
     */
    function symbol() public pure override returns (string memory) {
        return unicode"π";
    }

    /**
     * @inheritdoc ERC721
     */
    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        // forgefmt: disable-next-item
        return string.concat(
            "data:application/json;utf8,{",
                '"name": "",', // TODO
                unicode'"description": "π%20≈%20', piString(tokenId),'",',
                '"image": "', '', '"', // TODO
            "}"
        );
    }
}
