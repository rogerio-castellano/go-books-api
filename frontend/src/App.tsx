import "./App.css";
import Book from "./Book.model";
import { BooksEdit } from "./Components/BooksEdit";
import { BooksForm } from "./Components/BooksForm";
import { BooksList } from "./Components/BooksList";
import { useState, useEffect } from "react";

const apiUrl: string = "http://localhost:8081/api/books";

function App() {
  const [books, setBooks] = useState<Book[]>([]);
  const [bookToUpdate, setBookForUpdate] = useState<Book | null>(null);

  useEffect(() => {
    const initializeBooks = async () => {
      const fetchedBooks = await getBooks();
      setBooks(fetchedBooks);
    };
    initializeBooks();
  }, []);

  const getBooks = async (): Promise<Book[]> => {
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
      setBooks([...books, addedBook]);
      console.log("Book added successfully!", addedBook);
    } catch (error) {
      console.error("Error attempting to add book:", error);
    }
  };

  const deleteBook = async (bookId: number): Promise<void> => {
    try {
      console.log("Deleting book with ID:", bookId);
      const response = await fetch(`${apiUrl}/${bookId}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        console.log("Error deleting book:", response);
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      setBooks(books.filter((book) => book.id !== bookId));
      console.log("Book deleted successfully!");
    } catch (error) {
      console.error("Error attempting to delete book:", error);
    }
  };

  const updateBook = async (book: Book): Promise<void> => {
    try {
      console.log("Editing book", book);
      const response = await fetch(`${apiUrl}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(book),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      setBooks(await getBooks());
      setBookForUpdate(null);
    } catch (error) {
      console.error("Error attempting to edit book:", error);
    }
  };

  return (
    <div className="App">
      <main>
        <BooksForm onSubmit={(book) => addBook(book)} />
        <BooksList books={books} onEdit={(book) => setBookForUpdate(book)} onDelete={(id) => deleteBook(id)} />
        {bookToUpdate !== null && (
          <BooksEdit
            book={bookToUpdate}
            onSubmit={(book) => updateBook(book)}
            onCancel={() => setBookForUpdate(null)}
          />
        )}
      </main>
    </div>
  );
}

export default App;
