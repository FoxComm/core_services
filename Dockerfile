FROM scratch

MAINTAINER Foxcomm Team <team@foxcommerce.com>

COPY config.toml /config.toml
COPY router /router

EXPOSE 8000 8181 8182

ENTRYPOINT ["/router"]
CMD ["-tomlPath=/config", "-tomlWatch=true"]
