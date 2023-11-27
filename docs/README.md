# GenZ Documentation

> This documentation is a work in progress and will be updated as the project progresses.

Welcome to the GenZ documentation. 

## What is GenZ?

GenZ is generic golang generator. It is a tool that can be used to generate code from templates.
The templates are written in golang's template language by users. The templates can be used to generate any kind of code.

## Why GenZ?

### The problem

> [Clear is better than clever](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=14m35s)

Clear code is easier to understand and modify, which leads to fewer errors and less time spent debugging.
This is essential in industry where projects are often long-term and involve multiple developers.
It also improves communication and collaboration.

But we have observed that, sometimes, particularly on large data/domain models,
simple code tasks (like declaring constructor, setters, check model invariants or provide tooling/utility code around models)
can become repetitive and [Toil](https://sre.google/sre-book/eliminating-toil/),
impacting your velocity and maintainability as your model grows.

### Automate modeling with declaration and generation

Large projects like Kubernetes already automate utility code code by leveraging on static analysis and code generation
(e.g. [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and [controller-gen](https://book.kubebuilder.io/reference/controller-gen)).
For example, [adding Go "markers" comments on interfaces declaration](https://book.kubebuilder.io/reference/markers)
indicates business rules to the generator that will created the associated code.
Why wouldn't we create our own markers and generator ?

#### Parenthesis on reflection
Reflection is one way to achieve code introspection and access your model attributes, types and methods
but it's [hard to manipulate and operates at runtime](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=15m22s), which is very unsafe.

#### Parenthesis on tags
Tags is a very common way to add metadata to your model attributes. Thus, they can be introspected at execution time.
There are many libraries that use tags to generate code (e.g. [gorm](https://gorm.io/)).

The issues with tags are:
- they use reflection under the hood
- they only work on structs and not on other types (e.g. interfaces)

**We wanted a simple, generic and testable solution to introspect our models and generate code at compile time.**

We've built GenZ as a simple tool to help eliminate developer [Toil](https://sre.google/sre-book/eliminating-toil/) tied to your models.

## How to use GenZ?

GenZ is a CLI tool. It can be installed using the following command:

```bash
go get -u github.com/Joffref/genz
```

or 

```bash
go install github.com/Joffref/genz@latest
```

Once installed, you can use the `genz` command to generate code from templates.

Either using go generate:

```go
package main

//go:generate genz -type Human -template ./human.tmpl -output human_validator.gen.go
type Human struct {
    Firstname string
}
```

or directly:

```bash
genz -type Human -template ./human.tmpl -output human_validator.gen.go
```

## How does GenZ work?

GenZ parses the selected type and generates a struct that is injected into the template.

> NOTE: At the moment, GenZ only supports generating code from structs.

