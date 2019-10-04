---
title: "Types"
date: 2019-10-02T16:09:11+02:00
weight: 1
---

`types.yaml` is a file where you define your custom definition for your objects that you want to use within pligos.yaml. 
For example, you know that when we want to setup service for kubernetes, we need to know the service type and port.
So, the service can be defined as:
```
service:
  type: string
  port: numeric
```

Another example could be for `Image` that is used within deployment. So, you know that image has following attributes:

- registry
- repository
- tag
- pullPolicy


So, it can be defined in `types.yaml` as :
```
image:
  registry: string
  repository: string
  tag: string
  pullPolicy: string
```

#### Supported Types

Pligos supports the following basic types: 

- string
- numeric
- bool
- object

Want to see types.yaml in action? Have a look on our [getting-started](/getting-started/#types) guide

**Additionally** the language supports the following meta types which can be applied to any custom, or basic types:

- repeated
- mapped
- embedded
- embedded mapped

These types are explained with example in [schema compiler](/pligos-components/schema/#use-of-repeated-instance)