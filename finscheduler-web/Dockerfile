FROM node:22-alpine AS build

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@11.0.8 --activate

COPY finscheduler-web/package.json finscheduler-web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY finscheduler-web/ .

ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
ARG VITE_FARO_COLLECT_URL=/collect
ENV VITE_FARO_COLLECT_URL=$VITE_FARO_COLLECT_URL
ARG VITE_APP_ENVIRONMENT=LocalKubernetes
ENV VITE_APP_ENVIRONMENT=$VITE_APP_ENVIRONMENT

RUN pnpm build

FROM nginx:1.27-alpine

ENV API_PROXY_PASS=http://host.docker.internal:8081/api

COPY finscheduler-web/nginx/default.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
