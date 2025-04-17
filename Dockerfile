FROM golang:latest

# Installe Chromium et les dépendances
RUN apt-get update && apt-get install -y \
  chromium \
  ca-certificates \
  fonts-liberation \
  libappindicator3-1 \
  libasound2 \
  libatk-bridge2.0-0 \
  libatk1.0-0 \
  libcups2 \
  libdbus-1-3 \
  libgdk-pixbuf2.0-0 \
  libnspr4 \
  libnss3 \
  libx11-xcb1 \
  libxcomposite1 \
  libxdamage1 \
  libxrandr2 \
  xdg-utils \
  && apt-get clean
 
# Création dossier app
WORKDIR /app

# Copie le code
COPY . .

# Build
RUN go mod vendor
RUN go build -o main .

# Expose le port
EXPOSE 8080

# Commande de démarrage
CMD ["/app/main"]