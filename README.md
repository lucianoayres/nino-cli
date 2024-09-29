# nino CLI

![nino-banner-github](https://github.com/user-attachments/assets/4a2c91af-e212-46f9-9f07-e68ca42d65f7)

### Run LLMs from the Command Line (Always Free)

[About](#about-nino) · [Practical Examples](#practical-examples) · [Screenshot](#screenshot) · [Features](#features) · [Ollama Dependency](#ollama-dependency) · [Requirements](#requirements) · [Installation](#installation) · [Usage](#usage) · [Using Env Vars](#using-environment-variables) · [Command-line Flags](#command-line-flags) · [Makefile](#makefile-usage) · [TODOs](#todos) · [Acknowledgements](#acknowledgements) · [License](#license) · [Contribution](#contribution)

## About Nino

Nino is a Golang command-line tool that simplifies interaction with local language models served by [Ollama](https://github.com/jmorganca/ollama). It allows you to send prompts to models, receive real-time streaming responses directly in your terminal, and configure models using straightforward command-line arguments.

Nino enhances the basic interaction provided by Ollama by displaying full model responses in the terminal and enabling you to save outputs to a file, offering a seamless experience for working with language models.

## Practical Examples

### Example 1: Analyzing Live Bitcoin Data

Using the Nino CLI to request the AI to generate an investment strategy based on live bitcoin performance data:

```bash
./nino "Analyze Bitcoin's performance data and develop a long-term investment strategy: $(btcq -all-data)"
```

This command uses Bash's native command substitution to pull Bitcoin historical data through [btcq](https://github.com/lucianoayres/btcq-cli), another CLI tool. The analysis is conducted using the Llama 3.2 model.

![nino-cli-screenshot-bitcoin-live-data](https://github.com/user-attachments/assets/3b431013-cfbb-49cb-bc0a-a1f4b7b1017f)

### Example 2: Utilizing Optional Arguments

Discover how to enhance Nino CLI's functionality with optional arguments.

![nino-cli-screenshot](https://github.com/user-attachments/assets/49cb338b-098a-4789-bd8e-e349681b0de4)

## Features

Enhance command-line workflows with Nino CLI:

-   💎 Pipe outputs to the AI for real-time analysis.
-   🐶 Pass file contents as arguments.
-   🤖 Save AI responses to text files.
-   💻 Seamlessly integrate with command-line tools.

💖 Best of all, it's completely free, forever!

## Ollama Dependency

Nino relies on the [Ollama CLI tool](https://github.com/jmorganca/ollama) to interact with local language models. Ollama must be installed and running on your machine or server for nino to function properly.

### Install Ollama

Follow the instructions on the [Ollama GitHub repository](https://github.com/jmorganca/ollama) to install and set up Ollama. Ensure that Ollama is available in your system’s `$PATH`.

### Start the Ollama Server

Once Ollama is installed, you need to start the server:

1. **Start the Ollama Server**: This command will run the Ollama server on `http://localhost:11434/api/generate` (default URL and port).

    ```bash
    ollama serve
    ```

2. **Run the Model**: To run the desired model (e.g., `llama3.2`), execute the following command in a separate terminal window:

    ```bash
    ollama run llama3.2
    ```

> **Note:** The `-model` parameter in nino **must match** the model that you run on Ollama. For example, if you start `llama3.2` in Ollama, you must pass `llama3.2` as the `-model` in nino. Otherwise, nino will not be able to communicate with the correct model.

Ollama should now be running, and nino can interact with it by sending prompts.

## Requirements

-   Go 1.23+ installed on your system
-   [Ollama](https://github.com/jmorganca/ollama) installed and running locally or on your server
-   Ensure that the Ollama server is running via `ollama serve`

## Installation

1. Clone this repository:

    ```bash
    git clone https://github.com/lucianoayres/nino.git
    cd nino
    ```

2. Build the project:

    ```bash
    make build
    ```

## Usage

After building the project and ensuring that the Ollama server is running, you can run nino with the following commands:

### Using Default Model and URL:

You can use nino with just a prompt as the only argument. By default, it will use the `llama3.2` model and connect to the default URL and port for the local Ollama server:

```bash
./nino "Who said the quote, 'A person who never made a mistake never tried anything new'?"
```

To prevent unintended line breaks or splitting of arguments in the shell, it's recommended to enclose the prompt in double quotes.

```bash
./nino "What's the typical temperature range for a CPU while gaming?"
```

### Using `-model` and `-prompt` Arguments:

```bash
./nino -model llama3.2 -prompt "Which country has the most time zones?"
```

### Using `-prompt-file` Argument:

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

```bash
./nino -model llava -prompt "What's in this image? /home/luciano/img-9872.png"
```

### Using an Alternative Model

This example uses all parameters with the `mistral` model. Ensure Ollama is running with `mistral`:

```bash
./nino -model mistral -prompt "What is the capital of Australia?" -url http://localhost:55555/api/generate -output result.txt
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

Additionally, file content can be passed as input using shell redirection:

```bash
./nino "Review the following project details and provide feedback:" < project.txt
```

## Using Environment Variables

You can set environment variables to use as defaults for the `-model` and `-url` parameters if they are not passed on the command line.

### Setting `NINO_MODEL` and `NINO_URL`

-   Set a default model:

    ```bash
    export NINO_MODEL="llama3.2"
    ```

-   Set a default URL:

    ```bash
    export NINO_URL="http://localhost:11434/api/generate"
    ```

When the environment variables are set, nino will use them as default values. You can still override them by passing `-model` and `-url` flags at runtime.

### Clearing Environment Variables

To clear an environment variable, use:

```bash
unset NINO_MODEL
unset NINO_URL
```

## Command-line Flags:

-   `-model` : The model to use (default: "llama3.2").
    -   **This must match the model that is currently running on Ollama.**
-   `-prompt` : The prompt to send to the language model (required unless `-prompt-file` is used).
-   `-prompt-file` : The path to a text file containing the prompt (optional).
    -   If both `-prompt` and `-prompt-file` are provided, `-prompt` takes precedence.
-   `-url` : The host and port where the Ollama server is running (optional).
    -   **The default `http://localhost:11434/api/generate` will be used if no URL is passed.**
-   `-output` : Specifies the filename where the model output will be saved (optional).

## Makefile

The `Makefile` in the nino project automates several key tasks like installing dependencies, building, testing, and cleaning the project.

## TODOs

-   [x] Launch v1.0
-   [ ] Create GitHub Actions Recipes
-   [ ] Add More Unit Tests
-   [ ] Add Run With Docker Method
-   [ ] Add Chat Mode Option

## Acknowledgements

I would like to thank the developers of [Ollama](https://github.com/jmorganca/ollama) for providing the core tools that nino relies on. Additionally, a big thanks to the open-source community for creating the resources that made this project possible.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contribution

Contributions are welcome! Please fork the repository and submit a pull request if you'd like to propose any changes.
