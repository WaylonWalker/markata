# Documentation is built once in CI; the cluster serves only immutable static files.
FROM python:3.12-alpine AS build
WORKDIR /src
RUN python -m pip install --no-cache-dir "Markdown==3.8.2"
COPY docs/ ./docs/
COPY docs-site/ ./docs-site/
COPY scripts/build_docs_site.py ./scripts/build_docs_site.py
ENV DOCS_BASE_PATH=/
RUN python scripts/build_docs_site.py \
 && test -s docs-site/_site/index.html \
 && test -s docs-site/_site/docs/quickstart/index.html \
 && test -s docs-site/_site/docs/guides/authoring/index.html

FROM nginxinc/nginx-unprivileged:stable-alpine
COPY docs-site/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/docs-site/_site/ /usr/share/nginx/html/
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -q -O /dev/null http://127.0.0.1:8080/ || exit 1
