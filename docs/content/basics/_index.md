---
title: "Basics"
date: 2019-10-02T15:55:39+02:00
weight: 2
---

#### Installation

Pligos installation is super easy. If you have already installed the helm (to install helm follow the instructions [here](https://helm.sh/docs/using_helm/#installing-helm) ), then you can simply install pligos as `helm plugin` by following commands:

For **OSX**
```bash
helm plugin install https://github.com/real-digital/pligos/releases/latest/download/darwin_amd64_pligos.tar.gz
```

For **Linux**
```bash
helm plugin install https://github.com/real-digital/pligos/releases/latest/download/linux_amd64_pligos.tar.gz
```

#### Command

Pligos can be used by following command:

```bash
#It will render the pligos.yaml and generate the values.yaml
#containing the configuration for <context-name> context
#And it will also create templates folder with required
#files (from flavor) automatically
helm pligos <context-name> -c <path-to-pligos.yaml>

#Example: Running pligos for local context
helm pligos local -c /my-deployment

```

If you need some more help for using pligos, you can use the following command to get help about pligos. 

```bash
helm pligos --help
```