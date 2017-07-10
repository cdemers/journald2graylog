FROM fedora:25
RUN dnf update -qy && dnf upgrade -qy

##############################################################################
# Specific container layers:
ARG BUILD_DATE
ARG BUILD_HOST
ARG VCS_REF

# Asserting that all the necessary variables are defined for building a
# proper Docker container.
RUN if [ "$BUILD_DATE" == "" ]; then echo -e "\n\n\tERROR: BUILD_DATE MUST be defined, use docker build --build-arg BUILD_DATE=XXX\n\n"; exit 1; fi; \
    if [ "$BUILD_HOST" == "" ]; then echo "\n\n\tERROR: BUILD_HOST MUST be defined, use docker build --build-arg BUILD_HOST=XXX\n\n"; exit 1; fi; \
    if [ "$VCS_REF" == "" ]; then echo "\n\n\tERROR: VCS_REF MUST be defined, use docker build --build-arg VCS_REF=XXX\n\n"; exit 1; fi;

LABEL name="journald2graylog" \
      build-date=$BUILD_DATE \
      build-host=$BUILD_HOST \
      vcs-ref=$VCS_REF \
      maintainer="Charle Demers"

# The MAINTAINER key is deprecated, but still necessary for backward
# compatibility with Harbor and other registry browsers.
MAINTAINER Charle Demers

ADD ./journald2graylog /bin/journald2graylog
ADD ./docker/start.sh /bin/start.sh

CMD ["/bin/start.sh"]
