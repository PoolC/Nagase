FROM node:12.14-slim

WORKDIR /usr/src/app
COPY package.json .
COPY yarn.lock .

RUN apt update && apt install -y g++ make python3

RUN yarn --frozen-lockfile

COPY . .
RUN yarn build
RUN yarn --production
RUN yarn cache clean

EXPOSE 8080
CMD ["node", "dist/index.js"]
