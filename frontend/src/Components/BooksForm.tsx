import Book from "../Book.model";
import { useForm } from "react-hook-form";

interface Props {
  onSubmit: (book: Book) => void;
}

export const BooksForm = ({ onSubmit: onSubmittedBook }: Props) => {
  const { register, handleSubmit, reset } = useForm<Book>();

  return (
    <section id="book-form">
      <h2>Add a Book</h2>
      <form
        onSubmit={handleSubmit((data) => {
          onSubmittedBook(data);
          reset();
        })}
      >
        <label htmlFor="title">Title:</label>
        <input type="text" id="title" placeholder="Enter book title" required {...register("title")} />

        <label htmlFor="author">Author:</label>
        <input type="text" id="author" placeholder="Enter author's name" required {...register("author")} />

        <label htmlFor="pages">Pages:</label>
        <input
          type="number"
          id="pages"
          placeholder="Enter page count"
          required
          {...register("pages", { valueAsNumber: true })}
        />
        <button type="submit">Add Book</button>
      </form>
    </section>
  );
};
