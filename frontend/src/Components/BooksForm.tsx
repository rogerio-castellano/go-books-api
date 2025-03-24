export const BooksForm = () => {
  return (
    <section id="book-form">
      <h2>Add a Book</h2>
      <form id="addBookForm">
        <label htmlFor="title">Title:</label>
        <input type="text" id="title" name="title" placeholder="Enter book title" required />

        <label htmlFor="author">Author:</label>
        <input type="text" id="author" name="author" placeholder="Enter author's name" required />

        <label htmlFor="pages">Pages:</label>
        <input type="number" id="pages" name="pages" placeholder="Enter page count" required />

        <button type="submit">Add Book</button>
      </form>
    </section>
  );
};
