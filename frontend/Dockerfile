FROM node:20-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

RUN npm install sonner

COPY . .

RUN npm run build


EXPOSE 5551

CMD ["npm", "run", "dev"]

#CMD ["serve", "-s", "build", "-l", "5551"]