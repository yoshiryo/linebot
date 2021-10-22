DROP DATABASE IF EXISTS linebot;
CREATE DATABASE linebot;

DROP TABLE IF EXISTS linebot.manga;

CREATE TABLE linebot.manga
(
    id              INTEGER             NOT NULL PRIMARY KEY,
    cve_id          VARCHAR(100)        NOT NULL,
    name            VARCHAR(64)         NOT NULL,
    version         VARCHAR(20)                 ,
    cve_score       VARCHAR(50)                 ,
    priority        VARCHAR(30)                                
);