# The workflow renders docs before building this image; runtime needs only nginx.
FROM nginxinc/nginx-unprivileged:stable-alpine
COPY docs-site/nginx.conf /etc/nginx/conf.d/default.conf
COPY docs-site/_site/ /usr/share/nginx/html/
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -q -O /dev/null http://127.0.0.1:8080/ || exit 1
