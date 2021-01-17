FROM node:lts-alpine

WORKDIR /app

COPY . ./

RUN chmod +x run.sh
CMD ["./run.sh"]

RUN npm install

ENTRYPOINT node server.js