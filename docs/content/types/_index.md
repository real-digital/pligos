---
title: "Types"
date: 2019-10-01T12:58:37+02:00
weight: 6
---

##### Author: Abdul Jabbar
##### Created at: 19 Sep 2019

### types.yaml
types.yaml is a file where you define your custom definitions for your objects that you want to use within pligos.yaml. 
For example, you know that when we want to setup service for kubernetes, we need to know that service type and port.
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

### Supported Types

Pligos supports the following basic types: 

- string
- numeric
- bool
- object

**Additionally** the language supports the following meta types which can be applied to any custom, or basic types:

- repeated
- mapped
- embedded
- embedded mapped