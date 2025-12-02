![Fishline logo](https://i.imgur.com/51D1iLH.png)
Fishline is a GitHub/GitLab webhook receiver that executes commands based on incoming webhook payloads. It lets you run server-side commands, such as `git pull` in response to events from GitHub or GitLab. Fishline starts a server on a specified port that listens for these webhook requests, and you can configure different command sets for different projects and branches.

### :zap: Get started
You can download the pre-compiled binaries from the Github [releases](https://github.com/hyvr-official/Fishline/releases) page and copy them to the desired location. After that you can follow the below steps in order.

#### Create a `config.json` file in the root folder where you but the Fishline binary. Here is the format of the JSON file. Fill all the commands and other properties as needed also.
You can find more details about the paramters in config file in below sections.
`````json
{
    "port": "3000",
    "logPath": "./",
    "debug": false,
    "commands": {
        "project-name": {
            "main": [
                "cd /var/www/projects"
                "git pull"
            ]
        }
    }
}
`````

#### Setup a service for running Fishline in the background. Below given are the steps to create them on a Linux distro.

#### Create a service file called `fishline_kafka.service` in the directory `/etc/systemd/system` using the following commands.

`````bash
cd /etc/systemd/system
nano fishline.service
`````

#### Copy and paste the below contents to the above created service file `fishline.service`.

`````s
[Unit]
Description=Fishline Service
After=network.target

[Service]
Type=simple
ExecStart=fishline

[Install]
WantedBy=multi-user.target
`````

> Please note that path in `ExecStart` needs to change while creating the service file.

#### Now you can start the service and also check the status of the service.

`````bash
systemctl start fishline.service
systemctl status fishline.service
`````
