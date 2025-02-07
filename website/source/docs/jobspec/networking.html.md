---
layout: "docs"
page_title: "Nomad Networking"
sidebar_current: "docs-jobspec-networking"
description: |-
  Learn how to configure networking and ports for Nomad tasks.
---

# Networking

When scheduling jobs in Nomad they are provisioned across your fleet of
machines along with other jobs and services. Because you don't know in advance
what host your job will be provisioned on, Nomad will provide your task with
network configuration when they start up.

Note that this document only applies to services that want to _listen_
on a port. Batch jobs or services that only make outbound connections do not
need to allocate ports, since they will use any available interface to make an
outbound connection.

## IP Address

Hosts in Nomad may have multiple network interfaces attached to them. This
allows you to have a higher density of services, or bind to interfaces on
different subnets (for example public vs. private subnets).

Each task will receive port allocations on a single interface. The IP is passed
to your job via the `NOMAD_IP` environment variable.

## Ports

In addition to allocating an interface, Nomad can allocate static or dynamic
ports to your task.

### Dynamic Ports

Dynamic ports are allocated in a range from `20000` to `60000`.

Most services run in your cluster should use dynamic ports. This means that the
port will be allocated dynamically by the scheduler, and your service will have
to read an environment variable (see below) to know which port to bind to at
startup.

```
task "webservice" {
    port "http" {}
    port "https" {}
}
```

### Static Ports

Static ports bind your job to a specific port on the host they're placed on.
Since multiple services cannot share a port, the port must be open in order to
place your task.

```
task "dnsservice" {
    port "dns" {
        static = 53
    }        
}
```

We recommend _only_ using static ports for [system
jobs](/docs/jobspec/schedulers.html) or specialized jobs like load balancers.

### Labels and Environment Variables

The label assigned to the port is used to identify the port in service
discovery, and used for the name of the environment variable that indicates
which port your application should bind to. For example, we've labeled this
port `http`:

```
port "http" {}
```

When the task is started, it is passed an environment variable named
`NOMAD_PORT_http` which indicates the port.

```
NOMAD_PORT_http=53423 ./start-command
```

### Mapped Ports

Some drivers (such as Docker and QEMU) allow you to map ports. A mapped port
means that your application can listen on a fixed port (it does not need to
read the environment variable) and the dynamic port will be mapped to the port
in your container or VM.

```
driver = "docker"

port "http" {}

config {
    port_map = {
        http = 8080
    }
}
```

The above example is for the Docker driver. The service is listening on port
`8080` inside the container. The driver will automatically map the dynamic port
to this service.

Please refer to the [Docker](/docs/drivers/docker.html) and [QEMU](/docs/drivers/qemu.html) drivers for additional information.
