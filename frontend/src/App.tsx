import "./App.css";
import Book from "./Book.model";
import { BooksForm } from "./Components/BooksForm";
import { BooksList } from "./Components/BooksList";
import { useState, useEffect } from "react";

function App() {
  const [books, setBooks] = useState<Book[]>([]);

  useEffect(() => {
    const initializeBooks = async () => {
      const fetchedBooks = await fetchBooks();
      setBooks(fetchedBooks);
    };
    initializeBooks();
  }, []);

  function addBook(book: Book): void {
    books.push(book);
    setBooks([...books]);
    console.log(books);
  }

  return (
    <div className="App">
      <main>
        <BooksForm onSubmittedBook={(book) => addBook(book)} />
        <BooksList books={books} />
      </main>
    </div>
  );
}

const apiUrl: string = "http://localhost:8081/api/books";

// Fetch books from the API
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

export default App;
