# 🐾 nino CLI

![nino-banner-github](https://github.com/user-attachments/assets/4a2c91af-e212-46f9-9f07-e68ca42d65f7)

## Run LLMs from the Command Line (Always Free)

[About](#about-nino) · [What's New?](#whats-new) · [Features](#features) · [Practical Examples](#practical-examples) · [Ollama Dependency](#ollama-dependency) · [Requirements](#requirements) · [Installation](#installation) · [Usage](#usage) · [Context History](#context-history) · [Using Env Vars](#using-environment-variables) · [Command-line Flags](#command-line-flags) · [Makefile](#makefile-usage) · [GitHub Actions](#github-actions) · [TODOs](#todos) · [Acknowledgements](#acknowledgements) · [License](#license) · [Contribution](#contribution)

## About Nino

Nino is a Golang command-line tool that simplifies interaction with local language models served by [Ollama](https://github.com/jmorganca/ollama). It allows you to send prompts to models, receive real-time streaming responses directly in your terminal, and configure models using straightforward command-line arguments.

Nino enhances the basic interaction provided by Ollama by displaying full model responses in the terminal and enabling you to save outputs to a file, offering a seamless experience for working with language models.

🎧 [Listen to the Nino CLI Audio Overview](https://notebooklm.google.com/notebook/43c94b77-3ee3-475d-a2a5-478ae3112068/audio)

### What's New?

-   🤗 **Support for Hugging Face GGUF models**: You can now run any GGUF models hosted on Hugging Face using `Nino` with Ollama's infrastructure.

### Example: Running Hugging Face GGUF Models

You can now run Hugging Face models with `nino` by simply providing the repository URL. First pull the model using Ollama:

```bash
ollama pull hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF
```

Then use the Hugging Face model with Nino:

```bash
./nino -m "hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF" -p "Explain me the core concepts of Linear Algebra"
```

You can now explore a wide range of models and tailor the performance for your needs using **nino**.
For more advanced customizations, refer to the [Hugging Face and Ollama documentation](https://huggingface.co/docs/hub/ollama).

## Features

Enhance command-line workflows with Nino CLI:

-   💎 Pipe outputs to the AI for real-time analysis.
-   🤗 Run any Hugging Face GGUF models
-   💡 Remembers your last interaction for a more conversational experience.
-   🖼️ Analyze images with multimodal models.
-   🐶 Pass file contents as prompts and save AI responses to text files.
-   💻 Seamlessly integrate with command-line tools.
-   🔒 Data never leaves your computer, ensuring privacy.

💖 Best of all, it's completely free, forever!

## Practical Examples

### Example 1: Analyzing Live Bitcoin Data

Using the Nino CLI to request the AI to generate an investment strategy based on live bitcoin performance data:

```bash
./nino "Analyze Bitcoin's performance data and develop a long-term investment strategy: $(btcq -all-data)"
```

This command uses Bash's native command substitution to pull Bitcoin historical data through [btcq](https://github.com/lucianoayres/btcq-cli), another CLI tool. The analysis is conducted using the Llama 3.2 model.

![nino-cli-screenshot-bitcoin-live-data-analysis](https://github.com/user-attachments/assets/bc96f1f6-d528-4461-b67a-2c6c83104b74)

### Example 2: Advanced Image Analysis

Nino CLI supports multimodal models like `llava`, allowing analysis beyond text by processing images:

```bash
./nino -m "llava" -p "Describe this image in a short sentence" -image ./assets/images/sample-01.png
```

![nino-cli-screenshot-advanced-image-analysis](https://github.com/user-attachments/assets/1833cabc-973f-41bb-a716-5be69ba18e7a)

Users can pass an image file to analyze visual data, recognize patterns, or extract insights from images.

### Example 3: Utilizing Optional Arguments

Discover how to enhance Nino CLI's functionality with optional arguments.

![nino-cli-screenshot-optional-arguments](https://github.com/user-attachments/assets/f774cc0f-47e5-4579-8efb-9fbb86b2d682)

### Example 4: Enabling Verbose Mode

Use the `-verbose` or `-v` flag to enable detailed logging for debugging and performance validation:

```bash
./nino -m llama3.2 -p "Explain the concept of chemical equilibrium." -verbose
```

This will display detailed logs of the request payload, response status, and operation timings, aiding in troubleshooting and performance assessment.

## Ollama Dependency

Nino relies on the [Ollama CLI tool](https://github.com/jmorganca/ollama) to interact with local language models. Ollama must be installed and running on your machine or server for nino to function properly.

### Install Ollama

Follow the instructions on the [Ollama GitHub repository](https://github.com/jmorganca/ollama) to install and set up Ollama. Ensure that Ollama is available in your system’s `$PATH`.

### Pull the Model

To pull the desired model (e.g., `llama3.2`), execute the following command:

```bash
ollama pull llama3.2
```

### Start the Ollama Server

Once Ollama is installed and the chosen model pulled, you need to start the server. This command will run the Ollama server on `http://localhost:11434/api/generate` (default URL and port):

```bash
ollama serve
```

> **Note:** The `-model` parameter in nino **must match** the model that you run on Ollama. For example, if you start `llama3.2` in Ollama, you must pass `llama3.2` as the `-model` in nino. Otherwise, nino will not be able to communicate with the correct model.

Ollama should now be running, and nino can interact with it by sending prompts.

## Requirements

-   Go 1.23+ installed on your system
-   [Ollama](https://github.com/jmorganca/ollama) installed and running locally or on your server
-   Ensure that the Ollama server is running via `ollama serve`

## Installation

### Method 1: Download the Binary from GitHub Releases

1. Download the `nino` binary from the [GitHub release page](https://github.com/lucianoayres/nino-cli/releases).

2. Add execution permission:

    ```bash
    chmod +x ./nino
    ```

3. Optionally, move the binary to your local binary folder to make it accessible from anywhere:

    ```bash
    sudo mv ./nino /usr/local/bin/
    ```

### Method 2: Clone and Build from Source

1. Clone this repository:

    ```bash
    git clone https://github.com/lucianoayres/nino.git
    cd nino
    ```

2. Build the project:

    ```bash
    make
    ```

## Usage

After building the project and ensuring that the Ollama server is running, you can run nino with the following commands:

### Using Default Model and URL

You can use nino with just a prompt as the only argument. By default, it will use the `llama3.2` model and connect to the default URL and port for the local Ollama server:

```bash
./nino "Who said the quote, 'A person who never made a mistake never tried anything new'?"
```

To prevent unintended line breaks or splitting of arguments in the shell, it's recommended to enclose the prompt in double quotes.

```bash
./nino "What's the typical temperature range for a CPU while gaming?"
```

### Using `-model` and `-prompt` Arguments

```bash
./nino -model llama3.2 -prompt "Which country has the most time zones?"
```

### Using `-prompt-file` Argument

You can pass a text file containing the prompt using the `-prompt-file` flag:

```bash
./nino -model llama3.2 -prompt-file ./prompts/question.txt
```

This will read the contents of `question.txt` and send it as the prompt to the language model.

### Using Multiline Input

Wrap the prompt text with """:

```bash
./nino """Hey!
> Explain me:
> - Neural Networks
> - How LLM Works
> """
```

### Using Multimodal Models

For models that support image inputs (like `llava`), you can include images using the `-image` or `-i` flag:

```bash
./nino -model llava -prompt "What's in this image?" -image ./assets/images/sample-01.png
```

You can pass multiple images as arguments:

```bash
./nino -model llava -prompt "Describe each image in a single word." -image ./assets/images/sample-01.png -image ./assets/images/sample-02.png
```

### Using an Alternative Model

This example uses all parameters with the `mistral` model. Ensure Ollama is running with `mistral`:

```bash
./nino -model mistral -prompt "What is the capital of Australia?" -url http://localhost:55555/api/generate -output result.txt
```

### Using JSON Format Responses

To get a JSON response, use the `-format "json"` flag and ensure your prompt explicitly requests a JSON response:

```bash
./nino -model llama3.2 -prompt "What are the top 5 most abundant chemical elements on Earth? Respond using JSON." -format "json"
```

### Using an Output File

You can optionally save the model's output to a file while still printing it to the console with the following command:

```bash
./nino -model llama3.2 -prompt "What's the Japanese word for 'Thank you'?" -output answer.txt
```

### Using Command Substitution

You can dynamically generate input for nino by using shell command substitution with the $(...) syntax. This allows the output of a shell command to be used as a prompt input:

```bash
./nino "Analyze my project directory and suggest maintenance improvements: $(ls -la)"
```

Additionally, you can pass a shell script output as input:

```bash
./nino "$(./prompts/generate_commit_message.sh)"
```

### Disabling the Loading Animation

Use the `-no-loading` flag to disable the loading animation for a cleaner output:

```bash
./nino -no-loading "Explain the concept of chemical equilibrium."
```

### Using Silent Mode

You can suppress the model output and loading animation and only save the output to a file:

```bash
./nino -model llama3.2 -prompt "What color models are available in CSS?" -silent -output answer.txt
```

### Enabling Verbose Mode

Use the `-verbose` or `-v` flag to enable detailed logging for debugging and performance validation:

```bash
./nino -model llama3.2 -prompt "Explain the concept of chemical equilibrium." -verbose
```

This will display detailed logs of the request payload, response status, and operation timings, aiding in troubleshooting and performance assessment.

## Context History

### ⚠️ Feature temporariry disabled due to performance issues

Nino automatically maintains context between requests for the same model, allowing for more coherent and conversational interactions. To disable context for a particular request, use the `no-context` flag:

```bash
./nino -model llama3.2 -no-context -prompt "What's the Linux command to list hidden files in a directory?"
```

Note: The context is limited to the last interaction, not the entire conversation history.

### Reseting Context History

If you wish to reset the context entirely, you can delete the `context.json` file for the specific model. The context files are stored in the following directory:

-   If `XDG_DATA_HOME` is set:
    -   Context files are located at `$XDG_DATA_HOME/nino/models/MODEL_NAME/context.json`
-   If `XDG_DATA_HOME` is not set:
    -   Context files are located at `~/.local/share/nino/models/MODEL_NAME/context.json`

Replace `MODEL_NAME` with the name of the model you're using (e.g., `llama3.2`). Deleting this file will remove the saved context for that model.

## Using Environment Variables

Nino allows you to configure default settings through environment variables. These include the model and URL for requests, a system prompt that automatically prefixes all user prompts, and the keep-alive duration for how long the model stays active after a request. Below are details on how to configure each of these options.

### 1. Default Model and URL

You can set default values for the model and URL used in requests, so you don't need to pass them every time via the command line.

-   **Set a default model**:

    ```bash
    export NINO_MODEL="llama3.2"
    ```

-   **Set a default URL**:

    ```bash
    export NINO_URL="http://localhost:11434/api/generate"
    ```

If these environment variables are set, Nino will use them as defaults. You can still override these defaults by passing the `-model` and `-url` flags at runtime.

### 2. Keep-Alive Duration

The `NINO_KEEP_ALIVE` variable controls how long the model stays active after a request before shutting down. By default, this value is **60 minutes** (`60m`).

-   **Set a custom keep-alive duration**:

    ```bash
    export NINO_KEEP_ALIVE="90m"
    ```

In this example, the model will remain active for 90 minutes after a request.

### 3. System Prompt

You can set a default system prompt to be automatically added to every user prompt. This is useful for ensuring consistent instructions across all interactions.

-   **Set a default system prompt**:

    ```bash
    export NINO_SYSTEM_PROMPT="Do not use markdown in your answer:"
    ```

Once set, this system prompt cannot be overridden in individual prompts. You must clear it to change it.

### 4. Clearing Environment Variables

To clear any of the environment variables mentioned above, use:

```bash
unset NINO_MODEL
unset NINO_URL
unset NINO_KEEP_ALIVE
unset NINO_SYSTEM_PROMPT
```

## Command-line Flags

-   `-model` or `-m` : The model to use (default: "llama3.2").
    -   Note: This must match the model that is currently running on Ollama.
-   `-prompt` or `-p` : The prompt to send to the language model (required unless `-prompt-file` is used).
-   `-prompt-file` or `-pf` : The path to a text file containing the prompt (optional).
    -   Note: If both `-prompt` and `-prompt-file` are provided, `-prompt` takes precedence.
-   `-image` or `-i`: Path to local image file to include in the request (optional).
    -   Note: This flag is compatible only with multimodal models that support image inputs. It can be used multiple times to include multiple images in a single request.
-   `-url` or `-u` : The host and port where the Ollama server is running (optional).
    -   Note: The default `http://localhost:11434/api/generate` will be used if no URL is passed.
-   `-format` or `-f` : Specifies the format of the response from the model.
    -   Note: Currently, the only supported value is `json`. This flag also requires that your prompt explicitly instructs the model to respond in JSON format.
-   `-output` or `-o`: Specifies the filename where the model output will be saved (optional).
-   `-no-loading` or `-nl` : Disable the loading animation (optional).
-   `-no-stream` or `-ns`: Disables streaming mode, displaying the entire response at once instead of progressively showing it on the screen.
    -   Note: This may result in a longer wait time before the response is displayed.
-   `-no-context` or `-nc` : Disable the context from the previous request (optional).
    -   Note: Previous context won't be used for this response, but the new context will be cached.
-   `-silent` or `-s` : Suppresses model output and loading animation (optional).
    -   Note: Requires `-output` flag.
-   `-verbose` or `-v` : Enables verbose logging for debugging and performance validation (optional).
    -   Note: When enabled, detailed logs including request payloads and operation timings are displayed to aid in troubleshooting and performance assessment.

## Makefile

The `Makefile` in the nino project automates several key tasks like installing dependencies, building, testing, and cleaning the project.

## GitHub Actions

[Sample workflows](https://github.com/lucianoayres/nino-cli/tree/main/.github/workflows) using Nino CLI for AI-Generated content integration:

-   [Generate Daily Quote](https://github.com/lucianoayres/nino-cli/actions/workflows/generate-daily-quote.yml): Generate a quote, export it to a file and save it as artifact on GitHub daily at midnight (00:00 UTC).

-   [Save Output to File](https://github.com/lucianoayres/nino-cli/actions/workflows/save-output-to-file.yml): Dispatch the workflow with selected inputs by the user, export the model response to a file, commit it, and then push the changes to the remote repository.

### Triggering the Workflow via REST API

You can trigger the GitHub Actions workflow with a REST API call using the following example. Be sure to replace placeholders with your actual `GitHub Token`, `Username`, `Repository name`, and `Workflow filename`. Example:

```bash
curl -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer YOUR_GITHUB_TOKEN" \
  https://api.github.com/repos/lucianoayres/nino-cli/actions/workflows/save-output-to-file.yml/dispatches \
  -d '{"ref":"main", "inputs": {"model": "llama3.2", "prompt": "Explain me the BM25 ranking algorithm", "output_filename": "result.txt"}}'
```

#### Steps to Generate a Personal Access Token (PAT) on GitHub

To trigger workflows via the API, you’ll need a GitHub personal access token. Follow these steps to generate one:

1. Click on your profile photo in GitHub, go to **Settings**, and navigate to **Developer Settings**.
2. Under **Personal Access Tokens**, click [Generate a new token](https://github.com/settings/tokens?type=beta).
3. Set the **Expiration** time and select a **Repository** as the scope.
4. In **Repository Permissions**, ensure `Actions` and `Workflows` have `Read & Write` access.
5. Generate and copy the token for use in your API call.

## TODOs

-   [x] Launch v1.0
-   [x] Create GitHub Actions Recipes
-   [x] Add Multimodal Model Support
-   [x] Add JSON format Argument
-   [x] Add Stream Mode Argument
-   [ ] Add Context Support
-   [ ] Fix the Performance Issue with Context Data Loading
-   [ ] Increase Test Coverage
-   [ ] Add Custom Modelfiles
-   [ ] Add Run With Docker Method

### 🦖 Create Custom AI Models with Modelzilla

Looking to build your own AI models? Use [**Modelzilla**](https://github.com/lucianoayres/modelzilla) 🦖 to effortlessly generate customized Modelfiles.

## Acknowledgements

I would like to thank the developers of [Ollama](https://github.com/jmorganca/ollama) for providing the core tools that nino relies on. Additionally, a big thanks to the open-source community for creating the resources that made this project possible.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contribution

Contributions are welcome! Please fork the repository and submit a pull request if you'd like to propose any changes.
