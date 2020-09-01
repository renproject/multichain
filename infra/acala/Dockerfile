FROM ubuntu:xenial

RUN apt-get update && apt-get install --yes --fix-missing software-properties-common curl git clang
RUN apt-get install --yes --fix-missing --no-install-recommends build-essential

# Install Rust
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Clone repository
RUN git clone https://github.com/AcalaNetwork/Acala.git

RUN mv Acala /app
WORKDIR /app

# TEMPORARY: use the branch that has a good reference to the submodules
# TODO: remove when the `master` branch of Acala is updated
RUN git fetch
RUN git checkout update-orml
RUN git pull

# Make sure submodule.recurse is set to true to make life with submodule easier.
RUN git config --global submodule.recurse true

# Build
RUN make init
RUN make build

WORKDIR /
COPY run.sh /root/
RUN chmod +x /root/run.sh

# rpc port
EXPOSE 9933
# ws port
EXPOSE 9944

ENTRYPOINT ["./root/run.sh"]
