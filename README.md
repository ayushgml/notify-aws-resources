<a name="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/ayushgml/notify-aws-resources">
    <!-- You can insert a logo image link here -->
    <img src="image.webp" alt="Logo" width="50%" height="auto">
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

3. Compile the Go script (if necessary)
   ```sh
   go build -o monitorAWSResourcesProgram
   ```

4. Change the location of the directory in the `notify_aws_resources.sh` script to the project directory
   ```sh
   cd /path/to/notify-aws-resources
   ```

5. Make the `notify_aws_resources.sh` script executable
   ```sh
   chmod +x notify_aws_resources.sh
   ```

6. Set up a cron job or a scheduled task to run the notification script regularly
   ```sh
    crontab -e
    ```

7. Add the following line to the crontab file to run the script(You can change the schedule as per your requirement)
   ```sh
   0 15 * * * /path-to-your-directory/notify_aws_resources.sh
   ```

## Contribution
The project is open to contributions. Feel free to open a pull request or an issue if you find a bug or want to add a feature.


## Contact

Ayush Gupta - [@itsayush\_\_](https://twitter.com/itsayush__) - ayushgml@gmail.com