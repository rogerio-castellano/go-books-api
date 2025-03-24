import "./App.css";
import { BooksForm } from "./Components/BooksForm";
import { BooksList } from "./Components/BooksList";

function App() {
  return (
    <div className="App">
      <main>
        <BooksForm />
        <BooksList />
      </main>
    </div>
  );
}

export default App;
