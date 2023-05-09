FROM golang:latest

WORKDIR /opt/quepasa-source/

COPY ./quepasa-source/ /opt/quepasa-source/

RUN chmod 755 /opt/quepasa-source/docker/docker-entrypoint.sh
RUN sed -i 's/\r//' /opt/quepasa-source/docker/docker-entrypoint.sh

RUN ls -la /opt/quepasa-source/
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 31000

ENTRYPOINT ["/opt/quepasa-source/docker/docker-entrypoint.sh"]
CMD ["sh"]