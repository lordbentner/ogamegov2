FROM golang:latest
 
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