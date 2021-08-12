/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: "0.7.3",
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {
      chainId: 5,
      blockGasLimit: 15000000,
      hardfork: "london",
      mining: {
        auto: false,
        interval: 1000
      },
      accounts: {
        mnemonic: "clutch captain shoe salt awake harvest setup primary inmate ugly among become",
        count: 105
      }
    }
  }
}
