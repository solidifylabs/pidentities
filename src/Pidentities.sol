// SPDX-License-Identifier: MIT
// Copyright 2024 Arran Schlosberg (@divergencearran)
pragma solidity 0.8.24;

import {ERC721} from "solady/tokens/ERC721.sol";
import {Ownable} from "solady/auth/Ownable.sol";
import {FixedPointMathLib} from "solady/utils/FixedPointMathLib.sol";
import {LibString} from "solady/utils/LibString.sol";
import {SSTORE2Map, Bytecode} from "sstore2/SSTORE2Map.sol";

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
    using SSTORE2Map for bytes32;

    /// @dev Errors
    error DeploymentFailed();
    error NonPiPrefix(uint160);
    error NonExistentToken(uint256);

    /// @dev Whether to require that salts result in a pi-like prefix. Disabled for early testing.
    bool private immutable requirePiPrefix;

    /// @dev Equivalent to baseTokenURI but only for the image.
    string private _baseImageURI;

    constructor(address owner, bool _requirePiPrefix) {
        _initializeOwner(owner);
        requirePiPrefix = _requirePiPrefix;
    }

    /// @dev A `Pidentity` contract to deploy.
    struct Contract {
        string name;
        bytes deployedCode;
        uint256 salt;
    }

    /// @dev Coupling of a `Contract` and the recipient of the NFT.
    struct Mint {
        address to;
        Contract c;
    }

    /**
     * @dev Deploys all `Contracts` and mints an associated NFT to the `Mint.to` address.
     * @return Addresses of all deployed contracts; NFT token IDs are simply the addresses cast as uint256.
     */
    function mint(Mint[] memory mints) external onlyOwner returns (address[] memory) {
        address[] memory deployed = new address[](mints.length);
        for (uint256 i = 0; i < mints.length; ++i) {
            deployed[i] = _mint(mints[i].to, mints[i].c);
        }
        return deployed;
    }

    /// @dev Internal implementation of `mint()` for a single token, without authorisation check.
    function _mint(address to, Contract memory c) internal returns (address) {
        address addr = _deploy(c);

        _storageKey(addr).write(bytes(c.name));
        _mint(to, uint256(uint160(addr)));

        return addr;
    }

    /// @dev CREATE2 deploys the `Contract`, validating the pi-like prefix if required.
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
        if (requirePiPrefix && uint160(addr) / 1e40 != 31415926) {
            revert NonPiPrefix(uint160(addr));
        }

        return addr;
    }

    /// @dev Sets the URI prefix for `tokenURI()` images.
    function setBaseImageURI(string memory base) external onlyOwner {
        _baseImageURI = base;
    }

    /// @dev Sets the metadata name of a specific token.
    function setTokenName(uint256 tokenId, string memory name_) external onlyOwner {
        uint96 data = _getExtraData(tokenId) + 1;
        _setExtraData(tokenId, data);
        _storageKey(tokenId).write(bytes(name_));
    }

    /**
     * @dev Calls the NFT's associated contract, returning pi as a fraction.
     * @return numerator Numerator of the pi fraction.
     * @return denominator Denominator of the pi fraction. Always 1 << `bitsOfPrecision`, returned as a convenience.
     * @return bitsOfPrecision Number of bits used in approximating pi. This does not guarantee correctness to this precision.
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
        (numerator, bits) = approximator(tokenId).pi();
        return (numerator, 1 << bits, bits);
    }

    /// @dev Converts `piFrac(tokenId)` into a (decimal) string representation of pi.
    function piString(uint256 tokenId) public view returns (string memory) {
        (uint256 num, uint256 denom, uint256 bits) = piFrac(tokenId);
        uint256 hartleys = (1000 * bits) / 3322 - 1;

        uint256 integer = num / denom;
        uint256 decimal = FixedPointMathLib.fullMulDiv(num % denom, 10 ** hartleys, denom);

        return string.concat(integer.toString(), ".", decimal.toString());
    }

    /// @dev The pi-approximating contract associated with the NFT.
    function approximator(uint256 tokenId) public pure returns (Pidentity) {
        return Pidentity(_tokenToAddress(tokenId));
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
        return string.concat(
            "data:application/json;utf8,{",
            "\"name\": \"",
            string(_storageKey(tokenId).read()),
            "\",",
            unicode"\"description\": \"π%20≈%20",
            piString(tokenId),
            "\",",
            "\"image\": \"",
            _baseImageURI,
            LibString.toHexStringNoPrefix(uint256(_tokenToAddress(tokenId).codehash), 32),
            ".png\"}"
        );
    }

    /// @dev An SSTORE2Map storage key for the `tokenId`.
    function _storageKey(uint256 tokenId) internal view returns (bytes32) {
        return bytes32(abi.encodePacked(_tokenToAddress(tokenId), _getExtraData(tokenId)));
    }

    /// @dev An SSTORE2Map storage key for the deployed pi approximator. Equivalent to the uint256 variant.
    function _storageKey(address addr) internal view returns (bytes32) {
        return bytes32(abi.encodePacked(addr, _getExtraData(_addressToToken(addr))));
    }

    /// @dev Returns the address of the contract associated with the NFT. Simply strips the leading 96 bits.
    function _tokenToAddress(uint256 tokenId) internal pure returns (address) {
        return address(uint160(tokenId));
    }

    /// @dev Inverse of _tokenToAddress.
    function _addressToToken(address addr) internal pure returns (uint256) {
        return uint256(uint160(addr));
    }
}
