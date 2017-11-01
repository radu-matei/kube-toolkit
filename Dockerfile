FROM alpine

COPY bin/gotham-linux .

EXPOSE 10000

ENTRYPOINT ["./gotham-linux"]
CMD ["start"]