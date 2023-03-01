module.exports = {
  apps: [
    {
      script: "chainlink node start --password=.password",
      name: "chainlink",

      // @NOTE: node options for development
      env: {
        CHAINLINK_TLS_PORT: "0",
        SECURE_COOKIES: "false",
        CHAINLINK_DEV: "true",

        LINK_CONTRACT_ADDRESS: "0x326C977E6efc84E512bB9C30f76E30c160eD06FB0",
        ETH_URL: "wss://goerli.infura.io/ws/v3/KEY",
        ETH_CHAIN_ID: "5",
        JSON_CONSOLE: "true",

        DATABASE_URL: "postgresql://user:password@localhost:5432/chainlink_dev?sslmode=disable",
        FEATURE_EXTERNAL_INITIATORS: "true",
        DATABASE_TIMEOUT: "0",

        MIN_INCOMING_CONFIRMATIONS: "1",
        MIN_OUTGOING_CONFIRMATIONS: "1",
        FEATURE_FLUX_MONITOR: "true",
        SECURE_COOKIES: "false",
        LOG_LEVEL: "debug",
        ALLOW_ORIGINS: "*",
      },
    },
  ],
};
