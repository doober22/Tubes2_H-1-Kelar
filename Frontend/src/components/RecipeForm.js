import { useState } from "react";
import { searchRecipe } from "@/utils/api"; // <<-- ini bagian fetch-nya

export default function RecipeForm({ onResult }) {
  const [target, setTarget] = useState("");
  const [method, setMethod] = useState("bfs");
  const [mode, setMode] = useState("single");
  const [limit, setLimit] = useState(3);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = await searchRecipe({ target, method, mode, limit });
    onResult(data); // lempar ke parent
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input
        type="text"
        value={target}
        onChange={(e) => setTarget(e.target.value)}
        placeholder="Target Element"
        className="border p-2 rounded w-full"
        required
      />
      <select value={method} onChange={(e) => setMethod(e.target.value)} className="border p-2 rounded w-full">
        <option value="bfs">BFS</option>
        <option value="dfs">DFS</option>
      </select>
      <select value={mode} onChange={(e) => setMode(e.target.value)} className="border p-2 rounded w-full">
        <option value="single">Single Recipe</option>
        <option value="multiple">Multiple Recipes</option>
      </select>
      {mode === "multiple" && (
        <input
          type="number"
          min="1"
          value={limit}
          onChange={(e) => setLimit(e.target.value)}
          className="border p-2 rounded w-full"
        />
      )}
      <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">Search</button>
    </form>
  );
}
