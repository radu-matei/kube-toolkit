FROM golang:1.9.2 as gateway-builder

WORKDIR /go/src/github.com/radu-matei/kube-toolkit
COPY . .

RUN ["chmod", "+x", "prerequisites.sh"]
RUN ./prerequisites.sh

RUN make gateway-linux


FROM node:8-alpine as web-builder

COPY gateway/web /app
WORKDIR /app/web

RUN npm install -g typescript
RUN npm install -g @angular/cli
RUN npm install

RUN ng build --prod --base-href=/ui --deploy-url=/ui


# starting from ubuntu right now, there's an issue starting from alpine/scratch
FROM ubuntu

COPY --from=web-builder /app/dist /app/web
COPY --from=gateway-builder /go/src/github.com/radu-matei/kube-toolkit/bin /app

EXPOSE 8080

WORKDIR /app
CMD ["./gateway-linux"]