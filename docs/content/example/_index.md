---
title: "Example"
date: 2019-10-01T13:26:58+02:00
weight: 5
---

#### Example Using Pligos

You can find a basic example of pligos at [this](https://github.com/real-digital/pligos/tree/master/examples/pligos/configs/golang-hello-world) repo 

Clone the repo
```
git clone https://github.com/real-digital/pligos.git
```

Get into the example directory
```
cd pligos/examples/pligos/configs/golang-hello-world
```

The structure of this directory is as following:
![]({{< resource url="example-directory-before-running-pligos.png" >}})

**Note:** You can see that this directory is missing the `values.yaml` file and `templates` folder. Because these two will be created by pligos. 

So, Run the pligos command. 

```
#helm pligos <context-name> -c <path-to-pligos.yaml>
helm pligos default -c .
```
**Note:** Please keep in mind that we're using here `default` because our `pligos.yaml` contains the context named default.

Now, you'll see that there is auto-generated `values.yaml` file and `templates` folder.
![]({{< resource url="example-directory-after-running-pligos.png" >}})

Then you can run your usual helm command to deploy this hello-world example.
```
helm upgrade --install hello-world .
```