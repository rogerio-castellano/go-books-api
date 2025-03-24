import "./App.css";
import Book from "./Book.model";
import { BooksForm } from "./Components/BooksForm";
import { BooksList } from "./Components/BooksList";

function App() {
  const books: Book[] = [
    { id: 1, title: "The Go Programming Language", author: "Alan Donovan", pages: 400 },
    { id: 2, title: "Clean Code", author: "Robert C. Martin", pages: 464 },
    { id: 3, title: "You Don't Know JS", author: "Kyle Simpson", pages: 278 },
    { id: 4, title: "Design Patterns", author: "Erich Gamma", pages: 395 },
  ];

  return (
    <div className="App">
      <main>
        <BooksForm />
        <BooksList books={books} />
      </main>
    </div>
  );
}

export default App;
