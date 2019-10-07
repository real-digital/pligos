---
title: "Types"
date: 2019-10-02T16:09:11+02:00
weight: 2
---

You define your custom definition for your objects that you want to use within pligos.yaml. 
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

Want to see types.yaml in action? Have a look on our [getting-started]({{< ref "/getting-started/_index.md#types" >}}) guide

**Additionally** the language supports the following meta types which can be applied to any custom, or basic types:

- repeated
- mapped
- embedded
- embedded mapped

These types are explained with example in [schema compiler]({{< ref "/pligos-components/schema/_index.md#use-of-repeated-instance" >}})

#### Where these types can be defined?

You can define your context structure and your types both in `schema.yaml` just like the examples given at [schema compiler]({{< ref "/pligos-components/schema/_index.md" >}}) . Another short example could be as following:

```
#schema.yaml
route:
  port: string

container:
 route: route

context:
  container: container
```

```
# pligos.yaml
pligos:
  version: '1'
  # As we don't have separate types file in this case, so we can leave this empty
  types: []
  
contexts:
  dev:
    # This flavor folder contains schema.yaml and templates
    flavor: ./flavor 
    spec:
      container: gowebservice
    
values:
  route:
   - name: http
     port: 80

  container:
   - name: gowebservice
     route: http
```
But if you want to separate the generic types and context structure of your specific flavor (in this way, you can reuse your types across different flavors), then you can have separate `types.yaml` and `schema.yaml`

```
# types.yaml
route:
  port: string

container:
  route: route
```

```
# schema.yaml
context:
  container: container
```

```
# pligos.yaml
pligos:
  version: '1'
  # Path to types.yaml
  types: [./types.yaml]
  
contexts:
  dev:
    # This flavor folder contains schema.yaml and templates
    flavor: ./flavor
    spec:
      container: gowebservice
    
values:
  route:
   - name: http
     port: 80

  container:
   - name: gowebservice
     route: http
```