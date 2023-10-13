COMMAND=$1

cd sql/schema
goose postgres $DB_URL $COMMAND