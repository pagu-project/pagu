# Installing and Running Pagu

This document provides detailed instructions on how to install and run Pagu for development purposes.

## Prerequisites

Before proceeding, ensure that your system meets the following requirements:

- **Go**: Pagu is developed using the Go programming language.
  You can find installation instructions [here](https://go.dev/doc/install).
- **Database**: Pagu uses [MySQL](https://dev.mysql.com/downloads/workbench/) as its primary database in production.
  For local development, you can use [SQLite](https://www.sqlite.org/).

## Installation Steps

Follow these steps to install and configure Pagu on your local machine:

### 1. Clone the Repository

Clone the Pagu repository to your local machine:

```bash
git clone https://github.com/pagu-project/pagu.git
cd pagu
```

### 2. Install Development Tools

Install the necessary development tools by running:

```bash
make devtools
```

### 3. Running Local Pactus Nodes (Optional)

You can run local Pactus nodes and configure them in your `config.yml`.
Refer to the [Pactus Daemon documentation](https://docs.pactus.org/get-started/pactus-daemon/).
Alternatively, Pagu can fetch information from public nodes without requiring a local node.

### 4. Wallet Requirements (Optional)

Pagu requires a Pactus wallet to manage transactions.
If you don't have a wallet, follow the [instructions to create one](https://docs.pactus.org/tutorials/pactus-wallet/#create-a-wallet).
A wallet is essential for sending transactions through Pagu.

### 5. Discord Setup (Optional)

If you plan to run Pagu on a Discord server, you will need a Guild ID and a Discord application token.
These can be obtained by following the [Discord Developer Guide](https://discord.com/developers/docs/quick-start/getting-started).

### 6. Telegram Setup (Optional)

If you plan to run Pagu on Telegram, you will need a Telegram Bot Token.

### 7. WhatsApp Setup (Optional)

If you plan to run Pagu on WhatsApp, you will need to create a WhatsApp app,
set up the webhook, and generate the access token.

## Running Pagu

Run Pagu using the Command-Line Interface (CLI) without the need for integration into Discord or Telegram.
Use the following command:

```bash
go run ./internal/cmd run --config ./internal/config/config.sample.yml
```

Now, you can interact with Pagu:

```bash
calculator reward --stake=1000 --days=1
```

Check the version of Pagu:

```bash
about
```
