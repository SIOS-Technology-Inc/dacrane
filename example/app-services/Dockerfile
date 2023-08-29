FROM node:16.20.1-alpine3.18

COPY ./index.js ./index.js
COPY ./package.json ./package.json
COPY ./package-lock.json ./package-lock.json

RUN npm ci

CMD [ "npm", "start" ]
