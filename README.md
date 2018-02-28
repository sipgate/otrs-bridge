# OTRS Trello bridge

## Introduction

> Keep a trello board in sync with your OTRS


## Building

To build the otrs trello bridge you need a [go](https://golang.org/doc/install) runtime
and the [dep](https://golang.github.io/dep/docs/installation.html) dependency manager

After installing the prerequisites simply run
```bash
./build.sh
```

This will produce the `otrs-trello-bridge` statically linked binary

## Running

Before running, make sure you copy the `config.toml.dist` file to `config.toml` in your working directory
Then modify the settings to match your setup.

After that you should be able to run the binary:
```bash
./otrs-trello-bridge
```

The bridge defaults to port 8080 but can be overridden via the `PORT` environment variable.

For more information please refer to the [gin docs](https://gin-gonic.github.io/gin/)

## OTRS Setup

Import the webservice definitions from the [otrs-webservices](otrs-webservices) directory, make sure to use the correct URL of the bridge service

## Hacking

If you are using intellij or GoLand, simply open this directory as existing project.
If using something else, just open the directory and start hacking.

## License

```text
Copyright 2018 sipgate GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```