FROM node:lts as build-deps
WORKDIR /usr/src/app
ARG REACT_APP_SIGN_IN_URL 
ARG REACT_APP_SIGN_UP_URL 
ARG REACT_APP_IDM_URL
ENV REACT_APP_SIGN_IN_URL=${REACT_APP_SIGN_IN_URL} REACT_APP_SIGN_UP_URL=${REACT_APP_SIGN_UP_URL} REACT_APP_IDM_URL=${REACT_APP_IDM_URL}
COPY package.json yarn.lock ./
RUN yarn
COPY . ./
RUN yarn build

FROM nginx:1.17-alpine
COPY --from=build-deps /usr/src/app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]