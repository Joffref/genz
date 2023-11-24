# GenZ Documentation

> This documentation is a work in progress and will be updated as the project progresses.

Welcome to the GenZ documentation. 

## What is GenZ?

GenZ is generic golang generator. It is a tool that can be used to generate code from templates.
The templates are written in golang's template language by users. The templates can be used to generate any kind of code.

## Why GenZ?

The rationals behind GenZ are:
- To reduce the amount of boilerplate code that needs to be written.
- To reduce the amount of code that needs to be maintained.
- To reduce the amount of code that needs to be tested.

## How to use GenZ?

GenZ is a CLI tool. It can be installed using the following command:

```bash
go get -u github.com/utkarsh-pro/genz
```

or 

```bash
go install github.com/utkarsh-pro/genz@latest
```

Once installed, you can use the `genz` command to generate code from templates.

Either using go generate:

```go
//go:generate genz -t templates -o output -d data.json
```

or directly:

```bash
genz -t templates -o output -d data.json
```

## How does GenZ work?

GenZ parses the selected type and generates a struct that is injected into the template.

> NOTE: At the moment, GenZ only supports generating code from structs.

