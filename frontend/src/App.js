import logo from "./logo.svg";
import "./App.css";

function App() {
  return (
    <div className="App">
      <main>
        <section id="book-form">
          <h2>Add a Book</h2>
          <form id="addBookForm">
            <label for="title">Title:</label>
            <input type="text" id="title" name="title" placeholder="Enter book title" required />

            <label for="author">Author:</label>
            <input type="text" id="author" name="author" placeholder="Enter author's name" required />

            <label for="pages">Pages:</label>
            <input type="number" id="pages" name="pages" placeholder="Enter page count" required />

            <button type="submit">Add Book</button>
          </form>
        </section>
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
              <tr>
                <td>1</td>
                <td>Book Title</td>
                <td>Book Author</td>
                <td>350</td>
                <td>
                  <button>Edit</button>
                  <button>Delete</button>
                </td>
              </tr>
            </tbody>
          </table>
        </section>
      </main>
    </div>
  );
}

export default App;
