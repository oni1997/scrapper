FROM golang:alpine

# Install Chromium and necessary dependencies
RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    harfbuzz \
    nss \
    freetype \
    ttf-freefont

# Create a symlink for `google-chrome`
RUN ln -s /usr/bin/chromium-browser /usr/bin/google-chrome

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o app

CMD ["./app"]
