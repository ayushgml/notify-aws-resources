<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/ayushgml/notify-aws-resources">
    <!-- You can insert a logo image link here -->
  </a>

  <h1 align="center">AWS Resources Monitor & Notifier</h1>

  <p align="center">
    An efficient solution for monitoring and notifying the status of AWS resources!
    <br />
    <a href="https://github.com/ayushgml/notify-aws-resources"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/ayushgml/notify-aws-resources">View Demo</a>
    ·
    <a href="https://github.com/ayushgml/notify-aws-resources/issues">Report Bug</a>
    ·
    <a href="https://github.com/ayushgml/notify-aws-resources/issues">Request Feature</a>
  </p>
</div>


## About The Project

This project offers a comprehensive solution for monitoring AWS resources and sending notifications based on predefined metrics or status changes. It is designed to help developers and system administrators maintain optimal performance and uptime for their AWS infrastructure. The solution leverages AWS SDKs, custom scripts, and integration with notification services.

### Built With

* [Go](https://golang.org/)
* [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/)
* [Bash](https://www.gnu.org/software/bash/)

## Project Files

* `main.go`: The main Go script that implements AWS resources monitoring logic, including API calls to AWS services.
* `monitorAWSResourcesProgram`: A compiled version or additional script supporting the main monitoring functionality.
* `notify_aws_resources.sh`: Bash script responsible for sending notifications based on the monitoring results.

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

- Go installed on your system
- AWS CLI configured with your credentials
- Bash environment (Linux/MacOS or WSL for Windows)

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/ayushgml/notify-aws-resources.git
    ```

2. Navigate to the project directory
   ```sh
   cd notify-aws-resources
    ```