FROM node:24.11.1-alpine AS build
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:24.11.1-alpine AS runtime
ENV NODE_ENV=production
WORKDIR /app
USER node
COPY --from=build --chown=node:node /app/.output ./.output
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
