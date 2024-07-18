# Pump Fun Py

Pump Fun Bot in python

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file at the root of your project.

`WEBSOCKET_NODE_URL`
`HTTP_NODE_URL`

> **Solana public RPC endpoint:** https://api.devnet.solana.com

##Linting rules 

Golint has been deprecated as per the go documentation https://pkg.go.dev/golang.org/x/lint however i used it on the code base as most people online say its still good, we can potentially change linter however from what i used of it, it seems to work.

to install : go install golang.org/x/lint/golint@latest
Run at project path : golint ./...

there are a few alternatives that are more configurable such as golangci-lint(https://github.com/golangci/golangci-lint) this supports yaml custom config etc should we need it however i think its overkill for what we are doing.

## Run Locally

Setting up a virtual environment

> It is recommended to use a [virtual environment](https://www.geeksforgeeks.org/python-virtual-environment/).<br>
> This is a one-time setup. See below for running a virtual environment.

```bash
  python3 -m venv venv
```

Running a virtual environment

```bash
  source venv/bin/activate
```

Install dependencies

```bash
  pip install -r requirements.txt
```

Start the application

> **TO BE UPDATED**
