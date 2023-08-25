

# discord-dero-bot

This example demonstrates how to utilize the secret-discord-server's secret_pong_bot with DERO functions built in.

Such as message components, buttons, modals and pongs to user !pings; while also demonstrating ChatGPT's helpful bot, and DERO's wallet API functionality for its permissionless, private by default, cryptocurrency.

**Join [Secret Discord Server](https://discord.gg/GM5mY2t7Wg)
Discord chat channel for support.**

This project is under active development, so have fun with that.

## Build

This example assumes the following about your set up:

1. working Go environment setup
2. registered DERO wallet running in `--rpc-server` mode
3. DERO node running
4. Monero node running
5. ChatGPT API key
6. Discord API key

So first you will need to

```sh
git clone https://github.com/secretnamebasis/discord-dero-bot
```

Afterwords you are going to want to set up your `.env` file for loading enviroment dependent variables.

```sh
cd discord-dero-bot
touch .env
vim .env
```

Then copy the following template to your new `.env`

```
# DISCORD
# Bare in mind that you will want
# to make sure that your application
# has proper permissions
BOT_TOKEN=yourAppToken
GUILD_ID=yourServerID
APP_ID=yourDiscordAppID

# DERO
# The example has your node and wallet
# on the same device. Obviously, you could
# simple set up the DERO_NODE_IP & DERO_WALLET_IP
# as seperate; but I is lazy
DERO_SERVER_IP=ipaddress of you dero node and wallet
DERO_WALLET_PORT=10103
DERO_NODE_PORT=10102
USER=username
PASS=passwor

# ChatGPT
# This is a paid feature, so that kind of sucks;
# But users love it.
OPEN_AI_TOKEN=yourChatGPTToken
```

You will need to `source` your new `.env` file

```sh
source .env
```

And then you will need to

```sh
go run .```
