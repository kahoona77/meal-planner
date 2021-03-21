FROM alpine

WORKDIR /app
# copy binary into image
COPY web/tmpl/ /app/web/tmpl/
COPY web/assets/ /app/web/assets/
COPY meal-planner /app/

ENTRYPOINT ["/app/meal-planner"]

EXPOSE 8080