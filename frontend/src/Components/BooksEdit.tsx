import { useForm } from "react-hook-form";
import Book from "../Book.model";

interface Props {
  book: Book;
  onSubmit: (book: Book) => void;
  onCancel: () => void;
}

export const BooksEdit = ({ book, onSubmit, onCancel }: Props) => {
  const { register, handleSubmit } = useForm<Book>();

  const handleCancel = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    onCancel();
  };

  return (
    <div>
      <h2>Edit a Book</h2>
      <form
        onSubmit={handleSubmit((data) => {
          onSubmit(data);
        })}
      >
        <input type="hidden" id="id" defaultValue={book.id} required {...register("id", { valueAsNumber: true })} />
        <label htmlFor="title">Title:</label>
        <input type="text" id="title" defaultValue={book.title} required {...register("title")} />

        <label htmlFor="author">Author:</label>
        <input type="text" id="author" defaultValue={book.author} required {...register("author")} />

        <label htmlFor="pages">Pages:</label>
        <input
          type="number"
          id="pages"
          defaultValue={book.pages}
          required
          {...register("pages", { valueAsNumber: true })}
        />
        <button type="submit">Save Book</button>
        <button onClick={handleCancel}>Cancel</button>
      </form>
    </div>
  );
};
