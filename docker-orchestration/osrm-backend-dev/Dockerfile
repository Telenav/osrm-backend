
FROM debian:stretch-slim

# Install necessary packages for proper system state
RUN DEBIAN_FRONTEND=noninteractive apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y \
  build-essential \
  cmake \
  curl \
  wget \
  vim \
  git \
  zip \
  unzip \
  tree \
  clang-format-3.8 \
  libboost-all-dev \
  libbz2-dev \
  libtbb-dev \
  libxml2-dev \
  libzip-dev \
  lua5.2 \
  liblua5.2-dev \
  libluabind-dev \
  pkg-config \
  && rm -rf /var/lib/apt/lists/*


# install nvm and nodejs
# https://github.com/nvm-sh/nvm#install--update-script
RUN wget -qO- https://raw.githubusercontent.com/nvm-sh/nvm/v0.36.0/install.sh | bash && \
  export NVM_DIR="$HOME/.nvm" && \
  # This loads nvm
  [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  && \
  # we can only use node 10 for now, see https://github.com/Project-OSRM/osrm-backend/issues/5720
  nvm install 10 && \
  # nvm install node && \
  node --version && \
  npm --version

# install go
RUN wget --progress=dot:mega --no-check-certificate https://golang.org/dl/go1.15.3.linux-amd64.tar.gz && \
  tar -zxf go1.15.3.linux-amd64.tar.gz && \
  mv go /usr/local/ && \
  rm -f go1.15.3.linux-amd64.tar.gz && \
  mkdir -p /workspace/go/bin /workspace/go/src /workspace/go/pkg

# set env
ENV GO_HOME=/usr/local/go
ENV PATH=${GO_HOME}/bin/:$PATH
ENV GOPATH=/workspace/go


LABEL maintainer="wangyoucao577@gmail.com"

