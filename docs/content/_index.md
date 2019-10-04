---
title: "Pligos Documentation"
date: 2019-09-30T17:32:22+02:00
description: This article describes how we can use pligos to manage helm charts easily.
weight: 1
---

<h1>Pligos Documentation</h1>

#### Overview

Pligos allows you to navigate multiple services by making kubernetes
infrastructure configuration scalable. Without pligos, the usual
approach is to create a seperate helm chart for each service. While
this definitely can scale for a small amount of services, maintaining
`deployment, service, ingress, ...` templates for more than 5 helm
charts can be **burdensome** and error prone.

We observed that services, in it's core, often don't differ that
much. You can find a set of configurations that need to be
individually defined for each service, such as `images, routes,
mounts, ...`, so why not standardize around these configuration types,
while beeing disconnected from the underlying templates? This is why
pligos let's you define these configuration types (`image, route,
...`) and adds a schema language that allows you to compile those
configs into any form necessary.

So, what you will end up with is a small set of helm starters (in
pligos lingua franca they are called flavors) and a pligos
configuration for each service that map to these flavors.

Pligos helps you to manage your deployment in multiple environments by specifying your configurations in a single file called `pligos.yaml`. By using pligos, you ***don't need to play*** with the files within *templates* folder. ***Neither*** you have to dig into the multiple `values.yaml` files. You'll design a generic `templates` folder (called as pligos flavor) in a way to meet most of your requirements for your multiple micro-services.

Then for deploying a new micro-service, all you have to do is specify your values for local, staging, production or whatever, in pligos.yaml and you are good to go.

It makes much easier for you to manage cloud infrastructure. For example, if you want to update the hostname for your ingress in future, you don't have to find and replace it from bunch of files. You just need to modify your pligos.yaml accordingly.

![]({{< resource url="what-is-pligos.jpg" >}})