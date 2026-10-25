# syntax=docker/dockerfile:1.19
# 開発用イメージ。本番用の multi-stage（standalone 出力）は、デプロイ先が決まってから足す。
FROM node:24.22.1-alpine3.24

ENV NEXT_TELEMETRY_DISABLED=1

WORKDIR /app

# node 公式イメージには uid 1000 の node ユーザーが既にいるので、それを使う。
# 依存を先に入れてレイヤーキャッシュを効かせる。
# ここで入れた node_modules が、compose の named volume の初期値になる。
COPY --chown=node:node package.json package-lock.json ./
RUN npm ci

# .next も named volume にするため、node 所有で先に作っておく。
RUN mkdir -p /app/.next && chown node:node /app/.next

USER node
EXPOSE 3000
CMD ["npm", "run", "dev"]
