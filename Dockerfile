FROM golang:latest

RUN apt-get update \
 && apt-get install -y openssh-server \
 && cp /etc/ssh/sshd_config /etc/ssh/sshd_config-original \
 && sed -i 's/^#\s*Port.*/Port 23234/' /etc/ssh/sshd_config \
 && sed -i 's/^#\s*PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config \
 && mkdir -p /root/.ssh \
 && chmod 700 /root/.ssh \
 && mkdir /var/run/sshd \
 && chmod 755 /var/run/sshd \
 && rm -rf /var/lib/apt/lists /var/cache/apt/archives

# Set the working directory
WORKDIR /battlesship

# Copy the Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go application (for linux)
#RUN make build
RUN GOOS=linux GOARCH=amd64 make

# Create a minimal image for the application
FROM gcr.io/distroless/base

# Copy the built application
COPY /cmd/main /usr/bin/battlesship

# Copy the systemd file
COPY /battlesship.service /etc/systemd/system/battlesship.service

# Copy the built application and the host key
COPY battlesship_hostkey /battlesship_hostkey

# need to run this every time you change the unit file
#RUN systemctl daemon-reload
#RUN systemctl start battlesship

# Set the command to run the application
CMD ["battlesship"]
