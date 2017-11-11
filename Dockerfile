FROM scratch
ADD goSafeChatBackend /
CMD ["/goSafeChatBackend"]
EXPOSE 8090