FROM alpine:latest

RUN apk update && apk upgrade && apk add --no-cache \
    dbus \
    libdrm \
    mesa-gbm \
    nss \
    libxkbcommon \
    gnupg \
    x11vnc \
    xvfb \
    fluxbox \
    wget \
    wmctrl \
    curl \
    nodejs \
    npm \
    chromium \
    && rm -rf /var/cache/apk/*

RUN adduser -D apps
RUN mkdir -p /home/apps && chown apps:apps /home/apps

COPY image/alpine/bootstrap.sh /home/apps/bootstrap.sh
RUN chown -R apps:apps /home/apps
RUN chmod +x /home/apps/bootstrap.sh
RUN chmod 777 /home/apps/bootstrap.sh

# Aggiorna npm alla versione più recente
RUN npm install -g npm

WORKDIR /home/apps/crawler
COPY package*.json ./
RUN npm ci --only=production \
    && npm prune --production

RUN npx playwright install-deps

COPY . /home/apps/crawler
RUN chown -R apps:apps /home/apps

CMD ["bash", "/home/apps/bootstrap.sh"]
