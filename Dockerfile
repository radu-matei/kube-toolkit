FROM scratch

# it assumes the name you gave to the linux binary is ktkd-linux
# TODO - create a build time argument for the binary name
COPY bin/ktkd-linux .

EXPOSE 10000

ENTRYPOINT ["./ktkd-linux"]
CMD ["start", "--debug"]