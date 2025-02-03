# lgtm-poc

The goal here is going to be to learn how to use an LGTM stack.

I should be able to self-host all these applications. 

## iterations

- [ ] set up containers for services
- [ ] create a docker compose application
- [ ] create a kubernetes cluster

## LGTM

Loki - logs
Grafana - dashboards and visualization
Tempo - traces
Mimir - metrics

- where would prometheus come in? would it replace mimir?

## development

- ~~make a data generator and a data collector~~
- ~~containerize the applications~~
- ~~implement a healthcheck for collector and generator~~
- implement shutdown method for collector
- add metrics to golang servers - depending on how tempo/mimir need these, though
- implement stack:
  - loki - added to stack. seems to be connected to grafana. need to figure out how to ship data to it
    - promtail might need to be installed in the golang servers
  - grafana - added to stack. dashboard loads and shows loki connection
  - tempo
  - mimir