import Book from "../Book.model";

export const BooksList = () => {
  const Books: Book[] = [
    { id: 1, title: "The Go Programming Language", author: "Alan Donovan", pages: 400 },
    { id: 2, title: "Clean Code", author: "Robert C. Martin", pages: 464 },
    { id: 3, title: "You Don't Know JS", author: "Kyle Simpson", pages: 278 },
    { id: 4, title: "Design Patterns", author: "Erich Gamma", pages: 395 },
  ];

  return (
    <section id="book-list">
      <h2>Book List</h2>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Title</th>
            <th>Author</th>
            <th>Pages</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {Books.map((book) => (
            <tr key={book.id}>
              <td>{book.id}</td>
              <td>{book.title}</td>
              <td>{book.author}</td>
              <td>{book.pages}</td>
              <td>
                <button>Edit</button>
                <button>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
};
