# 2019-03-06 (cc) <paul.houghton.ywi9@statefarm.com>
#
FROM registry.sfgitlab.opr.statefarm.org/registry/sf/centos:7.5.1804

COPY agate /bin/agate

RUN mkdir -p /etc/agate /opt/agate
RUN chown -R 1000:1000 /etc/agate /opt/agate

USER 1000
EXPOSE 4464
VOLUME [ "/opt/agate" ]
WORKDIR /opt/agate
ENTRYPOINT [ "/bin/agate" ]
CMD	    [ "--config=/etc/agate/agate.yml", \
	      "--data /opt/agate" ]
