
FROM ubuntu:18.04

EXPOSE 7001
ADD docker /srv/install
RUN /srv/install/init.sh
CMD start-app
