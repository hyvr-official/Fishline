![Fishline logo](https://github.com/hyvr-official/Fishline/blob/main/art/readme-header.png?raw=true)
Fishline is a GitHub/GitLab webhook receiver that executes commands based on incoming webhook payloads. It lets you run server-side commands, such as `git pull` in response to events from GitHub or GitLab. Fishline starts a server on a specified port that listens for these webhook requests, and you can configure different command sets for different projects and branches.

### :zap: Get started
You can download the pre-compiled binaries from the Github [releases](https://github.com/hyvr-official/Fishline/releases) page and copy them to the desired location. After that you can follow the below steps in order.

<details name="install">
<summary>Install on Docker container</summary>
<br>


</details>

<details name="install">
<summary>Install on Linux distro</summary>
<br>
    
You can download the pre-compiled binaries from the Github [releases](https://github.com/hyvr-official/Fishline/releases) page and copy them to the desired location. After that you can follow the below steps in order.

#### 1. Create a `config.json` file in the root folder where you but the Fishline binary. Here is the format of the JSON file. Fill all the commands and other properties as needed also.
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

#### 2. Create a service file called `fishline_kafka.service` in the directory `/etc/systemd/system` using the following commands.

`````bash
cd /etc/systemd/system
nano fishline.service
`````

#### 3. Copy and paste the below contents to the above created service file `fishline.service`.

`````s
[Unit]
Description=Fishline Service
After=network.target

[Service]
Type=simple
ExecStart=fishline --config=/var/www/fishline/config.json

[Install]
WantedBy=multi-user.target
`````

> Please note that path in `ExecStart` needs to change while creating the service file. Also change the path for the file in `--config` to the `config.json` file that you created in the above step.

#### 4. Now you can start the service and also check the status of the service.

`````bash
systemctl start fishline.service
systemctl status fishline.service
`````
</details>

<details name="install">
<summary>Install on Windows machine</summary>
<br>

You can download the pre-compiled binaries from the Github [releases](https://github.com/hyvr-official/Fishline/releases) page and copy them to the desired location. After that you can follow the below steps in order.

#### 1. Create a `config.json` file in the root folder where you but the Fishline binary. Here is the format of the JSON file. Fill all the commands and other properties as needed also.
You can find more details about the paramters in config file in below sections.
`````json
{
    "port": "3000",
    "logPath": "./",
    "debug": false,
    "commands": {
        "project-name": {
            "main": [
                "cd C:\\Projects",
                "git pull"
            ]
        }
    }
}
`````

#### 2. Download and install [NSSM](https://nssm.cc/download)
We will be using NSSM to install and manage the service for Fishline.

#### 3. Create a service for Fishline using NSSM
* After installing [NSSM](https://nssm.cc/download)
* Run `nssm install fishline` to start the service building process.
* In the GUI add the details give below.
* `Application path` will the folder where you keep the downloaded `fishline.exe` binary
* `Aruguments` will be the `--config=[path to the config.json file]` (eg: `--config=C:\Fishline\config.json`)
* Then install and start the service
</details> 

### :rocket: Using in Gitlab and Github
Use the URL `http://[server pulic ip]:[port fishline is running]/[project name]` (eg. http://9.9.9.9:3000/project-name). You can give this URL as the webhook URL in Github and Gitlab. `[project name]` is given in the `config.json` commands array.

### :gear: Configuring Athena
Fishline can be configured using the `config.json` file created on the root the Fishline binary. Here are the details of the configuration keys and what they do in table format.

| Option | Description | Example |
| --- | --- | --- |
| `port` | Port number where Fishline server will be running | `3000` |
| `logPath` | Path to the log file where Fishline will save the logs | `/var/www/fishline` |
| `debug` | This option defines wheather debug mode is enabled or not. Debug will print the logs in terminal | `true` or `false` |
| `commands` | This will be object where we save the project names, branches and commands that should be run | Examples are give below |

#### Examples of how to configure `commands` option in `config.json` 
Here shows an example project called `chat-app` and shows the commands to run when webhook is called with `main` and `development` branch.
`````json
{
    "...": "..."
    "commands": {
        "chat-app": {
            "main": [
                "cd /var/www/chat-app-prod",
                "git pull"
            ],
            "development": [
                "cd /var/www/chat-app-dev",
                "git pull"
            ]
        }
    }
}
`````
> You can use this URL as webhook for this example `http://[pulic ip]:[port]/chat-app` (eg. http://9.9.9.9:3000/chat-app).

Here shows an example project called `blog-app` and shows the commands to run when webhook is called with `features/signup-form` branch.
`````json
{
    "...": "..."
    "commands": {
        "blog-app": {
            "features/signup-form": [
                "cd /var/www/login-app",
                "docker restart login-app"
            ]
        }
    }
}
`````
> You can use this URL as webhook for this example `http://[pulic ip]:[port]/blog-app` (eg. http://9.9.9.9:3000/blog-app).

Here shows an example with two projects called `chat-app` and `blog-app`.
`````json
{
    "...": "..."
    "commands": {
        "chat-app": {
            "main": [
                "cd /var/www/chat-app-prod",
                "git pull"
            ],
            "development": [
                "cd /var/www/chat-app-dev",
                "git pull"
            ]
        },
        "blog-app": {
            "features/signup-form": [
                "cd /var/www/login-app",
                "docker restart login-app"
            ]
        }
    }
}
`````
> You can use these URLs as webhook for this example:
> * `http://[pulic ip]:[port]/chat-app` (eg. http://9.9.9.9:3000/chat-app).
> * `http://[pulic ip]:[port]/blog-app` (eg. http://9.9.9.9:3000/blog-app).

### :hammer_and_wrench: How to build
You can build the binaries or do development of Fishline by following the below steps. Fishline is build fully on Golang. So you should install latest version of Go from [here](https://go.dev/doc/install).

* Clone that project from Github.
* Run `go mod download` command to install all mods.
* Run `build.bash` if you are building from Linux.
* Run `build.bat` if you are building from Windows.
* Build will be generated in the `./build` folder

### :page_with_curl: License
Fishline is licensed under the [MIT License](https://github.com/hyvr-official/Fishline/blob/master/LICENSE).
