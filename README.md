# CAFENETCHI

# TEST
run this for creating cafenetchi-test database (testing)

```Bash
docker exec -it cafenetchi-postgres psql \
    -U cafenetchi-user \
    -d postgres
```

then

```SQL
    CREATE DATABASE "cafenetchi-test";
```

## TODOS
- for running project on vps on docker we should use iptables-tables or iptables-legacy
- JWT validation tests
- auth middleware tests
- token expiration handling
- refresh token strategy
