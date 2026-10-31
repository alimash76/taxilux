# syntax=docker/dockerfile:1
FROM node:20-alpine AS builder
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
RUN apk add --no-cache curl \
    && addgroup -S app \
    && adduser -S app -G app
COPY --from=builder /src/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
RUN chown -R app:app /usr/share/nginx/html /var/cache/nginx /var/run /run /var/log/nginx \
    && chmod -R u+w /usr/share/nginx/html /var/cache/nginx /var/run
USER app
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -fsS http://localhost:8080/ || exit 1
CMD ["nginx", "-g", "daemon off;"]
