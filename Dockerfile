FROM golang:1.14-alpine as builder
WORKDIR /usr/src/globalserver
COPY ./globalserver ./
RUN apk add --no-cache tzdata upx
RUN upx --best globalserver -o _upx_globalserver && \
mv -f _upx_globalserver globalserver

FROM scratch
WORKDIR /opt/globalserver
COPY --from=builder /usr/src/globalserver/globalserver ./
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/
ENV TZ=Asia/Shanghai
CMD ["./globalserver"]