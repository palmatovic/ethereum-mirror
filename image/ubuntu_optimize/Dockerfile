FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y && apt-get install -y \
    libdbus-1-3 \
    libdrm2 \
    libgbm1 \
    libnss3 \
    libxkbcommon0 \
    gnupg \
    x11vnc \
    xvfb \
    fluxbox \
    wget \
    wmctrl \
    curl \
    && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list \
    #&& apt-get update && apt-get -y install google-chrome-stable \
    && curl -sL https://deb.nodesource.com/setup_14.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && apt-get remove -y curl \ 
    && apt-get remove -y wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get -y install google-chrome-stable
RUN useradd apps
RUN mkdir -p /home/apps && chown apps:apps /home/apps

COPY image/ubuntu_optimize/bootstrap_node.sh /home/apps/bootstrap.sh
RUN chown -R apps:apps /home/apps
RUN chmod +x /home/apps/bootstrap.sh
RUN chmod 777 /home/apps/bootstrap.sh

# Aggiorna npm alla versione più recente
RUN npm install -g npm

WORKDIR /home/apps/crawler
COPY package*.json ./
#RUN npm ci --production
RUN npm ci --only=production \
    && npm prune --production

RUN npx playwright install-deps
#RUN npx playwright install chrome

COPY . /home/apps/crawler
RUN chown -R apps:apps /home/apps

CMD ["bash", "/home/apps/bootstrap.sh"]
