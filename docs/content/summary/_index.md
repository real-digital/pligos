---
title: "What is Pligos?"
date: 2019-09-30T17:32:22+02:00
weight: 1
---

##### Author: Abdul Jabbar
##### Created at: 19 Sep 2019


Pligos enables scalable management of cloud infrastructures based on Helm. It helps you to manage your deployment in multiple environments by specifying your configurations in a single file called ***pligos.yaml***. By using pligos, you ***don't need to play*** with the files within *templates* folder. ***Neither*** you have to dig into the multiple *values.yaml* files. We have a generic *templates* folder designed in a way to meet almost all of your needs for deploying micro-services in *real-platform*.

All you have to do is specify your values for *local*, *staging*, *production* or whatever, in pligos.yaml and you are good to go.

It makes much easier for you to manage cloud infrastructure. For example, if you want to update the hostname for your ingress in future, you don't have to find and replace it from bunch of files. You just need to modify your pligos.yaml accordingly.

![]({{< resource url="what-is-pligos.jpg" >}})