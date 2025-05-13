export default function RecipeTree({ data }) {
  if (!data || data.length === 0) return null;

  return (
    <div className="p-4 border rounded">
      <h2 className="font-semibold mb-2">Recipe Steps</h2>
      <ul className="list-disc ml-6 space-y-1">
        {data.map((step, idx) => (
          <li key={idx}>{step}</li>
        ))}
      </ul>
    </div>
  );
}
