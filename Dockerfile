# 只使用二进制可执行文件
# Only use binary executable file.

FROM ubuntu:16.04
MAINTAINER AceDarkknight

WORKDIR /app

COPY . /app

EXPOSE 8090/tcp

ENTRYPOINT ./GoProxyCollector
