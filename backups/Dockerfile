FROM ubuntu:14.04

MAINTAINER Foxcomm Team <team@foxcommerce.com>

RUN apt-get update
RUN apt-get install -y postgresql-client
RUN apt-get install -y mongodb-clients
ADD backups /bin/backups
EXPOSE 5004

ENTRYPOINT ["backups"]
