DROP TABLE IF EXISTS PERSON CASCADE;
CREATE TABLE Person
(
    Person_ID     BIGSERIAL NOT NULL PRIMARY KEY,
    name_en       varchar(100),
    name_rus      varchar(100),
    picture_url   varchar(100),
    career        varchar(100),
    height        integer,
    age           integer,
    birthday      date,
    death         date,
    birthplace    varchar(100),
    death_place   varchar(100),
    gender        varchar(10),
    family_status varchar(100),
    film_number   integer
);

DROP TABLE IF EXISTS FILM CASCADE;
CREATE TABLE FILM
(
    Film_ID        BIGSERIAL NOT NULL PRIMARY KEY,
    title          varchar(100),
    title_original varchar(100),
    rating         double precision,
    description    varchar(2000),
    poster_url     varchar(100),
    trailer_url    varchar(100),
    total_revenue  money,
    release_year   integer,
    duration       integer,
    screenwriter   integer default -1,
    director       integer default -1,
    content_type   varchar(30),
    CONSTRAINT film_screenwriter FOREIGN KEY (screenwriter) REFERENCES Person (Person_ID) ON DELETE SET DEFAULT,
    CONSTRAINT film_director FOREIGN KEY (director) REFERENCES Person (Person_ID) ON DELETE SET DEFAULT
);

DROP TABLE IF EXISTS Profile CASCADE;
CREATE TABLE Profile
(
    User_ID       BIGSERIAL NOT NULL PRIMARY KEY,
    first_name    varchar(100),
    surname       varchar(100),
    email         varchar(100),
    password      varchar(100),
    picture_url   varchar(100),
    gender        varchar(10),
    register_date timestamp
);

DROP TABLE IF EXISTS Review CASCADE;
CREATE TABLE Review
(
    Review_ID   BIGSERIAL NOT NULL PRIMARY KEY,
    Film_ID     BIGINT,
    Author_ID   BIGINT,
    review_text varchar(2000),
    type        integer,
    stars       double precision,
    review_date timestamp,
    CONSTRAINT to_film FOREIGN KEY (Film_ID) REFERENCES FILM (Film_ID) ON DELETE CASCADE,
    CONSTRAINT to_user FOREIGN KEY (Author_ID) REFERENCES Profile (User_ID) ON DELETE CASCADE
);

DROP TABLE IF EXISTS Collection CASCADE;
CREATE TABLE Collection
(
    Collection_ID   BIGSERIAL NOT NULL PRIMARY KEY,
    Author_ID       BIGINT,
    collection_name varchar(60),
    description     varchar(200),
    creation_time   timestamp,
    picture_url     varchar(100),
    CONSTRAINT to_author FOREIGN KEY (Author_ID) REFERENCES Profile (User_ID) ON DELETE CASCADE
);

DROP TABLE IF EXISTS Country CASCADE;
CREATE TABLE Country
(
    Country_ID   SERIAL NOT NULL PRIMARY KEY,
    Country_name varchar(50)
);

DROP TABLE IF EXISTS Genre CASCADE;
CREATE TABLE Genre
(
    Genre_ID   SERIAL NOT NULL PRIMARY KEY,
    Genre_name varchar(50)
);

DROP TABLE IF EXISTS Subscription CASCADE;
CREATE TABLE Subscription
(
    Subscriber_ID bigint,
    Author_ID     bigint,
    CONSTRAINT to_sub FOREIGN KEY (Subscriber_ID) REFERENCES Profile (User_ID) ON DELETE CASCADE,
    CONSTRAINT to_auth FOREIGN KEY (Author_ID) REFERENCES Profile (User_ID) ON DELETE CASCADE,
    CONSTRAINT subscription_ID PRIMARY KEY (Subscriber_ID, Author_ID)
);

DROP TABLE IF EXISTS FilmCast CASCADE;
CREATE TABLE FilmCast
(
    Film_ID   bigint NOT NULL,
    Person_ID bigint NOT NULL,
    CONSTRAINT to_film FOREIGN KEY (Film_ID) REFERENCES FILM (Film_ID) ON DELETE CASCADE,
    CONSTRAINT to_person FOREIGN KEY (Person_ID) REFERENCES Person (Person_ID) ON DELETE CASCADE,
    CONSTRAINT Cast_ID PRIMARY KEY (Film_ID, Person_ID)
);

DROP TABLE IF EXISTS CountryProduction CASCADE;
CREATE TABLE CountryProduction
(
    Film_ID    bigint NOT NULL,
    Country_ID bigint NOT NULL,
    CONSTRAINT to_film FOREIGN KEY (Film_ID) REFERENCES FILM (Film_ID) ON DELETE CASCADE,
    CONSTRAINT to_country FOREIGN KEY (Country_ID) REFERENCES Country (Country_ID) ON DELETE CASCADE,
    CONSTRAINT Production_ID PRIMARY KEY (Film_ID, Country_ID)
);

DROP TABLE IF EXISTS FilmGenres CASCADE;
CREATE TABLE FilmGenres
(
    Film_ID  bigint NOT NULL,
    Genre_ID bigint NOT NULL,
    CONSTRAINT to_film FOREIGN KEY (Film_ID) REFERENCES FILM (Film_ID) ON DELETE CASCADE,
    CONSTRAINT to_country FOREIGN KEY (Genre_ID) REFERENCES Genre (Genre_ID) ON DELETE CASCADE,
    CONSTRAINT FilmGenre_ID PRIMARY KEY (Film_ID, Genre_ID)
);

DROP TABLE IF EXISTS Bookmark CASCADE;
CREATE TABLE Bookmark
(
    Film_ID bigint NOT NULL,
    User_ID bigint NOT NULL,
    CONSTRAINT to_film FOREIGN KEY (Film_ID) REFERENCES FILM (Film_ID) ON DELETE CASCADE,
    CONSTRAINT to_user FOREIGN KEY (User_ID) REFERENCES Profile (User_ID) ON DELETE CASCADE,
    CONSTRAINT Bookmark_ID PRIMARY KEY (Film_ID, User_ID)
);

DROP TABLE IF EXISTS CollectionConnection CASCADE;
CREATE TABLE CollectionConnection
(
    Film_ID       bigint NOT NULL,
    Collection_ID bigint NOT NULL,
    CONSTRAINT to_film FOREIGN KEY (Film_ID) REFERENCES FILM (Film_ID) ON DELETE CASCADE,
    CONSTRAINT to_coll FOREIGN KEY (Collection_ID) REFERENCES Collection (Collection_ID) ON DELETE CASCADE,
    CONSTRAINT Connect_ID PRIMARY KEY (Film_ID, Collection_ID)
);

DROP TABLE IF EXISTS Recommended CASCADE;
CREATE TABLE Recommended
(
    Film_ID bigint NOT NULL,
    Recommended_ID bigint NOT NULL,
    CONSTRAINT to_film FOREIGN KEY (Film_ID) REFERENCES FILM(Film_ID) ON DELETE CASCADE,
    CONSTRAINT to_filmRec FOREIGN KEY (Film_ID) REFERENCES  FILM(Film_ID) ON DELETE CASCADE,
    CONSTRAINT Recommended_ID PRIMARY KEY (Film_ID, Recommended_ID)
);