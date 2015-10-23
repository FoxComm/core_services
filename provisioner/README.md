# Provisioning

We use ansible playbooks to provision and configure the FoxComm
architecture hosted at Google Compute Engine.

## Requirements

### Prerequisites

    `ansible-galaxy install -r requirements.txt`

### Google compute dynamic inventory

First of all, make sure gcloud is installed & configured

Follow [this guide](https://cloud.google.com/sdk/gcloud/)

```bash
# Install the inventory script
$ go install -a ./provisioner/inventory
```

This installs the `$GOPATH/bin/inventory` binary which spits out GCE
inventory for a specific env (defaults to *staging*).

```bash
$ $GOPATH/bin/inventory --env production | jq 'keys'
[
  "_meta",
  "app",
  "bkdb",
  "db",
  "fleetctl-tunnel",
  "http-server",
  "https-server",
  "nat",
  "no-ip",
  "sidekiq",
  "spree-api",
  "spree-api-balancer"
]
```

```bash
# test the staging inventory script (replace http-server with a valid instance tag)
$ ansible -i provisioner/hosts/staging/inventory http-server -m ping
```

## Common operations

### Deploy user

```bash
$ ansible-playbook -i provisioner/hosts/staging/inventory provisioner/create-deploy-user.yml --ask-vault-pass
```

### HAProxy balancer
```bash
# Provision/Configure/Reconfigure a HAProxy balancer
$ ansible-playbook -i provisioner/hosts/staging/inventory provisioner/fc-balancer.yml --ask-vault-pass
```

### PostgreSQL
```bash
$ ansible-playbook -i provisioner/hosts/staging/inventory provisioner/postgresql.yml --ask-vault-pass
```

### Solr cluster

```bash
# Provision/Configure/Reconfigure a brand new solr-cluster
$ ansible-playbook -i provisioner/hosts/staging/inventory provisioner/solr-cluster.yml --ask-vault-pass

# Provision/Configure/Reconfigure existing solr-slaves
$ ansible-playbook -i provisioner/hosts/staging/inventory -l solr-slave provisioner/solr-cluster.yml --ask-vault-pass

# Provision/Configure/Reconfigure existing solr-master
$ ansible-playbook -i provisioner/hosts/staging/inventory -l solr-master provisioner/solr-cluster.yml --ask-vault-pass
```

### Fullstack

```bash
$ ansible-playbook -i provisioner/hosts/staging/inventory provisioner/fullstack.yml --ask-vault-pass
```

### Kibana

```bash
$ ansible-playbook -i provisioner/hosts/staging/inventory provisioner/kibana.yml --ask-vault-pass
```

### rsyslog

To deploy rsyslog configuration for kibana

```bash
$ ansible-playbook -i provisioner/hosts/staging/inventory provisioner/rsyslog.yml --ask-vault-pass
```

### Update FoxComm router config

```bash
WITH_COREOS=1 ansible-playbook -i hosts/staging/inventory foxcomm.yml --ask-vault-pass
```

## Tips

- You can use `--vault-password-file=$abs_filepath_with_pw` instead of `--ask-vault-pass`
