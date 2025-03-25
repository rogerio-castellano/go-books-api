import Book from "../Book.model";

interface Props {
  books: Book[];
  onDelete: (id: number) => void;
  onEdit: (book: Book) => void;
}

export const BooksList = ({ books, onEdit, onDelete }: Props) => {
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
          {books.map((book) => (
            <tr key={book.id}>
              <td>{book.id}</td>
              <td>{book.title}</td>
              <td>{book.author}</td>
              <td>{book.pages}</td>
              <td>
                <button className="edit" onClick={() => onEdit(book)}>
                  Edit
                </button>
                <button className="delete" onClick={() => onDelete(book.id)}>
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
};
