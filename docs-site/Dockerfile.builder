FROM python:3.12-alpine
RUN apk add --no-cache git ca-certificates \
 && python -m pip install --no-cache-dir "Markdown==3.8.2"
WORKDIR /app
COPY docs-site/builder_admin.py /app/builder_admin.py
ENV PYTHONDONTWRITEBYTECODE=1 PYTHONUNBUFFERED=1 HOME=/tmp
USER 1000:1000
EXPOSE 8080
CMD ["python", "-u", "/app/builder_admin.py"]
