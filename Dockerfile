# Utilise l'image officielle Go
FROM golang:1.24

# Création dossier app
WORKDIR /app

# Copie le code
COPY . .

# Build
RUN go build -o main .

# Expose le port
EXPOSE 8080

# Commande de démarrage
CMD ["/app/main"]
