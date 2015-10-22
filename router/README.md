# Router

Router based on top of `foxcomm/vulcand` server with additional middlewares.
`foxcomm/vulcand` is a fork of `mailgun/vulcand` with minimal customization. Mostly it's TOML engine instead of etcd for service management.

Documentation with architecture of vulcand avialable at http://www.vulcanproxy.com/

## Common tips, errors & troubleshooting

* To check router config (at runtime) match deployment configuration run following script:

```bash
# make sure FLEETCTL_TUNNEL is set
# for staging: 
export FLEETCTL_TUNNEL=$(scripts/fleetctl_tunnel staging)
#for production:
export FLEETCTL_TUNNEL=$(scripts/fleetctl_tunnel production)

FC_ENV=staging ./bin/router_check_endpoints.rb
# or 
FC_ENV=production ./bin/router_check_endpoints.rb
```

## Routerctl
This is vctl tool compiled with router for support of our middlewares and features.


### Install

```bash
go install ./router/routerctl
```

### Example usage

#### For first, list of commands

```bash
routerctl -h
# or w/o args.
routerctl
```

#### List of subcommands

```bash
routerctl frontend -h
# or w/o args.
routerctl frontend
```

#### Help for subcommand (with arguments description)

```bash
routerctl frontend h upsert
routerctl frontend help upsert
```

#### Set/get logging severity

```bash
routerctl log get_severity # get severity
routerctl --vulcan "http://${router_ip}:8182" log set_severity -s=INFO # set to INFO
routerctl --vulcan "http://${router_ip}:8182" log set_severity -s=WARN # set to WARN
```

#### Operations with frontends

```bash
# list of frontends
routerctl --vulcan "http://${router_ip}:8182" frontend ls

# Update or create frontend (backend must be exists)
routerctl frontend upsert -id origin_frontend -b backupServer1 --route='PathRegexp(`^/.*`)'

#  See available options for upsert
routerctl frontend upsert -h
```

#### Operations with backends

```bash
# show backend list
routerctl --vulcan "http://${router_ip}:8182" backend ls

# Add https backend without ssl verify (self-signed or invalid certs. for example), other settings are default
routerctl --vulcan "http://${router_ip}:8182" backend upsert --id backupServer1 --tlsSkipVerify
```

#### Operations with servers (endpoints)

```bash
# See list of servers of "origin_frontend" backend 
routerctl --vulcan "http://${router_ip}:8182" server ls -b=origin_frontend

# Upsert https endpoint to "origin_frontend" backend
routerctl server upsert --id "hotswappedserver11" --url https://localhost:8443 -b=origin_frontend

# Delete server from backend
routerctl server rm -b=origin_frontend --id hotswappedserver11
```

#### Middlewares

```bash
# Add cbreaker plugin
routerctl --vulcan "http://${router_ip}:8182" cbreaker upsert --id cbreaker2 --priority "0" --frontend origin_frontend --condition 'NetworkErrorRatio() > 0.5' --fallbackDuration 10s --checkPeriod 100ms --fallback '{"Type": "response", "Action": {"StatusCode": 400, "Body": "Come back later!"}}'

# Add feature_validator plugin
routerctl --vulcan "http://${router_ip}:8182" feature_validator upsert --id fv1 --priority "1" --frontend origin_frontend

# remove feature_validator plugin
routerctl feature_validator rm --id fv1 --frontend origin_frontend
```

#### See top

```bash
routerctl --vulcan "http://${router_ip}:8182" top
```

#### Using curl
Actually *routerctl* just do http api calls to API endpoint. So, you can use curl for do the same things:

```bash
curl -X GET 'http://localhost:8182/v2/log/severity' 
curl -X PUT --data "severity=INFO" 'http://localhost:8182/v2/log/severity'
```

