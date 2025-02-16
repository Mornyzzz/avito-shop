-- migrations/001_init_tables.up.sql

-- создание таблиц

CREATE TABLE Balance (
    Username VARCHAR(255) PRIMARY KEY,
    Coins INT NOT NULL
);

CREATE TABLE Inventory (
    Username VARCHAR(255) NOT NULL,
    Item VARCHAR(255) NOT NULL,
    Quantity INT NOT NULL,
    PRIMARY KEY (Username, Item)
);

CREATE TABLE Item (
    Name VARCHAR(255) PRIMARY KEY,
    Price INT NOT NULL
);

CREATE TABLE CoinTransaction (
    ID SERIAL PRIMARY KEY,
    FromUser VARCHAR(255) NOT NULL,
    ToUser VARCHAR(255) NOT NULL,
    Amount INT NOT NULL
);

CREATE TABLE Users (
    Username VARCHAR(255) PRIMARY KEY,
    Password VARCHAR(255) NOT NULL
);