FROM scratch

COPY bin/ktkd-linux .

EXPOSE 10000

ENTRYPOINT ["./ktkd-linux"]
CMD ["start", "--debug"]