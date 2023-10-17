## Backend Mid Developer Test

### Guidelines

Prerequisite:
- Docker
- Make (CLI)

In this source code has refactored and written to adopt Clean Architecture with some modified code and directory styles.
Every file has separated to its functionality. And the flow of this API application started from Router -> Handler -> Service -> Repository (DB).


#### Installation & Running
1. Copy file `config.yml.example` to `config.yml` in the directory `/configs` (if not existed) for application variables
2. Copy file `.env.example` to `.env` (if not existed) for docker environment variables
3. Fill out the DB connection config based on your Dockerized Postgres DB 
4. Run the app containers using Docker compose
```cmd
docker-compose up -d
```
5. For the default, The service will be available on `localhost:8081` on your local (address & port based in `config.yml` file)


### API Documentation
There are file `./apispec.json` and `./docs.postman_collection.json` to look up the API docs
