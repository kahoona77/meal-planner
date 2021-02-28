FROM alpine

WORKDIR /app
# copy binary into image
COPY web/tmpl/ /app/web/tmpl/
COPY web/assets/build/ /app/web/assets/build/
COPY meal-planner /app/

ENTRYPOINT ["/app/meal-planner"]

EXPOSE 8080