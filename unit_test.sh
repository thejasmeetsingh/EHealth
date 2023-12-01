# Run the testing database container
docker run --name test_db \
        -e POSTGRES_USER=test_db_user \
        -e POSTGRES_PASSWORD=1234 \
        -e POSTGRES_DB=ehealth_test_db \
        -p 5432:5432 \
        --health-cmd='pg_isready -d ehealth_test_db -U test_db_user' \
        --health-interval=10s \
        --health-timeout=5s \
        --health-retries=5 \
        -d postgis/postgis`

echo "Waiting for DB..."

while true; do
    # Check the health of the database container
    if docker inspect --format '{{json .State.Health.Status}}' test_db | grep -q "healthy"
    then
        break
    fi
    sleep 1
done

sleep 30

# Run migrations
goose -dir sql/schema postgres postgres://test_db_user:1234@localhost:5432/ehealth_test_db up

# Run test case
go test ./tests/