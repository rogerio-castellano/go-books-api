import "./App.css";
import Book from "./Book.model";
import { BooksForm } from "./Components/BooksForm";
import { BooksList } from "./Components/BooksList";
import { useState } from "react";

const initialBooks: Book[] = [
  { id: 1, title: "The Go Programming Language", author: "Alan Donovan", pages: 400 },
  { id: 2, title: "Clean Code", author: "Robert C. Martin", pages: 464 },
  { id: 3, title: "You Don't Know JS", author: "Kyle Simpson", pages: 278 },
  { id: 4, title: "Design Patterns", author: "Erich Gamma", pages: 395 },
];

function App() {
  const [books, setBooks] = useState<Book[]>(initialBooks);

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

export default App;
