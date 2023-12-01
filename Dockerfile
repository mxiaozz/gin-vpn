FROM openvpn:2.5.8-sqlite

ADD app.yaml /app/app.yaml
ADD gin-vpn  /app/gin-vpn

CMD ["/app/gin-vpn"]
