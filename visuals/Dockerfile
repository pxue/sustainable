# build environment
FROM node:12.13.0-alpine as build

WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json /app/package.json

RUN yarn --silent
RUN yarn global add webpack --silent

ARG MapboxAccessToken pk.eyJ1IjoicHh1ZSIsImEiOiJjazNqYmZxNXcwaGZsM3BvM203dndoYm5sIn0.-pf2_7SDCpWeMGgRw9UHGw
ENV MapboxAccessToken $MapboxAccessToken
ARG ApiUrl https://sustainable-api.mannr.co
ENV ApiUrl $ApiUrl

COPY . /app
RUN yarn build

# production environment
FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx/nginx.conf /etc/nginx/conf.d

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
