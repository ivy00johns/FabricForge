# Fabric Forge

Fabric CLI is an interactive command-line interface for the Fabric AI project, designed to streamline the process of selecting and executing AI patterns.

## Overview

This CLI tool provides a user-friendly interface to browse, filter, and execute Fabric AI patterns. It integrates with the [Fabric](https://github.com/danielmiessler/fabric/) project to offer an efficient way of utilizing various AI patterns for different tasks and workflows.

## Features

-   Interactive pattern selection
-   Advanced filtering options (Global Search, Tags, Categories, Directories)
-   Real-time pattern list updates as you type
-   Command preview and confirmation before execution
-   Integration with clipboard for easy input handling

## Requirements

-   Go 1.16 or later
-   Fabric AI project installed and configured
-   Node.js and npm (for Prettier)

## Installation

1. Clone this repository:

    ```
    git clone https://github.com/ivy00johns/FabricForge.git
    cd FabricForge
    ```

2. Install dependencies:

    ```
    make deps
    ```

3. Build the project:

    ```
    make build
    ```

4. Merge pattern metadata:

    ```
    make merge
    ```

    This command combines individual pattern metadata files into a single JSON file that Fabric Forge uses. Run this command whenever you add or update patterns in the Fabric project.

## Configuration

Before running the CLI, make sure to set up your `.env` file with the following variables:

1. Copy the example `.env` file:

    ```
    cp .env.example .env
    ```

2. Update the following variables, if necessary:

-   `CLI_WIDTH`: Width of the CLI interface
-   `CLI_HEIGHT`: Height of the CLI interface
-   `CLI_TITLE`: Title displayed in the CLI
-   `CLI_PLACEHOLDER`: Placeholder text for the filter input
-   `MERGED_PATTERNS_METADATA_PATH`: Path to the JSON file containing pattern metadata
-   `ALPHA_SORT`: Set to "true" to enable alphabetical sorting of patterns
-   `SORT_BY_DIR_NAME`: Set to "true" to sort patterns by directory name
-   `OUTPUT_DIR`: Directory to save command output files
-   `STREAM_RESULTS`: Set to "true" to stream results in real-time
-   `OUTPUT_RESULTS`: Set to "true" to save command output to files

## Usage

1. Run the CLI:

    ```
    make run
    ```

2. Use the arrow keys to navigate through the list of patterns.

3. Press `/` to access the filter menu. You can filter by:

    - Global Search (search across all fields)
    - Tags
    - Categories
    - Directories

4. In the Global Search mode, type to filter patterns in real-time.

5. Press Enter to select a pattern or apply a filter.

6. When a pattern is selected, you'll see a command preview. Confirm to execute the command.

7. The selected pattern will be executed using the Fabric AI project, with input taken from your clipboard.

## Development

This project uses a Makefile to streamline development tasks. Here are some useful commands:

-   `make build`: Format code and build the binary
-   `make run`: Build and run the binary
-   `make dev`: Run the application without building a binary
-   `make clean`: Remove built binary
-   `make rebuild`: Clean, build, and run
-   `make merge`: Run the JSON merge script to update pattern metadata
-   `make update_json`: Update metadata JSON
-   `make test`: Run tests
-   `make fmt`: Format Go code
-   `make format`: Format code using Prettier
-   `make deps`: Ensure dependencies are up to date
-   `make build-all`: Build for multiple platforms (Linux, macOS, Windows)

For a full list of available commands, run `make help`.

## Acknowledgments

This project is built to enhance the usability of the Fabric AI project by Daniel Miessler. Special thanks to the creators and maintainers of the libraries used in this CLI, including Bubble Tea, Lip Gloss, and others.

## Links

### Repos

-   [Fabric](https://github.com/danielmiessler/fabric)
-   [PromptUI](https://github.com/manifoldco/promptui)

### API Keys

-   [Anthropic](https://console.anthropic.com/settings/keys)
-   [OpenAI](https://platform.openai.com/api-keys)
-   [Groq](https://console.groq.com/keys)
