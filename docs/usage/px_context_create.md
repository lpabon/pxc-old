## px context create

Create a context

### Synopsis

A context is the information needed to connect to
Portworx and any other system. This information will be saved
to a file called config.yml in a directory called .px under
your home directory.

```
px context create [NAME] [flags]
```

### Examples

```

  # To create a context called k8s which will communicate to a Kubernetes
  # cluster and to one of the nodes running Portworx:
  px context create mycluster --kubeconfig=/path/to/kubeconfig --endpoint=123.456.1.10:9020

  # To create a context called mycluster which will point to one of the Portworx
  # nodes on the cluster:
  px context create mycluster --endpoint=123.456.1.10:9020
```

### Options

```
      --cafile string       Path to client CA certificate if needed
      --endpoint string     Portworx service endpoint. Ex. 127.0.0.1:9020
  -h, --help                help for create
      --kubeconfig string   Path to Kubeconfig file if any
      --secure              Use secure connection
      --token string        Token for use in this context
```

### Options inherited from parent commands

```
      --config string    config file (default is $HOME/.px/config.yml)
      --context string   Force context name for the command
      --show-labels      Show labels in the last column of the output
```

### SEE ALSO

* [px context](px_context.md)	 - Manage connections to Portworx and other systems

###### Auto generated by spf13/cobra on 12-Aug-2019
