CREATE TABLE books (
    id SERIAL PRIMARY KEY,        -- Auto-incrementing ID, equivalent to "int"
    title VARCHAR(255) NOT NULL,  -- Title of the book, string with a max length of 255
    author VARCHAR(255) NOT NULL, -- Author of the book, string with a max length of 255
    pages INT NOT NULL            -- Number of pages, integer
);
