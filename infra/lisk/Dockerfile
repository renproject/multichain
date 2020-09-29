FROM ubuntu:xenial

RUN apt-get update && apt-get install --yes --fix-missing software-properties-common curl git clang
RUN apt-get install --yes --fix-missing --no-install-recommends build-essential
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install --yes nodejs

RUN chmod 777 /root

RUN adduser --help
RUN adduser --system --group --home /home/lisk --shell /bin/bash --uid 1100 --disabled-password lisk

# Clone repository
RUN git clone https://github.com/LiskHQ/lisk-core.git

RUN mv lisk-core /root
WORKDIR /root/lisk-core

# TEMPORARY: use the branch that has a good reference to the submodules
# TODO: remove when the `master` branch of Lisk is updated
RUN git fetch
RUN git checkout master
RUN git pull

# Make sure submodule.recurse is set to true to make life with submodule easier.
RUN git config --global submodule.recurse true

# Build
RUN npm ci
RUN npm run build

COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh
ENTRYPOINT ["/root/run.sh"]
