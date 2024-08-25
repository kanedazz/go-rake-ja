# Use the official Golang image as the base image
FROM golang:latest

# change directory to /app
WORKDIR /app

# Set environment variables to avoid prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# Install required packages and MeCab dependencies
RUN apt-get update && apt-get install -y \
    mecab \
    libmecab-dev \
    mecab-ipadic-utf8 \
    git \
    make \
    curl \
    xz-utils \
    file \
    sudo \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Optional: Install additional dictionaries (e.g., mecab-ipadic-utf8)
RUN apt-get update && apt-get install -y \
    mecab-ipadic-utf8 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Optional: Install MeCab Python bindings if needed
# RUN pip install mecab-python3

# Verify the installation
RUN mecab --version

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the source code into the container
# COPY . .


# # Set the environment variables for compiling and linking MeCab with Go
# RUN export CGO_FLAGS="`mecab-config --inc-dir`" && \
#     export CGO_LDFLAGS="`mecab-config --libs`" && \
#     go build -mod=vendor -o main .

# Set the command to run the executable when the container starts
# CMD ["./main"]
CMD sh -c 'export CGO_FLAGS="`mecab-config --inc-dir`" && \
           export CGO_LDFLAGS="`mecab-config --libs`" && \
           go test -mod=vendor ./...'

