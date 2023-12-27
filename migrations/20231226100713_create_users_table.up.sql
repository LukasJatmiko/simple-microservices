CREATE TABLE IF NOT EXISTS users(
   `id` int NOT NULL AUTO_INCREMENT,
   `name` VARCHAR (50) NOT NULL,
   `email` VARCHAR (255) UNIQUE NOT NULL,
   `password` VARCHAR (255) NOT NULL,
   PRIMARY KEY(`id`),
   CONSTRAINT `unique_email` UNIQUE (`email`)
);