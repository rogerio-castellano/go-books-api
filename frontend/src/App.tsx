import "./App.css";
import Book from "./Book.model";
import { BooksForm } from "./Components/BooksForm";
import { BooksList } from "./Components/BooksList";
import { useState, useEffect } from "react";

const apiUrl: string = "http://localhost:8081/api/books";

function App() {
  const [books, setBooks] = useState<Book[]>([]);

  useEffect(() => {
    const initializeBooks = async () => {
      const fetchedBooks = await fetchBooks();
      setBooks(fetchedBooks);
    };
    initializeBooks();
  }, []);

  const fetchBooks = async (): Promise<Book[]> => {
    try {
      const response = await fetch(apiUrl);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const books: Book[] = await response.json();
      console.log("Books fetched successfully:", books);
      books.forEach((book) => {
        console.log(`Title: ${book.title}, Author: ${book.author}, Pages: ${book.pages}`);
      });
      return books;
    } catch (error) {
      console.error("Error fetching books:", error);
      return [];
    }
  };

  const addBook = async (book: Book): Promise<void> => {
    try {
      console.log("Adding book:", JSON.stringify(book));
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(book),
      });

      if (!response.ok) {
        console.log("Error adding book:", response);
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const addedBook: Book = await response.json();
      setBooks([...books, addedBook]); // Update state with the new book
      console.log("Book added successfully!", addedBook);
      // console.log("Book added successfully!", await response.json());
    } catch (error) {
      console.error("Error attempting to add book:", error);
    }
  };

  return (
    <div className="App">
      <main>
        <BooksForm onSubmittedBook={(book) => addBook(book)} />
        <BooksList books={books} />
      </main>
    </div>
  );
}

export default App;
