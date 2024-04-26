-- CREATE TABLE IF NOT EXISTS Users (
--     ID SERIAL PRIMARY KEY,
--     Email TEXT NOT NULL UNIQUE,
--     Username TEXT NOT NULL,
--     Phone TEXT,
--     Password_hash TEXT NOT NULL,
--     Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     LastVisitAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     Verification TEXT,
--     Photo TEXT,
--     Gender TEXT,
--     Price NUMERIC(10, 2),
--     Сommunal BOOLEAN,
--     CommunalPrice NUMERIC(10, 2),
--     Contact TEXT[]
-- );
--
-- CREATE TABLE IF NOT EXISTS Posts (
--     ID SERIAL PRIMARY KEY,
--     UserID INT,
--     Title TEXT NOT NULL,
--     Content TEXT,
--     Type TEXT,
--     Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     Visibility TEXT,
--     Likes INT,
--     LocationID INT
-- );
--
-- CREATE TABLE IF NOT EXISTS Comments (
--     ID SERIAL PRIMARY KEY,
--     PostID INT,
--     UserID INT,
--     Content TEXT,
--     Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
--
-- CREATE TABLE IF NOT EXISTS Locations (
--     ID SERIAL PRIMARY KEY,
--     Country TEXT,
--     City TEXT,
--     Street TEXT,
--     HouseNumber TEXT,
--     Floor INT,
--     ApartmentNumber TEXT,
--     Latitude DECIMAL(10, 8),
--     Longitude DECIMAL(11, 8)
-- );
-- ALTER TABLE Posts ADD CONSTRAINT fk_user FOREIGN KEY (UserID) REFERENCES Users(ID);
--
-- ALTER TABLE Posts ADD CONSTRAINT fk_location FOREIGN KEY (LocationID) REFERENCES Locations(ID);
--
-- ALTER TABLE Comments ADD CONSTRAINT fk_post FOREIGN KEY (PostID) REFERENCES Posts(ID);
--
-- ALTER TABLE Comments ADD CONSTRAINT fk_user_comment FOREIGN KEY (UserID) REFERENCES Users(ID);

    CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_visit_at TIMESTAMP,
    verification_code VARCHAR(255),
    verification_verified BOOLEAN
);
