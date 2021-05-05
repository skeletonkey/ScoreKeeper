# Score Keeper Backend API

[[_TOC_]]

## Docker

```bash
# build a local image
docker build -t scorekeeperapi .

# start up the container - needs to be run in the direcotry where your db file is
docker run -d -p 8080:8080 -v $PWD/ScoreKeeper.db:/code/ScoreKeeper.db --name my_scorekeeperapi scorekeeperapi

# container can be started and stopped as need be
docker stop my_scorekeeprapi
docker start my_scorekeeperapi

# delete the local image if you're starting over 
docker rm my_scorekeeperapi

# watch the logs of the program
docker logs --follow my_scorekeeperapi
```

### Developing

Run the build command from the above steps

```bash
docker run -it -p 8080:8080 -v $PWD/code:/code --entrypoint /bin/bash --rm scorekeeperapi
```

Once inside the running container you can edit main.go using your favorite IDE and run:

```bash
go run main.go
```

## SQLite

This is the schema of the database that backs ScoreKeeper's API.  `ScoreKeeper.db.empty` is a clean copy ready to use - I'd change the name first though.

```bash
sqlite3 ScoreKeeper.db

CREATE TABLE user (
   id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
   first_name TEXT                              NOT NULL,
   last_name  INTEGER                           NOT NULL,
   active     INTEGER                           NOT NULL DEFAULT 1,
   UNIQUE(first_name, last_name)
);

CREATE TABLE game (
    id   INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT                              NOT NULL,
    UNIQUE(name)
);

# date_played needs be in YYYY-MM-DD (HH:MM:SS.SSS) format
CREATE TABLE score (
    id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id     INTEGER                           NOT NULL,
    game_id     INTEGER                           NOT NULL,
    date_played TEXT                              NOT NULL,
    score       INTEGER                           NOT NULL DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES user(id),
    FOREIGN KEY(game_id) REFERENCES game(id)
);

```

## API

All communication is done using JSON.

Any 500 HTTP Responses are unrecoverable and you need to check the logs by attaching to the docker container.

### Game

Represent the game that you are playing

#### Get All

#### Get Specific

#### Add

##### Endpoint

```html
/api/game
```

POST Body:

```json
{
    "name": "Rummikub"
}
```

##### Response

###### 201

Score has been inserted - body will have the new entity in it

```json
{
    "id": 1,
    "name": "Rummikub"
}
```

###### 404

Something went wrong - see error message in body.

```json
{
    "error": "Error while adding a game: UNIQUE constraint failed: game.name"
}
```

### Score

#### Get All

#### Get Specific

#### Add

##### Endpoint

```html
/api/score
```

POST Body:

date_played is a string and NEEDS to be in the format: yyyy-mm-dd

```json
{
    "user_id": 1,
    "game_id": 1,
    "date_played": "2021-01-01",
    "score": -5
}
```

##### Response

###### 201

Score has been inserted - body will have the new entity in it

```json
{
    "id": 8,
    "user_id": 1,
    "game_id": 1,
    "date_played": "2021-01-01",
    "score": -5
}
```

###### 404

Something went wrong - see error message in body.

```json
{
    "error": "Date Played needs to be in the correct format: YYYY-MM-DD"
}
```

### User

#### Get All

##### Endpoint

```html
/api/users
```

##### Response

###### 200

If there are no game entered an empty list will be returned.

Example:

```json
[
    {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "active": 1
    }
]
```

#### Get Specific

##### Endpoint

```html
/api/user/123
```

##### Response

###### 200

Found user

Example:

```json
{
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "active": 1
}
```

###### 204

User not found

###### 404

Something went wrong - check error message:

```json
{
    "error": "Error while attempting to convert user id: strconv.Atoi: parsing \"abc\": invalid syntax"
}
```

#### Add

##### Endpoint

```html
/api/user
```

POST Body:

```json
{
    "first_name": "John",
    "last_name": "Doe"
}
```

##### Response

###### 201

User create - body will have the new entity in it

```json
{
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "active": 1
}
```

###### 404

Something went wrong - see error message in body.

```json
{
    "error": "Error while adding a user: UNIQUE constraint failed: user.first_name, user.last_name"
}
```
